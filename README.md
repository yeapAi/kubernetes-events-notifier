# kubernetes-events-notifier

[![Build & Release](https://github.com/damienjacinto/kubernetes-events-notifier/actions/workflows/build_push.yml/badge.svg?branch=master)](https://github.com/damienjacinto/kubernetes-events-notifier/actions/workflows/build_push.yml)

Anlayze Kuebernetes event at cluster level, for now only events targeting pod are handle.
Every pod events will be processed with every processor activated.

## Design

Controller listens to Kubernetes API for Events and changes. For eah event received, it checks whether this event is a "ContainerStarted" event, based on the event and the kind of the involved object all the process are called on this event.

## Processor

Each event filtered by the controller are sent to active processor.

- oomkilled (Send a slack event on a oomkilled pod event)

## Usage

    Usage:
      kubernetes-event-notifier [OPTIONS]

    Application Options:
      -debug    active debug logging level

    Processor
      -oom      active oomkill notifier (needs WEBHOOKURL env vars)

    Help Options:
      -h, --help     Show this help message

## Processor usage

### oomkilled

Send a slack message with webhook (WEBHOOKURL env vars mandatory) on an event for a pod resttart with his last status as OomKilled.
You can add text before the detail of the pod oomkilled with the environment variable CONTEXTINFO.
