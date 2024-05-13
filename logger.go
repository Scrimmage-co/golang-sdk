package scrimmage

import "github.com/sirupsen/logrus"

type LogLevel int

const (
	LogLevel_Debug LogLevel = iota + 1
	LogLevel_Info
	LogLevel_Log
	LogLevel_Warn
	LogLevel_Error
)

type Logger interface {
	Log(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

type defaultLogger struct {
}

func newDefaultLogger() Logger {
	return &defaultLogger{}
}

func (d *defaultLogger) Log(args ...interface{}) {
	logrus.Info(args...)
}

func (d *defaultLogger) Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func (d *defaultLogger) Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func (d *defaultLogger) Info(args ...interface{}) {
	logrus.Info(args...)
}

func (d *defaultLogger) Error(args ...interface{}) {
	logrus.Error(args...)
}
