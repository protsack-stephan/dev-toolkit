package log

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetReportCaller(true)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Println(args ...interface{}) {
	logger.Println(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Exit(code int) {
	logger.Error(code)
}
