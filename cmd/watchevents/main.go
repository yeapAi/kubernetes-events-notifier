package watchevents

import (
    "time"
    "reflect"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    core "k8s.io/api/core/v1"
    "github.com/damienjacinto/internal/utils"
    "github.com/damienjacinto/pkg/eventprocessor"
    "github.com/rs/zerolog/log"
)

type Controller struct {
    k8sClient      kubernetes.Interface
    k8sFactory     informers.SharedInformerFactory
    eventAddedCh   chan *core.Event
    eventUpdatedCh chan *eventUpdateGroup
    startTime      time.Time
    stopChan       chan struct{}
    processors     []eventprocessor.EventProcessor
}

const (
    // informerSyncMinute defines how often the cache is synced from Kubernetes
    informerSyncMinute = 2
    startedEvent = "Started"
    podKind = "Pod"
)

type eventUpdateGroup struct {
    oldEvent *core.Event
    newEvent *core.Event
}

func isContainerStartedEvent(event *core.Event) bool {
    return (event.Reason == startedEvent &&
        event.InvolvedObject.Kind == podKind)
}

func isSameEventOccurrence(g *eventUpdateGroup) bool {
    return (g.oldEvent.InvolvedObject == g.newEvent.InvolvedObject &&
        g.oldEvent.Count == g.newEvent.Count)
}

func Run(k8sClient kubernetes.Interface, eventProcessors []eventprocessor.EventProcessor) *Controller {

    k8sFactory := informers.NewSharedInformerFactory(k8sClient, informerSyncMinute * time.Minute)
    stopChan := make(chan struct{})

    controller := &Controller{
        k8sClient:      k8sClient,
        k8sFactory:     k8sFactory,
        eventAddedCh:   make(chan *core.Event),
        eventUpdatedCh: make(chan *eventUpdateGroup),
        startTime:      time.Now(),
        stopChan:       stopChan,
        processors:     eventProcessors,
    }

    informer := k8sFactory.Core().V1().Events().Informer()
    utils.InstallSignalHandler(stopChan)

    informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            controller.eventAddedCh <- obj.(*core.Event)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            controller.eventUpdatedCh <- &eventUpdateGroup{
                oldEvent: oldObj.(*core.Event),
                newEvent: newObj.(*core.Event),
            }
        }})

    return controller
}

func (c *Controller) HandleEvents() error {

    c.k8sFactory.Start(c.stopChan)

    stop := make(chan struct{})
    utils.InstallSignalHandler(stop)
    c.k8sFactory.WaitForCacheSync(stop)

    for {
        select {
            case event := <-c.eventAddedCh:
                c.evaluateEvent(event)
            case eventUpdate := <-c.eventUpdatedCh:
                c.evaluateEventUpdate(eventUpdate)
            case <-c.stopChan:
                log.Info().Msg("channel stopped")
                return nil
        }
    }
}

func (c *Controller) evaluateEvent(event *core.Event) {
    log.Debug().Msgf("got event %s/%s (count: %d), reason: %s, involved object: %s",
        event.ObjectMeta.Namespace, event.ObjectMeta.Name, event.Count, event.Reason, event.InvolvedObject.Kind)
    if isContainerStartedEvent(event) {
        c.evaluatePodStatus(event)
    }
}

func (c *Controller) evaluateEventUpdate(eventUpdate  *eventUpdateGroup) {
    event := eventUpdate.newEvent
    switch {
        case (eventUpdate.oldEvent == nil):
            log.Debug().Msgf("No old event present for event %s/%s (count: %d), reason: %s, involved object: %s, skipping processing",
                event.ObjectMeta.Namespace, event.ObjectMeta.Name, event.Count, event.Reason, event.InvolvedObject.Kind)
        case reflect.DeepEqual(eventUpdate.oldEvent, eventUpdate.newEvent):
            log.Debug().Msgf("Event %s/%s (count: %d), reason: %s, involved object: %s, did not change: skipping processing",
                event.ObjectMeta.Namespace, event.ObjectMeta.Name, event.Count, event.Reason, event.InvolvedObject.Kind)
        case !isContainerStartedEvent(event):
            log.Debug().Msgf("Event %s/%s (count: %d), reason: %s, involved object: %s, is not a container started event",
                event.ObjectMeta.Namespace, event.ObjectMeta.Name, event.Count, event.Reason, event.InvolvedObject.Kind)
        case isSameEventOccurrence(eventUpdate):
            log.Debug().Msgf("Event %s/%s (count: %d), reason: %s, involved object: %s, did not change wrt. to restart count: skipping processing",
                eventUpdate.newEvent.ObjectMeta.Namespace, eventUpdate.newEvent.ObjectMeta.Name, eventUpdate.newEvent.Count, eventUpdate.newEvent.Reason, eventUpdate.newEvent.InvolvedObject.Kind)
        default:
            c.evaluatePodStatus(event)
    }
}

func (c *Controller) evaluatePodStatus(event *core.Event) {
    for _, p := range c.processors {
        log.Debug().Msgf("Processing event with %s", p.GetName())
        p.Process(c.k8sClient, event, c.startTime)
    }
}
