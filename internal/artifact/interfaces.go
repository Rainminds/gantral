package artifact

import (
	"context"

	"github.com/Rainminds/gantral/pkg/models"
)

// ArtifactEmitter defines the contract for emitting immutable commitment artifacts.
// It is the primary interface for the "Evidence Emitter" component required by Phase 6.
type ArtifactEmitter interface {
	// EmitArtifact generates, seals, and persists a commitment artifact for a state transition.
	// It requires the previous artifact's hash to form a chain (Merkle Log).
	//
	// Parameters:
	//   - instanceID: The ID of the execution instance.
	//   - prevHash: The hash of the antecedent artifact. Use "" for Genesis.
	//   - state: The target authority state (APPROVED, REJECTED, OVERRIDDEN).
	//   - policyVer: The version of the policy applied.
	//   - contextHash: The SHA256 digest of the execution context.
	//   - actorID: The ID of the authority (human or system).
	//
	// Returns:
	//   - *models.CommitmentArtifact: The sealed artifact.
	//   - error: If any step fails (serialization, IO, etc.). Fail-Closed.
	EmitArtifact(
		ctx context.Context,
		instanceID string,
		prevHash string,
		state string,
		policyVer string,
		contextHash string,
		actorID string,
	) (*models.CommitmentArtifact, error)
}
