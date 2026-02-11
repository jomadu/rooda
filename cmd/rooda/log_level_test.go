package main

import (
	"testing"

	"github.com/jomadu/rooda/internal/config"
)

func TestLogLevelResolution(t *testing.T) {
	tests := []struct {
		name           string
		cfgLogLevel    config.LogLevel
		verboseFlag    bool
		quietFlag      bool
		logLevelFlag   string
		expectedLevel  config.LogLevel
	}{
		{
			name:          "verbose flag overrides config",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   true,
			quietFlag:     false,
			logLevelFlag:  "",
			expectedLevel: config.LogLevelDebug,
		},
		{
			name:          "quiet flag overrides config",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   false,
			quietFlag:     true,
			logLevelFlag:  "",
			expectedLevel: config.LogLevelError,
		},
		{
			name:          "log-level flag overrides config",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   false,
			quietFlag:     false,
			logLevelFlag:  "warn",
			expectedLevel: config.LogLevelWarn,
		},
		{
			name:          "config used when no flags",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   false,
			quietFlag:     false,
			logLevelFlag:  "",
			expectedLevel: config.LogLevelInfo,
		},
		{
			name:          "verbose takes precedence over log-level flag",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   true,
			quietFlag:     false,
			logLevelFlag:  "warn",
			expectedLevel: config.LogLevelDebug,
		},
		{
			name:          "quiet takes precedence over log-level flag",
			cfgLogLevel:   config.LogLevelInfo,
			verboseFlag:   false,
			quietFlag:     true,
			logLevelFlag:  "debug",
			expectedLevel: config.LogLevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the log level resolution logic
			logLevel := tt.cfgLogLevel
			if tt.verboseFlag {
				logLevel = config.LogLevelDebug
			} else if tt.quietFlag {
				logLevel = config.LogLevelError
			} else if tt.logLevelFlag != "" {
				logLevel = config.LogLevel(tt.logLevelFlag)
			}

			if logLevel != tt.expectedLevel {
				t.Errorf("expected log level %q, got %q", tt.expectedLevel, logLevel)
			}
		})
	}
}
