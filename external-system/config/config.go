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

// The `config` package provides service level configuration
package config

import (
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

/* Config type is used to store common configuration
   of an external system such as database configuration */
type Config struct {
	DB             *DBConfig  `yaml:"db,flow"`
	AuthSettings   AuthConfig `yaml:"authSettings"` // configs to enabling auth support
	WebUIDirectory string     `yaml:"webUIDirectory"`
}

// DBConn used to access database throughout the service
var DBConn *gorm.DB

// Used to hold database specific metadata
type DBConfig struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

// AuthConfig defines the authentication configurations
type AuthConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ClientID      string `yaml:"clientID"`
	ClientSecret  string `yaml:"clientSecret"`
	IDPAddress    string `yaml:"idpAddress"`
	RedirectURL   string `yaml:"redirectURL"`
	SessionSecret string `yaml:"sessionSecret"`
}

// Return MySQL Database metadata to connect with local db instance
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     "postgres",
			Port:     5432,
			Username: "",
			Password: "",
			Database: "external",
		},
		AuthSettings: AuthConfig{
			Enabled:      false,
			ClientID:     "",
			ClientSecret: "",
			IDPAddress:   "https://localhost:8080",
			RedirectURL:  "http://externalsystem/auth/acronis/callback",
		},
		WebUIDirectory: "web",
	}
}

func (c *Config) Validate() error {
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
	// DB User
	envVal, ok := os.LookupEnv("DB_USER")
	if ok {
		c.DB.Username = envVal
	}

	// DB Password
	envVal, ok = os.LookupEnv("DB_PASSWORD")
	if ok {
		c.DB.Password = envVal
	}

	// SSO Auth Client id
	envVal, ok = os.LookupEnv("SSO_AUTH_CLIENT_ID")
	if ok {
		c.AuthSettings.ClientID = envVal
	}

	// SSO Auth Client Secret
	envVal, ok = os.LookupEnv("SSO_AUTH_CLIENT_SECRET")
	if ok {
		c.AuthSettings.ClientSecret = envVal
	}

	// SSO Auth Session Secret
	envVal, ok = os.LookupEnv("SSO_AUTH_SESSION_SECRET")
	if ok {
		c.AuthSettings.SessionSecret = envVal
	}
}
