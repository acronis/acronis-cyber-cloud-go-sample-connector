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
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	// ContextID is used by logrus to maintain session logging
	ContextID contextKey = "contextID"
)

// StackFrameSkip defines number of stacks to ascend when printing log line
const StackFrameSkip = 2

// LogrusLogger is the implementation of Logger interface that utilizes logrus library
type LogrusLogger struct {
	Logger *logrus.Logger
	Ctx    context.Context
}

var levelToLogrusLevel = [LevelDebug + 1]logrus.Level{
	logrus.FatalLevel,
	logrus.FatalLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

// SetupLogrusLogger initializes LogrusLogger as the default logger the application.
// It applies config as the log configuration.
func SetupLogrusLogger(config *LogConfig) {
	if config.WithJSONFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "@timestamp",
			},
		})
	}
	logrus.SetOutput(os.Stdout)
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatalf("Failed to parse level: %v", err)
	}
	logrus.SetLevel(logLevel)
	DefaultLoggingLib = logrusLogger
}

// getLogrusLogger returns a logrus object
func getLogrusLogger(ctx context.Context) Logger {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
	}).Logger
	return &LogrusLogger{
		Logger: logger,
		Ctx:    ctx,
	}
}

// LoggerDetails prepares some common fields that will be logged in every log.
// It will take contextID value from context
func (a *LogrusLogger) LoggerDetails(stackSkip int) *logrus.Entry {
	pc, filePath, line, ok := runtime.Caller(stackSkip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		filePathParts := strings.Split(filePath, "/")
		return a.Logger.WithFields(logrus.Fields{
			"ctx":       a.Ctx.Value(ContextID),
			"func":      details.Name(),
			"file_name": filePathParts[len(filePathParts)-1],
			"line":      line,
		})
	}
	return a.Logger.WithFields(logrus.Fields{
		"log_details": "undefined",
	})
}

// JSON performs logging with JSON format.
func (a *LogrusLogger) JSON(level int, key string, value interface{}) {
	switch v := value.(type) {
	case string:
		a.LoggerDetails(StackFrameSkip).WithField(key, v).Log(levelToLogrusLevel[level], v)
	default:
		a.LoggerDetails(StackFrameSkip).WithField(key, v).Log(levelToLogrusLevel[level], "json")
	}
}

// Fatal logs the message with Fatal severity level
func (a *LogrusLogger) Fatal(message string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Fatal(message)
}

// Fatalf logs the message that contains placeholder with Fatal severity level
func (a *LogrusLogger) Fatalf(format string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Fatalf(format, args...)
}

// Error logs the message with Error severity level
func (a *LogrusLogger) Error(message string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Error(message)
}

// Errorf logs the message that contains placeholder with Error severity level
func (a *LogrusLogger) Errorf(format string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Errorf(format, args...)
}

// Warn logs the message with Warning severity level
func (a *LogrusLogger) Warn(message string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Warn(message)
}

// Warn logs the message that contains placeholder with Warning severity level
func (a *LogrusLogger) Warnf(format string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Warnf(format, args...)
}

// Info logs the message with Info severity level
func (a *LogrusLogger) Info(message string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Info(message)
}

// Infof logs the message that contains placeholder with Info severity level
func (a *LogrusLogger) Infof(format string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Infof(format, args...)
}

// Debug logs the message with Debug severity level
func (a *LogrusLogger) Debug(message string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Debug(message)
}

// Debugf logs the message that contains placeholder with Debug severity level
func (a *LogrusLogger) Debugf(format string, args ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Debugf(format, args...)
}

// Println logs the message with Info severity level followed by newline
func (a *LogrusLogger) Println(v ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Println(v...)
}

// Printf logs the message that contains placeholder with Info severity level
func (a *LogrusLogger) Printf(format string, v ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Printf(format, v...)
}

// Print logs the message with Info severity level followed by newline
func (a *LogrusLogger) Print(v ...interface{}) {
	a.LoggerDetails(StackFrameSkip).Print(v...)
}
