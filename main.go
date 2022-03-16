package main

import (
    "os"
    "github.com/rs/zerolog/log"
    "github.com/damienjacinto/internal/kubeclient"
    "github.com/damienjacinto/internal/utils"
    "github.com/damienjacinto/pkg/eventprocessor"
    "github.com/damienjacinto/watchevents"
    "github.com/damienjacinto/oomkilled"
)

func main() {

    utils.InitLog()

    webhookURL, ok := os.LookupEnv("WEBHOOKURL")
    if !ok {
        log.Fatal().Msgf("%s not set\n", "WEBHOOKURL")
    }

    k8sClient, err := kubeclient.CreateClient()
    if err != nil {
        log.Fatal().Err(err)
    }

    log.Info().Msg("Starting kubernetes watcher...")
    contextInfo := "Cluster : sandbox, env : eu\n"
    oomkilled := &oomkilled.Oomkilled{Name: "oomkill", WebhookUrl: webhookURL, AdditionnalText: contextInfo}
    controller := watchevents.Run(k8sClient, []eventprocessor.EventProcessor{oomkilled})

    log.Info().Msg("Starting event handler...")
    controller.HandleEvents()
}
