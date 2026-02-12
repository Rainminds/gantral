package replay_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/internal/replay"
	"github.com/Rainminds/gantral/internal/storage/local"
	"github.com/Rainminds/gantral/pkg/models"
)

func Test_Replay_Gaslight(t *testing.T) {
	// Setup: Real LocalStore in Temp Dir
	tmpDir := t.TempDir()
	store, err := local.NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	guard := replay.NewReplayGuard(store)
	ctx := context.Background()

	// 1. Establish Ground Truth (Legitimate Artifact)
	validArt := models.NewCommitmentArtifact(
		"inst-1", "prev", "APPROVED", "v1", "ctx-hash", "admin",
	)
	if err := validArt.CalculateHashAndSetID(); err != nil {
		t.Fatal(err)
	}
	if err := store.Write(ctx, validArt); err != nil {
		t.Fatal(err)
	}

	// 2. Scenario A: Honest Replay (Matches Ground Truth)
	// The workflow history claims "Artifact X" execution. Guard checks Store.
	if err := guard.ValidateReplay(ctx, validArt); err != nil {
		t.Errorf("Honest replay failed: %v", err)
	}

	// 3. Scenario B: Gaslight Attack (Fabricated Event)
	// Attacker injects an event into history saying "Attacker Approved",
	// but NO such artifact exists in the store.
	fakeArt := models.NewCommitmentArtifact(
		"inst-1", "prev", "APPROVED", "v1", "ctx-hash", "attacker",
	)
	_ = fakeArt.CalculateHashAndSetID()

	err = guard.ValidateReplay(ctx, fakeArt)
	if err == nil {
		t.Error("Guard accepted fabricated artifact (should fail-closed on missing evidence)")
	} else if !strings.Contains(err.Error(), "missing from store") {
		t.Errorf("Expected 'missing from store' error, got: %v", err)
	}

	// 4. Scenario C: Tampered Substitution (Integrity Failure)
	// Attacker tries to reuse a valid ID but changes the Claim (e.g. State).
	// Since ID is content-addressed, this MUST fail the self-integrity check first.
	tamperedArt := *validArt
	tamperedArt.AuthorityState = "REJECTED" // Content changed, but ID kept same

	err = guard.ValidateReplay(ctx, &tamperedArt)
	if err == nil {
		t.Error("Guard accepted tampered artifact (integrity check failed)")
	} else if !strings.Contains(err.Error(), "history integrity compromised") {
		t.Errorf("Expected 'history integrity compromised' error, got: %v", err)
	}
}
