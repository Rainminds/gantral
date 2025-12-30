package engine

import (
	"context"
	"fmt"
	"log/slog"
	"time"
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
	InstanceID    string
	Type          DecisionType
	ActorID       string
	Justification string
}

// RecordDecision records a human decision and updates the instance state accordingly.
func (e *Engine) RecordDecision(ctx context.Context, cmd RecordDecisionCmd) (*Instance, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. Fetch Instance
	instance, ok := e.instances[cmd.InstanceID]
	if !ok {
		return nil, fmt.Errorf("instance not found: %s", cmd.InstanceID)
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

	// 3. Update State
	instance.State = nextState
	instance.UpdatedAt = time.Now()

	// 4. Log Decision
	slog.Info("decision_recorded",
		"instance_id", cmd.InstanceID,
		"type", cmd.Type,
		"actor", cmd.ActorID,
		"next_state", nextState,
	)

	return instance, nil
}
