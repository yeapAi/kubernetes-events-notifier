package eventprocessor

import (
    core "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
    "time"
)

type EventProcessor interface {
    GetName() string
    Process(k8sClient kubernetes.Interface, event *core.Event, controllerStartTime time.Time)
}
