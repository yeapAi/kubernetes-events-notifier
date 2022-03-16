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

type Config struct {
    Debug bool
    Oomkilled bool
}

func InitLog() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func GetFlag() *Config {
    debug := flag.Bool("debug", false, "sets log level to debug")
    oomkilled := flag.Bool("oom", false, "sets events handler oomkilled")
    flag.Parse()

    config := &Config{Oomkilled: false, Debug: false}

    if *debug {
        config.Debug = true
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }

    if *oomkilled {
        config.Oomkilled = true
    }
    return config
}
