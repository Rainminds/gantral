package artifact

import (
	"context"
	"errors"

	"github.com/Rainminds/gantral/pkg/models"
)

var (
	// ErrImmutableViolation indicates a violation of WORM (Write-Once-Read-Many).
	// Evidence cannot be overwritten.
	ErrImmutableViolation = errors.New("artifact immutability violated: artifact already exists")

	// ErrArtifactNotFound indicates the requested artifact does not exist on the medium.
	ErrArtifactNotFound = errors.New("artifact not found")
)

// ImmutableStore defines the interface for the authoritative append-only storage layer.
// Implementations MUST ensure that writes are atomic, fail-closed, and immutable.
type ImmutableStore interface {
	// WriteArtifact persists the artifact to the underlying WORM storage.
	// It MUST fail synchronously with ErrImmutableViolation if an artifact with the same ID already exists.
	// It MUST be atomic (no partial writes).
	WriteArtifact(ctx context.Context, teamID string, artifact *models.CommitmentArtifact) error

	// GetArtifact retrieves an artifact by its hash.
	GetArtifact(ctx context.Context, teamID string, artifactHash string) (*models.CommitmentArtifact, error)
}
