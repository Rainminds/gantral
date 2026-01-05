package tests

import (
	"context"
	"os"
	"testing"

	"github.com/Rainminds/gantral/adapters/secondary/postgres"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"github.com/Rainminds/gantral/infra"
	"github.com/joho/godotenv"
)

func TestIntegration_HITL_Flow(t *testing.T) {
	// 0. Setup
	_ = godotenv.Load("../.env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Run Migrations ensure schema exists
	if err := infra.RunMigrations(dbURL); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	ctx := context.Background()
	store, err := postgres.NewStore(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer store.Close()

	// Clean up tables
	// Note: In a real environment, we'd use a separate test DB or schema.
	// For this demo context, assuming dev DB is fine to truncate or specific IDs used.
	// Let's rely on unique IDs to avoid collision, validation is based on ID.

	// 1. Initialize Engine with Real Store
	eng := engine.NewEngine(store)

	// 2. Create High Materiality Instance (expect PAUSE)
	pol := policy.Policy{ID: "demo-policy", Materiality: policy.MaterialityHigh}

	inst, err := eng.CreateInstance(ctx, "integration-wf", map[string]interface{}{"key": "val"}, pol)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	t.Logf("Created Instance: %s", inst.ID)

	if inst.State != engine.StateWaitingForHuman {
		t.Errorf("expected WAITING_FOR_HUMAN, got %s", inst.State)
	}

	// 3. Verify Persistence (Fetch from Store directly)
	fetched, err := store.GetInstance(ctx, inst.ID)
	if err != nil {
		t.Fatalf("failed to fetch instance from store: %v", err)
	}
	if fetched.State != engine.StateWaitingForHuman {
		t.Errorf("store persisted wrong state: %s", fetched.State)
	}

	// 4. Record Decision (Approve)
	cmd := engine.RecordDecisionCmd{
		InstanceID:    inst.ID,
		Type:          engine.DecisionApprove,
		ActorID:       "integration-tester",
		Justification: "Integration Test Approval",
	}

	updated, err := eng.RecordDecision(ctx, cmd)
	if err != nil {
		t.Fatalf("RecordDecision failed: %v", err)
	}

	if updated.State != engine.StateApproved {
		t.Errorf("expected APPROVED, got %s", updated.State)
	}

	// 5. Verify Decision Persistence
	// We don't have GetDecisions exposed on Engine or Store public interface for the test yet,
	// unless we use the Store implementation specific query or check side effects on instance state.
	// Instance state update confirms the transaction succeeded.

	// Double check fetching again
	finalFetch, _ := store.GetInstance(ctx, inst.ID)
	if finalFetch.State != engine.StateApproved {
		t.Errorf("final persisted state incorrect: %s", finalFetch.State)
	}
}
