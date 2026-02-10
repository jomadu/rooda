package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ValidateConfig validates the merged configuration.
// Returns error with actionable messages if validation fails.
func ValidateConfig(config *Config) error {
	if err := validateLoopConfig(&config.Loop); err != nil {
		return err
	}

	for name, proc := range config.Procedures {
		if err := validateProcedure(name, &proc); err != nil {
			return err
		}
	}

	return nil
}

func validateLoopConfig(loop *LoopConfig) error {
	// Validate DefaultMaxIterations
	if loop.DefaultMaxIterations != nil && *loop.DefaultMaxIterations < 1 {
		return fmt.Errorf("loop.default_max_iterations must be >= 1, got %d", *loop.DefaultMaxIterations)
	}

	// Validate IterationTimeout
	if loop.IterationTimeout != nil && *loop.IterationTimeout < 1 {
		return fmt.Errorf("loop.iteration_timeout must be >= 1 second, got %d", *loop.IterationTimeout)
	}

	// Validate MaxOutputBuffer
	if loop.MaxOutputBuffer < 1024 {
		return fmt.Errorf("loop.max_output_buffer must be >= 1024 bytes, got %d", loop.MaxOutputBuffer)
	}

	// Validate FailureThreshold
	if loop.FailureThreshold < 1 {
		return fmt.Errorf("loop.failure_threshold must be >= 1, got %d", loop.FailureThreshold)
	}

	// Validate LogLevel
	if err := validateLogLevel(loop.LogLevel); err != nil {
		return err
	}

	// Validate LogTimestampFormat
	if err := validateTimestampFormat(loop.LogTimestampFormat); err != nil {
		return err
	}

	// Validate IterationMode
	if err := validateIterationMode(loop.IterationMode); err != nil {
		return err
	}

	// Validate AI command if set
	if loop.AICmd != "" {
		if err := validateAICommand(loop.AICmd); err != nil {
			return fmt.Errorf("loop.ai_cmd: %w", err)
		}
	}

	return nil
}

func validateProcedure(name string, proc *Procedure) error {
	// Validate DefaultMaxIterations
	if proc.DefaultMaxIterations != nil && *proc.DefaultMaxIterations < 1 {
		return fmt.Errorf("procedure %q: default_max_iterations must be >= 1, got %d", name, *proc.DefaultMaxIterations)
	}

	// Validate IterationTimeout
	if proc.IterationTimeout != nil && *proc.IterationTimeout < 1 {
		return fmt.Errorf("procedure %q: iteration_timeout must be >= 1 second, got %d", name, *proc.IterationTimeout)
	}

	// Validate MaxOutputBuffer
	if proc.MaxOutputBuffer != nil && *proc.MaxOutputBuffer < 1024 {
		return fmt.Errorf("procedure %q: max_output_buffer must be >= 1024 bytes, got %d", name, *proc.MaxOutputBuffer)
	}

	// Validate IterationMode
	if proc.IterationMode != "" {
		if err := validateIterationMode(proc.IterationMode); err != nil {
			return fmt.Errorf("procedure %q: %w", name, err)
		}
	}

	// Validate AI command if set
	if proc.AICmd != "" {
		if err := validateAICommand(proc.AICmd); err != nil {
			return fmt.Errorf("procedure %q: ai_cmd: %w", name, err)
		}
	}

	return nil
}

func validateLogLevel(level LogLevel) error {
	switch level {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return nil
	default:
		return fmt.Errorf("invalid log_level %q, must be one of: debug, info, warn, error", level)
	}
}

func validateTimestampFormat(format TimestampFormat) error {
	switch format {
	case TimestampTime, TimestampTimeMs, TimestampRelative, TimestampISO, TimestampNone:
		return nil
	default:
		return fmt.Errorf("invalid log_timestamp_format %q, must be one of: time, time-ms, relative, iso, none", format)
	}
}

func validateIterationMode(mode IterationMode) error {
	switch mode {
	case ModeMaxIterations, ModeUnlimited:
		return nil
	case "":
		return nil // Empty is valid (means inherit)
	default:
		return fmt.Errorf("invalid iteration_mode %q, must be one of: max-iterations, unlimited", mode)
	}
}

func validateAICommand(cmd string) error {
	// Parse command to extract binary path
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	binary := parts[0]

	// Check if it's an absolute path
	if strings.HasPrefix(binary, "/") || strings.HasPrefix(binary, "~") {
		return validateBinaryPath(binary)
	}

	// Check if binary exists in PATH
	path, err := exec.LookPath(binary)
	if err != nil {
		return fmt.Errorf("command %q not found in PATH. Suggestions:\n  - Install the tool\n  - Use absolute path\n  - Add to PATH", binary)
	}

	return validateBinaryPath(path)
}

func validateBinaryPath(path string) error {
	// Expand ~ if present
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot expand ~: %w", err)
		}
		path = strings.Replace(path, "~", home, 1)
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("binary %q does not exist", path)
		}
		return fmt.Errorf("cannot access binary %q: %w", path, err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("binary %q is not a regular file", path)
	}

	// Check if it's executable
	if info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("binary %q is not executable. Suggestion: chmod +x %q", path, path)
	}

	return nil
}
