package engine

import (
	"context"
	"fmt"
	"strings"
)

// DecisionType defines the type of decision made.
type DecisionType string

const (
	DecisionApprove  DecisionType = "APPROVE"
	DecisionReject   DecisionType = "REJECT"
	DecisionOverride DecisionType = "OVERRIDE"
)

// CalculateNextState determines the target state based on the decision type.
// This logic is shared between the Engine (transactional update) and Activities (Artifact emission).
func CalculateNextState(decisionType DecisionType) (State, error) {
	switch decisionType {
	case DecisionApprove:
		return StateApproved, nil
	case DecisionReject:
		return StateRejected, nil
	case DecisionOverride:
		return StateOverridden, nil
	default:
		return "", fmt.Errorf("invalid decision type: %s", decisionType)
	}
}

// RecordDecisionCmd is the input for recording a decision.
type RecordDecisionCmd struct {
	InstanceID      string
	Type            DecisionType
	ActorID         string
	Justification   string
	Role            string
	ContextSnapshot map[string]interface{}
	ContextDelta    map[string]interface{}
	PolicyVersionID string
	NewArtifactHash string // The hash of the artifact emitted for this decision (for chain linking)
}

// RecordDecision records a human decision and updates the instance state accordingly.
func (e *Engine) RecordDecision(ctx context.Context, cmd RecordDecisionCmd) (*Instance, error) {
	// 1. Fetch Instance to validate state (optimized: could be done in store, but logic here is safer)
	instance, err := e.store.GetInstance(ctx, cmd.InstanceID)
	if err != nil {
		return nil, err
	}

	// 2. Validate State Transition logic
	// Only allow decisions if waiting for human
	if instance.State != StateWaitingForHuman {
		return nil, fmt.Errorf("instance is not waiting for human decision (state: %s)", instance.State)
	}

	// 3. Enforce Invariants
	if cmd.Type == DecisionApprove || cmd.Type == DecisionOverride {
		if len(strings.TrimSpace(cmd.Justification)) == 0 {
			return nil, fmt.Errorf("justification is required for decision type %s", cmd.Type)
		}
	}

	nextState, err := CalculateNextState(cmd.Type)
	if err != nil {
		return nil, err
	}

	// 4. Delegate to Store for Transactional Update
	return e.store.RecordDecision(ctx, cmd, nextState)
}
