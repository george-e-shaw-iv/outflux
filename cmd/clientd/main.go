package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/config"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/mechanism"
	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/sync"
	"github.com/george-e-shaw-iv/outflux/internal/grpc"
	"github.com/george-e-shaw-iv/outflux/internal/log"
)

// Variable block for command line flags.
var (
	// serverInfluxToken stores a token used to authenticate to the Outflux server.
	// This token is an InfluxDB token that the Outflux server uses to ping the
	// backing InfluxDB instance on the server.
	serverInfluxToken string

	// configFile stores the location of the outflux client configuration file.
	configFile string
)

// init is where we register and parse command line flags.
func init() {
	flag.StringVar(&serverInfluxToken, "token", "", "A valid InfluxDB token for the InfluxDB instance used on the server which is used to authenticate the Outflux Client to the Outflux server when establishing a connection. This can also be passed via the OUTFLUX_TOKEN environment variable.")
	flag.StringVar(&configFile, "config", "/etc/outflux/config.yaml", "The location of the outflux config file. For information on the structure of this file, see the outflux README.")
	flag.Parse()

	// Fallback to the OUTFLUX_SERVER_INFLUX_TOKEN environment variable if -token was empty.
	if serverInfluxToken == "" {
		serverInfluxToken = os.Getenv("OUTFLUX_SERVER_INFLUX_TOKEN")
	}
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
		"file": cfg.File,
		"server": log.Fields{
			"address": cfg.Server.Address,
		},
	})

	outfluxClient, err := grpc.NewClient(ctx, cfg.Server.Address)
	if err != nil {
		logger.Error("create outflux grpc client", log.WithError(err))
		return
	}
	defer outfluxClient.Close()

	// TODO(george-e-shaw-iv): Before we get any further here we need authenticate
	// with the server using the token given. This will require adding an RPC to the
	// server to do so.

	syncRunner, err := sync.NewRunner(ctx, outfluxClient)
	if err != nil {
		log.Error("create sync runner", log.WithError(err))
		return
	}
	defer syncRunner.Close(ctx)

	// TODO(george-e-shaw-iv): Now that I'm looking at this code it's weird that I'm passing
	// a runner into the thing that it runs as opposed to passing the things that the runner
	// runs... into the runner. Fix this.
	mechanisms := mechanism.NewRunner(cfg, syncRunner)
	defer mechanisms.Close(ctx)

	mechanismErrs := make(chan error, 1)
	go func() {
		if err := mechanisms.RunAll(ctx); err != nil {
			mechanismErrs <- err
		}
	}()

	logger.Info("client running")

	// Blocking main and waiting for shutdown of the daemon.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Waiting for an osSignal or a fatal mechanism running error.
	select {
	case e := <-mechanismErrs:
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

	if err := mechanisms.Close(shutdownCtx); err != nil {
		logger.Error("error during client cleanup process", log.WithError(err))
	}
}
