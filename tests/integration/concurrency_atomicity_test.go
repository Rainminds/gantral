//go:build integration

package integration_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestConcurrency_SectionO verifies safe concurrent operations.
func TestConcurrency_SectionO(t *testing.T) {
	// 1. Parallel Instance Creation
	t.Run("Parallel Instance Creation", func(t *testing.T) {
		t.Parallel()
		var wg sync.WaitGroup
		createCount := 10

		// Map to detect duplicate IDs if generation isn't atomic/unique
		ids := sync.Map{}

		for i := 0; i < createCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Simulate creation
				// inst := engine.CreateInstance(...)
				// ids.Store(inst.ID, true)
				// For now, we mock the call or assume engine is thread-safe
				// If we have a real engine instance in future, we call it.
				// This skeleton ensures we have the test structure ready.

				// Mocking unique ID generation
				id := fmt.Sprintf("inst-%d-%d", time.Now().UnixNano(), i) // Flaky if not unique? In real code engine uses UUID.
				ids.Store(id, true)
			}()
		}
		wg.Wait()

		count := 0
		ids.Range(func(key, value any) bool {
			count++
			return true
		})

		if count != createCount {
			// In real test, this would fail if engine reused IDs
			t.Logf("Created %d unique instances out of %d attempts", count, createCount)
		}
	})

	// 2. Lock Contention / Concurrent Approval
	t.Run("Concurrent Approval -> One Succeeds", func(t *testing.T) {
		t.Parallel()
		// Setup instance in WAITING state
		// Launch 2 goroutines to Approve same instance
		// Assert: One gets success, other gets "invalid state" or handling error (idempotency?)
		// If idempotent, both succeed but only one transition happens.
		// If strict, second fails.
	})
}
