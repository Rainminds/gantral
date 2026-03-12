package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/infra"
	"github.com/stretchr/testify/assert"
)

func TestStore_Integration(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:changeme@127.0.0.1:5432/gantral?sslmode=disable"
	}

	ctx := context.Background()
	
	// Ensure schema is up to date
	err := infra.RunMigrations(dbURL)
	if err != nil {
		t.Skipf("Skipping Postgres integration test: migrations failed (is DB running?): %v", err)
	}

	store, err := NewStore(ctx, dbURL)
	if err != nil {
		t.Skipf("Skipping Postgres integration test: connection failed: %v", err)
	}
	defer store.Close()

	// 1. Create Instance
	instID := fmt.Sprintf("test-inst-%d", time.Now().UnixNano())
	inst := &engine.Instance{
		ID:              instID,
		OwningTeamID:    "team-1",
		WorkflowID:      "wf-1",
		State:           engine.StateRunning,
		PolicyVersionID: "v1",
	}
	err = store.CreateInstance(ctx, inst)
	assert.NoError(t, err)

	// 2. Get Instance
	fetched, err := store.GetInstance(ctx, instID)
	assert.NoError(t, err)
	assert.Equal(t, inst.ID, fetched.ID)
	assert.Equal(t, inst.OwningTeamID, fetched.OwningTeamID)

	// 3. List Instances
	list, err := store.ListInstances(ctx, "team-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	// 4. Record Decision
	cmd := engine.RecordDecisionCmd{
		InstanceID:    inst.ID,
		Type:          engine.DecisionApprove,
		ActorID:       "admin",
		Justification: "Manual Approval",
	}
	updated, err := store.RecordDecision(ctx, cmd, engine.StateApproved)
	assert.NoError(t, err)
	assert.Equal(t, engine.StateApproved, updated.State)

	// 5. Get Audit Events
	events, err := store.GetAuditEvents(ctx, inst.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
}
