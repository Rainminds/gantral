package unit

import (
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
)

func Test_Hitl_Enforcement_SectionB(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		instance    *engine.Instance
		cmd         engine.RecordDecisionCmd
		expectError bool
		errContains string
	}{
		// 1. Valid Approval
		{
			name: "Valid Approval",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionApprove,
				ActorID:       "user-1",
				Justification: "Looks good",
			},
			expectError: false,
		},
		// 2. Invalid State
		{
			name: "Invalid State (Running)",
			instance: &engine.Instance{
				State: engine.StateRunning,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionApprove,
				ActorID:       "user-1",
				Justification: "Looks good",
			},
			expectError: true,
			errContains: "not waiting for human",
		},
		// 3. Missing Justification (Approve)
		{
			name: "Missing Justification (Approve)",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionApprove,
				ActorID:       "user-1",
				Justification: "",
			},
			expectError: true,
			errContains: "justification is required",
		},
		// 4. Whitespace Justification
		{
			name: "Whitespace Justification",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionApprove,
				ActorID:       "user-1",
				Justification: "   ",
			},
			expectError: true,
			errContains: "justification is required",
		},
		// 5. Valid Override
		{
			name: "Valid Override",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionOverride,
				ActorID:       "admin-1",
				Justification: "Fixing issue",
				ContextDelta:  map[string]interface{}{"foo": "bar"},
			},
			expectError: false,
		},
		// 6. Override Missing ContextDelta
		{
			name: "Override Missing ContextDelta",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionOverride,
				ActorID:       "admin-1",
				Justification: "Fixing issue",
				ContextDelta:  nil, // Empty
			},
			expectError: true,
			errContains: "context_delta is required",
		},
		// 7. Missing Identity
		{
			name: "Missing Identity",
			instance: &engine.Instance{
				State: engine.StateWaitingForHuman,
			},
			cmd: engine.RecordDecisionCmd{
				Type:          engine.DecisionReject,
				ActorID:       "", // Empty
				Justification: "Bad",
			},
			expectError: true,
			errContains: "missing actor identity",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := engine.ValidateDecision(tt.instance, tt.cmd)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errContains)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
