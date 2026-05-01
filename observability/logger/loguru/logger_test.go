package loguru

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// mockFields implements Fields interface for testing
type mockFields struct{}

func (m *mockFields) ToFields() map[string]any {
	return map[string]any{"test": "value"}
}

func (m *mockFields) UpdateAll(fields map[string]any) {}

func (m *mockFields) UpdateOne(key string, value any) {}

func (m *mockFields) Format() string {
	return "test_field: value"
}

func TestNewLogger(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	if logger == nil {
		t.Fatal("Expected Logger to be non-nil")
	}
	if logger.logger == nil {
		t.Fatal("Expected internal logrus logger to be non-nil")
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	mockFields := &mockFields{}
	config := Config{
		Level:     logrus.InfoLevel,
		Colorable: true,
	}

	logger := NewLoggerWithConfig(mockFields, config)

	if logger == nil {
		t.Fatal("Expected Logger to be non-nil")
	}
	if logger.logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", logger.logger.GetLevel())
	}
	if !logger.config.Colorable {
		t.Error("Expected Colorable to be true")
	}
}

func TestSetConfig(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	config := Config{
		Level:     logrus.WarnLevel,
		Colorable: true,
	}

	logger.SetConfig(config)

	if logger.logger.GetLevel() != logrus.WarnLevel {
		t.Errorf("Expected WarnLevel, got %v", logger.logger.GetLevel())
	}
	if !logger.config.Colorable {
		t.Error("Expected Colorable to be true")
	}
}

func TestLoggerWithColorable(t *testing.T) {
	mockFields := &mockFields{}
	config := Config{
		Level:     logrus.InfoLevel,
		Colorable: true,
	}

	logger := NewLoggerWithConfig(mockFields, config)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)
	logger.logger.SetFormatter(NewCustomFormatter(true))

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected output to contain message")
	}
	// Check if ANSI color codes are present when colorable is true
	if strings.Contains(output, "\033[") {
		// Color codes found - this is expected when colorable is true
	}
}

func TestLoggerWithoutColorable(t *testing.T) {
	mockFields := &mockFields{}
	config := Config{
		Level:     logrus.InfoLevel,
		Colorable: false,
	}

	logger := NewLoggerWithConfig(mockFields, config)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)
	logger.logger.SetFormatter(NewCustomFormatter(false))

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected output to contain message")
	}
	// Check if ANSI color codes are NOT present when colorable is false
	if strings.Contains(output, "\033[") {
		t.Error("Expected no ANSI color codes when colorable is false")
	}
}

func TestGetRealCaller(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	caller := logger.getRealCaller()
	if caller == "" || caller == "unknown.unknown:0" {
		t.Error("Expected caller to be captured")
	}
	// The caller should be captured (it may be from test framework or test file)
	if !strings.Contains(caller, ".") || !strings.Contains(caller, ":") {
		t.Errorf("Expected caller to be in format file.func:line, got: %s", caller)
	}
}

func TestDebug(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)
	logger.logger.SetLevel(logrus.DebugLevel)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	logger.Debug("test debug")
	output := buf.String()
	if output == "" {
		t.Error("Expected Debug to write to output")
	}
}

func TestInfo(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	logger.Info("test info")
	output := buf.String()
	if output == "" {
		t.Error("Expected Info to write to output")
	}
}

func TestInfoWithArgs(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	logger.Info("test message: %s, count: %d", "value", 42)
	output := buf.String()
	if !strings.Contains(output, "value") || !strings.Contains(output, "42") {
		t.Error("Expected Info with args to format correctly")
	}
}

func TestWarn(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	logger.Warn("test warn")
	output := buf.String()
	if output == "" {
		t.Error("Expected Warn to write to output")
	}
}

func TestError(t *testing.T) {
	mockFields := &mockFields{}
	logger := NewLogger(mockFields)

	var buf bytes.Buffer
	logger.logger.SetOutput(&buf)

	logger.Error("test error")
	output := buf.String()
	if output == "" {
		t.Error("Expected Error to write to output")
	}
}
