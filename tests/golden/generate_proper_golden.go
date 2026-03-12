//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Rainminds/gantral/pkg/models"
)

func main() {
	art := models.CommitmentArtifact{
		ArtifactVersion:     "v1",
		ArtifactID:          "00000000-0000-0000-0000-000000000000",
		InstanceID:          "inst-golden-001",
		WorkflowVersionID:   "golden-wf-version",
		PrevArtifactHash:    nil,
		AuthorityState:      "APPROVED",
		PolicyVersionID:     "pol-golden-v1",
		ContextSnapshotHash: "ctx-hash-golden-123",
		HumanActorID:        "actor-golden-456",
		Justification:       "golden justification",
		CryptoProfile:       "standard-v1",
		TimestampToken:      "2023-01-01T00:00:00Z",
		TimestampAlgorithm:  "dummy-algo",
		ArtifactSignature:   "mock-sig",
		SignatureAlgorithm:  "mock-sig-algo",
		Attestations:        nil,
	}

	err := art.CalculateHash()
	if err != nil {
		panic(err)
	}

	// 1. Write the canonical_artifact_v1.json
	fullBytes, _ := json.MarshalIndent(art, "", "  ")
	err = os.WriteFile("canonical_artifact_v1.json", fullBytes, 0644)
	if err != nil {
		panic(err)
	}

	// 2. Write expected_sha256_v1.txt
	err = os.WriteFile("expected_sha256_v1.txt", []byte(art.ArtifactHash), 0644)
	if err != nil {
		panic(err)
	}

	// 3. Write replay format
	replayResult := `{"valid":true,"artifact_id":"` + art.ArtifactID + `","calculated_hash":"` + art.ArtifactHash + `"}`
	os.WriteFile("expected_replay_result_v1.txt", []byte(replayResult), 0644)

	fmt.Printf("Generated Golden Vectors.\nArtifact ID: %s\nHash: %s\n", art.ArtifactID, art.ArtifactHash)
}
