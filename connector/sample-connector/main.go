// Copyright (c) 2021 Acronis International GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core/updater"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/logs"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/sample-connector/external"
	extclient "github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/client"
)

func main() {
	fmt.Println("Sample Connector")

	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Init config
	config := NewDefaultConfig()
	if err := config.LoadConfigFromFile(*configFile); err != nil {
		log.Fatalf("Failed to load config with error: %v", err)
	}

	// Init logger
	logs.SetupLogrusLogger(&config.UpdaterSettings.LogSettings)
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.ContextID, "init")
	logger := logs.GetDefaultLogger(ctx)

	logger.Info("Init Complete.")

	// Setup client for external system
	extClient := extclient.NewClient(http.DefaultClient, config.ExternalSystemURL)
	externalClient := external.NewExternalSystem(extClient)

	// initialize Updater to pull information from Acronis cloud
	coreUpdater, err := updater.NewUpdater(config.UpdaterSettings, externalClient)
	if err != nil {
		log.Fatalf("Failed to initialize updater: %v", err)
	}

	// Run updater
	if err := coreUpdater.Start(); err != nil {
		logger.Errorf("Failed to start updater with error: %v", err)
		os.Exit(1)
	}

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}
