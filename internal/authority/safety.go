package authority

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/internal/artifact"
)

var (
	// ErrStateAmbiguous indicates a critical disconnect between execution state and evidence.
	// This usually implies a phantom write or a replays attack.
	ErrStateAmbiguous = errors.New("CRITICAL: state ambiguity detected (phantom artifact)")
)

// ConsistencyGuard safeguards transitions by verifying evidence existence.
type ConsistencyGuard struct {
	store artifact.Store
}

// NewConsistencyGuard creates a guard with access to the immutable evidence store.
func NewConsistencyGuard(store artifact.Store) *ConsistencyGuard {
	return &ConsistencyGuard{store: store}
}

// EnsureStateConsistency verifies that a claimed artifact actually exists in the WORM store
// before allowing any side-effect or state transition dependent on it.
func (g *ConsistencyGuard) EnsureStateConsistency(ctx context.Context, instanceID string, artifactID string) error {
	if artifactID == "" {
		// If no artifact is claimed (e.g. Genesis), strictly validate if that is allowed.
		// For now, if caller claims "no artifact", we assume they know what they are doing (e.g. start)
		// But if they claim an ID, it MUST exist.
		return nil
	}

	// 1. Check Store
	art, err := g.store.Get(ctx, artifactID)

	// 2. Strict Error Handling
	if err != nil {
		if errors.Is(err, artifact.ErrArtifactNotFound) {
			slog.Error("SECURITY ALERT: Phantom Artifact Detected",
				"instance_id", instanceID,
				"claimed_artifact_id", artifactID)
			return fmt.Errorf("%w: artifact %s not found", ErrStateAmbiguous, artifactID)
		}
		// Other errors (e.g. connection) -> Fail Closed
		return fmt.Errorf("consistency check failed: %w", err)
	}

	// 3. Verify Binding (Defense against ID reuse or collision)
	if art.InstanceID != instanceID {
		slog.Error("SECURITY ALERT: Cross-Instance Contamination",
			"expected_instance", instanceID,
			"artifact_instance", art.InstanceID)
		return fmt.Errorf("%w: instance mismatch", ErrStateAmbiguous)
	}

	return nil
}
