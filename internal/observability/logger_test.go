package observability

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/jomadu/rooda/internal/config"
)

func TestLogEvent_Format(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampTime, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Info("Test message", map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	})
	
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO level, got: %s", output)
	}
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected message, got: %s", output)
	}
	if !strings.Contains(output, "key1=value1") {
		t.Errorf("Expected key1=value1, got: %s", output)
	}
	if !strings.Contains(output, "key2=42") {
		t.Errorf("Expected key2=42, got: %s", output)
	}
}

func TestLogLevel_Filtering(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelWarn, config.TimestampTime, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Debug("debug message", nil)
	logger.Info("info message", nil)
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)
	
	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should be filtered")
	}
	if strings.Contains(output, "info message") {
		t.Error("Info message should be filtered")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be logged")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged")
	}
}

func TestTimestampFormat_Time(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampTime, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Info("test", nil)
	
	output := buf.String()
	// Should have [HH:MM:SS.mmm] format
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected bracketed timestamp, got: %s", output)
	}
}

func TestTimestampFormat_Relative(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampRelative, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	time.Sleep(10 * time.Millisecond)
	logger.Info("test", nil)
	
	output := buf.String()
	// Should have [+X.XXXs] format
	if !strings.Contains(output, "[+") || !strings.Contains(output, "s]") {
		t.Errorf("Expected relative timestamp [+X.XXXs], got: %s", output)
	}
}

func TestTimestampFormat_ISO(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampISO, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Info("test", nil)
	
	output := buf.String()
	// Should have ISO 8601 format (contains T and either timezone offset or Z)
	if !strings.Contains(output, "T") {
		t.Errorf("Expected ISO timestamp with T separator, got: %s", output)
	}
	if !strings.Contains(output, "+") && !strings.Contains(output, "-") && !strings.Contains(output, "Z") {
		t.Errorf("Expected ISO timestamp with timezone, got: %s", output)
	}
}

func TestTimestampFormat_None(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampNone, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Info("test", nil)
	
	output := buf.String()
	// Should not have timestamp
	if strings.HasPrefix(output, "[") {
		t.Errorf("Expected no timestamp, got: %s", output)
	}
	if !strings.HasPrefix(output, "INFO") {
		t.Errorf("Expected to start with INFO, got: %s", output)
	}
}

func TestLogfmtFormatting(t *testing.T) {
	startTime := time.Now()
	logger := NewLogger(config.LogLevelInfo, config.TimestampNone, startTime)
	
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	
	logger.Info("test", map[string]interface{}{
		"string":  "value with spaces",
		"number":  42,
		"boolean": true,
	})
	
	output := buf.String()
	// Multi-word values should be quoted
	if !strings.Contains(output, `string="value with spaces"`) {
		t.Errorf("Expected quoted string value, got: %s", output)
	}
	// Numbers should not be quoted
	if !strings.Contains(output, "number=42") {
		t.Errorf("Expected unquoted number, got: %s", output)
	}
	// Booleans should not be quoted
	if !strings.Contains(output, "boolean=true") {
		t.Errorf("Expected unquoted boolean, got: %s", output)
	}
}
