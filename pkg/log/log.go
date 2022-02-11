// Package log is a structured logger.
// The package can be extended with logger library API.
package log

import (
	"github.com/sirupsen/logrus"
)

type WithFielder interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Formatter interface {
	SetFormatter(formatter logrus.Formatter)
}

type Printer interface {
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Panic(args ...interface{})
	Error(args ...interface{})
}

type Logger interface {
	WithFielder
	Formatter
	Printer
}

var logger Logger

func init() {
	logger = logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
}

// SetLogger sets the logger instance
func SetLogger(l Logger) {
	logger = l
}

// WithFields call logrus WithFields
func WithFields(fields map[string]interface{}) Printer {
	return logger.WithFields(fields)
}

// Printf call logrus Printf
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Println call logrus Println
func Println(args ...interface{}) {
	logger.Println(args...)
}

// Debug call logrus Debug
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info call logrus Info
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn call logrus Warn
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Panic call logrus Panic
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Error call logrus Error
func Error(args ...interface{}) {
	logger.Error(args...)
}
