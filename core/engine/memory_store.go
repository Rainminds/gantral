package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// MemoryStore implements InstanceStore in memory for testing.
type MemoryStore struct {
	mu        sync.RWMutex
	instances map[string]*Instance
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		instances: make(map[string]*Instance),
	}
}

func (s *MemoryStore) CreateInstance(ctx context.Context, inst *Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Deep copy to simulate storage boundary
	s.instances[inst.ID] = copyInstance(inst)
	return nil
}

func (s *MemoryStore) GetInstance(ctx context.Context, id string) (*Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	inst, ok := s.instances[id]
	if !ok {
		return nil, fmt.Errorf("instance not found: %s", id)
	}
	return copyInstance(inst), nil
}

func (s *MemoryStore) ListInstances(ctx context.Context) ([]*Instance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Instance
	for _, inst := range s.instances {
		result = append(result, copyInstance(inst))
	}
	return result, nil
}

func (s *MemoryStore) RecordDecision(ctx context.Context, cmd RecordDecisionCmd, nextState State) (*Instance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	inst, ok := s.instances[cmd.InstanceID]
	if !ok {
		return nil, fmt.Errorf("instance not found: %s", cmd.InstanceID)
	}

	inst.State = nextState
	// In a real store, we would also save the decision record

	s.instances[cmd.InstanceID] = copyInstance(inst)
	return copyInstance(inst), nil
}

func copyInstance(src *Instance) *Instance {
	dst := *src
	// Helper to copy inner maps if needed, but for now shallow copy of maps is risky if tests mutate them
	// Ideally we deep copy maps
	dst.TriggerContext = deepCopyMap(src.TriggerContext)
	dst.PolicyContext = deepCopyMap(src.PolicyContext)
	return &dst
}

func deepCopyMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	bytes, _ := json.Marshal(src)
	var dst map[string]interface{}
	json.Unmarshal(bytes, &dst)
	return dst
}
