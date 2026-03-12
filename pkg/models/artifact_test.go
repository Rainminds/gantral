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

	art := NewCommitmentArtifact(instanceID, "wf-v1", prevHash, state, policy, ctxHash, actor, "justification")

	if art.ArtifactVersion != CurrentArtifactVersion {
		t.Errorf("expected version %s, got %s", CurrentArtifactVersion, art.ArtifactVersion)
	}
	if art.InstanceID != instanceID {
		t.Errorf("expected instance %s, got %s", instanceID, art.InstanceID)
	}
	// Note: Timestamp is removed in v7.0, replaced by TimestampToken generated later.
}

func TestArtifact_PayloadForHashing(t *testing.T) {
	art := &CommitmentArtifact{
		ArtifactVersion:     CurrentArtifactVersion,
		InstanceID:          "inst-1",
		PrevArtifactHash:    GenesisHash,
		AuthorityState:      "APPROVED",
		PolicyVersionID:     "v1",
		ContextSnapshotHash: "ctx-1",
		HumanActorID:        "user-1",
		CryptoProfile:       "standard-v1",
		// Timestamp is no longer in the payload
	}

	t.Run("Valid Payload", func(t *testing.T) {
		payload, err := art.PayloadForHashing()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if payload["instance_id"] != "inst-1" {
			t.Errorf("expected inst-1, got %v", payload["instance_id"])
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
			{"Missing ContextSnapshotHash", func(a *CommitmentArtifact) { a.ContextSnapshotHash = "" }, "missing context_snapshot_hash"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				cp := *art
				tc.modify(&cp)
				_, err := cp.PayloadForHashing()
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
	art := NewCommitmentArtifact("inst-1", "wfv1", GenesisHash, "APPROVED", "v1", "ctx-1", "user-1", "just")
	err := art.CalculateHash()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(art.ArtifactHash) != 64 {
		t.Errorf("expected 64-char hex hash, got len %d", len(art.ArtifactHash))
	}
}

func TestArtifact_MarshalJSON(t *testing.T) {
	art := NewCommitmentArtifact("inst-1", "wf1", GenesisHash, "APPROVED", "v1", "ctx-1", "user-1", "just")
	if err := art.CalculateHash(); err != nil {
		t.Fatalf("failed to calculate hash: %v", err)
	}

	data, err := json.Marshal(art)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	required := []string{"artifact_version", "artifact_id", "instance_id"}
	for _, f := range required {
		if m[f] == "" || m[f] == nil {
			t.Errorf("missing field %s in JSON", f)
		}
	}
}
