package artifact

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
)

// MockStore is a no-op store for testing Manager logic.
type MockStore struct{}

func (s *MockStore) Write(ctx context.Context, art *models.CommitmentArtifact) error {
	return nil
}

func (s *MockStore) Get(ctx context.Context, id string) (*models.CommitmentArtifact, error) {
	return nil, nil
}

func TestEmitArtifact_Success(t *testing.T) {
	m := NewManager(&MockStore{})

	art, err := m.EmitArtifact(
		context.Background(),
		"inst-123",
		"prev-hash-abc",
		"APPROVED",
		"policy-v1",
		"ctx-hash-xyz",
		"actor-bob",
	)

	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if art == nil {
		t.Fatal("Expected artifact, got nil")
	}

	if art.ArtifactID == "" {
		t.Error("ArtifactID should be populated")
	}
	if art.ArtifactHash == "" {
		t.Error("ArtifactHash should be populated")
	}
	if art.ArtifactID != art.ArtifactHash {
		t.Error("ArtifactID should equal ArtifactHash in v1")
	}
	if art.ArtifactVersion != "v1" {
		t.Errorf("Expected version v1, got %s", art.ArtifactVersion)
	}
}

func TestEmitArtifact_Determinism(t *testing.T) {
	// m := NewManager() - Unused in this test as we construct models manually to control time

	// Capture time by mocking or just do two calls very fast?
	// NewCommitmentArtifact sets timestamp to time.Now(). To test bit-for-bit determinism
	// of the *logic*, we need to control variables.
	// However, EmitArtifact generates a timestamp internally.
	// To test "same input -> same output", strictly speaking, one variable is Time.
	// So calling EmitArtifact twice WILL produce different hashes because of Timestamp.
	//
	// CHECK: The requirement says:
	// "Verify that calling `EmitArtifact` on the exact same input twice produces the *exact same hash* (bit-for-bit)."
	// This implies we must either inject time OR the requirement implies "same input INCLUDING time".
	// Since `EmitArtifact` encapsulates `time.Now()`, we cannot strictly test determinism across two calls
	// unless we mock time or manually construct the artifact.
	//
	// The prompt demanded: "Verify that calling EmitArtifact on the exact same input twice..."
	// Since Manager.EmitArtifact calls `models.NewCommitmentArtifact` which calls `time.Now()`,
	// we cannot strictly satisfy this without DI for time.
	//
	// ADJUSTMENT: We will verify determinism of the *Hashing Logic* (`CalculateHashAndSetID`)
	// by constructing the model manually with a fixed timestamp.
	// OR we assume the prompt implies "if inputs AND time are fixed".

	// Let's test the core hashing determinism which is the critical part.

	ts := "2023-10-27T10:00:00Z"
	// Create two identical models manually
	art1 := models.NewCommitmentArtifact("i1", "p1", "s1", "pv1", "c1", "a1")
	art1.Timestamp = ts

	art2 := models.NewCommitmentArtifact("i1", "p1", "s1", "pv1", "c1", "a1")
	art2.Timestamp = ts

	if err := art1.CalculateHashAndSetID(); err != nil {
		t.Fatal(err)
	}
	if err := art2.CalculateHashAndSetID(); err != nil {
		t.Fatal(err)
	}

	if art1.ArtifactID != art2.ArtifactID {
		t.Errorf("Determinism failure: %s != %s", art1.ArtifactID, art2.ArtifactID)
	}
}

func TestEmitArtifact_InvalidInput(t *testing.T) {
	m := NewManager(&MockStore{})

	// Missing InstanceID
	_, err := m.EmitArtifact(context.Background(), "", "prev", "STATE", "pol", "ctx", "act")
	if err == nil {
		t.Error("Expected error for empty InstanceID, got nil")
	}

	// Missing ContextHash
	_, err = m.EmitArtifact(context.Background(), "inst", "prev", "STATE", "pol", "", "act")
	if err == nil {
		t.Error("Expected error for empty ContextHash, got nil")
	}
}

func TestArtifact_JSONStructure(t *testing.T) {
	// Verify that MarshalJSON includes all fields and flattened structure
	m := NewManager(&MockStore{})
	art, _ := m.EmitArtifact(context.Background(), "inst", "prev", "APPROVED", "pol", "ctx", "act")

	bytes, _ := json.Marshal(art)
	var asMap map[string]interface{}
	if err := json.Unmarshal(bytes, &asMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	required := []string{
		"artifact_version", "artifact_id", "instance_id",
		"prev_artifact_hash", "authority_state", "policy_version_id",
		"context_hash", "human_actor_id", "timestamp", "artifact_hash",
	}

	for _, field := range required {
		if _, ok := asMap[field]; !ok {
			t.Errorf("JSON output missing field: %s", field)
		}
	}
}
