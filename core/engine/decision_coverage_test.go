package engine

import (
	"context"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/policy"
)

func TestValidateDecision_Unit(t *testing.T) {
	tests := []struct {
		name        string
		instance    *Instance
		cmd         RecordDecisionCmd
		expectError bool
		errMsg      string
	}{
		{
			name: "Valid Approval",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:          DecisionApprove,
				ActorID:       "admin",
				Justification: "Looks good",
			},
			expectError: false,
		},
		{
			name: "Missing Justification for Approval",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:          DecisionApprove,
				ActorID:       "admin",
				Justification: "",
			},
			expectError: true,
			errMsg:      "justification is required",
		},
		{
			name: "Valid Override",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:          DecisionOverride,
				ActorID:       "admin",
				Justification: "Emergency",
				ContextDelta:  map[string]interface{}{"key": "val"},
			},
			expectError: false,
		},
		{
			name: "Missing ContextDelta for Override",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:          DecisionOverride,
				ActorID:       "admin",
				Justification: "Emergency",
			},
			expectError: true,
			errMsg:      "context_delta is required",
		},
		{
			name: "Missing Identity",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:          DecisionApprove,
				ActorID:       " ",
				Justification: "Docs",
			},
			expectError: true,
			errMsg:      "missing actor identity",
		},
		{
			name: "Reject - No justification needed",
			instance: &Instance{
				State: StateWaitingForHuman,
			},
			cmd: RecordDecisionCmd{
				Type:    DecisionReject,
				ActorID: "admin",
			},
			expectError: false,
		},
		{
			name: "Wrong Initial State",
			instance: &Instance{
				State: StateRunning,
			},
			cmd: RecordDecisionCmd{
				Type:    DecisionApprove,
				ActorID: "admin",
			},
			expectError: true,
			errMsg:      "is not waiting for human decision",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDecision(tt.instance, tt.cmd)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errMsg)
				} else if tt.errMsg != "" && !anyContains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing '%s', got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCalculateNextState(t *testing.T) {
	tests := []struct {
		decision DecisionType
		expected State
		err      bool
	}{
		{DecisionApprove, StateApproved, false},
		{DecisionReject, StateRejected, false},
		{DecisionOverride, StateOverridden, false},
		{DecisionType("KILL"), "", true},
	}

	for _, tt := range tests {
		res, err := CalculateNextState(tt.decision)
		if tt.err {
			if err == nil {
				t.Errorf("expected error for decision %s", tt.decision)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.decision, err)
			}
			if res != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, res)
			}
		}
	}
}

func TestEngine_CreateInstance_PolicyBranching(t *testing.T) {
	ms := NewMemoryStore()
	e := NewEngine(ms)
	ctx := context.Background()

	// 1. Auto-Run Policy
	polAuto := policy.Policy{ID: "p-auto"} // Default policy engine returns ShouldPause=false
	inst, err := e.CreateInstance(ctx, "team-1", "wf-1", nil, polAuto)
	if err != nil {
		t.Fatal(err)
	}
	if inst.State != StateRunning {
		t.Errorf("expected RUNNING, got %s", inst.State)
	}

	// 2. Missing Team
	_, err = e.CreateInstance(ctx, "", "wf-1", nil, polAuto)
	if err != ErrMissingTeamID {
		t.Errorf("expected ErrMissingTeamID, got %v", err)
	}
}

func anyContains(s string, sub string) bool {
	return strings.Contains(s, sub)
}
