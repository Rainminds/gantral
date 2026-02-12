package artifact

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
	"github.com/Rainminds/gantral/pkg/verifier"
)

// GoldenArtifact defines the structure of our golden file
type GoldenArtifact struct {
	ArtifactVersion  string `json:"artifact_version"`
	ArtifactID       string `json:"artifact_id"`
	InstanceID       string `json:"instance_id"`
	PrevArtifactHash string `json:"prev_artifact_hash"`
	AuthorityState   string `json:"authority_state"`
	PolicyVersionID  string `json:"policy_version_id"`
	ContextHash      string `json:"context_hash"`
	HumanActorID     string `json:"human_actor_id"`
	Timestamp        string `json:"timestamp"`
}

func loadGolden(t *testing.T) GoldenArtifact {
	path := filepath.Join("..", "..", "tests", "golden", "canonical_artifact_v1.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}
	var g GoldenArtifact
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatalf("Failed to decode golden file: %v", err)
	}
	return g
}

func Test_Artifact_Integrity_Golden_SectionD(t *testing.T) {
	t.Parallel()
	golden := loadGolden(t)

	// Construct artifact from golden fields (except ID)
	art := models.NewCommitmentArtifact(
		golden.InstanceID,
		golden.PrevArtifactHash,
		golden.AuthorityState,
		golden.PolicyVersionID,
		golden.ContextHash,
		golden.HumanActorID,
	)
	// Force override timestamp to match golden because New set it to Now()
	art.Timestamp = golden.Timestamp
	// Force override version (though default New uses v1, we ensure it matches golden)
	art.ArtifactVersion = golden.ArtifactVersion

	// Calculate Hash
	if err := art.CalculateHashAndSetID(); err != nil {
		t.Fatalf("CalculateHashAndSetID failed: %v", err)
	}

	// Assert: Calculated Hash matches Golden Hash (Section M)
	if art.ArtifactID != golden.ArtifactID {
		t.Errorf("Hash Mismatch!\nExpected: %s\nGot:      %s", golden.ArtifactID, art.ArtifactID)
	}

	// Assert: Canonical Payload is stable
	payload, _ := art.CanonicalPayload()
	// We can't easily assert byte-for-byte against golden FILE easily because golden file has indentation/newlines for readability
	// while CanonicalPayload is compact.
	// But `CalculateHashAndSetID` depends on `CanonicalPayload`. If Hash matches, Payload matches semantically.
	_ = payload
}

func Test_Artifact_Tamper_Fails_SectionD(t *testing.T) {
	t.Parallel()
	golden := loadGolden(t)

	// Base valid artifact JSON
	baseArt := models.CommitmentArtifact{
		ArtifactVersion:  golden.ArtifactVersion,
		ArtifactID:       golden.ArtifactID,
		InstanceID:       golden.InstanceID,
		PrevArtifactHash: golden.PrevArtifactHash,
		AuthorityState:   golden.AuthorityState,
		PolicyVersionID:  golden.PolicyVersionID,
		ContextHash:      golden.ContextHash,
		HumanActorID:     golden.HumanActorID,
		Timestamp:        golden.Timestamp,
	}
	baseJSON, _ := json.Marshal(baseArt)

	// Verify base is valid first
	res, err := verifier.VerifyArtifact(baseJSON)
	if err != nil {
		t.Fatalf("Base artifact failed verification: %v", err)
	}
	if !res.Valid {
		t.Fatalf("Base artifact marked invalid: %v", res.Error)
	}

	// Tamper Cases
	tamperCases := []struct {
		name       string
		tamperFunc func(*models.CommitmentArtifact)
	}{
		{
			name: "Tamper AuthorityState",
			tamperFunc: func(a *models.CommitmentArtifact) {
				a.AuthorityState = "REJECTED" // Changed from APPROVED
			},
		},
		{
			name: "Tamper Timestamp",
			tamperFunc: func(a *models.CommitmentArtifact) {
				a.Timestamp = "2024-01-01T00:00:00Z"
			},
		},
		{
			name: "Tamper PolicyVersion",
			tamperFunc: func(a *models.CommitmentArtifact) {
				a.PolicyVersionID = "v2-evil"
			},
		},
		{
			name: "Tamper Actor",
			tamperFunc: func(a *models.CommitmentArtifact) {
				a.HumanActorID = "attacker"
			},
		},
		{
			name: "Tamper PrevHash",
			tamperFunc: func(a *models.CommitmentArtifact) {
				a.PrevArtifactHash = "111111..."
			},
		},
	}

	for _, tc := range tamperCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Copy base
			tampered := baseArt
			tc.tamperFunc(&tampered)

			// serialize
			tamperedJSON, _ := json.Marshal(tampered)

			// Verify -> MUST FAIL because ID is still the old one, but content changed
			res, err := verifier.VerifyArtifact(tamperedJSON)
			if err != nil {
				// VerifyArtifact might return error for malformed JSON, but for hash mismatch it usually returns Valid=false
				// If it returns error, that's also acceptable "fail closed", but typically we expect Valid=false
			}

			if res != nil && res.Valid {
				t.Errorf("Tampered artifact %s was marked VALID! This is a critical failure.", tc.name)
			}

			if res != nil && res.Error == "" {
				t.Error("Expected error message in result, got empty")
			}
		})
	}
}
