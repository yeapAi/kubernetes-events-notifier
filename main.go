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
    config := utils.GetFlag()

    contextInfo, _ := os.LookupEnv("CONTEXTINFO")
    eventPocessors := []eventprocessor.EventProcessor{}

    if config.Oomkilled {
        webhookURL, ok := os.LookupEnv("WEBHOOKURL")
        if !ok {
            log.Fatal().Msgf("%s not set\n", "WEBHOOKURL")
        }
        oomkilled := &oomkilled.Oomkilled{Name: "oomkill", WebhookUrl: webhookURL, AdditionnalText: contextInfo}
        eventPocessors = append(eventPocessors, oomkilled)
    }

    k8sClient, err := kubeclient.CreateClient(config.NotInCluster)
    if err != nil {
        log.Fatal().Err(err)
    }

    log.Info().Msg("Starting kubernetes watcher...")
    controller := watchevents.Run(k8sClient, eventPocessors)

    log.Info().Msg("Starting event handler...")
    controller.HandleEvents()
}
