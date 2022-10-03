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

package updater

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/logs"
)

// Config defines the configuration structure of the sample-connector
type Config struct {
	LogSettings            logs.LogConfig  `yaml:"logSettings,flow"`       // logging config
	AuthSettings           AuthConfig      `yaml:"authSettings,flow"`      // configs to enabling auth support
	APIServerSettings      APIServerConfig `yaml:"apiServerSettings,flow"` // configs to connect to api server
	UpdateInterval         uint            `yaml:"updateInterval"`         // update interval, in seconds
	ReconciliationInterval uint            `yaml:"reconciliationInterval"` // reconciliation interval, in seconds
	UsageReportInterval    uint            `yaml:"usageReportInterval"`    // usage report interval, in seconds
}

// AuthConfig defines the authentication configurations
type AuthConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

// APIServerConfig contains the information to connect to Acronis Cyber Cloud Platform
type APIServerConfig struct {
	BaseURL string `yaml:"baseURL"`
}

// NewDefaultConfig returns the default configuration values
func NewDefaultConfig() *Config {
	return &Config{
		LogSettings: logs.LogConfig{
			LoggingLib:        "logrus",
			WithJSONFormatter: true,
			LogLevel:          "debug",
		},
		AuthSettings: AuthConfig{
			ClientID:     "",
			ClientSecret: "",
		},
		APIServerSettings: APIServerConfig{
			BaseURL: "",
		},
		UpdateInterval:         5,
		ReconciliationInterval: 86400,
		UsageReportInterval:    21600,
	}
}

// Validate function is used to validate config values
func (c *Config) Validate() error {
	switch c.LogSettings.LogLevel {
	case "trace", "debug", "info", "warn", "error", "fatal", "panic":
	default:
		return fmt.Errorf("invalid logging level: %v", c.LogSettings.LogLevel)
	}

	// Remove trailing slash
	c.APIServerSettings.BaseURL = strings.TrimRight(c.APIServerSettings.BaseURL, "/")

	if _, err := url.ParseRequestURI(c.APIServerSettings.BaseURL); err != nil {
		return fmt.Errorf("error validating API server base url: %s", c.APIServerSettings.BaseURL)
	}

	return nil
}
