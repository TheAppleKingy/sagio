package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sagio/internal/infra/configs"
	"sagio/internal/infra/transports"

	"github.com/caarlos0/env/v11"
)

type SetupConfig struct {
	HTTP bool
	AMQP bool
}

func main() {
	if err := env.Parse(&configs.AppConfig); err != nil {
		slog.Error("Unable to parse configs", "error", err)
		os.Exit(1)
	}
	var conf SetupConfig
	var resource string
	flag.BoolVar(&conf.HTTP, "http", true, "Enable HTTP transport")
	flag.BoolVar(&conf.AMQP, "amqp", false, "Enable AMQP transport")
	flag.StringVar(&resource, "resource", "init", "Resource name for saga initialization (endpoint/queue)")
	flag.Parse()
	core := transports.NewSagioCore(resource)
	// if err != nil {
	// 	slog.Error("Unable to create sagio core", "error", err)
	// 	os.Exit(1)
	// }
	registry := transports.NewTransportRegistry(core)
	typesMap := map[transports.ConnectionProtocol]bool{
		transports.HTTP: conf.HTTP,
		transports.AMQP: conf.AMQP,
	}
	for connType, ok := range typesMap {
		if ok {
			if err := registry.Register(connType); err != nil {
				slog.Error(fmt.Sprintf("Unable to register %s transport: %v", connType, err))
				registry.Stop()
				os.Exit(1)
			}
		}
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	registry.Stop()
	slog.Info("Sagio stopped.")
}
