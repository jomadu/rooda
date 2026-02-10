package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig_Valid(t *testing.T) {
	maxIter := 5
	timeout := 300
	maxBuf := 2048
	config := &Config{
		Loop: LoopConfig{
			IterationMode:        ModeMaxIterations,
			DefaultMaxIterations: &maxIter,
			IterationTimeout:     &timeout,
			MaxOutputBuffer:      10485760,
			FailureThreshold:     3,
			LogLevel:             LogLevelInfo,
			LogTimestampFormat:   TimestampTime,
			ShowAIOutput:         false,
			AICmd:                "/bin/sh",
		},
		Procedures: map[string]Procedure{
			"test": {
				DefaultMaxIterations: &maxIter,
				IterationTimeout:     &timeout,
				MaxOutputBuffer:      &maxBuf,
			},
		},
	}

	if err := ValidateConfig(config); err != nil {
		t.Errorf("Expected valid config to pass, got error: %v", err)
	}
}

func TestValidateConfig_InvalidDefaultMaxIterations(t *testing.T) {
	zero := 0
	config := &Config{
		Loop: LoopConfig{
			DefaultMaxIterations: &zero,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for DefaultMaxIterations < 1")
	}
}

func TestValidateConfig_InvalidIterationTimeout(t *testing.T) {
	zero := 0
	config := &Config{
		Loop: LoopConfig{
			IterationTimeout: &zero,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for IterationTimeout < 1")
	}
}

func TestValidateConfig_InvalidMaxOutputBuffer(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			MaxOutputBuffer: 512,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for MaxOutputBuffer < 1024")
	}
}

func TestValidateConfig_InvalidFailureThreshold(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			FailureThreshold: 0,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for FailureThreshold < 1")
	}
}

func TestValidateConfig_InvalidLogLevel(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			LogLevel: "invalid",
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid LogLevel")
	}
}

func TestValidateConfig_InvalidTimestampFormat(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			LogTimestampFormat: "invalid",
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid LogTimestampFormat")
	}
}

func TestValidateConfig_InvalidIterationMode(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			IterationMode: "invalid",
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid IterationMode")
	}
}

func TestValidateConfig_NonExistentAICommand(t *testing.T) {
	config := &Config{
		Loop: LoopConfig{
			AICmd: "/nonexistent/binary",
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for non-existent AI command")
	}
}

func TestValidateConfig_NonExecutableAICommand(t *testing.T) {
	tmpDir := t.TempDir()
	nonExecFile := filepath.Join(tmpDir, "nonexec")
	if err := os.WriteFile(nonExecFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &Config{
		Loop: LoopConfig{
			AICmd: nonExecFile,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for non-executable AI command")
	}
}

func TestValidateConfig_ProcedureInvalidMaxIterations(t *testing.T) {
	zero := 0
	config := &Config{
		Procedures: map[string]Procedure{
			"test": {
				DefaultMaxIterations: &zero,
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for procedure DefaultMaxIterations < 1")
	}
}

func TestValidateConfig_ProcedureInvalidTimeout(t *testing.T) {
	zero := 0
	config := &Config{
		Procedures: map[string]Procedure{
			"test": {
				IterationTimeout: &zero,
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for procedure IterationTimeout < 1")
	}
}

func TestValidateConfig_ProcedureInvalidMaxBuffer(t *testing.T) {
	small := 512
	config := &Config{
		Procedures: map[string]Procedure{
			"test": {
				MaxOutputBuffer: &small,
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for procedure MaxOutputBuffer < 1024")
	}
}

func TestValidateConfig_ProcedureInvalidIterationMode(t *testing.T) {
	config := &Config{
		Procedures: map[string]Procedure{
			"test": {
				IterationMode: "invalid",
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid procedure IterationMode")
	}
}
