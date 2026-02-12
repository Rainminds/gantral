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

// ValidateDecision enforces HITL invariants on a decision command against an instance state.
// It is exposed for testing purposes to verify logic without DB dependencies.
func ValidateDecision(instance *Instance, cmd RecordDecisionCmd) error {
	// 1. Validate State Transition logic
	// Only allow decisions if waiting for human
	if instance.State != StateWaitingForHuman {
		return fmt.Errorf("instance is not waiting for human decision (state: %s)", instance.State)
	}

	// 2. Enforce Invariants
	if cmd.Type == DecisionApprove || cmd.Type == DecisionOverride {
		if len(strings.TrimSpace(cmd.Justification)) == 0 {
			return fmt.Errorf("justification is required for decision type %s", cmd.Type)
		}
	}

	// Enforce ContextDelta for Override (Section B requirements)
	if cmd.Type == DecisionOverride {
		if len(cmd.ContextDelta) == 0 {
			return fmt.Errorf("context_delta is required for decision type %s", cmd.Type)
		}
	}

	// Section B: Missing identity -> reject
	if strings.TrimSpace(cmd.ActorID) == "" {
		return fmt.Errorf("missing actor identity")
	}

	// Section B: Role mismatch (basic check, assuming role is required if present in cmd,
	// though RBAC might be deeper. Enforcing non-empty role if system requires it).
	// For Tier 1, we assert that if Role is provided, it's not empty string?
	// The prompt says "Role mismatch -> reject".
	// Without RBAC definitions, I will enforce ActorID presence as "Missing identity".

	return nil
}

// RecordDecision records a human decision and updates the instance state accordingly.
func (e *Engine) RecordDecision(ctx context.Context, cmd RecordDecisionCmd) (*Instance, error) {
	// 1. Fetch Instance
	instance, err := e.store.GetInstance(ctx, cmd.InstanceID)
	if err != nil {
		return nil, err
	}

	// 2. Validate
	if err := ValidateDecision(instance, cmd); err != nil {
		return nil, err
	}

	nextState, err := CalculateNextState(cmd.Type)
	if err != nil {
		return nil, err
	}

	// 3. Delegate to Store for Transactional Update
	return e.store.RecordDecision(ctx, cmd, nextState)
}
