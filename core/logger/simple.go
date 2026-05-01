package logger

import "fmt"

type SimpleLogger struct {
	Level string
}

func NewSimpleLogger(level string) *SimpleLogger {
	return &SimpleLogger{Level: level}
}

func (l *SimpleLogger) Debug(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelDebug.String()) {
		fmt.Printf(message, args...)
	}
}

func (l *SimpleLogger) Info(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelInfo.String()) {
		fmt.Printf(message, args...)
	}
}

func (l *SimpleLogger) Warn(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelWarn.String()) {
		fmt.Printf(message, args...)
	}
}

func (l *SimpleLogger) Error(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelError.String()) {
		fmt.Printf(message, args...)
	}
}

func (l *SimpleLogger) Fatal(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelFatal.String()) {
		fmt.Printf(message, args...)
	}
}

func (l *SimpleLogger) Panic(message string, args ...any) {
	if IsLoggableLevel(l.Level, LevelPanic.String()) {
		fmt.Printf(message, args...)
	}
}
