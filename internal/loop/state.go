package loop

import (
	"math"
	"time"
)

// LoopStatus represents the current state of the iteration loop
type LoopStatus string

const (
	StatusRunning     LoopStatus = "running"
	StatusSuccess     LoopStatus = "success"     // AI signaled SUCCESS
	StatusMaxIters    LoopStatus = "max-iters"   // Max iterations reached
	StatusAborted     LoopStatus = "aborted"     // Failure threshold exceeded
	StatusInterrupted LoopStatus = "interrupted" // User pressed Ctrl+C (SIGINT/SIGTERM)
)

// IterationState tracks the state of the iteration loop
type IterationState struct {
	Iteration           int            // Current iteration number (0-indexed)
	MaxIterations       *int           // Termination threshold (nil = unlimited)
	IterationTimeout    *int           // Per-iteration timeout in seconds (nil = no timeout)
	MaxOutputBuffer     int            // Max AI CLI output buffer size in bytes (default: 10485760 = 10MB)
	ConsecutiveFailures int            // Consecutive AI CLI failures
	FailureThreshold    int            // Max consecutive failures before abort (default: 3)
	StartedAt           time.Time      // When the loop started
	Status              LoopStatus     // running, completed, aborted, interrupted
	ProcedureName       string         // Name of the procedure being executed
	Stats               IterationStats // Running statistics for iteration timing
}

// IterationStats tracks iteration timing statistics using Welford's online algorithm
// for constant memory usage regardless of iteration count
type IterationStats struct {
	Count     int           // Total iterations completed
	TotalTime time.Duration // Sum of all iteration durations
	MinTime   time.Duration // Fastest iteration (0 if no iterations)
	MaxTime   time.Duration // Slowest iteration (0 if no iterations)
	M2        float64       // Sum of squared differences from mean (for variance calculation)
}

// updateStats updates iteration statistics using Welford's online algorithm
// This maintains constant memory O(1) regardless of iteration count
func (s *IterationStats) updateStats(duration time.Duration) {
	s.Count++
	s.TotalTime += duration

	// Update min/max
	if s.Count == 1 || duration < s.MinTime {
		s.MinTime = duration
	}
	if duration > s.MaxTime {
		s.MaxTime = duration
	}

	// Welford's online algorithm for variance
	// See: https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
	durationSeconds := duration.Seconds()
	
	// Calculate old mean BEFORE updating count
	oldMean := 0.0
	if s.Count > 1 {
		oldMean = (s.TotalTime - duration).Seconds() / float64(s.Count-1)
	}
	
	// Calculate new mean AFTER updating count
	newMean := s.TotalTime.Seconds() / float64(s.Count)
	
	// M2 accumulates the squared distance from the mean
	// M2 = M2 + (x - old_mean) * (x - new_mean)
	s.M2 += (durationSeconds - oldMean) * (durationSeconds - newMean)
}

// getMean returns the mean iteration duration
func (s *IterationStats) getMean() time.Duration {
	if s.Count == 0 {
		return 0
	}
	return s.TotalTime / time.Duration(s.Count)
}

// getStdDev returns the standard deviation of iteration durations
func (s *IterationStats) getStdDev() time.Duration {
	if s.Count < 2 {
		return 0
	}
	
	// Variance = M2 / count (population variance)
	variance := s.M2 / float64(s.Count)
	stddevSeconds := math.Sqrt(variance)
	
	return time.Duration(stddevSeconds * float64(time.Second))
}
