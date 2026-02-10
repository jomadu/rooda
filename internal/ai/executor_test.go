package ai

import (
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/jomadu/rooda/internal/config"
)

func TestExecuteAICLI_Success(t *testing.T) {
	aiCmd := config.AICommand{
		Command: "echo hello",
		Source:  "test",
	}
	result := ExecuteAICLI(aiCmd, "", false, nil, 1024, nil)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got: %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", result.Output)
	}
	if result.Truncated {
		t.Error("expected output not truncated")
	}
	if result.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestExecuteAICLI_NonZeroExit(t *testing.T) {
	aiCmd := config.AICommand{
		Command: "sh -c 'exit 42'",
		Source:  "test",
	}
	result := ExecuteAICLI(aiCmd, "", false, nil, 1024, nil)

	if result.Error != nil {
		t.Fatalf("expected no error for non-zero exit, got: %v", result.Error)
	}
	if result.ExitCode != 42 {
		t.Errorf("expected exit code 42, got: %d", result.ExitCode)
	}
}

func TestExecuteAICLI_Timeout(t *testing.T) {
	timeout := 1 // 1 second
	aiCmd := config.AICommand{
		Command: "sleep 10",
		Source:  "test",
	}
	result := ExecuteAICLI(aiCmd, "", false, &timeout, 1024, nil)

	if result.Error == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(result.Error.Error(), "timeout") {
		t.Errorf("expected timeout error, got: %v", result.Error)
	}
	if result.Duration < time.Second || result.Duration > 2*time.Second {
		t.Errorf("expected duration ~1s, got: %v", result.Duration)
	}
}

func TestExecuteAICLI_OutputTruncation(t *testing.T) {
	// Generate output larger than buffer
	aiCmd := config.AICommand{
		Command: "sh -c 'for i in $(seq 1 100); do echo \"line $i with some padding text\"; done'",
		Source:  "test",
	}
	maxBuffer := 100 // Small buffer to force truncation
	result := ExecuteAICLI(aiCmd, "", false, nil, maxBuffer, nil)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if !result.Truncated {
		t.Error("expected output to be truncated")
	}
	if len(result.Output) > maxBuffer {
		t.Errorf("expected output <= %d bytes, got: %d", maxBuffer, len(result.Output))
	}
	// Should keep most recent output (end of output)
	if !strings.Contains(result.Output, "100") {
		t.Error("expected truncated output to contain last line (100)")
	}
}

func TestExecuteAICLI_InvalidCommand(t *testing.T) {
	aiCmd := config.AICommand{
		Command: "nonexistent-binary-xyz",
		Source:  "test",
	}
	result := ExecuteAICLI(aiCmd, "", false, nil, 1024, nil)

	if result.Error == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestExecuteAICLI_WithPrompt(t *testing.T) {
	aiCmd := config.AICommand{
		Command: "cat", // cat reads stdin and outputs it
		Source:  "test",
	}
	prompt := "test prompt content"
	result := ExecuteAICLI(aiCmd, prompt, false, nil, 1024, nil)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if !strings.Contains(result.Output, prompt) {
		t.Errorf("expected output to contain prompt, got: %s", result.Output)
	}
}

func TestScanOutputForSignals_Success(t *testing.T) {
	output := "Some output\n<promise>SUCCESS</promise>\nMore output"
	hasSuccess, hasFailure := ScanOutputForSignals(output)

	if !hasSuccess {
		t.Error("expected SUCCESS signal to be detected")
	}
	if hasFailure {
		t.Error("expected no FAILURE signal")
	}
}

func TestScanOutputForSignals_Failure(t *testing.T) {
	output := "Some output\n<promise>FAILURE</promise>\nMore output"
	hasSuccess, hasFailure := ScanOutputForSignals(output)

	if hasSuccess {
		t.Error("expected no SUCCESS signal")
	}
	if !hasFailure {
		t.Error("expected FAILURE signal to be detected")
	}
}

func TestScanOutputForSignals_Both(t *testing.T) {
	output := "Some output\n<promise>SUCCESS</promise>\n<promise>FAILURE</promise>\n"
	hasSuccess, hasFailure := ScanOutputForSignals(output)

	if !hasSuccess {
		t.Error("expected SUCCESS signal to be detected")
	}
	if !hasFailure {
		t.Error("expected FAILURE signal to be detected")
	}
}

func TestScanOutputForSignals_None(t *testing.T) {
	output := "No signals here"
	hasSuccess, hasFailure := ScanOutputForSignals(output)

	if hasSuccess {
		t.Error("expected no SUCCESS signal")
	}
	if hasFailure {
		t.Error("expected no FAILURE signal")
	}
}

func TestExecuteAICLI_SignalInterrupt(t *testing.T) {
	aiCmd := config.AICommand{
		Command: "sleep 10",
		Source:  "test",
	}
	sigChan := make(chan os.Signal, 1)
	
	// Send signal after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()
	
	result := ExecuteAICLI(aiCmd, "", false, nil, 1024, sigChan)

	if result.Error != ErrInterrupted {
		t.Errorf("expected ErrInterrupted, got: %v", result.Error)
	}
	if result.Duration < 100*time.Millisecond {
		t.Errorf("expected duration >= 100ms, got: %v", result.Duration)
	}
}
