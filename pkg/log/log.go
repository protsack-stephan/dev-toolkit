// Package log is a structured logger.
// The package can be extended with logger library API.
package log

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	SetFormatter(formatter logrus.Formatter)
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Panic(args ...interface{})
	Error(args ...interface{})
	Exit(code int)
}

var logger Logger

func init() {
	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// SetLogger sets the logger instance
func SetLogger(l Logger) {
	logger = l
}

// Printf call logrus Printf
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Println call logrus Println
func Println(args ...interface{}) {
	logger.Println(args...)
}

// Panic call logrus Panic
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Error call logrus Error
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Exit call logrus Exit
func Exit(code int) {
	logger.Exit(code)
}
