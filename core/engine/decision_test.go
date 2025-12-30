package engine

import (
	"context"
	"testing"
)

func TestRecordDecision(t *testing.T) {
	engine := NewEngine()
	ctx := context.Background()

	tests := []struct {
		name          string
		cmd           RecordDecisionCmd
		expectedState State
		expectError   bool
	}{
		{
			name: "Approve Decision",
			cmd: RecordDecisionCmd{
				InstanceID:    "inst-waiting",
				Type:          DecisionApprove,
				ActorID:       "human-1",
				Justification: "Looks good",
			},
			expectedState: StateApproved,
			expectError:   false,
		},
		{
			name: "Reject Decision",
			cmd: RecordDecisionCmd{
				InstanceID:    "inst-waiting",
				Type:          DecisionReject,
				ActorID:       "human-1",
				Justification: "Bad inputs",
			},
			expectedState: StateRejected,
			expectError:   false,
		},
		{
			name: "Override Decision",
			cmd: RecordDecisionCmd{
				InstanceID:    "inst-waiting",
				Type:          DecisionOverride,
				ActorID:       "admin",
				Justification: "Emergency",
			},
			expectedState: StateOverridden,
			expectError:   false,
		},
		{
			name: "Invalid Decision Type",
			cmd: RecordDecisionCmd{
				InstanceID:    "inst-waiting",
				Type:          "INVALID_TYPE",
				ActorID:       "human-1",
				Justification: "Typo",
			},
			expectedState: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Seed the instance
			instStub := &Instance{
				ID:    tt.cmd.InstanceID,
				State: StateWaitingForHuman,
			}
			engine.instances[instStub.ID] = instStub

			inst, err := engine.RecordDecision(ctx, tt.cmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if inst.State != tt.expectedState {
				t.Errorf("expected state %s, got %s", tt.expectedState, inst.State)
			}
		})
	}
}
