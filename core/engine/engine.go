package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Rainminds/gantral/core/policy"
)

// Engine is the core component that manages execution lifecycles.
type Engine struct {
	policyEngine *policy.Engine

	// In-memory store for now.
	// WARNING: Data is lost on server restart. Replace with DB in Phase 3.
	mu        sync.RWMutex
	instances map[string]*Instance
}

// NewEngine creates a new instance of the Engine.
func NewEngine() *Engine {
	return &Engine{
		policyEngine: policy.NewEngine(),
		instances:    make(map[string]*Instance),
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

	// 4. Store in memory
	e.mu.Lock()
	e.instances[instance.ID] = instance
	e.mu.Unlock()

	return instance, nil
}

// GetInstance retrieves an instance by ID.
func (e *Engine) GetInstance(ctx context.Context, id string) (*Instance, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	inst, ok := e.instances[id]
	if !ok {
		return nil, fmt.Errorf("instance not found: %s", id)
	}
	return inst, nil
}

// ListInstances retrieves all instances.
func (e *Engine) ListInstances(ctx context.Context) ([]*Instance, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]*Instance, 0, len(e.instances))
	for _, inst := range e.instances {
		result = append(result, inst)
	}
	return result, nil
}
