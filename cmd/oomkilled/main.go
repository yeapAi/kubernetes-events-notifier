package oomkilled

import (
    "context"
    "time"
    "fmt"
    "k8s.io/client-go/kubernetes"
    core "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "github.com/rs/zerolog/log"
    "github.com/damienjacinto/slack"
)

const (
    TerminationReasonOOMKilled = "OOMKilled"
)

type Oomkilled struct {
    Name string
    WebhookUrl string
    AdditionnalText string
}

func (o Oomkilled) Process(k8sClient kubernetes.Interface, event *core.Event, controllerStartTime time.Time) {
    pod, err := k8sClient.CoreV1().Pods(event.InvolvedObject.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
    if err != nil {
        log.Info().Msgf("Failed to retrieve pod %s/%s, due to: %v", event.InvolvedObject.Namespace, event.InvolvedObject.Name, err)
    } else {
        for _, s := range pod.Status.ContainerStatuses {
            if s.LastTerminationState.Terminated == nil || s.LastTerminationState.Terminated.Reason != TerminationReasonOOMKilled {
                log.Info().Msgf("container %s in %s/%s was not oomkilled, event ignored", s.Name, pod.Namespace, pod.Name)
                continue
            }

            if s.LastTerminationState.Terminated.FinishedAt.Time.Before(controllerStartTime) {
                log.Info().Msgf("container '%s' in '%s/%s' was terminated before this controller started", s.Name, pod.Namespace, pod.Name)
                continue
            }

            msg := fmt.Sprintf("%s Container '%s' in '%s/%s' (%s) was OOMKilled", o.AdditionnalText, s.Name, pod.Namespace, pod.Name, s.ContainerID)
            log.Info().Msg(msg)
            o.sendSlackMessage(msg)
        }
    }
}

func (o Oomkilled) GetName() string {
    return o.Name
}

func (o Oomkilled) sendSlackMessage(msg string) {
    errSlack := slack.SendSlackNotification(o.WebhookUrl, msg)
    if errSlack != nil {
        log.Fatal().Err(errSlack)
    }
}
