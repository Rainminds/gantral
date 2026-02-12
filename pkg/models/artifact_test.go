package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewCommitmentArtifact(t *testing.T) {
	instanceID := "inst-123"
	prevHash := GenesisHash
	state := "APPROVED"
	policy := "v1.1"
	ctxHash := "ctx-hash-123"
	actor := "human-1"

	art := NewCommitmentArtifact(instanceID, prevHash, state, policy, ctxHash, actor)

	if art.ArtifactVersion != SchemaVersionV1 {
		t.Errorf("expected version %s, got %s", SchemaVersionV1, art.ArtifactVersion)
	}
	if art.InstanceID != instanceID {
		t.Errorf("expected instance %s, got %s", instanceID, art.InstanceID)
	}
	if art.Timestamp == "" {
		t.Fatal("expected timestamp to be set")
	}
}

func TestArtifact_CanonicalPayload(t *testing.T) {
	art := &CommitmentArtifact{
		ArtifactVersion:  SchemaVersionV1,
		InstanceID:       "inst-1",
		PrevArtifactHash: GenesisHash,
		AuthorityState:   "APPROVED",
		PolicyVersionID:  "v1",
		ContextHash:      "ctx-1",
		HumanActorID:     "user-1",
		Timestamp:        "2024-01-01T00:00:00Z",
	}

	t.Run("Valid Payload", func(t *testing.T) {
		payload, err := art.CanonicalPayload()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var m map[string]string
		if err := json.Unmarshal(payload, &m); err != nil {
			t.Fatalf("failed to unmarshal payload: %v", err)
		}
		if m["instance_id"] != "inst-1" {
			t.Errorf("expected inst-1, got %s", m["instance_id"])
		}
	})

	t.Run("Missing Fields", func(t *testing.T) {
		cases := []struct {
			name     string
			modify   func(*CommitmentArtifact)
			expected string
		}{
			{"Missing InstanceID", func(a *CommitmentArtifact) { a.InstanceID = "" }, "missing instance_id"},
			{"Missing AuthorityState", func(a *CommitmentArtifact) { a.AuthorityState = "" }, "missing authority_state"},
			{"Missing ContextHash", func(a *CommitmentArtifact) { a.ContextHash = "" }, "missing context_hash"},
			{"Missing Timestamp", func(a *CommitmentArtifact) { a.Timestamp = "" }, "missing timestamp"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				cp := *art
				tc.modify(&cp)
				_, err := cp.CanonicalPayload()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.expected) {
					t.Errorf("expected error containing %q, got %q", tc.expected, err.Error())
				}
			})
		}
	})
}

func TestArtifact_CalculateHashAndSetID(t *testing.T) {
	art := NewCommitmentArtifact("inst-1", GenesisHash, "APPROVED", "v1", "ctx-1", "user-1")
	err := art.CalculateHashAndSetID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(art.ArtifactID) != 64 {
		t.Errorf("expected 64-char hex hash, got len %d", len(art.ArtifactID))
	}
}

func TestArtifact_MarshalJSON(t *testing.T) {
	art := NewCommitmentArtifact("inst-1", GenesisHash, "APPROVED", "v1", "ctx-1", "user-1")
	art.CalculateHashAndSetID()

	data, err := json.Marshal(art)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	required := []string{"artifact_version", "artifact_id", "instance_id", "timestamp"}
	for _, f := range required {
		if m[f] == "" {
			t.Errorf("missing field %s in JSON", f)
		}
	}
}
