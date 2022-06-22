package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/george-e-shaw-iv/outflux/cmd/clientd/internal/config"
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

	// Fallback to the OUTFLUX_TOKEN environment variable if -token was empty.
	if serverInfluxToken == "" {
		serverInfluxToken = os.Getenv("OUTFLUX_TOKEN")
	}
}

func main() {
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
		logger.Error("parse config", log.Fields{
			"error": err,
		})
		return
	}

	logger.Info("parsed config", log.Fields{
		"file": cfg.File,
		"server": log.Fields{
			"address": cfg.Server.Address,
		},
	})

	exitCode = 0
}
