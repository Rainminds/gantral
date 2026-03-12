package golden_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
)

func TestCanonicalSerialization_MatchesGolden(t *testing.T) {
	// Read the canonical golden struct
	goldenJSON, err := os.ReadFile(filepath.Join("canonical_artifact_v1.json"))
	if err != nil {
		t.Fatalf("Failed to read golden JSON: %v", err)
	}

	// Read expected hash
	expectedHashBytes, err := os.ReadFile(filepath.Join("expected_sha256_v1.txt"))
	if err != nil {
		t.Fatalf("Failed to read expected hash: %v", err)
	}
	expectedHash := string(expectedHashBytes)

	var art models.CommitmentArtifact
	if err := json.Unmarshal(goldenJSON, &art); err != nil {
		t.Fatalf("Failed to unmarshal golden JSON: %v", err)
	}

	err = art.CalculateHash()
	if err != nil {
		t.Fatalf("Calculate hash failed: %v", err)
	}

	if art.ArtifactHash != expectedHash {
		t.Errorf("Golden vector mismatch! \nExpected: %s\nGot: %s", expectedHash, art.ArtifactHash)
	}
}

func TestCanonicalSerialize_Rules(t *testing.T) {
	// TRD A.10.2: Lexicographic ascending order
	// TRD A.10.3: Compact JSON
	// TRD A.10.4: Stable Field Inclusion (null)
	// TRD A.10.5: String normalization (no HTML escape)
	// TRD A.10.6: Boolean lower case

	payload := map[string]interface{}{
		"z_bool_true":  true,
		"y_bool_false": false,
		"x_null":       nil,
		"w_string":     "value&<>",
		"v_number":     42.5,
		"u_int":        123,
	}

	canonicalBytes, err := models.CanonicalSerialize(payload)
	if err != nil {
		t.Fatalf("CanonicalSerialize failed: %v", err)
	}

	expectedCanonical := `{"u_int":123,"v_number":42.5,"w_string":"value&<>","x_null":null,"y_bool_false":false,"z_bool_true":true}`
	if string(canonicalBytes) != expectedCanonical {
		t.Errorf("Canonical mismatch.\nGot: %s\nExp: %s", string(canonicalBytes), expectedCanonical)
	}
}

func TestArtifact_CalculateHash_Determinism(t *testing.T) {
	art := models.NewCommitmentArtifact("i1", "wfv1", "prev", "APPROVED", "pv1", "ctx1", "act1", "j")

	err := art.CalculateHash()
	if err != nil {
		t.Fatal(err)
	}

	expectedHash := art.ArtifactHash

	// Re-calculate should result in the exact same hash
	art.ArtifactHash = ""
	err = art.CalculateHash()
	if err != nil {
		t.Fatal(err)
	}

	if art.ArtifactHash != expectedHash {
		t.Errorf("Hashing is non-deterministic")
	}
}

func TestArtifact_PayloadExcludesRestrictedFields(t *testing.T) {
	art := models.NewCommitmentArtifact("i1", "wfv1", "prev", "APPROVED", "pv1", "ctx1", "act1", "j")
	art.ArtifactSignature = "sig"
	art.TimestampToken = "token"

	payload, err := art.PayloadForHashing()
	if err != nil {
		t.Fatal(err)
	}

	// Verify exclusions
	if _, ok := payload["artifact_signature"]; ok {
		t.Error("artifact_signature should not be in hash payload")
	}
	if _, ok := payload["timestamp_token"]; ok {
		t.Error("timestamp_token should not be in hash payload")
	}
	if _, ok := payload["attestations"]; ok {
		t.Error("attestations should not be in hash payload")
	}

	// Verify prev_artifact_hash explicitly null if empty
	artGenesis := models.NewCommitmentArtifact("i1", "wfv1", models.GenesisHash, "APPROVED", "pv1", "ctx1", "act1", "j")
	genesisPayload, _ := artGenesis.PayloadForHashing()

	if val, ok := genesisPayload["prev_artifact_hash"]; !ok || val != nil {
		t.Errorf("Genesis prev_artifact_hash must be nil: %v", val)
	}
}
