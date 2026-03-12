package verifier

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Rainminds/gantral/pkg/constants"
	"github.com/Rainminds/gantral/pkg/models"
)

// VerificationStatus reflects TRD A.13 classification.
type VerificationStatus string

const (
	StatusValid        VerificationStatus = "VALID"
	StatusInvalid      VerificationStatus = "INVALID"
	StatusInconclusive VerificationStatus = "INCONCLUSIVE"
)

// VerificationResult contains the outcome of an artifact check.
type VerificationResult struct {
	Valid          bool               `json:"valid"` // maintained for backwards compatibility
	Status         VerificationStatus `json:"status"`
	ArtifactID     string             `json:"artifact_id"`
	CalculatedHash string             `json:"calculated_hash"`
	Error          string             `json:"error,omitempty"`
}

// ChainResult contains the outcome of a chain verification.
type ChainResult struct {
	Valid        bool               `json:"valid"` // maintained for backwards compatibility
	Status       VerificationStatus `json:"status"`
	BrokenIndex  int                `json:"broken_index"`
	BrokenReason string             `json:"broken_reason,omitempty"`
}

// Verifier implements TRD A.13 normative verification.
type Verifier struct {
	sigVerifier models.SignatureVerifier
	tsVerifier  models.TimestampVerifier
}

// New creates a new verifier with strictly required dependencies.
func New(sigVerifier models.SignatureVerifier, tsVerifier models.TimestampVerifier) *Verifier {
	return &Verifier{
		sigVerifier: sigVerifier,
		tsVerifier:  tsVerifier,
	}
}

// isTerminalState verifies if a state ends an authority chain.
func isTerminalState(s string) bool {
	return s == constants.StateTerminated || s == constants.StateCompleted || s == constants.StateRejected
}

// checkTransition validates the Section 4.2 sequence.
func checkTransition(from, to string) bool {
	// Genesis artifacts don't have a from state, they typically start at CREATED or RUNNING.
	if from == "" {
		return to == constants.StateCreated || to == constants.StateRunning
	}

	transitions := map[string]map[string]bool{
		constants.StateCreated:         {constants.StateRunning: true},
		constants.StateRunning:         {constants.StateWaitingForHuman: true, constants.StateCompleted: true, constants.StateTerminated: true},
		constants.StateWaitingForHuman: {constants.StateApproved: true, constants.StateRejected: true, constants.StateOverridden: true},
		constants.StateApproved:        {constants.StateResumed: true},
		constants.StateOverridden:      {constants.StateResumed: true},
		constants.StateResumed:         {constants.StateRunning: true},
	}

	if allowed, ok := transitions[from]; ok {
		return allowed[to]
	}
	return false
}

