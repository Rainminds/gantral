package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	// ArtifactVersionV1 represents the initial v1 schema.
	ArtifactVersionV1 = "v1"
	// ArtifactVersionV1_0 represents the v1.0 schema version for backward compatibility.
	ArtifactVersionV1_0 = "v1.0"
	// CurrentArtifactVersion defines the current schema version for new commitment artifacts.
	CurrentArtifactVersion = ArtifactVersionV1
)

// IsSupportedVersion checks if the given artifact version is supported by the system.
func IsSupportedVersion(version string) bool {
	switch version {
	case ArtifactVersionV1, ArtifactVersionV1_0:
		return true
	default:
		return false
	}
}

// ErrUnsupportedArtifactVersion indicates that the artifact version is unknown or no longer supported.
var ErrUnsupportedArtifactVersion = errors.New("unsupported artifact version")

// GenesisHash is the SHA-256 hash of comparable length (64 chars) consisting of zeros.
const GenesisHash = "0000000000000000000000000000000000000000000000000000000000000000"

// Attestation represents a cryptographic extension to an artifact.
type Attestation struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// CommitmentArtifact represents the immutable proof of an execution authority transition.
type CommitmentArtifact struct {
	ArtifactVersion     string        `json:"artifact_version"`
	ArtifactID          string        `json:"artifact_id"`
	InstanceID          string        `json:"instance_id"`
	WorkflowVersionID   string        `json:"workflow_version_id"`
	PolicyVersionID     string        `json:"policy_version_id"`
	AuthorityState      string        `json:"authority_state"`
	ContextSnapshotHash string        `json:"context_snapshot_hash"`
	HumanActorID        string        `json:"human_actor_id"`
	Justification       string        `json:"justification"`
	PrevArtifactHash    interface{}   `json:"prev_artifact_hash"` // Explicitly handle null for genesis
	ArtifactHash        string        `json:"artifact_hash"`
	CryptoProfile       string        `json:"crypto_profile"`
	ArtifactSignature   string        `json:"artifact_signature"`
	SignatureAlgorithm  string        `json:"signature_algorithm"`
	TimestampToken      string        `json:"timestamp_token"`
	TimestampAlgorithm  string        `json:"timestamp_algorithm"`
	Attestations        []Attestation `json:"attestations"` // Explicitly handle null if empty
}

// NewCommitmentArtifact creates a new artifact with the given fields.
// It automatically generates a UUIDv4 for ArtifactID and sets Timestamp.
func NewCommitmentArtifact(
	instanceID string,
	workflowVersionID string,
	prevArtifactHash string,
	authorityState string,
	policyVersionID string,
	contextSnapshotHash string,
	humanActorID string,
	justification string,
) *CommitmentArtifact {
	var prevHash interface{} = prevArtifactHash
	if prevArtifactHash == GenesisHash || prevArtifactHash == "" {
		prevHash = nil
	}

	return &CommitmentArtifact{
		ArtifactVersion:     CurrentArtifactVersion,
		ArtifactID:          uuid.New().String(),
		InstanceID:          instanceID,
		WorkflowVersionID:   workflowVersionID,
		PolicyVersionID:     policyVersionID,
		AuthorityState:      authorityState,
		ContextSnapshotHash: contextSnapshotHash,
		HumanActorID:        humanActorID,
		Justification:       justification,
		PrevArtifactHash:    prevHash,
		CryptoProfile:       "standard-v1",
		Attestations:        nil,
	}
}

// CalculateHash computes the SHA256 hash of the canonical payload
// as defined in TRD A.3 and A.11, and sets artifacts.ArtifactHash.
// It explicitly excludes signature and timestamp fields.
func (a *CommitmentArtifact) CalculateHash() error {
	payload, err := a.PayloadForHashing()
	if err != nil {
		return err
	}

	canonicalBytes, err := CanonicalSerialize(payload)
	if err != nil {
		return fmt.Errorf("canonical serialization failed: %w", err)
	}

	var hash [32]byte
	if a.PrevArtifactHash != nil && a.PrevArtifactHash != "" {
		prevHashStr, ok := a.PrevArtifactHash.(string)
		if ok && prevHashStr != "" {
			// TRD A.11: SHA256(canonical_json_bytes || prev_artifact_hash)
			combined := append(canonicalBytes, []byte(prevHashStr)...)
			hash = sha256.Sum256(combined)
		} else {
			hash = sha256.Sum256(canonicalBytes)
		}
	} else {
		hash = sha256.Sum256(canonicalBytes)
	}

	a.ArtifactHash = hex.EncodeToString(hash[:])
	return nil
}

// PayloadForHashing returns the strictly deterministic map used for hashing,
// omitting the signature, timestamp, attestations, and artifact_hash per TRD A.3.
func (a *CommitmentArtifact) PayloadForHashing() (map[string]interface{}, error) {
	// 1. Validate required fields (Fail-Closed)
	if !IsSupportedVersion(a.ArtifactVersion) {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedArtifactVersion, a.ArtifactVersion)
	}
	if a.InstanceID == "" {
		return nil, errors.New("canonical payload: missing instance_id")
	}
	if a.AuthorityState == "" {
		return nil, errors.New("canonical payload: missing authority_state")
	}
	if a.ContextSnapshotHash == "" {
		return nil, errors.New("canonical payload: missing context_snapshot_hash")
	}

	// 2. Construct map per schema
	msg := map[string]interface{}{
		"artifact_version":      a.ArtifactVersion,
		"artifact_id":           a.ArtifactID,
		"instance_id":           a.InstanceID,
		"workflow_version_id":   a.WorkflowVersionID,
		"policy_version_id":     a.PolicyVersionID,
		"authority_state":       a.AuthorityState,
		"context_snapshot_hash": a.ContextSnapshotHash,
		"human_actor_id":        a.HumanActorID,
		"justification":         a.Justification,
		"prev_artifact_hash":    a.PrevArtifactHash,
		"crypto_profile":        a.CryptoProfile,
	}

	return msg, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (a *CommitmentArtifact) MarshalJSON() ([]byte, error) {
	msg := map[string]interface{}{
		"artifact_version":      a.ArtifactVersion,
		"artifact_id":           a.ArtifactID,
		"instance_id":           a.InstanceID,
		"workflow_version_id":   a.WorkflowVersionID,
		"policy_version_id":     a.PolicyVersionID,
		"authority_state":       a.AuthorityState,
		"context_snapshot_hash": a.ContextSnapshotHash,
		"human_actor_id":        a.HumanActorID,
		"justification":         a.Justification,
		"prev_artifact_hash":    a.PrevArtifactHash,
		"artifact_hash":         a.ArtifactHash,
		"crypto_profile":        a.CryptoProfile,
		"artifact_signature":    a.ArtifactSignature,
		"signature_algorithm":   a.SignatureAlgorithm,
		"timestamp_token":       a.TimestampToken,
		"timestamp_algorithm":   a.TimestampAlgorithm,
	}

	// Handle optional attestations
	if len(a.Attestations) == 0 {
		msg["attestations"] = nil
	} else {
		msg["attestations"] = a.Attestations
	}

	return json.Marshal(msg)
}
