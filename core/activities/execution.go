package activities

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/ports"
	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"go.temporal.io/sdk/activity"
)

// ExecutionActivities provides activities for persisting execution state.
type ExecutionActivities struct {
	DB              ports.InstanceStore
	ArtifactEmitter artifact.ArtifactEmitter
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
	InstanceID      string                 `json:"instance_id"`
	DecisionType    engine.DecisionType    `json:"decision_type"`
	ActorID         string                 `json:"actor_id"`
	Justification   string                 `json:"justification"`
	Role            string                 `json:"role"`
	ContextSnapshot map[string]interface{} `json:"context_snapshot"`
	ContextDelta    map[string]interface{} `json:"context_delta"`
	PolicyVersionID string                 `json:"policy_version_id"`
	EvidenceHash    string                 `json:"evidence_hash"` // Hash of tool execution evidence (Phase 5.5)
}

// RecordDecision persists a human decision.
func (a *ExecutionActivities) RecordDecision(ctx context.Context, input RecordDecisionInput) (*models.CommitmentArtifact, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Recording decision", "instance_id", input.InstanceID, "type", input.DecisionType)

	// 1. Fetch Current Instance State (to get Previous Hash)
	// We need 'prevHash' to be the LastArtifactHash of the instance.
	instance, err := a.DB.GetInstance(ctx, input.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch instance for chaining: %w", err)
	}

	// 2. Prepare Metadata
	// Calculate next state using shared engine logic to ensure Artifact matches DB state.
	nextState, err := engine.CalculateNextState(input.DecisionType)
	if err != nil {
		return nil, err
	}

	// 3. Emit Commitment Artifact (Evidence)
	// Compute Deterministic Context Hash
	// If EvidenceHash is provided (Tool Mediation), use it.
	// Otherwise, fallback to hashing the ContextSnapshot (Human Decision).
	var contextHash string
	if input.EvidenceHash != "" {
		contextHash = input.EvidenceHash
	} else {
		var err error
		contextHash, err = artifact.HashContext(input.ContextSnapshot)
		if err != nil {
			return nil, fmt.Errorf("failed to hash context: %w", err)
		}
	}

	art, err := a.ArtifactEmitter.EmitArtifact(
		ctx,
		input.InstanceID,
		instance.LastArtifactHash, // Chain Link
		string(nextState),
		input.PolicyVersionID,
		contextHash,
		input.ActorID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to emit artifact: %w", err)
	}

	// 4. Persist to DB (State + Chain Link)
	cmd := engine.RecordDecisionCmd{
		InstanceID:      input.InstanceID,
		Type:            input.DecisionType,
		ActorID:         input.ActorID,
		Justification:   input.Justification,
		Role:            input.Role,
		ContextSnapshot: input.ContextSnapshot,
		ContextDelta:    input.ContextDelta,
		PolicyVersionID: input.PolicyVersionID,
		NewArtifactHash: art.ArtifactID, // Persist the new link
	}

	if _, err := a.DB.RecordDecision(ctx, cmd, nextState); err != nil {
		// CONSISTENCY MODEL: At-Least-Once Artifacts
		// If DB fails here, we have an "Orphaned Artifact" in WORM storage.
		// This is ACCEPTABLE for an audit system: it represents a valid attempt
		// that failed to operationalize.
		// The reverse (DB updated, Artifact missing) is NEVER allowed.
		// Therefore, we emit first, then write to DB.
		return nil, fmt.Errorf("failed to record decision in DB: %w", err)
	}

	return art, nil
}
