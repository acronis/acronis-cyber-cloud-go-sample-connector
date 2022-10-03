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

package logs

import (
	"context"
)

// DefaultLoggingLib is the current logging library chosen as the default logger
var DefaultLoggingLib = "logrus"

const logrusLogger = "logrus"

// RFC5424 log message levels, should be mapped by each logger implementation
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// Logger interface. This would provide a way to switch logger implementations as needed
type Logger interface {
	JSON(level int, key string, value interface{})

	Fatal(message string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Error(message string, args ...interface{})
	Errorf(format string, args ...interface{})

	Warn(message string, args ...interface{})
	Warnf(format string, args ...interface{})

	Info(message string, args ...interface{})
	Infof(format string, args ...interface{})

	Debug(message string, args ...interface{})
	Debugf(format string, args ...interface{})

	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// LogConfig holds possible configuration setting for logger
type LogConfig struct {
	LoggingLib        string `yaml:"loggingLib"` // possible value: logrus
	LogLevel          string `yaml:"logLevel"`   // possible values: trace, debug, info, warn, error, fatal, panic
	WithJSONFormatter bool   `yaml:"withJSONFormatter"`
}

// GetDefaultLogger returns default logger depending on the logging library chosen during setup
func GetDefaultLogger(ctx context.Context) Logger {
	if DefaultLoggingLib == logrusLogger {
		return getLogrusLogger(ctx)
	}

	return getLogrusLogger(ctx) // default
}
