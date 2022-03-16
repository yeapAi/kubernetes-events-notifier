package utils

import (
    "os"
    "flag"
    "os/signal"
    "syscall"
    "github.com/rs/zerolog"
)

func InstallSignalHandler(stop chan struct{}) {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        stop <- struct{}{}
        close(stop)
    }()
}

func InitLog() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    debug := flag.Bool("debug", false, "sets log level to debug")

    flag.Parse()
    zerolog.SetGlobalLevel(zerolog.InfoLevel)

    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
}
