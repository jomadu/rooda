package loop

import (
	"testing"
	"time"

	"github.com/jomadu/rooda/internal/config"
)

func TestRunLoop_MaxIterationsReached(t *testing.T) {
	maxIters := 2
	state := &IterationState{
		Iteration:        0,
		MaxIterations:    &maxIters,
		FailureThreshold: 3,
		Status:           StatusRunning,
		ProcedureName:    "test",
		StartedAt:        time.Now(),
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "echo 'test'", Source: "test"}

	status := RunLoop(state, cfg, aiCmd, "", false)

	if status != StatusMaxIters {
		t.Errorf("Expected status %s, got %s", StatusMaxIters, status)
	}

	if state.Iteration != 2 {
		t.Errorf("Expected 2 iterations, got %d", state.Iteration)
	}
}

func TestRunLoop_FailureThresholdExceeded(t *testing.T) {
	maxIters := 10
	state := &IterationState{
		Iteration:           0,
		MaxIterations:       &maxIters,
		ConsecutiveFailures: 0,
		FailureThreshold:    2,
		Status:              StatusRunning,
		ProcedureName:       "test",
		StartedAt:           time.Now(),
		MaxOutputBuffer:     config.DefaultMaxOutputBuffer,
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "sh -c 'exit 1'", Source: "test"}

	status := RunLoop(state, cfg, aiCmd, "", false)

	if status != StatusAborted {
		t.Errorf("Expected status %s, got %s", StatusAborted, status)
	}

	if state.ConsecutiveFailures < state.FailureThreshold {
		t.Errorf("Expected consecutive failures >= %d, got %d", state.FailureThreshold, state.ConsecutiveFailures)
	}
}

func TestRunLoop_SuccessSignal(t *testing.T) {
	maxIters := 10
	state := &IterationState{
		Iteration:        0,
		MaxIterations:    &maxIters,
		FailureThreshold: 3,
		Status:           StatusRunning,
		ProcedureName:    "test",
		StartedAt:        time.Now(),
		MaxOutputBuffer:  config.DefaultMaxOutputBuffer,
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "echo '<promise>SUCCESS</promise>'", Source: "test"}

	status := RunLoop(state, cfg, aiCmd, "", false)

	if status != StatusSuccess {
		t.Errorf("Expected status %s, got %s", StatusSuccess, status)
	}

	if state.Iteration != 1 {
		t.Errorf("Expected 1 iteration, got %d", state.Iteration)
	}
}

func TestRunLoop_FailureSignalIncrementsCounter(t *testing.T) {
	maxIters := 5
	state := &IterationState{
		Iteration:           0,
		MaxIterations:       &maxIters,
		ConsecutiveFailures: 0,
		FailureThreshold:    5,
		Status:              StatusRunning,
		ProcedureName:       "test",
		StartedAt:           time.Now(),
		MaxOutputBuffer:     config.DefaultMaxOutputBuffer,
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "echo '<promise>FAILURE</promise>'", Source: "test"}

	status := RunLoop(state, cfg, aiCmd, "", false)

	if status != StatusMaxIters {
		t.Errorf("Expected status %s, got %s", StatusMaxIters, status)
	}

	if state.ConsecutiveFailures != 5 {
		t.Errorf("Expected 5 consecutive failures, got %d", state.ConsecutiveFailures)
	}
}

func TestRunLoop_SuccessResetsFailureCounter(t *testing.T) {
	maxIters := 3
	state := &IterationState{
		Iteration:           0,
		MaxIterations:       &maxIters,
		ConsecutiveFailures: 2,
		FailureThreshold:    5,
		Status:              StatusRunning,
		ProcedureName:       "test",
		StartedAt:           time.Now(),
		MaxOutputBuffer:     config.DefaultMaxOutputBuffer,
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "echo 'success'", Source: "test"}

	status := RunLoop(state, cfg, aiCmd, "", false)

	if status != StatusMaxIters {
		t.Errorf("Expected status %s, got %s", StatusMaxIters, status)
	}

	if state.ConsecutiveFailures != 0 {
		t.Errorf("Expected consecutive failures reset to 0, got %d", state.ConsecutiveFailures)
	}
}

func TestRunLoop_StatisticsUpdated(t *testing.T) {
	maxIters := 2
	state := &IterationState{
		Iteration:        0,
		MaxIterations:    &maxIters,
		FailureThreshold: 3,
		Status:           StatusRunning,
		ProcedureName:    "test",
		StartedAt:        time.Now(),
		MaxOutputBuffer:  config.DefaultMaxOutputBuffer,
	}

	cfg := config.Config{
		Procedures: map[string]config.Procedure{
			"test": {
				Observe: []config.FragmentAction{{Content: "observe"}},
				Orient:  []config.FragmentAction{{Content: "orient"}},
				Decide:  []config.FragmentAction{{Content: "decide"}},
				Act:     []config.FragmentAction{{Content: "act"}},
			},
		},
	}

	aiCmd := config.AICommand{Command: "echo 'test'", Source: "test"}

	RunLoop(state, cfg, aiCmd, "", false)

	if state.Stats.Count != 2 {
		t.Errorf("Expected 2 iterations in stats, got %d", state.Stats.Count)
	}

	if state.Stats.TotalTime == 0 {
		t.Errorf("Expected non-zero total time")
	}

	if state.Stats.MinTime == 0 {
		t.Errorf("Expected non-zero min time")
	}

	if state.Stats.MaxTime == 0 {
		t.Errorf("Expected non-zero max time")
	}
}
