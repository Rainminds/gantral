package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/Rainminds/gantral/core/policy"
)

// InstanceStore defines the persistence layer requirements.
type InstanceStore interface {
	CreateInstance(ctx context.Context, inst *Instance) error
	GetInstance(ctx context.Context, id string) (*Instance, error)
	ListInstances(ctx context.Context) ([]*Instance, error)
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
func (e *Engine) CreateInstance(ctx context.Context, workflowID string, triggerContext map[string]interface{}, pol policy.Policy) (*Instance, error) {
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

// GetInstance retrieves an instance by ID.
func (e *Engine) GetInstance(ctx context.Context, id string) (*Instance, error) {
	return e.store.GetInstance(ctx, id)
}

// ListInstances retrieves all instances.
func (e *Engine) ListInstances(ctx context.Context) ([]*Instance, error) {
	return e.store.ListInstances(ctx)
}
