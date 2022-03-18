package kubeclient

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "github.com/damienjacinto/internal/kubeconfig"
)

func getKubeconfigPath(forceRunningInCluster bool) (string) {
    config, runningInCluster, err := kubeconfig.DefaultKubeconfigPath(forceRunningInCluster)
    if err != nil {
        log.Fatal().Err(err)
    }
    log.Info().Msgf("KubeconfigPath: %s, RunningInCluster: %t", config, runningInCluster)
    return config
}

func loadKubeconfigPath(kubeconfigPath string) {
    kubeconfig.LoadKubeconfig(kubeconfigPath)
}

func CreateClient(forceNotInCluster bool) (kubernetes.Interface, error) {
    var kubeconfig *rest.Config

    kubeconfigPath := getKubeconfigPath(forceNotInCluster)

    if kubeconfigPath != "" {
        config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
        if err != nil {
            return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
        }
        kubeconfig = config
    } else {
        config, err := rest.InClusterConfig()
        if err != nil {
            return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
        }
        kubeconfig = config
    }

    loadKubeconfigPath(kubeconfigPath)
    client, err := kubernetes.NewForConfig(kubeconfig)

    if err != nil {
        return nil, fmt.Errorf("unable to create a client: %v", err)
    }

    return client, nil
}
