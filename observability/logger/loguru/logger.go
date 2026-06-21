package loguru

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level     logrus.Level
	Colorable bool
}

// Logger is a logrus-based implementation of logger.ILogger.
// It provides structured logging with support for hooks and custom formatters.
// It uses Config to manage context fields (trace_id, caller_id, method, client_ip, user_id)
// that are automatically included in all log entries.
type Logger struct {
	logger *logrus.Logger
	fields Fields
	config Config
	mu     sync.RWMutex
}

// NewLogger creates a new Logger with the provided Fields.
// Uses default Config (Debug level, no colors).
func NewLogger(fields Fields) *Logger {
	return NewLoggerWithConfig(fields, Config{
		Level:     logrus.DebugLevel,
		Colorable: false,
	})
}

// NewLoggerWithConfig creates a new Logger with the provided Fields and Config.
func NewLoggerWithConfig(fields Fields, config Config) *Logger {
	l := logrus.New()
	l.SetLevel(config.Level)
	l.SetReportCaller(true) // Enable caller information

	// Configure custom formatter with the required format and colorable setting
	l.SetFormatter(NewCustomFormatter(config.Colorable))

	return &Logger{
		logger: l,
		fields: fields,
		config: config,
	}
}

// SetConfig updates the logger's configuration (level and colorable).
func (l *Logger) SetConfig(config Config) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.config = config
	l.logger.SetLevel(config.Level)
	l.logger.SetFormatter(NewCustomFormatter(config.Colorable))
}

// getRealCaller captures the actual caller by skipping wrapper functions.
// Skip 2 to skip: runtime.Caller -> getRealCaller -> Logger.Info/Debug/etc -> user code
func (l *Logger) getRealCaller() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown.unknown:0"
	}

	file = filepath.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	// Extract just the function name (without package path)
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		funcName = funcName[idx+1:]
	}

	return fmt.Sprintf("%s.%s:%d", file, funcName, line)
}

func (l *Logger) Debug(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Debug(msg)
	} else {
		entry.Debug(message)
	}
}

func (l *Logger) Info(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Info(msg)
	} else {
		entry.Info(message)
	}
}

func (l *Logger) Warn(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Warn(msg)
	} else {
		entry.Warn(message)
	}
}

func (l *Logger) Error(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Error(msg)
	} else {
		entry.Error(message)
	}
}

func (l *Logger) Fatal(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Fatal(msg)
	} else {
		entry.Fatal(message)
	}
}

func (l *Logger) Panic(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := l.logger.WithField("caller", l.getRealCaller())

	if l.fields != nil {
		// Add fields to entry for hooks (e.g., Loki labels)
		if fieldsMap := l.fields.ToFields(); len(fieldsMap) > 0 {
			entry = entry.WithFields(fieldsMap)
		}
		// Format fields in message
		fields := l.fields.Format()
		msg := fmt.Sprintf("%s | %s", fields, message)
		entry.Panic(msg)
	} else {
		entry.Panic(message)
	}
}

// AddHook adds a hook to the underlying logrus logger.
// Hooks can be used to send logs to external services (Loki, Sentry, etc.).
func (l *Logger) AddHook(hook logrus.Hook) {
	l.logger.AddHook(hook)
}

// GetLogrusLogger returns the underlying logrus.Logger instance.
// This allows access to advanced logrus features if needed.
func (l *Logger) GetLogrusLogger() *logrus.Logger {
	return l.logger
}
