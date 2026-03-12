package artifact

import (
	"context"
	"encoding/hex"
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
	store       ImmutableStore
	signer      models.Signer
	timeStamper models.TimeStamper
}

// NewManager creates a new artifact manager with the given persistence store, signer, and time stamper.
func NewManager(store ImmutableStore, signer models.Signer, timeStamper models.TimeStamper) *Manager {
	return &Manager{
		store:       store,
		signer:      signer,
		timeStamper: timeStamper,
	}
}

// EmitArtifact generates, seals, and calculates the ID for a new commitment artifact.
// It implements strict validation and fail-closed logic.
//
// Security Property: "EmitArtifact generates a non-repudiable proof of authorization bound to execution state."
// EmitArtifact generates, seals, and persists a new commitment artifact.
// It implements strict validation, fail-closed logic, and atomic persistence.
//
// Security Property: "EmitArtifact generates a non-repudiable proof of authorization bound to execution state."
func (m *Manager) EmitArtifact(
	ctx context.Context,
	teamID string,
	instanceID string,
	workflowVersionID string,
	prevHash string,
	state string,
	policyVer string,
	contextHash string,
	actorID string,
	justification string,
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
	if teamID == "" {
		return nil, fmt.Errorf("%w: teamID required", ErrInvalidInput)
	}
	if actorID == "" {
		return nil, fmt.Errorf("%w: actorID required", ErrInvalidInput)
	}
	// prevHash can be empty for genesis, so we don't strictly block it,
	// but we might want to enforce "0000..." for genesis in future iterations.

	// 2. Instantiate Model
	art := models.NewCommitmentArtifact(
		instanceID,
		workflowVersionID,
		prevHash,
		state,
		policyVer,
		contextHash,
		actorID,
		justification,
	)

	// 3. Calculate Canonical Hash (The "Seal")
	// If this fails, strict fail-closed: we simply return error and NO artifact.
	if err := art.CalculateHash(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrArtifactSerialization, err)
	}

	hashBytes, err := hex.DecodeString(art.ArtifactHash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash for signing: %w", err)
	}

	// 4. Sign the Artifact
	if m.signer == nil {
		return nil, fmt.Errorf("fail-closed: strict signer configuration required")
	}
	sig, alg, err := m.signer.Sign(hashBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign artifact: %w", err)
	}
	art.ArtifactSignature = hex.EncodeToString(sig)
	art.SignatureAlgorithm = alg

	// 5. Timestamp the Artifact
	if m.timeStamper == nil {
		return nil, fmt.Errorf("fail-closed: strict timestamper configuration required")
	}
	token, alg, err := m.timeStamper.Token(hashBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to timestamp artifact: %w", err)
	}
	art.TimestampToken = hex.EncodeToString(token)
	art.TimestampAlgorithm = alg

	// 6. Persistence
	// We persist to the WORM storage before returning.
	// If persistence fails, we MUST fail the operation (Atomicity).
	if err := m.store.WriteArtifact(ctx, teamID, art); err != nil {
		return nil, fmt.Errorf("persistence failure: %w", err)
	}

	return art, nil
}
