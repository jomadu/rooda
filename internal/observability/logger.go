package observability

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jomadu/rooda/internal/config"
)

// Logger provides structured logging with configurable levels and timestamps.
type Logger struct {
	level         config.LogLevel
	timestampFmt  config.TimestampFormat
	startTime     time.Time
	output        io.Writer
}

// NewLogger creates a new logger with the specified level and timestamp format.
func NewLogger(level config.LogLevel, timestampFmt config.TimestampFormat, startTime time.Time) *Logger {
	return &Logger{
		level:        level,
		timestampFmt: timestampFmt,
		startTime:    startTime,
		output:       os.Stderr,
	}
}

// SetOutput sets the output destination for log messages.
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// Debug logs a debug-level message.
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log(config.LogLevelDebug, message, fields)
}

// Info logs an info-level message.
func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(config.LogLevelInfo, message, fields)
}

// Warn logs a warning-level message.
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(config.LogLevelWarn, message, fields)
}

// Error logs an error-level message.
func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(config.LogLevelError, message, fields)
}

func (l *Logger) log(level config.LogLevel, message string, fields map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	var parts []string

	// Timestamp
	if ts := l.formatTimestamp(); ts != "" {
		parts = append(parts, ts)
	}

	// Level
	parts = append(parts, levelString(level))

	// Message
	parts = append(parts, message)

	// Fields
	if len(fields) > 0 {
		parts = append(parts, formatLogfmt(fields))
	}

	fmt.Fprintln(l.output, strings.Join(parts, " "))
}

func (l *Logger) shouldLog(level config.LogLevel) bool {
	return levelValue(level) >= levelValue(l.level)
}

func levelValue(level config.LogLevel) int {
	switch level {
	case config.LogLevelDebug:
		return 0
	case config.LogLevelInfo:
		return 1
	case config.LogLevelWarn:
		return 2
	case config.LogLevelError:
		return 3
	default:
		return 1
	}
}

func levelString(level config.LogLevel) string {
	switch level {
	case config.LogLevelDebug:
		return "DEBUG"
	case config.LogLevelInfo:
		return "INFO"
	case config.LogLevelWarn:
		return "WARN"
	case config.LogLevelError:
		return "ERROR"
	default:
		return "INFO"
	}
}

func (l *Logger) formatTimestamp() string {
	now := time.Now()

	switch l.timestampFmt {
	case config.TimestampTime, config.TimestampTimeMs:
		return fmt.Sprintf("[%s]", now.Format("15:04:05.000"))
	case config.TimestampRelative:
		elapsed := now.Sub(l.startTime).Seconds()
		return fmt.Sprintf("[+%.3fs]", elapsed)
	case config.TimestampISO:
		return now.Format(time.RFC3339Nano)
	case config.TimestampNone:
		return ""
	default:
		return fmt.Sprintf("[%s]", now.Format("15:04:05.000"))
	}
}

func formatLogfmt(fields map[string]interface{}) string {
	var parts []string
	for k, v := range fields {
		parts = append(parts, formatField(k, v))
	}
	return strings.Join(parts, " ")
}

func formatField(key string, value interface{}) string {
	switch v := value.(type) {
	case string:
		if needsQuoting(v) {
			return fmt.Sprintf("%s=%q", key, v)
		}
		return fmt.Sprintf("%s=%s", key, v)
	case bool:
		return fmt.Sprintf("%s=%t", key, v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%s=%d", key, v)
	case float32, float64:
		return fmt.Sprintf("%s=%v", key, v)
	default:
		return fmt.Sprintf("%s=%v", key, v)
	}
}

func needsQuoting(s string) bool {
	return strings.Contains(s, " ") || strings.Contains(s, "=") || strings.Contains(s, "\"")
}