// VerifyArtifact validates a single artifact bytes payload following TRD A.13 1-9.
func (v *Verifier) VerifyArtifact(data []byte) (*VerificationResult, error) {
	// Step 1: Deserialize the artifact.
	var art models.CommitmentArtifact
	if err := json.Unmarshal(data, &art); err != nil {
		return &VerificationResult{Valid: false, Status: StatusInvalid, Error: fmt.Sprintf("invalid json structure: %v", err)}, nil
	}

	// Step 2: Validate schema completeness.
	if art.ArtifactID == "" || art.InstanceID == "" || art.AuthorityState == "" {
		return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "missing critical schema fields"}, nil
	}

	if !models.IsSupportedVersion(art.ArtifactVersion) {
		return &VerificationResult{Valid: false, Status: StatusInvalid, Error: fmt.Sprintf("%v: %s", models.ErrUnsupportedArtifactVersion, art.ArtifactVersion)}, nil
	}

	claimedHash := art.ArtifactHash

	// Step 3 & 4 & 5: Reconstruct canonical JSON payload & Recompute Hash
	checkArt := art
	if err := checkArt.CalculateHash(); err != nil {
		return &VerificationResult{
			Valid:      false,
			Status:     StatusInvalid,
			ArtifactID: art.ArtifactID,
			Error:      fmt.Sprintf("calculation failed: %v", err),
		}, nil
	}

	// Step 6: Verify recomputed hash
	if checkArt.ArtifactHash != claimedHash {
		return &VerificationResult{
			Valid:          false,
			Status:         StatusInvalid,
			ArtifactID:     art.ArtifactID,
			CalculatedHash: checkArt.ArtifactHash,
			Error:          "hash mismatch: integrity compromised",
		}, nil
	}

	hashBytes, err := hex.DecodeString(checkArt.ArtifactHash)
	if err != nil {
		return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "invalid hash encoding"}, nil
	}

	// Step 7: Verify signature
	if v.sigVerifier != nil {
		if art.ArtifactSignature == "" {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "missing signature"}, nil
		}
		sigBytes, err := hex.DecodeString(art.ArtifactSignature)
		if err != nil {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "invalid signature hex"}, nil
		}
		if err := v.sigVerifier.Verify(hashBytes, sigBytes, art.SignatureAlgorithm); err != nil {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "signature verification failed"}, nil
		}
	}

	// Step 8: Verify timestamp
	if v.tsVerifier != nil {
		if art.TimestampToken == "" {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "missing timestamp token"}, nil
		}
		tsBytes, err := hex.DecodeString(art.TimestampToken)
		if err != nil {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "invalid timestamp hex"}, nil
		}
		if err := v.tsVerifier.Verify(hashBytes, tsBytes, art.TimestampAlgorithm); err != nil {
			return &VerificationResult{Valid: false, Status: StatusInvalid, Error: "timestamp verification failed"}, nil
		}
	}

	// Step 11 logic for a single artifact
	status := StatusValid
	if !isTerminalState(art.AuthorityState) {
		status = StatusInconclusive
	}

	return &VerificationResult{
		Valid:          true,
		Status:         status,
		ArtifactID:     art.ArtifactID,
		CalculatedHash: checkArt.ArtifactHash,
	}, nil
}

// VerifyChain validates a sequence of artifacts.
func (v *Verifier) VerifyChain(chain []models.CommitmentArtifact) *ChainResult {
	if len(chain) == 0 {
		return &ChainResult{Valid: true, Status: StatusValid}
	}

	// Step 10: Check each link in the chain
	for i := 1; i < len(chain); i++ {
		prev := chain[i-1]
		curr := chain[i]

		// Array representation linkage Check
		// Note that `curr.PrevArtifactHash` might be typed as interface{} internally
		// if it handled nulls, therefore convert/evaluate string equivalence correctly.
		prevHashStr, _ := curr.PrevArtifactHash.(string)
		if prevHashStr != prev.ArtifactHash {
			return &ChainResult{
				Valid:        false,
				Status:       StatusInvalid,
				BrokenIndex:  i,
				BrokenReason: fmt.Sprintf("linkage broken: artifact[%d].prev (%v) != artifact[%d].hash (%s)", i, curr.PrevArtifactHash, i-1, prev.ArtifactHash),
			}
		}

		// Step 9: State Transition checking
		if !checkTransition(prev.AuthorityState, curr.AuthorityState) {
			return &ChainResult{
				Valid:        false,
				Status:       StatusInvalid,
				BrokenIndex:  i,
				BrokenReason: fmt.Sprintf("invalid transition from %s to %s", prev.AuthorityState, curr.AuthorityState),
			}
		}
	}

	// Need to check the first transition as well
	first := chain[0]
	if !checkTransition("", first.AuthorityState) {
		// some chains might start mid-way, assuming we are verifying full history
		// for testing let's not strictly fail a chain that starts on "RUNNING" as that's mapped in transition from "".
		return &ChainResult{
			Valid:        false,
			Status:       StatusInvalid,
			BrokenIndex:  0,
			BrokenReason: fmt.Sprintf("invalid genesis transition, chain starts with %s", first.AuthorityState),
		}
	}

	// Step 11: Classify based on terminal state
	finalArt := chain[len(chain)-1]
	status := StatusValid
	if !isTerminalState(finalArt.AuthorityState) {
		status = StatusInconclusive
	}

	return &ChainResult{Valid: true, Status: status}
}
