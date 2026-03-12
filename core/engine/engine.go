package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rainminds/gantral/core/policy"
)

var (
	ErrMissingTeamID        = errors.New("missing team id")
	ErrCrossTenantViolation = errors.New("cross-tenant access rejected")
)

// InstanceStore defines the persistence layer requirements.
type InstanceStore interface {
	CreateInstance(ctx context.Context, inst *Instance) error
	GetInstance(ctx context.Context, id string) (*Instance, error)
	ListInstances(ctx context.Context, teamID string) ([]*Instance, error)
	RecordDecision(ctx context.Context, cmd RecordDecisionCmd, nextState State) (*Instance, error)
}

// Engine is the core component that manages execution lifecycles.
type Engine struct {
	policyEngine *policy.Engine
	store        InstanceStore
}

// NewEngine creates a new instance of the Engine.
func NewEngine(store InstanceStore) *Engine {
	return &Engine{
		policyEngine: policy.NewEngine(),
		store:        store,
	}
}

// CreateInstance starts a new execution instance.
func (e *Engine) CreateInstance(ctx context.Context, teamID string, workflowID string, triggerContext map[string]interface{}, pol policy.Policy) (*Instance, error) {
	if teamID == "" {
		return nil, ErrMissingTeamID
	}

	// 1. Evaluate Policy
	evalResult, err := e.policyEngine.Evaluate(ctx, pol)
	if err != nil {
		return nil, fmt.Errorf("policy evaluation failed: %w", err)
	}

	// 2. Determine Initial State
	initialState := StateRunning
	if evalResult.ShouldPause {
		initialState = StateWaitingForHuman
	}

	// 3. Create Instance Record
	instance := &Instance{
		ID:             fmt.Sprintf("inst-%d", time.Now().UnixNano()), // Simple unique ID
		OwningTeamID:   teamID,
		WorkflowID:     workflowID,
		State:          initialState,
		TriggerContext: triggerContext,
		PolicyContext: map[string]interface{}{
			"policy_id": pol.ID,
			"decision":  evalResult.NextState,
			"reason":    evalResult.Reason,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Store via Interface
	if err := e.store.CreateInstance(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return instance, nil
}

// GetInstance retrieves an instance by ID and checks tenant isolation.
func (e *Engine) GetInstance(ctx context.Context, teamID string, id string) (*Instance, error) {
	inst, err := e.store.GetInstance(ctx, id)
	if err != nil {
		return nil, err
	}
	if inst.OwningTeamID != teamID {
		return nil, ErrCrossTenantViolation
	}
	return inst, nil
}

// ListInstances retrieves all instances for a specific team.
func (e *Engine) ListInstances(ctx context.Context, teamID string) ([]*Instance, error) {
	if teamID == "" {
		return nil, ErrMissingTeamID
	}
	return e.store.ListInstances(ctx, teamID)
}
