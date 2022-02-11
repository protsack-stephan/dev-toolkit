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

// Mock logger WithFields method behaviour
func (l *loggerMock) WithFields(fields logrus.Fields) *logrus.Entry {
	args := l.Called(fields)

	return args.Get(0).(*logrus.Entry)
}

// Mock logger Printf method behaviour
func (l *loggerMock) Printf(format string, args ...interface{}) {
	l.Called(format, args)
}

// Mock logger Println method behaviour
func (l *loggerMock) Println(args ...interface{}) {
	l.Called(args)
}

// Mock logger Debug method behaviour
func (l *loggerMock) Debug(args ...interface{}) {
	l.Called(args)
}

// Mock logger Info method behaviour
func (l *loggerMock) Info(args ...interface{}) {
	l.Called(args)
}

// Mock logger Warn method behaviour
func (l *loggerMock) Warn(args ...interface{}) {
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
const logTestDebugLog = "run Debug"
const logTestInfoLog = "run Info"
const logTestWarnLog = "run Warn"
const logTestPanicLog = "run Panic"
const logTestErrorLog = "run Error"
const logTestExitCode = 1

func TestLogger(t *testing.T) {
	loggerObj := new(loggerMock)

	loggerObj.On("WithFields", logrus.Fields{"name": "field1"}).Return(new(logrus.Entry))
	loggerObj.On("Printf", logTestPrintfFormat, []interface{}{logTestPrintfLog})
	loggerObj.On("Println", []interface{}{logTestPrintlnLog})
	loggerObj.On("Debug", []interface{}{logTestDebugLog})
	loggerObj.On("Info", []interface{}{logTestInfoLog})
	loggerObj.On("Warn", []interface{}{logTestWarnLog})
	loggerObj.On("Panic", []interface{}{logTestPanicLog})
	loggerObj.On("Error", []interface{}{logTestErrorLog})
	loggerObj.On("Exit", logTestExitCode)

	SetLogger(loggerObj)

	WithFields(logrus.Fields{"name": "field1"})
	Printf(logTestPrintfFormat, logTestPrintfLog)
	Println(logTestPrintlnLog)
	Debug(logTestDebugLog)
	Info(logTestInfoLog)
	Warn(logTestWarnLog)
	Panic(logTestPanicLog)
	Error(logTestErrorLog)
	Exit(logTestExitCode)
}
