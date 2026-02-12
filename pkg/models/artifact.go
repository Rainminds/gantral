package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SchemaVersionV1 defines the current schema version for commitment artifacts.
const SchemaVersionV1 = "v1"

// GenesisHash is the SHA-256 hash of comparable length (64 chars) consisting of zeros.
// It is used as the PrevArtifactHash for the first artifact in a chain.
const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

// CommitmentArtifact represents the immutable proof of an execution authority transition.
// It serves as the root of trust for auditability and verification.
//
// "EmitArtifact generates a non-repudiable proof of authorization bound to execution state."
type CommitmentArtifact struct {
	// ArtifactVersion is the schema version (e.g., "v1").
	ArtifactVersion string `json:"artifact_version"`

	// ArtifactID is the unique identifier of this artifact (usually the hash of its content).
	ArtifactID string `json:"artifact_id"`

	// InstanceID is the UUID of the execution instance.
	InstanceID string `json:"instance_id"`

	// PrevArtifactHash is the hash of the previous artifact in the chain.
	// This links artifacts together into an immutable log.
	PrevArtifactHash string `json:"prev_artifact_hash"`

	// AuthorityState is the state being transitioned to (e.g., APPROVED, REJECTED).
	AuthorityState string `json:"authority_state"`

	// PolicyVersionID is the version of the policy used for evaluation.
	PolicyVersionID string `json:"policy_version_id"`

	// ContextHash is the SHA256 hash of the execution context snapshot.
	ContextHash string `json:"context_hash"`

	// HumanActorID is the identity of the human/system authorizing the transition.
	HumanActorID string `json:"human_actor_id"`

	// Timestamp is the exact time of emission (RFC3339).
	Timestamp string `json:"timestamp"`
}

// NewCommitmentArtifact creates a new artifact with the given fields.
// It automatically sets the SchemaVersion and Timestamp.
// It does NOT calculate the ID/Hash; that must be done via CalculateHashAndSetID.
func NewCommitmentArtifact(
	instanceID string,
	prevArtifactHash string,
	authorityState string,
	policyVersionID string,
	contextHash string,
	humanActorID string,
) *CommitmentArtifact {
	return &CommitmentArtifact{
		ArtifactVersion:  SchemaVersionV1,
		InstanceID:       instanceID,
		PrevArtifactHash: prevArtifactHash,
		AuthorityState:   authorityState,
		PolicyVersionID:  policyVersionID,
		ContextHash:      contextHash,
		HumanActorID:     humanActorID,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
	}
}

// CalculateHashAndSetID computes the SHA256 hash of the canonical payload
// and sets both ArtifactHash and ArtifactID.
// It returns an error if serialization fails.
func (a *CommitmentArtifact) CalculateHashAndSetID() error {
	payloadBytes, err := a.CanonicalPayload()
	if err != nil {
		return fmt.Errorf("failed to generate canonical payload: %w", err)
	}

	hash := sha256.Sum256(payloadBytes)
	hashString := hex.EncodeToString(hash[:])

	a.ArtifactID = hashString
	return nil
}

// CanonicalPayload returns the strictly deterministic JSON bytes used for hashing.
// It uses a map to ensure keys are sorted alphabetically by encoding/json,
// guaranteeing a stable hash regardless of struct field order.
func (a *CommitmentArtifact) CanonicalPayload() ([]byte, error) {
	// 1. Validate required fields (Fail-Closed)
	if a.InstanceID == "" {
		return nil, errors.New("canonical payload: missing instance_id")
	}
	if a.AuthorityState == "" {
		return nil, errors.New("canonical payload: missing authority_state")
	}
	if a.ContextHash == "" {
		return nil, errors.New("canonical payload: missing context_hash")
	}
	if a.Timestamp == "" {
		return nil, errors.New("canonical payload: missing timestamp")
	}
	// Note: PrevArtifactHash can be empty for the first artifact (Genesis),
	// but strictly speaking, explicit nil handling should be done by caller.
	// We allow empty string or GenesisHash. Ideally, use models.GenesisHash.
	// For now, we allow it but ensure it's included in the map.

	// 2. Construct map for sorted keys
	msg := map[string]string{
		"artifact_version":   a.ArtifactVersion,
		"instance_id":        a.InstanceID,
		"prev_artifact_hash": a.PrevArtifactHash,
		"authority_state":    a.AuthorityState,
		"policy_version_id":  a.PolicyVersionID,
		"context_hash":       a.ContextHash,
		"human_actor_id":     a.HumanActorID,
		"timestamp":          a.Timestamp,
	}

	// 3. Marshal with standard library (sorts map keys)
	return json.Marshal(msg)
}

// MarshalJSON implements the json.Marshaler interface to ensure
// the full object handles 0-values correctly and adheres to schema.
// We force the standard struct marshaling but could enforce order if needed.
// Since typical consumers parse JSON into objects, key order in the OUTER object
// is less critical than the HASH payload order. But to be safe, we can reuse
// the map approach for the outer object too, adding the ID fields.
func (a *CommitmentArtifact) MarshalJSON() ([]byte, error) {
	// Construct map to ensure complete representation including ID/Hash
	msg := map[string]string{
		"artifact_version":   a.ArtifactVersion,
		"artifact_id":        a.ArtifactID,
		"instance_id":        a.InstanceID,
		"prev_artifact_hash": a.PrevArtifactHash,
		"authority_state":    a.AuthorityState,
		"policy_version_id":  a.PolicyVersionID,
		"context_hash":       a.ContextHash,
		"human_actor_id":     a.HumanActorID,
		"timestamp":          a.Timestamp,
	}
	return json.Marshal(msg)
}
