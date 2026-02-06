package local

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
)

// safeIDRegex ensures artifact IDs contain only safe characters to prevent path traversal.
var safeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// Store implements artifact.Store using the local filesystem.
type Store struct {
	basePath string
	mu       sync.RWMutex // Protects concurrent access if needed (though FS handles locking mostly)
}

// NewStore creates a new local filesystem store.
// It ensures the base directory exists.
func NewStore(basePath string) (*Store, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create artifact storage directory: %w", err)
	}
	return &Store{basePath: basePath}, nil
}

// Write persists the artifact atomically and consistently (WORM).
// Steps:
// 1. Check if target exists (Fail if yes).
// 2. Write to temp file.
// 3. Sync to disk.
// 4. Atomic Rename.
func (s *Store) Write(ctx context.Context, art *models.CommitmentArtifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !safeIDRegex.MatchString(art.ArtifactID) {
		return fmt.Errorf("invalid artifact ID format: %s", art.ArtifactID)
	}

	targetPath := filepath.Join(s.basePath, art.ArtifactID+".json")

	// 1. Immutability Check (WORM)
	if _, err := os.Stat(targetPath); err == nil {
		return artifact.ErrArtifactAlreadyExists
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check artifact existence: %w", err)
	}

	// Serialize
	data, err := json.Marshal(art)
	if err != nil {
		return fmt.Errorf("serialization failed: %w", err)
	}

	// 2. Write to Temp File (Atomicity)
	// We create the temp file in the same directory to ensure atomic rename works (same partition).
	tmpFile, err := os.CreateTemp(s.basePath, "artifact-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()

	// Cleanup temp file on error
	defer func() {
		if tmpFile != nil {
			tmpFile.Close()    // Close if not already closed
			os.Remove(tmpName) // Best effort remove
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// 3. Fsync (Durability)
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tmpFile = nil // Prevent double close in defer

	// 4. Atomic Rename
	if err := os.Rename(tmpName, targetPath); err != nil {
		return fmt.Errorf("failed to rename artifact file: %w", err)
	}

	return nil
}

// Get retrieves an artifact from disk.
func (s *Store) Get(ctx context.Context, artifactID string) (*models.CommitmentArtifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !safeIDRegex.MatchString(artifactID) {
		return nil, fmt.Errorf("invalid artifact ID format: %s", artifactID)
	}

	targetPath := filepath.Join(s.basePath, artifactID+".json")

	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, artifact.ErrArtifactNotFound
		}
		return nil, fmt.Errorf("failed to read artifact file: %w", err)
	}

	var art models.CommitmentArtifact
	if err := json.Unmarshal(data, &art); err != nil {
		return nil, fmt.Errorf("failed to deserialize artifact: %w", err)
	}

	return &art, nil
}
