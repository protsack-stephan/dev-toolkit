package log

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock object for logger
type loggerMock struct {
	mock.Mock
}

// Mock logger SetFormatter method behaviour
func (l *loggerMock) SetFormatter(formatter logrus.Formatter) {
	l.Called(formatter)
}

// Mock logger Printf method behaviour
func (l *loggerMock) Printf(format string, args ...interface{}) {
	l.Called(format, args)
}

// Mock logger Println method behaviour
func (l *loggerMock) Println(args ...interface{}) {
	l.Called(args)
}

// Mock logger Panic method behaviour
func (l *loggerMock) Panic(args ...interface{}) {
	l.Called(args)
}

// Mock logger Error method behaviour
func (l *loggerMock) Error(args ...interface{}) {
	l.Called(args)
}

// Mock logger Exit method behaviour
func (l *loggerMock) Exit(code int) {
	l.Called(code)
}

const logTestPrintfFormat = "Format"
const logTestPrintfLog = "run Printf"
const logTestPrintlnLog = "run Println"
const logTestPanicLog = "run Panic"
const logTestErrorLog = "run Error"
const logTestExitCode = 1

func TestLogger(t *testing.T) {
	loggerObj := new(loggerMock)

	loggerObj.On("Printf", logTestPrintfFormat, []interface{}{logTestPrintfLog})
	loggerObj.On("Println", []interface{}{logTestPrintlnLog})
	loggerObj.On("Panic", []interface{}{logTestPanicLog})
	loggerObj.On("Error", []interface{}{logTestErrorLog})
	loggerObj.On("Exit", logTestExitCode)

	SetLogger(loggerObj)

	Printf(logTestPrintfFormat, logTestPrintfLog)
	Println(logTestPrintlnLog)
	Panic(logTestPanicLog)
	Error(logTestErrorLog)
	Exit(logTestExitCode)
}
