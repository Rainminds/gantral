package engine

import (
	"testing"
)

func TestTransition(t *testing.T) {
	tests := []struct {
		name          string
		initialState  State
		targetState   State
		expectError   bool
		expectedState State
	}{
		// Happy Paths
		{
			name:          "Created -> Running",
			initialState:  StateCreated,
			targetState:   StateRunning,
			expectError:   false,
			expectedState: StateRunning,
		},
		{
			name:          "Running -> WaitingForHuman",
			initialState:  StateRunning,
			targetState:   StateWaitingForHuman,
			expectError:   false,
			expectedState: StateWaitingForHuman,
		},
		{
			name:          "WaitingForHuman -> Approved",
			initialState:  StateWaitingForHuman,
			targetState:   StateApproved,
			expectError:   false,
			expectedState: StateApproved,
		},
		{
			name:          "Approved -> Resumed",
			initialState:  StateApproved,
			targetState:   StateResumed,
			expectError:   false,
			expectedState: StateResumed,
		},
		{
			name:          "Resumed -> Running",
			initialState:  StateResumed,
			targetState:   StateRunning,
			expectError:   false,
			expectedState: StateRunning,
		},
		{
			name:          "Running -> Completed",
			initialState:  StateRunning,
			targetState:   StateCompleted,
			expectError:   false,
			expectedState: StateCompleted,
		},

		// Invalid Transitions
		{
			name:         "Created -> Completed (Skip Running)",
			initialState: StateCreated,
			targetState:  StateCompleted,
			expectError:  true,
		},
		{
			name:         "WaitingForHuman -> Running (Skip Decision)",
			initialState: StateWaitingForHuman,
			targetState:  StateRunning,
			expectError:  true,
		},
		{
			name:         "Completed -> Running (Terminal State)",
			initialState: StateCompleted,
			targetState:  StateRunning,
			expectError:  true,
		},
		{
			name:         "UnknownState -> Created",
			initialState: "UNKNOWN_STATE",
			targetState:  StateCreated,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				State: tt.initialState,
			}

			err := Transition(instance, tt.targetState)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if instance.State != tt.expectedState {
					t.Errorf("expected state %s, got %s", tt.expectedState, instance.State)
				}
				if instance.UpdatedAt.IsZero() {
					t.Errorf("expected UpdatedAt to be set")
				}
			}
		})
	}
}
