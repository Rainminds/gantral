package verifier_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/Rainminds/gantral/pkg/verifier"
	"github.com/stretchr/testify/assert"
)

type MockSigVerifier struct{}

func (m *MockSigVerifier) Verify(hash, sig []byte, alg string) error {
	return nil
}

type MockTsVerifier struct{}

func (m *MockTsVerifier) Verify(hash, token []byte, alg string) error {
	return nil
}

func TestVerifyArtifact_Success(t *testing.T) {
	// 1. Create a valid artifact manually
	payload := map[string]interface{}{"foo": "bar"}
	payloadHash, _ := artifact.HashContext(payload)

	artPtr := models.NewCommitmentArtifact("inst-1", "wfv1", "", "COMPLETED", "v1", payloadHash, "user-1", "just")
	artPtr.ArtifactID = "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
	artPtr.TimestampToken = "746f6b656e" // valid hex mock
	artPtr.TimestampAlgorithm = "mock-alg"
	artPtr.ArtifactSignature = "736967" // valid hex mock
	artPtr.SignatureAlgorithm = "mock-alg"

	err := artPtr.CalculateHash()
	assert.NoError(t, err)

	art := *artPtr

	tmpDir, err := os.MkdirTemp("", "verifier-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, art.ArtifactID+".json")
	fullData, _ := json.Marshal(art)
	err = os.WriteFile(filePath, fullData, 0644)
	assert.NoError(t, err)

	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	v := verifier.New(&MockSigVerifier{}, &MockTsVerifier{})
	res, err := v.VerifyArtifact(data)
	assert.NoError(t, err)
	assert.True(t, res.Valid)
	assert.Equal(t, verifier.StatusValid, res.Status)
	assert.Equal(t, art.ArtifactID, res.ArtifactID)
}

func TestVerifyArtifact_Inconclusive(t *testing.T) {
	art := models.CommitmentArtifact{
		ArtifactVersion:     models.CurrentArtifactVersion,
		ArtifactID:          "a1",
		InstanceID:          "inst-1",
		WorkflowVersionID:   "w",
		AuthorityState:      "RUNNING", // Non-terminal
		ContextSnapshotHash: "c",
	}
	err := art.CalculateHash()
	assert.NoError(t, err)

	data, _ := json.Marshal(art)
	v := verifier.New(nil, nil)
	res, _ := v.VerifyArtifact(data)

	assert.True(t, res.Valid)
	assert.Equal(t, verifier.StatusInconclusive, res.Status)
}

func TestVerifyArtifact_InvalidJSON(t *testing.T) {
	v := verifier.New(nil, nil)
	res, err := v.VerifyArtifact([]byte("{invalid-json"))
	assert.NoError(t, err)
	assert.False(t, res.Valid)
	assert.Equal(t, verifier.StatusInvalid, res.Status)
}

func TestVerifyChain_Success(t *testing.T) {
	artA := models.CommitmentArtifact{
		ArtifactVersion:     models.CurrentArtifactVersion,
		InstanceID:          "inst",
		AuthorityState:      "RUNNING", // Transition "" -> RUNNING is ok
		ContextSnapshotHash: "hash-a",
		PrevArtifactHash:    nil,
	}
	errA := artA.CalculateHash()
	assert.NoError(t, errA)

	artB := models.CommitmentArtifact{
		ArtifactVersion:     models.CurrentArtifactVersion,
		InstanceID:          "inst",
		AuthorityState:      "WAITING_FOR_HUMAN", // RUNNING -> WAITING
		ContextSnapshotHash: "hash-b",
		PrevArtifactHash:    artA.ArtifactHash,
	}
	errB := artB.CalculateHash()
	assert.NoError(t, errB)

	chain := []models.CommitmentArtifact{artA, artB}

	v := verifier.New(nil, nil)
	report := v.VerifyChain(chain)
	assert.True(t, report.Valid)
	assert.Equal(t, verifier.StatusInconclusive, report.Status) // Ends on WAITING
}

func TestVerifyChain_BrokenLink(t *testing.T) {
	artA := models.CommitmentArtifact{
		ArtifactVersion:  models.CurrentArtifactVersion,
		ArtifactID:       "hash-some",
		PrevArtifactHash: "genesis",
	}
	artB := models.CommitmentArtifact{
		ArtifactVersion:  models.CurrentArtifactVersion,
		ArtifactID:       "hash-b",
		PrevArtifactHash: "WRONG-HASH",
	}

	chain := []models.CommitmentArtifact{artA, artB}
	v := verifier.New(nil, nil)
	report := v.VerifyChain(chain)
	assert.False(t, report.Valid)
	assert.Equal(t, 1, report.BrokenIndex)
}

func TestVerifyChain_InvalidTransition(t *testing.T) {
	artA := models.CommitmentArtifact{
		ArtifactVersion:     models.CurrentArtifactVersion,
		InstanceID:          "inst-1",
		AuthorityState:      "COMPLETED", // Terminal
		ContextSnapshotHash: "hash-a",
	}
	errA := artA.CalculateHash()
	assert.NoError(t, errA)

	artB := models.CommitmentArtifact{
		ArtifactVersion:  models.CurrentArtifactVersion,
		InstanceID:       "inst-1",
		AuthorityState:   "RUNNING", // Cannot RUN after COMPLETED
		PrevArtifactHash: artA.ArtifactHash,
		ContextSnapshotHash: "hash-b", // Added required field
	}
	errB := artB.CalculateHash()
	assert.NoError(t, errB)

	chain := []models.CommitmentArtifact{artA, artB}
	v := verifier.New(nil, nil)
	report := v.VerifyChain(chain)
	assert.False(t, report.Valid)
	assert.Equal(t, verifier.StatusInvalid, report.Status)
}

func TestVerifyChain_InvalidGenesis(t *testing.T) {
	artA := models.CommitmentArtifact{
		ArtifactVersion:     models.CurrentArtifactVersion,
		InstanceID:          "inst-1",
		AuthorityState:      "APPROVED", // Cannot start chain with APPROVED
		ContextSnapshotHash: "hash-a",
	}
	err := artA.CalculateHash()
	assert.NoError(t, err)

	chain := []models.CommitmentArtifact{artA}
	v := verifier.New(nil, nil)
	report := v.VerifyChain(chain)

	assert.False(t, report.Valid)
	assert.Equal(t, verifier.StatusInvalid, report.Status)
	assert.Contains(t, report.BrokenReason, "invalid genesis transition")
}

func TestReplayDeterminism(t *testing.T) {
	v := verifier.New(nil, nil)

	t.Run("Scenario 1: VALID Chain", func(t *testing.T) {
		// CREATED -> RUNNING -> COMPLETED
		art1 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "CREATED",
			ContextSnapshotHash: "h1",
		}
		if err := art1.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		art2 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "RUNNING",
			PrevArtifactHash:    art1.ArtifactHash,
			ContextSnapshotHash: "h2",
		}
		if err := art2.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		art3 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "COMPLETED",
			PrevArtifactHash:    art2.ArtifactHash,
			ContextSnapshotHash: "h3",
		}
		if err := art3.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		chain := []models.CommitmentArtifact{art1, art2, art3}
		res := v.VerifyChain(chain)
		assert.True(t, res.Valid)
		assert.Equal(t, verifier.StatusValid, res.Status)
	})

	t.Run("Scenario 2: INVALID - Hash Mismatch", func(t *testing.T) {
		art1 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			ArtifactID:          "id-1",
			InstanceID:          "inst-1",
			AuthorityState:      "CREATED",
			ContextSnapshotHash: "h1",
		}
		err := art1.CalculateHash()
		assert.NoError(t, err)
		originalHash := art1.ArtifactHash

		// Tamper with payload
		art1.ContextSnapshotHash = "tampared-hash"
		// DO NOT recalculate hash, just simulate a mismatched data on disk/storage
		art1.ArtifactHash = originalHash

		data, _ := json.Marshal(art1)
		res, err := v.VerifyArtifact(data)
		assert.NoError(t, err)
		assert.False(t, res.Valid)
		assert.Equal(t, verifier.StatusInvalid, res.Status)
		assert.Contains(t, res.Error, "hash mismatch")
	})

	t.Run("Scenario 3: INVALID - Illegal Transition", func(t *testing.T) {
		// APPROVED -> COMPLETED (Invalid, must be APPROVED -> RESUMED)
		art1 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "APPROVED",
			ContextSnapshotHash: "h1",
		}
		if err := art1.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		art2 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "COMPLETED",
			PrevArtifactHash:    art1.ArtifactHash,
			ContextSnapshotHash: "h2",
		}
		err2 := art2.CalculateHash()
		assert.NoError(t, err2)

		chain := []models.CommitmentArtifact{art1, art2}
		res := v.VerifyChain(chain)
		assert.False(t, res.Valid)
		assert.Equal(t, verifier.StatusInvalid, res.Status)
		assert.Contains(t, res.BrokenReason, "invalid transition")
	})

	t.Run("Scenario 4: INCONCLUSIVE - Non-Terminal State", func(t *testing.T) {
		// CREATED -> RUNNING -> WAITING_FOR_HUMAN
		art1 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "CREATED",
			ContextSnapshotHash: "h1",
		}
		if err := art1.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		art2 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "RUNNING",
			PrevArtifactHash:    art1.ArtifactHash,
			ContextSnapshotHash: "h2",
		}
		if err := art2.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		art3 := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			InstanceID:          "inst-1",
			AuthorityState:      "WAITING_FOR_HUMAN",
			PrevArtifactHash:    art2.ArtifactHash,
			ContextSnapshotHash: "h3",
		}
		if err := art3.CalculateHash(); err != nil {
			t.Fatalf("CalculateHash failed: %v", err)
		}

		chain := []models.CommitmentArtifact{art1, art2, art3}
		res := v.VerifyChain(chain)
		assert.True(t, res.Valid)
		assert.Equal(t, verifier.StatusInconclusive, res.Status)
	})
}
