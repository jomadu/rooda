package loop

import (
	"fmt"
	"time"

	"github.com/jomadu/rooda/internal/ai"
	"github.com/jomadu/rooda/internal/config"
	"github.com/jomadu/rooda/internal/observability"
	"github.com/jomadu/rooda/internal/prompt"
)

// RunLoop executes the OODA iteration loop until a termination condition is met.
// Returns the final loop status (success, max-iters, aborted, interrupted).
func RunLoop(state *IterationState, cfg config.Config, aiCmd config.AICommand, userContext string, verbose bool, logger *observability.Logger) LoopStatus {
	// Setup signal handler
	sigChan := SetupSignalHandler()

	procedure, ok := cfg.Procedures[state.ProcedureName]
	if !ok {
		logger.Error(fmt.Sprintf("Procedure %s not found", state.ProcedureName), nil)
		state.Status = StatusAborted
		return StatusAborted
	}

	// Log loop start
	maxIters := "unlimited"
	if state.MaxIterations != nil {
		maxIters = fmt.Sprintf("%d", *state.MaxIterations)
	}
	logger.Info("Starting loop", map[string]interface{}{
		"procedure":       state.ProcedureName,
		"max_iterations":  maxIters,
	})

	for {
		// Check termination: max iterations
		if state.MaxIterations != nil && state.Iteration >= *state.MaxIterations {
			state.Status = StatusMaxIters
			break
		}

		// Check termination: failure threshold
		if state.ConsecutiveFailures >= state.FailureThreshold {
			logger.Error("Aborting: consecutive failures exceeded threshold", map[string]interface{}{
				"consecutive_failures": state.ConsecutiveFailures,
				"threshold":           state.FailureThreshold,
			})
			state.Status = StatusAborted
			break
		}

		// Start iteration
		iterationStart := time.Now()
		iterNum := state.Iteration + 1
		maxItersDisplay := "unlimited"
		if state.MaxIterations != nil {
			maxItersDisplay = fmt.Sprintf("%d", *state.MaxIterations)
		}
		logger.Info(fmt.Sprintf("Starting iteration %d/%s", iterNum, maxItersDisplay), map[string]interface{}{
			"procedure": state.ProcedureName,
		})

		// Assemble prompt
		assembledPrompt, err := prompt.AssemblePrompt(procedure, userContext, "")
		if err != nil {
			logger.Error("Prompt assembly failed", map[string]interface{}{
				"error": err.Error(),
			})
			state.Status = StatusAborted
			break
		}

		// Execute AI CLI
		result := ai.ExecuteAICLI(aiCmd, assembledPrompt, verbose, state.IterationTimeout, state.MaxOutputBuffer, sigChan)

		// Handle interrupt
		if result.Error == ai.ErrInterrupted {
			logger.Info("Interrupted by signal", nil)
			state.Status = StatusInterrupted
			break
		}

		// Handle timeout
		if result.Error == ai.ErrTimeout {
			logger.Warn(fmt.Sprintf("Iteration %d: AI CLI exceeded timeout", iterNum), map[string]interface{}{
				"timeout": fmt.Sprintf("%ds", *state.IterationTimeout),
			})
			state.ConsecutiveFailures++
			elapsed := time.Since(iterationStart)
			state.Stats.updateStats(elapsed)
			state.Iteration++
			continue
		}

		// Handle execution error
		if result.Error != nil {
			logger.Error("AI CLI execution failed", map[string]interface{}{
				"error": result.Error.Error(),
			})
			state.Status = StatusAborted
			break
		}

		// Determine outcome per matrix
		outcome := DetectIterationFailure(IterationResult{
			ExitCode: result.ExitCode,
			Output:   result.Output,
		})

		// Scan for signals (for logging)
		hasSuccess, hasFailure := ai.ScanOutputForSignals(result.Output)
		_ = hasSuccess // Used for logging context

		elapsed := time.Since(iterationStart)

		switch outcome {
		case OutcomeJobDone:
			// SUCCESS signal - terminate loop
			logger.Info(fmt.Sprintf("Iteration %d: AI signaled SUCCESS", iterNum), nil)
			state.Status = StatusSuccess
			state.Stats.updateStats(elapsed)
			logger.Info(fmt.Sprintf("Completed iteration %d/%s", iterNum, maxItersDisplay), map[string]interface{}{
				"elapsed": formatDuration(elapsed),
				"status":  "success",
			})
			state.Iteration++
			
			// Log loop completion
			totalElapsed := time.Since(state.StartedAt)
			logger.Info("Loop completed", map[string]interface{}{
				"status":        string(StatusSuccess),
				"iterations":    state.Iteration,
				"total_elapsed": formatDuration(totalElapsed),
			})
			
			// Display statistics
			logIterationStats(logger, &state.Stats)
			
			return StatusSuccess

		case OutcomeFailure:
			// FAILURE signal or non-zero exit - increment failures
			state.ConsecutiveFailures++
			if hasFailure {
				logger.Warn(fmt.Sprintf("Iteration %d: AI signaled FAILURE", iterNum), map[string]interface{}{
					"consecutive": state.ConsecutiveFailures,
				})
			} else {
				logger.Warn(fmt.Sprintf("Iteration %d failed with exit code %d", iterNum, result.ExitCode), map[string]interface{}{
					"consecutive": state.ConsecutiveFailures,
				})
			}

		case OutcomeSuccess:
			// Exit 0, no signal - reset failures
			state.ConsecutiveFailures = 0
			logger.Info(fmt.Sprintf("Iteration %d succeeded", iterNum), nil)
		}

		// Record timing
		state.Stats.updateStats(elapsed)
		logger.Info(fmt.Sprintf("Completed iteration %d/%s", iterNum, maxItersDisplay), map[string]interface{}{
			"elapsed": formatDuration(elapsed),
			"status":  string(outcome),
		})

		// Increment iteration counter
		state.Iteration++
	}

	// Log loop completion
	totalElapsed := time.Since(state.StartedAt)
	logger.Info("Loop completed", map[string]interface{}{
		"status":        string(state.Status),
		"iterations":    state.Iteration,
		"total_elapsed": formatDuration(totalElapsed),
	})

	// Display statistics if iterations completed
	logIterationStats(logger, &state.Stats)

	return state.Status
}

// logIterationStats displays iteration timing statistics
func logIterationStats(logger *observability.Logger, stats *IterationStats) {
	if stats.Count == 0 {
		return
	}

	mean := stats.getMean()
	fields := map[string]interface{}{
		"count": stats.Count,
		"min":   formatDuration(stats.MinTime),
		"max":   formatDuration(stats.MaxTime),
		"mean":  formatDuration(mean),
	}

	if stats.Count >= 2 {
		stddev := stats.getStdDev()
		fields["stddev"] = formatDuration(stddev)
	}

	logger.Info("Iteration timing:", fields)
}

// formatDuration formats a duration in human-readable format (e.g., "1.23s", "2m 15s")
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", d.Seconds()*1000)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

// GetExitCode returns the appropriate exit code for a loop status.
func GetExitCode(status LoopStatus) int {
	switch status {
	case StatusSuccess:
		return 0
	case StatusAborted:
		return 1
	case StatusMaxIters:
		return 2
	case StatusInterrupted:
		return 130
	default:
		return 1
	}
}

// FormatLoopSummary returns a human-readable summary of loop completion.
func FormatLoopSummary(state *IterationState) string {
	elapsed := time.Since(state.StartedAt)
	return fmt.Sprintf("Loop completed status=%s iterations=%d total_elapsed=%v",
		state.Status, state.Iteration, elapsed)
}
