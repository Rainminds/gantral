package statemachine

import (
	"testing"

	"github.com/Rainminds/gantral/core/engine"
)

// Test_StateMachine_Matrix enforces the strict transition matrix defined in TRD.
// It exhaustively tests every possible From->To combination.
func Test_StateMachine_Matrix(t *testing.T) {
	t.Parallel()

	allStates := []engine.State{
		engine.StateCreated,
		engine.StateRunning,
		engine.StateWaitingForHuman,
		engine.StateApproved,
		engine.StateRejected,
		engine.StateOverridden,
		engine.StateResumed,
		engine.StateCompleted,
		engine.StateTerminated,
	}

	// We use the authoritative map from the codebase as the source of truth for "Expected Success".
	// However, to ensure the TEST doesn't just mirror the code, we explicitly define expected allowed transitions here
	// to catch regressions if the code map is accidentally modified.
	allowedTransitions := map[engine.State]map[engine.State]bool{
		engine.StateCreated: {
			engine.StateRunning: true,
		},
		engine.StateRunning: {
			engine.StateWaitingForHuman: true,
			engine.StateCompleted:       true,
			engine.StateTerminated:      true,
		},
		engine.StateWaitingForHuman: {
			engine.StateApproved:   true,
			engine.StateRejected:   true,
			engine.StateOverridden: true,
		},
		engine.StateApproved: {
			engine.StateResumed: true,
		},
		engine.StateRejected: {
			engine.StateTerminated: true,
		},
		engine.StateOverridden: {
			engine.StateResumed: true,
		},
		engine.StateResumed: {
			engine.StateRunning: true,
		},
		engine.StateCompleted:  {}, // Terminal
		engine.StateTerminated: {}, // Terminal
	}

	for _, from := range allStates {
		from := from // capture loop var
		for _, to := range allStates {
			to := to // capture loop var

			t.Run(string(from)+"->"+string(to), func(t *testing.T) {
				t.Parallel()

				// Setup instance in 'from' state
				instance := &engine.Instance{
					State: from,
				}

				// Check if this transition SHOULD be allowed
				expectAllowed := false
				if allowedMap, ok := allowedTransitions[from]; ok {
					if allowedMap[to] {
						expectAllowed = true
					}
				}

				// Execute transition
				err := engine.Transition(instance, to)

				if expectAllowed {
					if err != nil {
						t.Errorf("Expected valid transition %s->%s failed: %v", from, to, err)
					}
					if instance.State != to {
						t.Errorf("Expected state to be %s, got %s", to, instance.State)
					}
				} else {
					if err == nil {
						t.Errorf("Expected invalid transition %s->%s to fail, but it succeeded", from, to)
					}
					// Verify state didn't change on failure (though implementation might not guarantee strict immutability on error struct-wise,
					// logically it shouldn't change). The Transition function in machine.go updates state ONLY if valid.
					if instance.State != from {
						t.Errorf("State mutated despite invalid transition! Expected %s, got %s", from, instance.State)
					}
				}
			})
		}
	}
}

// Test_StateMachine_InvalidTransition_Panics verifies that nil pointer or invalid inputs
// are handled safely (fail-closed or panic as per Go semantics, though Transition returns error).
// The Blueprint says "Invalid transitions must panic OR terminate appropriately".
// Since the current implementation returns error, we enforce that.
func Test_StateMachine_UnknownState_Rejected(t *testing.T) {
	t.Parallel()

	instance := &engine.Instance{
		State: "UNKNOWN_STATE_999",
	}

	err := engine.Transition(instance, engine.StateRunning)
	if err == nil {
		t.Error("Expected error when transitioning from unknown state")
	}
}
