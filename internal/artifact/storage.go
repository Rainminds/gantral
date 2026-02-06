package artifact

import (
	"context"
	"errors"

	"github.com/Rainminds/gantral/pkg/models"
)

var (
	// ErrArtifactAlreadyExists indicates a violation of WORM (Write-Once-Read-Many).
	// Evidence cannot be overwritten.
	ErrArtifactAlreadyExists = errors.New("artifact already exists (immutable)")

	// ErrArtifactNotFound indicates the requested artifact does not exist on the medium.
	ErrArtifactNotFound = errors.New("artifact not found")
)

// Store defines the interface for the immutable persistence layer.
// Implementations MUST ensure that writes are atomic and immutable.
type Store interface {
	// Write persists the artifact to the underlying storage.
	// It MUST fail if an artifact with the same ID already exists.
	// It MUST be atomic (no partial writes).
	Write(ctx context.Context, artifact *models.CommitmentArtifact) error

	// Get retrieves an artifact by its ID.
	Get(ctx context.Context, artifactID string) (*models.CommitmentArtifact, error)
}
