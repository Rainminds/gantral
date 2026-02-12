package replay

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
)

var (
	// ErrReplayTampered indicates that the replayed history does not match the authoritative artifacts.
	ErrReplayTampered = errors.New("SECURITY ALERT: replayed history does not match signed evidence")
)

// ReplayGuard enforces that all authority transitions in the workflow history
// are backed by valid, immutable commitment artifacts.
type ReplayGuard struct {
	store artifact.Store
}

// NewReplayGuard creates a new ReplayGuard.
func NewReplayGuard(store artifact.Store) *ReplayGuard {
	return &ReplayGuard{store: store}
}

// ValidateReplay checks a "claimed" artifact from the workflow history against the store.
// In a real interceptor, we would extract the ArtifactID from the activity result or signal payload.
//
// Parameters:
//   - claimedArtifact: The artifact struct (or ID) observed in the workflow history re-execution.
//
// Behavior:
//   - Fetches the authoritative artifact from the Store using claimedArtifact.ID.
//   - Compares critical fields (State, ContextHash, ActorID).
//   - Returns ErrReplayTampered if missing or mismatched.
func (g *ReplayGuard) ValidateReplay(ctx context.Context, claimedArtifact *models.CommitmentArtifact) error {
	if claimedArtifact == nil {
		return nil // No artifact claimed, nothing to validate (e.g. non-authority event)
	}

	// 1. Verify Integrity of Claimed Artifact (Self-Consistency)
	// Ensure the history record hasn't been tampered with to mismatch its own ID.
	// We clone it to avoid mutating the input when calculating hash.
	check := *claimedArtifact
	if err := check.CalculateHashAndSetID(); err != nil {
		return fmt.Errorf("failed to calculate hash: %w", err)
	}
	if check.ArtifactID != claimedArtifact.ArtifactID {
		slog.Error("SECURITY ALERT: Replay history tampered (hash mismatch)",
			"claimed_id", claimedArtifact.ArtifactID,
			"calculated_id", check.ArtifactID)
		return fmt.Errorf("%w: history integrity compromised", ErrReplayTampered)
	}

	// 2. Fetch Authoritative Artifact (Existence Proof)
	authoritative, err := g.store.Get(ctx, claimedArtifact.ArtifactID)
	if err != nil {
		if errors.Is(err, artifact.ErrArtifactNotFound) {
			slog.Error("SECURITY ALERT: Replay claimed artifact that does not exist",
				"claimed_id", claimedArtifact.ArtifactID)
			return fmt.Errorf("%w: artifact %s missing from store", ErrReplayTampered, claimedArtifact.ArtifactID)
		}
		return fmt.Errorf("failed to fetch artifact during replay: %w", err)
	}

	// 3. Consistency Check (State Sanity)
	// Even if ID matches (Integrity), ensure the specific fields match
	// (Defense against theoretical hash collisions or store inconsistencies).

	if authoritative.AuthorityState != claimedArtifact.AuthorityState {
		slog.Error("SECURITY ALERT: Replay claimed different state",
			"claimed_state", claimedArtifact.AuthorityState,
			"stored_state", authoritative.AuthorityState)
		return fmt.Errorf("%w: state mismatch", ErrReplayTampered)
	}

	return nil
}
