package engine

import (
	"context"
	"testing"

	"github.com/Rainminds/gantral/core/policy"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Error("NewEngine returned nil")
	}
}

func TestEngineCRUD(t *testing.T) {
	e := NewEngine()
	ctx := context.Background()

	// 1. List Empty
	list, err := e.ListInstances(ctx)
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}

	// 2. Create Instance
	pol := policy.Policy{ID: "test-policy"}
	inst, err := e.CreateInstance(ctx, "wf-1", nil, pol)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	if inst.ID == "" {
		t.Error("expected instance ID")
	}
	if inst.State != StateRunning { // Default low materiality
		t.Errorf("expected RUNNING, got %s", inst.State)
	}

	// 3. Get Instance
	fetched, err := e.GetInstance(ctx, inst.ID)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}
	if fetched.ID != inst.ID {
		t.Errorf("fetched wrong ID: %s", fetched.ID)
	}

	// 4. Get Non-Existent
	_, err = e.GetInstance(ctx, "99999")
	if err == nil {
		t.Error("expected error for non-existent instance")
	}

	// 5. List Populated
	list, err = e.ListInstances(ctx)
	if err != nil {
		t.Fatalf("ListInstances failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 instance, got %d", len(list))
	}
}

func TestRecordDecision_NotFound(t *testing.T) {
	e := NewEngine()
	cmd := RecordDecisionCmd{
		InstanceID: "missing",
		Type:       DecisionApprove,
	}
	_, err := e.RecordDecision(context.Background(), cmd)
	if err == nil {
		t.Error("expected error for missing instance")
	}
}

func TestRecordDecision_InvalidState(t *testing.T) {
	e := NewEngine()
	// Create running instance (not waiting)
	inst, _ := e.CreateInstance(context.Background(), "wf-1", nil, policy.Policy{ID: "p1"})

	cmd := RecordDecisionCmd{
		InstanceID: inst.ID,
		Type:       DecisionApprove,
	}
	_, err := e.RecordDecision(context.Background(), cmd)
	if err == nil {
		t.Error("expected error for instance not in WAITING_FOR_HUMAN state")
	}
}
