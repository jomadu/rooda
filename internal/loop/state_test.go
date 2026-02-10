package loop

import (
	"math"
	"testing"
	"time"
)

func TestIterationStats_InitialState(t *testing.T) {
	stats := IterationStats{}

	if stats.Count != 0 {
		t.Errorf("expected Count=0, got %d", stats.Count)
	}
	if stats.TotalTime != 0 {
		t.Errorf("expected TotalTime=0, got %v", stats.TotalTime)
	}
	if stats.MinTime != 0 {
		t.Errorf("expected MinTime=0, got %v", stats.MinTime)
	}
	if stats.MaxTime != 0 {
		t.Errorf("expected MaxTime=0, got %v", stats.MaxTime)
	}
	if stats.M2 != 0 {
		t.Errorf("expected M2=0, got %f", stats.M2)
	}
}

func TestIterationStats_SingleIteration(t *testing.T) {
	stats := IterationStats{}
	duration := 5 * time.Second

	stats.updateStats(duration)

	if stats.Count != 1 {
		t.Errorf("expected Count=1, got %d", stats.Count)
	}
	if stats.TotalTime != duration {
		t.Errorf("expected TotalTime=%v, got %v", duration, stats.TotalTime)
	}
	if stats.MinTime != duration {
		t.Errorf("expected MinTime=%v, got %v", duration, stats.MinTime)
	}
	if stats.MaxTime != duration {
		t.Errorf("expected MaxTime=%v, got %v", duration, stats.MaxTime)
	}
	if stats.M2 != 0 {
		t.Errorf("expected M2=0 for single iteration, got %f", stats.M2)
	}

	mean := stats.getMean()
	if mean != duration {
		t.Errorf("expected mean=%v, got %v", duration, mean)
	}

	stddev := stats.getStdDev()
	if stddev != 0 {
		t.Errorf("expected stddev=0 for single iteration, got %v", stddev)
	}
}

func TestIterationStats_MultipleIterations(t *testing.T) {
	stats := IterationStats{}
	durations := []time.Duration{
		2 * time.Second,
		4 * time.Second,
		6 * time.Second,
	}

	for _, d := range durations {
		stats.updateStats(d)
	}

	if stats.Count != 3 {
		t.Errorf("expected Count=3, got %d", stats.Count)
	}

	expectedTotal := 12 * time.Second
	if stats.TotalTime != expectedTotal {
		t.Errorf("expected TotalTime=%v, got %v", expectedTotal, stats.TotalTime)
	}

	if stats.MinTime != 2*time.Second {
		t.Errorf("expected MinTime=2s, got %v", stats.MinTime)
	}

	if stats.MaxTime != 6*time.Second {
		t.Errorf("expected MaxTime=6s, got %v", stats.MaxTime)
	}

	expectedMean := 4 * time.Second
	mean := stats.getMean()
	if mean != expectedMean {
		t.Errorf("expected mean=%v, got %v", expectedMean, mean)
	}

	// For values [2, 4, 6], mean=4, variance=((2-4)^2 + (4-4)^2 + (6-4)^2)/3 = 8/3
	// stddev = sqrt(8/3) ≈ 1.633 seconds
	expectedStdDev := time.Duration(math.Sqrt(8.0/3.0) * float64(time.Second))
	stddev := stats.getStdDev()
	
	// Allow small floating point error
	diff := math.Abs(float64(stddev - expectedStdDev))
	if diff > float64(time.Millisecond) {
		t.Errorf("expected stddev≈%v, got %v (diff=%v)", expectedStdDev, stddev, diff)
	}
}

func TestIterationStats_WelfordAlgorithm(t *testing.T) {
	// Test that Welford's algorithm produces correct variance
	stats := IterationStats{}
	values := []float64{10, 20, 30, 40, 50}
	
	for _, v := range values {
		stats.updateStats(time.Duration(v * float64(time.Second)))
	}

	// Mean = 30
	expectedMean := 30 * time.Second
	mean := stats.getMean()
	if mean != expectedMean {
		t.Errorf("expected mean=%v, got %v", expectedMean, mean)
	}

	// Variance = ((10-30)^2 + (20-30)^2 + (30-30)^2 + (40-30)^2 + (50-30)^2) / 5
	//          = (400 + 100 + 0 + 100 + 400) / 5 = 1000 / 5 = 200
	// StdDev = sqrt(200) ≈ 14.142 seconds
	expectedStdDev := time.Duration(math.Sqrt(200.0) * float64(time.Second))
	stddev := stats.getStdDev()
	
	diff := math.Abs(float64(stddev - expectedStdDev))
	if diff > float64(time.Millisecond) {
		t.Errorf("expected stddev≈%v, got %v (diff=%v)", expectedStdDev, stddev, diff)
	}
}

func TestIterationStats_ZeroIterations(t *testing.T) {
	stats := IterationStats{}

	mean := stats.getMean()
	if mean != 0 {
		t.Errorf("expected mean=0 for zero iterations, got %v", mean)
	}

	stddev := stats.getStdDev()
	if stddev != 0 {
		t.Errorf("expected stddev=0 for zero iterations, got %v", stddev)
	}
}

func TestIterationStats_ConstantMemory(t *testing.T) {
	// Verify that statistics use constant memory regardless of iteration count
	stats := IterationStats{}
	
	// Run many iterations
	for i := 0; i < 10000; i++ {
		stats.updateStats(time.Duration(i) * time.Millisecond)
	}

	// Should still have only 5 fields (Count, TotalTime, MinTime, MaxTime, M2)
	// This is a compile-time guarantee, but we verify the stats are still correct
	if stats.Count != 10000 {
		t.Errorf("expected Count=10000, got %d", stats.Count)
	}

	mean := stats.getMean()
	if mean == 0 {
		t.Error("expected non-zero mean")
	}

	stddev := stats.getStdDev()
	if stddev == 0 {
		t.Error("expected non-zero stddev")
	}
}
