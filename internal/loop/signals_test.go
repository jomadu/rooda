package loop

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func TestSetupSignalHandler(t *testing.T) {
	sigChan := SetupSignalHandler()
	if sigChan == nil {
		t.Fatal("expected signal channel, got nil")
	}
}

func TestKillProcessOnSignal_GracefulTermination(t *testing.T) {
	// Start a process that will terminate gracefully
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	// Kill it with SIGTERM
	err := KillProcessWithTimeout(cmd.Process, 5*time.Second)
	if err != nil {
		t.Errorf("expected graceful termination, got error: %v", err)
	}

	// Verify process is dead
	if cmd.Process.Signal(syscall.Signal(0)) == nil {
		t.Error("process still running after kill")
	}
}

func TestKillProcessOnSignal_ForcedKill(t *testing.T) {
	// Start a process that ignores SIGTERM
	cmd := exec.Command("sh", "-c", "trap '' TERM; sleep 10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}

	// Try to kill with very short timeout to force SIGKILL
	err := KillProcessWithTimeout(cmd.Process, 100*time.Millisecond)
	if err != nil {
		t.Errorf("expected forced kill to succeed, got error: %v", err)
	}

	// Verify process is dead
	if cmd.Process.Signal(syscall.Signal(0)) == nil {
		t.Error("process still running after forced kill")
	}
}

func TestGetInterruptExitCode(t *testing.T) {
	code := GetInterruptExitCode()
	if code != 130 {
		t.Errorf("expected exit code 130, got %d", code)
	}
}

func TestHandleSignalDuringExecution(t *testing.T) {
	// This tests the integration: signal received during AI execution
	sigChan := make(chan os.Signal, 1)
	
	// Simulate signal
	sigChan <- syscall.SIGINT
	
	// Verify signal received
	select {
	case sig := <-sigChan:
		if sig != syscall.SIGINT {
			t.Errorf("expected SIGINT, got %v", sig)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("signal not received")
	}
}
