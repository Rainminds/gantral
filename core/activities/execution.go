package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/ports"
	"go.temporal.io/sdk/activity"
)

// ExecutionActivities provides activities for persisting execution state.
type ExecutionActivities struct {
	DB ports.InstanceStore
}

// PersistInstanceInput defines the input for PersistInstance activity.
type PersistInstanceInput struct {
	InstanceID      string
	WorkflowID      string
	TriggerContext  map[string]interface{}
	Policy          map[string]interface{} // Using generic map to avoid circular deps if needed, but engine.Policy is fine usually.
	PolicyVersionID string
	// Pre-evaluated policy result
	InitialState engine.State
	PolicyResult map[string]interface{}
}

// PersistInstance persists a new execution instance to the database.
func (a *ExecutionActivities) PersistInstance(ctx context.Context, input PersistInstanceInput) (*engine.Instance, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Persisting new instance", "workflow_id", input.WorkflowID)

	id := input.InstanceID
	if id == "" {
		id = fmt.Sprintf("inst-%s", activity.GetInfo(ctx).WorkflowExecution.RunID)
	}

	inst := &engine.Instance{
		ID:              id,
		WorkflowID:      input.WorkflowID,
		State:           input.InitialState, // Should be RUNNING or WAITING_FOR_HUMAN
		TriggerContext:  input.TriggerContext,
		PolicyVersionID: input.PolicyVersionID,
		PolicyContext:   input.PolicyResult,
	}
	err := a.DB.CreateInstance(ctx, inst)
	if err != nil {
		slog.Error("Failed to persist instance", "error", err)
		return nil, err
	}

	return inst, nil
}

// RecordDecisionInput defines input for decision recording.
type RecordDecisionInput struct {
	InstanceID      string
	DecisionType    engine.DecisionType
	ActorID         string
	Justification   string
	Role            string
	ContextSnapshot map[string]interface{}
	ContextDelta    map[string]interface{}
	PolicyVersionID string
}

// RecordDecision persists a human decision.
func (a *ExecutionActivities) RecordDecision(ctx context.Context, input RecordDecisionInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Recording decision", "instance_id", input.InstanceID, "type", input.DecisionType)

	cmd := engine.RecordDecisionCmd{
		InstanceID:      input.InstanceID,
		Type:            input.DecisionType,
		ActorID:         input.ActorID,
		Justification:   input.Justification,
		Role:            input.Role,
		ContextSnapshot: input.ContextSnapshot,
		ContextDelta:    input.ContextDelta,
		PolicyVersionID: input.PolicyVersionID,
	}

	// Calculate next state
	var nextState engine.State
	switch input.DecisionType {
	case engine.DecisionApprove:
		nextState = engine.StateApproved
	case engine.DecisionReject:
		nextState = engine.StateRejected
	case engine.DecisionOverride:
		nextState = engine.StateOverridden
	default:
		return fmt.Errorf("invalid decision type: %s", input.DecisionType)
	}

	_, err := a.DB.RecordDecision(ctx, cmd, nextState)
	return err
}

func generateInstanceID() string {
	// Placeholder, actual logic uses RunID in PersistInstance
	return ""
}

// Helper to marshal map to JSON bytes for logging/debug if needed
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
