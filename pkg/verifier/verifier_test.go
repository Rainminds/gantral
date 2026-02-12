package verifier_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/Rainminds/gantral/pkg/verifier"
	"github.com/stretchr/testify/assert"
)

func TestVerifyArtifact_Success(t *testing.T) {
	// 1. Create a valid artifact manually
	payload := map[string]interface{}{"foo": "bar"}
	payloadHash, _ := artifact.HashContext(payload)

	art := models.CommitmentArtifact{
		ArtifactVersion:  models.SchemaVersionV1,
		InstanceID:       "inst-1",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		AuthorityState:   "APPROVED",
		PolicyVersionID:  "v1",
		ContextHash:      payloadHash,
		PrevArtifactHash: "",
		HumanActorID:     "user-1",
	}

	// Calculate ID (Hash of content except ArtifactID)
	// We use the internal method to set it correctly
	err := art.CalculateHashAndSetID()
	assert.NoError(t, err)

	// 2. Write to temp file
	tmpDir, err := os.MkdirTemp("", "verifier-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, art.ArtifactID+".json")
	fullData, _ := json.Marshal(art)
	err = os.WriteFile(filePath, fullData, 0644)
	assert.NoError(t, err)

	// 3. Verify
	// VerifyArtifact takes []byte derived from file
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	res, err := verifier.VerifyArtifact(data)
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Equal(t, art.ArtifactID, res.ArtifactID)
}

func TestVerifyArtifact_InvalidJSON(t *testing.T) {
	res, err := verifier.VerifyArtifact([]byte("{invalid-json"))
	assert.NoError(t, err) // It returns a result with Valid=false, not an error
	assert.False(t, res.Valid)
	assert.Contains(t, res.Error, "invalid json")
}

func TestVerifyArtifact_InvalidHash(t *testing.T) {
	// 1. Create artifact with MISMATCHED ID
	art := models.CommitmentArtifact{
		ArtifactID:       "bad-hash",
		InstanceID:       "inst-1",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		AuthorityState:   "APPROVED",
		PolicyVersionID:  "v1",
		ContextHash:      "hash",
		PrevArtifactHash: "",
		HumanActorID:     "user-1",
	}

	// 3. Verify
	data, _ := json.Marshal(art)
	res, err := verifier.VerifyArtifact(data)
	assert.NoError(t, err)
	assert.False(t, res.Valid)
	assert.Contains(t, res.Error, "hash mismatch")
}

func TestVerifyChain_Success(t *testing.T) {
	// Chain: A -> B
	artA := models.CommitmentArtifact{
		ArtifactVersion:  models.SchemaVersionV1,
		InstanceID:       "inst-chain",
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
		AuthorityState:   "PENDING",
		ContextHash:      "hash-a",
		PrevArtifactHash: models.GenesisHash,
		HumanActorID:     "user-a",
	}
	_ = artA.CalculateHashAndSetID()

	artB := models.CommitmentArtifact{
		ArtifactVersion:  models.SchemaVersionV1,
		InstanceID:       "inst-chain",
		Timestamp:        time.Now().UTC().Add(1 * time.Second).Format(time.RFC3339),
		AuthorityState:   "APPROVED",
		ContextHash:      "hash-b",
		PrevArtifactHash: artA.ArtifactID, // Chain Link
		HumanActorID:     "user-b",
	}
	_ = artB.CalculateHashAndSetID()

	chain := []models.CommitmentArtifact{artA, artB}

	// Verify
	report := verifier.VerifyChain(chain)
	assert.True(t, report.Valid)
}

func TestVerifyChain_BrokenLink(t *testing.T) {
	// Chain: A -> B (but B points to wrong prev)
	artA := models.CommitmentArtifact{
		ArtifactID:       "hash-A",
		PrevArtifactHash: "genesis",
	}
	artB := models.CommitmentArtifact{
		ArtifactID:       "hash-B",
		PrevArtifactHash: "WRONG-HASH", // Link Broken
	}

	chain := []models.CommitmentArtifact{artA, artB}

	// Verify
	report := verifier.VerifyChain(chain)
	assert.False(t, report.Valid)
	assert.Equal(t, 1, report.BrokenIndex)
}

func toJson(v interface{}) []byte {
	d, _ := json.Marshal(v)
	return d
}
