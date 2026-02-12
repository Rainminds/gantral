package verifier

import (
	"encoding/json"
	"fmt"

	"github.com/Rainminds/gantral/pkg/models"
)

// VerificationResult contains the outcome of an artifact check.
type VerificationResult struct {
	Valid          bool   `json:"valid"`
	ArtifactID     string `json:"artifact_id"`
	CalculatedHash string `json:"calculated_hash"`
	Error          string `json:"error,omitempty"`
}

// ChainResult contains the outcome of a chain verification.
type ChainResult struct {
	Valid        bool   `json:"valid"`
	BrokenIndex  int    `json:"broken_index"`
	BrokenReason string `json:"broken_reason,omitempty"`
}

// VerifyArtifact validates the integrity of a single artifact blob.
// It checks:
// 1. JSON structure validity.
// 2. Hash consistency (Claimed ID == Calculated Hash).
func VerifyArtifact(data []byte) (*VerificationResult, error) {
	// 1. Deserialize
	var art models.CommitmentArtifact
	if err := json.Unmarshal(data, &art); err != nil {
		return &VerificationResult{
			Valid: false,
			Error: fmt.Sprintf("invalid json structure: %v", err),
		}, nil
	}

	// 2. Check for missing critical fields (sanity check)
	if art.ArtifactID == "" {
		return &VerificationResult{Valid: false, Error: "missing artifact_id"}, nil
	}

	// 3. Re-calculate Hash
	// We use the canonical payload logic from pkg/models to ensure strict adherence to the spec.
	// NOTE: We work observing the object *as loaded*.
	claimedID := art.ArtifactID

	// Create a copy to recalculate to avoid mutating the source if it differed
	checkArt := art
	if err := checkArt.CalculateHashAndSetID(); err != nil {
		return &VerificationResult{
			Valid:      false,
			ArtifactID: claimedID,
			Error:      fmt.Sprintf("calculation failed: %v", err),
		}, nil
	}

	// 4. Compare
	if checkArt.ArtifactID != claimedID {
		return &VerificationResult{
			Valid:          false,
			ArtifactID:     claimedID,
			CalculatedHash: checkArt.ArtifactID,
			Error:          "hash mismatch: integrity compromised",
		}, nil
	}

	return &VerificationResult{
		Valid:          true,
		ArtifactID:     claimedID,
		CalculatedHash: checkArt.ArtifactID,
	}, nil
}

// VerifyChain validates a sequence of artifacts.
// The artifacts MUST be sorted by the caller (usually by timestamp or linkage).
// This function verifies strict cryptographic linkage: art[N].PrevHash == art[N-1].Hash.
func VerifyChain(chain []models.CommitmentArtifact) *ChainResult {
	if len(chain) == 0 {
		return &ChainResult{Valid: true}
	}

	// Verify links
	for i := 1; i < len(chain); i++ {
		prev := chain[i-1]
		curr := chain[i]

		// Strict Linkage Check
		if curr.PrevArtifactHash != prev.ArtifactID {
			return &ChainResult{
				Valid:        false,
				BrokenIndex:  i,
				BrokenReason: fmt.Sprintf("linkage broken: artifact[%d].prev (%s) != artifact[%d].id (%s)", i, curr.PrevArtifactHash, i-1, prev.ArtifactID),
			}
		}
	}

	return &ChainResult{Valid: true}
}
