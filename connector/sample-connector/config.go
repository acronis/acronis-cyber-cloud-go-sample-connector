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
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core/updater"
)

// Config defines the configuration structure of the app
type Config struct {
	ExternalSystemURL string          `yaml:"externalSystemURL"`
	UpdaterSettings   *updater.Config `yaml:"updaterSettings,flow"`
}

// NewDefaultConfig returns the default configuration values
func NewDefaultConfig() *Config {
	return &Config{
		ExternalSystemURL: "",
		UpdaterSettings:   updater.NewDefaultConfig(),
	}
}

// Validate function is used to validate config values
func (c *Config) Validate() error {
	if _, err := url.ParseRequestURI(c.ExternalSystemURL); err != nil {
		return fmt.Errorf("error validating external system url: %s", c.ExternalSystemURL)
	}

	if err := c.UpdaterSettings.Validate(); err != nil {
		return err
	}

	return nil
}

// LoadConfigFromFile loads config file and updates values of config
func (c *Config) LoadConfigFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(c); err != nil {
		return err
	}

	c.LoadEnvVar()

	return c.Validate()
}

func (c *Config) LoadEnvVar() {
	// Load variables from env if exist
	// Auth Client id
	envVal, ok := os.LookupEnv("AUTH_CLIENT_ID")
	if ok {
		c.UpdaterSettings.AuthSettings.ClientID = envVal
	}

	// Auth Client Secret
	envVal, ok = os.LookupEnv("AUTH_CLIENT_SECRET")
	if ok {
		c.UpdaterSettings.AuthSettings.ClientSecret = envVal
	}
}
