package artifact

import (
	"errors"
	"fmt"

	"github.com/Rainminds/gantral/pkg/models"
)

var (
	// ErrArtifactSerialization indicates a failure to strictly encode the artifact.
	ErrArtifactSerialization = errors.New("artifact serialization failure")

	// ErrHashMismatch indicates a verification failure (used in verification logic).
	ErrHashMismatch = errors.New("artifact hash mismatch")

	// ErrInvalidInput indicates missing or malformed input arguments.
	ErrInvalidInput = errors.New("invalid artifact input argument")
)

// Manager implements the ArtifactEmitter interface.
// It manages the lifecycle of commitment artifacts.
type Manager struct {
	// In a real implementation, this would likely hold references to a storage backend
	// (e.g., ObjectStorage, KMS) to persist the artifact.
	// For Phase 6.1, we focus on the domain logic and in-memory generation.
}

// NewManager creates a new artifact manager.
func NewManager() *Manager {
	return &Manager{}
}

// EmitArtifact generates, seals, and calculates the ID for a new commitment artifact.
// It implements strict validation and fail-closed logic.
//
// Security Property: "EmitArtifact generates a non-repudiable proof of authorization bound to execution state."
func (m *Manager) EmitArtifact(
	instanceID string,
	prevHash string,
	state string,
	policyVer string,
	contextHash string,
	actorID string,
) (*models.CommitmentArtifact, error) {
	// 1. Fail-Closed Input Validation
	if instanceID == "" {
		return nil, fmt.Errorf("%w: instanceID required", ErrInvalidInput)
	}
	if state == "" {
		return nil, fmt.Errorf("%w: authority state required", ErrInvalidInput)
	}
	if contextHash == "" {
		return nil, fmt.Errorf("%w: context hash required", ErrInvalidInput)
	}
	// prevHash can be empty for genesis, so we don't strictly block it,
	// but we might want to enforce "0000..." for genesis in future iterations.

	// 2. Instantiate Model
	art := models.NewCommitmentArtifact(
		instanceID,
		prevHash,
		state,
		policyVer,
		contextHash,
		actorID,
	)

	// 3. Calculate Canonical Hash (The "Seal")
	// If this fails, strict fail-closed: we simply return error and NO artifact.
	if err := art.CalculateHashAndSetID(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArtifactSerialization, err)
	}

	// 4. In a full implementation, we would now:
	//    - Sign the artifact with KMS (Signature field).
	//    - Persist to WORM storage.
	// For Phase 6.1, returning the chemically pure domain object is the goal.

	return art, nil
}
