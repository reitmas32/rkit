package logger

import (
	"os"
	"strings"
	"testing"
)

func TestNewSimpleLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"fatal level", "fatal"},
		{"panic level", "panic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewSimpleLogger(tt.level)
			if logger == nil {
				t.Fatal("NewSimpleLogger() returned nil")
			}
			if logger.Level != tt.level {
				t.Errorf("NewSimpleLogger() Level = %v, want %v", logger.Level, tt.level)
			}
		})
	}
}

func TestSimpleLogger_Debug(t *testing.T) {
	tests := []struct {
		name           string
		loggerLevel    string
		message        string
		args           []any
		shouldLog      bool
	}{
		{"debug level logs debug", "debug", "test message", nil, true},
		{"info level logs debug", "info", "test message", nil, false},
		{"warn level logs debug", "warn", "test message", nil, false},
		{"error level logs debug", "error", "test message", nil, false},
		{"all level logs debug", "all", "test message", nil, true},
		{"debug with args", "debug", "test %s", []any{"arg"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Debug(tt.message, tt.args...)

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Debug() should have logged but didn't. Output: %s", output)
			}
			if !tt.shouldLog && output != "" {
				t.Errorf("Debug() should not have logged but did. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_Info(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		shouldLog   bool
	}{
		{"debug level logs info", "debug", true},
		{"info level logs info", "info", true},
		{"warn level logs info", "warn", false},
		{"error level logs info", "error", false},
		{"all level logs info", "all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Info("test message")

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Info() should have logged but didn't. Output: %s", output)
			}
			if !tt.shouldLog && output != "" {
				t.Errorf("Info() should not have logged but did. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_Warn(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		shouldLog   bool
	}{
		{"debug level logs warn", "debug", true},
		{"info level logs warn", "info", true},
		{"warn level logs warn", "warn", true},
		{"error level logs warn", "error", false},
		{"all level logs warn", "all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Warn("test message")

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Warn() should have logged but didn't. Output: %s", output)
			}
			if !tt.shouldLog && output != "" {
				t.Errorf("Warn() should not have logged but did. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_Error(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		shouldLog   bool
	}{
		{"debug level logs error", "debug", true},
		{"info level logs error", "info", true},
		{"warn level logs error", "warn", true},
		{"error level logs error", "error", true},
		{"all level logs error", "all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Error("test message")

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Error() should have logged but didn't. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_Fatal(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		shouldLog   bool
	}{
		{"debug level logs fatal", "debug", true},
		{"info level logs fatal", "info", true},
		{"warn level logs fatal", "warn", true},
		{"error level logs fatal", "error", true},
		{"fatal level logs fatal", "fatal", true},
		{"all level logs fatal", "all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Fatal("test message")

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Fatal() should have logged but didn't. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_Panic(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel string
		shouldLog   bool
	}{
		{"debug level logs panic", "debug", true},
		{"info level logs panic", "info", true},
		{"warn level logs panic", "warn", true},
		{"error level logs panic", "error", true},
		{"fatal level logs panic", "fatal", true},
		{"panic level logs panic", "panic", true},
		{"all level logs panic", "all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := NewSimpleLogger(tt.loggerLevel)
			logger.Panic("test message")

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.shouldLog && !strings.Contains(output, "test") {
				t.Errorf("Panic() should have logged but didn't. Output: %s", output)
			}
		})
	}
}

func TestSimpleLogger_WithArgs(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewSimpleLogger("debug")
	logger.Info("test %s %d", "message", 42)

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "test") {
		t.Errorf("Info() with args should have logged. Output: %s", output)
	}
}
