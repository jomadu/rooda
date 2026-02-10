package loop

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SetupSignalHandler creates a channel for SIGINT/SIGTERM and returns it.
func SetupSignalHandler() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}

// KillProcessWithTimeout kills a process with SIGTERM, waits for termination,
// then sends SIGKILL if still running after timeout.
func KillProcessWithTimeout(process *os.Process, timeout time.Duration) error {
	if process == nil {
		return nil
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		// Process may already be dead
		return nil
	}

	// Wait for graceful termination
	done := make(chan error, 1)
	go func() {
		_, err := process.Wait()
		done <- err
	}()

	select {
	case <-done:
		// Process terminated gracefully
		return nil
	case <-time.After(timeout):
		// Timeout - send SIGKILL
		log.Printf("WARN: Process %d did not terminate within %v, sending SIGKILL", process.Pid, timeout)
		if err := process.Kill(); err != nil {
			log.Printf("WARN: Failed to SIGKILL process %d: %v", process.Pid, err)
		}
		// Wait a bit more for SIGKILL to take effect
		select {
		case <-done:
			return nil
		case <-time.After(1 * time.Second):
			log.Printf("WARN: Process %d still running after SIGKILL", process.Pid)
			return nil
		}
	}
}

// GetInterruptExitCode returns the standard exit code for SIGINT (130).
func GetInterruptExitCode() int {
	return 130
}
