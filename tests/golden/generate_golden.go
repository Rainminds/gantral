//go:build ignore

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Minimal struct to match pkg/models/artifact.go WITHOUT depending on it
// to avoid import cycles or breakage if core changes.
// We want to FREEZE the Golden Vector format.
type GoldenArtifact struct {
	ArtifactVersion  string `json:"artifact_version"`
	ArtifactID       string `json:"artifact_id"`
	InstanceID       string `json:"instance_id"`
	PrevArtifactHash string `json:"prev_artifact_hash"`
	AuthorityState   string `json:"authority_state"`
	PolicyVersionID  string `json:"policy_version_id"`
	ContextHash      string `json:"context_hash"`
	HumanActorID     string `json:"human_actor_id"`
	Timestamp        string `json:"timestamp"`
}

func main() {
	// 1. Define Fixed Golden Data
	artifact := GoldenArtifact{
		ArtifactVersion:  "v1",
		InstanceID:       "inst-golden-001",
		PrevArtifactHash: "0000000000000000000000000000000000000000000000000000000000000000",
		AuthorityState:   "APPROVED",
		PolicyVersionID:  "pol-golden-v1",
		ContextHash:      "ctx-hash-golden-123",
		HumanActorID:     "actor-golden-456",
		Timestamp:        "2023-01-01T00:00:00Z",
		// ArtifactID is calculated below
	}

	// 2. Calculate Canonical Payload (Map Sort)
	// Matches pkg/models logic: map[string]string -> json.Marshal
	canonicalMap := map[string]string{
		"artifact_version":   artifact.ArtifactVersion,
		"instance_id":        artifact.InstanceID,
		"prev_artifact_hash": artifact.PrevArtifactHash,
		"authority_state":    artifact.AuthorityState,
		"policy_version_id":  artifact.PolicyVersionID,
		"context_hash":       artifact.ContextHash,
		"human_actor_id":     artifact.HumanActorID,
		"timestamp":          artifact.Timestamp,
	}

	canonicalBytes, err := json.Marshal(canonicalMap)
	if err != nil {
		panic(err)
	}

	// 3. Calculate Hash
	hash := sha256.Sum256(canonicalBytes)
	artifactID := hex.EncodeToString(hash[:])
	artifact.ArtifactID = artifactID

	// 4. Write Files
	outputDir := "tests/golden"

	// 4a. Canonical/Full JSON
	// We export the FULL object including ID, matching what VerifyArtifact expects as input
	fullMap := map[string]string{
		"artifact_version":   artifact.ArtifactVersion,
		"artifact_id":        artifact.ArtifactID,
		"instance_id":        artifact.InstanceID,
		"prev_artifact_hash": artifact.PrevArtifactHash,
		"authority_state":    artifact.AuthorityState,
		"policy_version_id":  artifact.PolicyVersionID,
		"context_hash":       artifact.ContextHash,
		"human_actor_id":     artifact.HumanActorID,
		"timestamp":          artifact.Timestamp,
	}

	fullJSON, _ := json.MarshalIndent(fullMap, "", "  ") // Indent for readability in repo

	err = os.WriteFile(filepath.Join(outputDir, "canonical_artifact_v1.json"), fullJSON, 0644)
	if err != nil {
		panic(err)
	}

	// 4b. Expected Hash
	err = os.WriteFile(filepath.Join(outputDir, "expected_sha256_v1.txt"), []byte(artifactID), 0644)
	if err != nil {
		panic(err)
	}

	// 4c. Expected Replay Result
	replayResult := `{"valid":true,"artifact_id":"` + artifactID + `","calculated_hash":"` + artifactID + `"}`
	err = os.WriteFile(filepath.Join(outputDir, "expected_replay_result_v1.txt"), []byte(replayResult), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated Golden Vectors in %s\n", outputDir)
	fmt.Printf("Artifact ID: %s\n", artifactID)
}
