package loop

import (
	"fmt"
	"log"
	"time"

	"github.com/jomadu/rooda/internal/ai"
	"github.com/jomadu/rooda/internal/config"
	"github.com/jomadu/rooda/internal/prompt"
)

// RunLoop executes the OODA iteration loop until a termination condition is met.
// Returns the final loop status (success, max-iters, aborted, interrupted).
func RunLoop(state *IterationState, cfg config.Config, aiCmd config.AICommand, userContext string, verbose bool) LoopStatus {
	// Setup signal handler
	sigChan := SetupSignalHandler()

	procedure, ok := cfg.Procedures[state.ProcedureName]
	if !ok {
		log.Printf("ERROR: Procedure %s not found", state.ProcedureName)
		state.Status = StatusAborted
		return StatusAborted
	}

	for {
		// Check termination: max iterations
		if state.MaxIterations != nil && state.Iteration >= *state.MaxIterations {
			state.Status = StatusMaxIters
			break
		}

		// Check termination: failure threshold
		if state.ConsecutiveFailures >= state.FailureThreshold {
			log.Printf("ERROR: Aborting: %d consecutive failures", state.ConsecutiveFailures)
			state.Status = StatusAborted
			break
		}

		// Start iteration
		iterationStart := time.Now()
		log.Printf("INFO: Starting iteration %d", state.Iteration+1)

		// Assemble prompt
		assembledPrompt, err := prompt.AssemblePrompt(procedure, userContext, "")
		if err != nil {
			log.Printf("ERROR: Prompt assembly failed: %v", err)
			state.Status = StatusAborted
			break
		}

		// Execute AI CLI
		result := ai.ExecuteAICLI(aiCmd, assembledPrompt, verbose, state.IterationTimeout, state.MaxOutputBuffer, sigChan)

		// Handle interrupt
		if result.Error == ai.ErrInterrupted {
			log.Printf("INFO: Interrupted by signal")
			state.Status = StatusInterrupted
			break
		}

		// Handle timeout
		if result.Error == ai.ErrTimeout {
			log.Printf("WARN: Iteration %d: AI CLI exceeded timeout (%ds)", state.Iteration+1, *state.IterationTimeout)
			state.ConsecutiveFailures++
			elapsed := time.Since(iterationStart)
			state.Stats.updateStats(elapsed)
			state.Iteration++
			continue
		}

		// Handle execution error
		if result.Error != nil {
			log.Printf("ERROR: AI CLI execution failed: %v", result.Error)
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
			log.Printf("INFO: Iteration %d: AI signaled SUCCESS", state.Iteration+1)
			state.Status = StatusSuccess
			state.Stats.updateStats(elapsed)
			log.Printf("INFO: Iteration %d completed in %v (SUCCESS)", state.Iteration+1, elapsed)
			state.Iteration++
			return StatusSuccess

		case OutcomeFailure:
			// FAILURE signal or non-zero exit - increment failures
			state.ConsecutiveFailures++
			if hasFailure {
				log.Printf("WARN: Iteration %d: AI signaled FAILURE (consecutive: %d)", state.Iteration+1, state.ConsecutiveFailures)
			} else {
				log.Printf("WARN: Iteration %d failed with exit code %d (consecutive: %d)", state.Iteration+1, result.ExitCode, state.ConsecutiveFailures)
			}

		case OutcomeSuccess:
			// Exit 0, no signal - reset failures
			state.ConsecutiveFailures = 0
			log.Printf("INFO: Iteration %d succeeded", state.Iteration+1)
		}

		// Record timing
		state.Stats.updateStats(elapsed)
		log.Printf("INFO: Iteration %d completed in %v", state.Iteration+1, elapsed)

		// Increment iteration counter
		state.Iteration++
	}

	// Display statistics if iterations completed
	if state.Stats.Count > 0 {
		mean := state.Stats.getMean()
		if state.Stats.Count == 1 {
			log.Printf("INFO: Iteration timing: count=%d min=%v max=%v mean=%v",
				state.Stats.Count, state.Stats.MinTime, state.Stats.MaxTime, mean)
		} else {
			stddev := state.Stats.getStdDev()
			log.Printf("INFO: Iteration timing: count=%d min=%v max=%v mean=%v stddev=%v",
				state.Stats.Count, state.Stats.MinTime, state.Stats.MaxTime, mean, stddev)
		}
	}

	return state.Status
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
