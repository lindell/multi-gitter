package multigitter

import (
	"sync"
	"testing"
	"time"
)

func TestRunInParallel_NoSleep(t *testing.T) {
	start := time.Now()
	count := 0
	var mu sync.Mutex

	runInParallel(func(i int) {
		mu.Lock()
		count++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}, 10, 5, 0)

	elapsed := time.Since(start)
	if count != 10 {
		t.Errorf("Expected 10 executions, got %d", count)
	}
	// Should complete in roughly 20ms (2 batches * 10ms)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Took too long without sleep: %v", elapsed)
	}
}

func TestRunInParallel_WithSleep(t *testing.T) {
	start := time.Now()
	count := 0
	var mu sync.Mutex
	sleepDuration := 50 * time.Millisecond

	runInParallel(func(i int) {
		mu.Lock()
		count++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}, 10, 5, sleepDuration)

	elapsed := time.Since(start)
	if count != 10 {
		t.Errorf("Expected 10 executions, got %d", count)
	}
	// Should sleep once between the two batches (5+5)
	// Total: ~20ms execution + 50ms sleep = 70ms minimum
	if elapsed < sleepDuration {
		t.Errorf("Should have slept at least %v, but took %v", sleepDuration, elapsed)
	}
}

func TestRunInParallel_SingleBatch(t *testing.T) {
	start := time.Now()
	sleepDuration := 50 * time.Millisecond

	runInParallel(func(i int) {
		time.Sleep(10 * time.Millisecond)
	}, 5, 10, sleepDuration)

	elapsed := time.Since(start)
	// Single batch, no sleep should occur
	if elapsed > 100*time.Millisecond {
		t.Errorf("Should not have slept for single batch, took %v", elapsed)
	}
}

func TestRunInParallel_MultipleBatches(t *testing.T) {
	batches := make([][]int, 0)
	var mu sync.Mutex
	currentBatch := make([]int, 0)
	sleepDuration := 30 * time.Millisecond

	runInParallel(func(i int) {
		mu.Lock()
		currentBatch = append(currentBatch, i)
		if len(currentBatch) == 3 {
			batches = append(batches, currentBatch)
			currentBatch = make([]int, 0)
		}
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}, 9, 3, sleepDuration)

	// Flush remaining
	mu.Lock()
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}
	mu.Unlock()

	// Should have 3 batches
	if len(batches) != 3 {
		t.Errorf("Expected 3 batches, got %d", len(batches))
	}
}
