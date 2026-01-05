package engine

import (
	"context"
	"fmt"
)

// DecisionType defines the type of decision made.
type DecisionType string

const (
	DecisionApprove  DecisionType = "APPROVE"
	DecisionReject   DecisionType = "REJECT"
	DecisionOverride DecisionType = "OVERRIDE"
)

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

	var nextState State
	switch cmd.Type {
	case DecisionApprove:
		nextState = StateApproved
	case DecisionReject:
		nextState = StateRejected
	case DecisionOverride:
		nextState = StateOverridden
	default:
		return nil, fmt.Errorf("invalid decision type: %s", cmd.Type)
	}

	// 3. Delegate to Store for Transactional Update
	return e.store.RecordDecision(ctx, cmd, nextState)
}
