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

// The `main` package triggers external system service
package main

import (
	"flag"
	"log"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/server"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/storage"
)

// Setup DB configuration, initiate routes and run the server
func main() {
	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Init config
	appConfig := config.GetConfig()
	if err := appConfig.LoadConfigFromFile(*configFile); err != nil {
		log.Fatalf("Failed to load config with error: %v", err)
	}

	app := &server.App{}
	storage.SetupDB(appConfig)
	if err := app.InitializeRoutes(appConfig); err != nil {
		log.Fatalf("Failed to initialize routes with error: %v", err)
	}
	app.Run(":8080")
}
