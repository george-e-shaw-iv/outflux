package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/george-e-shaw-iv/outflux/cmd/serverd/config"
	"github.com/george-e-shaw-iv/outflux/internal/grpc"
	"github.com/george-e-shaw-iv/outflux/internal/log"
)

// configFile stores the location of the outflux server configuration file.
var configFile string

// init is where we register and parse command line flags.
func init() {
	flag.StringVar(&configFile, "config", "/etc/outflux/config.yaml", "The location of the outflux config file. For information on the structure of this file, see the outflux README.")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	exitCode := 1
	defer func() {
		os.Exit(exitCode)
	}()

	logger, err := log.NewLogger()
	if err != nil {
		fmt.Printf("error creating logger: %v\n", err)
		return
	}

	cfg, err := config.Parse(configFile)
	if err != nil {
		logger.Error("parse config", log.WithError(err))
		return
	}

	logger.Info("parsed config", log.Fields{
		"host": cfg.Host,
		"port": cfg.Port,
	})

	var outflux Outflux
	server, err := grpc.NewServer(ctx, &outflux, grpc.WithHost(cfg.Host), grpc.WithPort(cfg.Port))
	if err != nil {
		logger.Error("create grpc server", log.WithError(err))
		return
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := server.Listen(ctx); err != nil {
			serverErrors <- err
		}
	}()

	logger.Info("server listening", log.Fields{
		"address": server.Address(),
	})

	// Blocking main and waiting for shutdown of the daemon.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a fatal server error.
	select {
	case e := <-serverErrors:
		logger.Error("fatal server error occurred", log.WithError(e))
	case <-osSignals:
		signal.Reset()
		logger.Info("received shutdown signal, attempting to gracefully shutdown")

		// We're shutting down through controlled means here, set the exit code
		// to 0 preemptively.
		exitCode = 0
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := server.Close(shutdownCtx); err != nil {
		logger.Error("error during server cleanup process", log.WithError(err))
	}
}
