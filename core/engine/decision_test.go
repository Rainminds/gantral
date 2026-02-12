package engine

import (
	"context"
	"testing"
)

func TestRecordDecision(t *testing.T) {
	// engine initialized per test case
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
				ContextDelta:  map[string]interface{}{"reason": "emergency"},
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
			store := NewMemoryStore()
			engine := NewEngine(store)

			// Seed the instance
			instStub := &Instance{
				ID:    tt.cmd.InstanceID,
				State: StateWaitingForHuman,
			}
			_ = store.CreateInstance(ctx, instStub)

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

func TestRecordDecision_MissingJustification(t *testing.T) {
	// Setup
	store := NewMemoryStore()
	e := NewEngine(store)
	ctx := context.Background()

	// 1. Create Instance in Waiting State
	inst := &Instance{
		ID:    "inst-justification-test",
		State: StateWaitingForHuman,
	}
	_ = store.CreateInstance(ctx, inst)

	// 2. Attempt APPROVE with empty justification
	cmdApprove := RecordDecisionCmd{
		InstanceID:    inst.ID,
		Type:          DecisionApprove,
		ActorID:       "user-1",
		Justification: "   ", // Whitespace only
	}
	_, err := e.RecordDecision(ctx, cmdApprove)
	if err == nil {
		t.Fatal("Expected error for APPROVE with empty justification, got nil")
	}

	// 3. Attempt OVERRIDE with empty justification
	cmdOverride := RecordDecisionCmd{
		InstanceID:    inst.ID,
		Type:          DecisionOverride,
		ActorID:       "admin-1",
		Justification: "",
	}
	_, err = e.RecordDecision(ctx, cmdOverride)
	if err == nil {
		t.Fatal("Expected error for OVERRIDE with empty justification, got nil")
	}

	// 4. Attempt REJECT with empty justification (Should be allowed? User implies check for Approve/Override. Usually Reject also needs it but user query specifically asked for Approve/Override)
	// "Is justification required (non-null, non-empty) for APPROVE decisions? Is justification required for OVERRIDE decisions?"
	// Let's enforce for all for consistency, or strictly follow user. User asked "If ANY of the above are false".
	// Let's stick to Approve/Override for now.
}
