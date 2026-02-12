package policy

import (
	"context"
	"testing"
)

func TestEvaluate(t *testing.T) {
	engine := NewEngine()
	ctx := context.Background()

	tests := []struct {
		name          string
		policy        Policy
		expectedState string
		shouldPause   bool
	}{
		{
			name: "Low Materiality - Auto Run",
			policy: Policy{
				ID:          "pol-1",
				Materiality: MaterialityLow,
			},
			expectedState: "RUNNING",
			shouldPause:   false,
		},
		{
			name: "High Materiality - Pause",
			policy: Policy{
				ID:          "pol-2",
				Materiality: MaterialityHigh,
			},
			expectedState: "WAITING_FOR_HUMAN",
			shouldPause:   true,
		},
		{
			name: "Medium Materiality with Approval - Pause",
			policy: Policy{
				ID:                    "pol-3",
				Materiality:           MaterialityMedium,
				RequiresHumanApproval: true,
			},
			expectedState: "WAITING_FOR_HUMAN",
			shouldPause:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Evaluate(ctx, tt.policy)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.NextState != tt.expectedState {
				t.Errorf("expected state %s, got %s", tt.expectedState, result.NextState)
			}

			if result.ShouldPause != tt.shouldPause {
				t.Errorf("expected pause %v, got %v", tt.shouldPause, result.ShouldPause)
			}
		})
	}
}

func TestEvaluatePure_Exhaustive(t *testing.T) {
	tests := []struct {
		name     string
		policy   Policy
		expected bool
	}{
		{"Low Materiality", Policy{Materiality: MaterialityLow}, false},
		{"Medium Materiality", Policy{Materiality: MaterialityMedium}, false},
		{"High Materiality", Policy{Materiality: MaterialityHigh}, true},
		{"Low with Human", Policy{Materiality: MaterialityLow, RequiresHumanApproval: true}, true},
		{"Medium with Human", Policy{Materiality: MaterialityMedium, RequiresHumanApproval: true}, true},
		{"High with Human", Policy{Materiality: MaterialityHigh, RequiresHumanApproval: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluatePure(tt.policy)
			if result.ShouldPause != tt.expected {
				t.Errorf("expected pause %v, got %v", tt.expected, result.ShouldPause)
			}
		})
	}
}
