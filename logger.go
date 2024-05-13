package main

import "context"

type LogLevel int

const (
	LogLevel_Log   LogLevel = iota + 1
	LogLevel_Debug LogLevel = iota + 1
	LogLevel_Warn  LogLevel = iota + 1
	LogLevel_Info  LogLevel = iota + 1
	LogLevel_Error LogLevel = iota + 1
)

type Logger interface {
	Log(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
}
