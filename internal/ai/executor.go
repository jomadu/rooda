package ai

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jomadu/rooda/internal/config"
	"github.com/kballard/go-shellquote"
)

var ErrTimeout = errors.New("AI CLI execution timeout")
var ErrInterrupted = errors.New("interrupted by signal")

type AIExecutionResult struct {
	Output    string
	ExitCode  int
	Duration  time.Duration
	Truncated bool
	Error     error
}

func ExecuteAICLI(aiCmd config.AICommand, prompt string, verbose bool, aiExecutionTimeout *int, maxBuffer int, sigChan <-chan os.Signal) AIExecutionResult {
	startTime := time.Now()

	parts, err := shellquote.Split(aiCmd.Command)
	if err != nil || len(parts) == 0 {
		return AIExecutionResult{
			Error:    errors.New("invalid AI command"),
			Duration: time.Since(startTime),
		}
	}

	binary := parts[0]
	args := parts[1:]

	cmd := exec.Command(binary, args...)
	cmd.Dir, _ = os.Getwd()
	cmd.Env = os.Environ()
	cmd.Stdin = strings.NewReader(prompt)

	var outputBuffer bytes.Buffer
	var outputWriter io.Writer = &outputBuffer

	if verbose {
		outputWriter = io.MultiWriter(&outputBuffer, os.Stdout)
	}

	cmd.Stdout = outputWriter
	cmd.Stderr = outputWriter

	if err := cmd.Start(); err != nil {
		return AIExecutionResult{
			Error:    err,
			Duration: time.Since(startTime),
		}
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	var waitErr error
	if aiExecutionTimeout != nil {
		select {
		case waitErr = <-done:
		case <-time.After(time.Duration(*aiExecutionTimeout) * time.Second):
			cmd.Process.Kill()
			<-done
			return AIExecutionResult{
				Output:   outputBuffer.String(),
				Duration: time.Since(startTime),
				Error:    ErrTimeout,
			}
		case <-sigChan:
			// Signal received - kill process
			cmd.Process.Kill()
			<-done
			return AIExecutionResult{
				Output:   outputBuffer.String(),
				Duration: time.Since(startTime),
				Error:    ErrInterrupted,
			}
		}
	} else {
		select {
		case waitErr = <-done:
		case <-sigChan:
			// Signal received - kill process
			cmd.Process.Kill()
			<-done
			return AIExecutionResult{
				Output:   outputBuffer.String(),
				Duration: time.Since(startTime),
				Error:    ErrInterrupted,
			}
		}
	}

	duration := time.Since(startTime)
	output := outputBuffer.String()
	truncated := false

	if len(output) > maxBuffer {
		truncated = true
		output = output[len(output)-maxBuffer:]
	}

	exitCode := 0
	if waitErr != nil {
		if exitError, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return AIExecutionResult{
				Output:    output,
				Duration:  duration,
				Truncated: truncated,
				Error:     waitErr,
			}
		}
	}

	return AIExecutionResult{
		Output:    output,
		ExitCode:  exitCode,
		Duration:  duration,
		Truncated: truncated,
		Error:     nil,
	}
}

func ScanOutputForSignals(output string) (hasSuccess bool, hasFailure bool) {
	hasSuccess = strings.Contains(output, "<promise>SUCCESS</promise>")
	hasFailure = strings.Contains(output, "<promise>FAILURE</promise>")
	return hasSuccess, hasFailure
}
