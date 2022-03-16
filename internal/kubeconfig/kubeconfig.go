package kubeconfig

import (
    "fmt"
    "os"
    "k8s.io/client-go/util/homedir"
    "path/filepath"
    "github.com/rs/zerolog/log"
    "k8s.io/client-go/tools/clientcmd"
    client "k8s.io/client-go/tools/clientcmd/api"
)

const serviceAccountTokenInClusterPath = "/var/run/secrets/kubernetes.io/serviceaccount"

type KubeConfig *client.Config

func DefaultKubeconfigPath(forceNotInCluster bool) (filePath string, runningInCluster bool, err error) {
    if !forceNotInCluster {
        if _, err = os.Stat(serviceAccountTokenInClusterPath); err == nil {
            return "", true, nil
        }
    }

    if filePath := os.Getenv("KUBECONFIG"); filePath != "" {
        return filePath, false, nil
    }

    homedirPath := homedir.HomeDir()
    if homedirPath == "" {
        return "", false, fmt.Errorf("failed to determine homedir path")
    }

    return filepath.Join(homedirPath, ".kube", "config"), false, nil
}

func LoadKubeconfig(configFilePath string) (KubeConfig, error) {
    log.Info().Msgf("Using kubeconfig from '%v'", configFilePath)

    kubeconfig, err := clientcmd.LoadFromFile(configFilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to load kubeconfig from '%v': %v", configFilePath, err)
    }

    return kubeconfig, nil
}
