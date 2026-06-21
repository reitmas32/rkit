package loguru

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// CustomFormatter implements logrus.Formatter with the required format:
// DATE_TIME | LEVEL | FILE.FUNCTION:LINE | [config fields] | MESSAGE
type CustomFormatter struct {
	TimestampFormat string
	Colorable       bool
}

// NewCustomFormatter creates a new CustomFormatter with default settings.
func NewCustomFormatter(colorable bool) *CustomFormatter {
	return &CustomFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		Colorable:       colorable,
	}
}

// Format formats the log entry according to the custom format.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// DATE_TIME
	timestamp := entry.Time.Format(f.TimestampFormat)
	b.WriteString(timestamp)
	b.WriteString(" | ")

	// LEVEL (with color if colorable)
	level := strings.ToUpper(entry.Level.String())
	levelColor := f.getLevelColor(entry.Level)
	f.writeColored(&b, level, levelColor, f.Colorable)
	b.WriteString(" | ")

	// FILE.FUNCTION:LINE
	callerInfo := f.getCallerInfo(entry)
	b.WriteString(callerInfo)
	b.WriteString(" | ")

	// MESSAGE (already contains the config fields from logger methods)
	b.WriteString(entry.Message)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

// getLevelColor returns the ANSI color code for a log level.
func (f *CustomFormatter) getLevelColor(level logrus.Level) string {
	if !f.Colorable {
		return ""
	}

	switch level {
	case logrus.DebugLevel:
		return "\033[36m" // Cyan
	case logrus.InfoLevel:
		return "\033[32m" // Green
	case logrus.WarnLevel:
		return "\033[33m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "\033[31m" // Red
	default:
		return ""
	}
}

// writeColored writes text with optional ANSI color codes.
func (f *CustomFormatter) writeColored(b *bytes.Buffer, text, color string, useColor bool) {
	if useColor && color != "" {
		b.WriteString(color)
		b.WriteString(text)
		b.WriteString("\033[0m") // Reset color
	} else {
		b.WriteString(text)
	}
}

// getCallerInfo extracts file, function, and line information.
// It first checks if a "caller" field was set manually (from logger methods),
// otherwise falls back to logrus's caller information.
func (f *CustomFormatter) getCallerInfo(entry *logrus.Entry) string {
	// First, try to get caller from entry's Data (set by logger methods)
	if caller, ok := entry.Data["caller"].(string); ok {
		return caller
	}

	// Fallback: try to get from entry's caller (logrus sets this when SetReportCaller(true))
	if entry.HasCaller() && entry.Caller != nil {
		file := filepath.Base(entry.Caller.File)
		function := entry.Caller.Function
		// Extract just the function name (without package path)
		if idx := strings.LastIndex(function, "."); idx >= 0 {
			function = function[idx+1:]
		}
		line := entry.Caller.Line
		return fmt.Sprintf("%s.%s:%d", file, function, line)
	}

	return "unknown.unknown:0"
}
