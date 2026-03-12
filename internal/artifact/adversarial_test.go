package artifact

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
)

// Adversarial Test Suite

func Test_TamperResistance(t *testing.T) {
	// 1. Create a valid artifact
	m := NewManager(&MockStore{}, &MockSigner{}, &MockTimeStamper{})
	validArt, err := m.EmitArtifact(context.Background(), "team-1", "inst-1", "wf-1", "prev-0", "APPROVED", "v1", "ctx-hash", "admin", "j")
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 2. Serialize to JSON
	originalJSON, _ := json.Marshal(validArt)

	// 3. Manually flip one bit/field (change role/actor) in the JSON payload
	// We'll replace "admin" with "hacker"
	tamperedJSONStr := strings.Replace(string(originalJSON), "admin", "hacker", 1)

	// 4. Attempt to verify logic (Simulate Verification)
	// Since we don't have a standalone Verify function yet strictly defined in the interface,
	// we simulate verification: Reconstruct object -> Recalculate Hash -> Compare with ID
	var tamperedArt models.CommitmentArtifact
	if err := json.Unmarshal([]byte(tamperedJSONStr), &tamperedArt); err != nil {
		t.Fatalf("Failed to unmarshal tampered JSON: %v", err)
	}

	// The ID in the JSON is still the ORIGINAL ID (validArt.ArtifactID)
	// But the content is now "hacker".
	// Recalculate hash of the content.
	recalcArt := tamperedArt // copy
	if err := recalcArt.CalculateHash(); err != nil {
		t.Fatalf("Recalculation failed: %v", err)
	}

	// 5. ASSERT that verification fails (Hash Mismatch)
	if recalcArt.ArtifactHash == validArt.ArtifactHash {
		t.Fatal("Tamper Resistance FAILED: Hash collision or logic ignored change.")
	}
	// Content change MUST change ID in v1
	if recalcArt.ArtifactHash == validArt.ArtifactHash {
		t.Fatal("Tamper Resistance FAILED: Content change did not change ID.")
	}
}

func Test_ChainBreakage(t *testing.T) {
	m := NewManager(&MockStore{}, &MockSigner{}, &MockTimeStamper{})

	// 1. Create Artifact A
	artA, _ := m.EmitArtifact(context.Background(), "team-1", "inst-1", "wf1", models.GenesisHash, "RUNNING", "v1", "ctx-A", "sys", "j")

	// 2. Create Artifact B pointing to Artifact A
	artB, _ := m.EmitArtifact(context.Background(), "team-1", "inst-1", "wf1", artA.ArtifactHash, "APPROVED", "v1", "ctx-B", "human", "j")

	// 3. "Corrupt" Artifact A (Simulate that the history log was altered)
	// Actually, we check that B *requires* A's ID.
	if artB.PrevArtifactHash != artA.ArtifactHash {
		t.Fatal("Chain Linkage FAILED: B not pointing to A")
	}

	// Let's simulate a verification check:
	// "Does ArtB point to the calculated hash of ArtA?"

	// Re-calculate A's hash from its content
	if err := artA.CalculateHash(); err != nil {
		t.Fatal(err)
	}

	// If someone changed ArtA's content:
	artA.HumanActorID = "malicious"
	if err := artA.CalculateHash(); err != nil {
		t.Fatal(err)
	} // Now A has a new ID

	// ASSERT that verification of the chain fails
	if artB.PrevArtifactHash == artA.ArtifactHash {
		t.Fatal("Chain Breakage FAILED: B still points to modified A")
	}
}

func Test_AmbiguityRefusal(t *testing.T) {
	m := NewManager(&MockStore{}, &MockSigner{}, &MockTimeStamper{})

	// 1. Attempt to emit with missing critical field (Empty Context Hash)
	// The Manager.EmitArtifact is expected to Fail-Closed (return error, no artifact).
	art, err := m.EmitArtifact(context.Background(), "team-1", "inst-1", "wf", "prev", "APPROVED", "v1", "", "actor", "j")

	// 2. ASSERT fatal error
	if err == nil {
		t.Fatal("Ambiguity Refusal FAILED: Emitted artifact despite missing context hash")
	}
	if art != nil {
		t.Fatal("Ambiguity Refusal FAILED: Returned partial artifact on error")
	}
}

func Test_ReplayDeterminism(t *testing.T) {
	// 1. Hardcode a known valid artifact JSON (generated "today")
	// This ensures that "Time" is not a hidden variable in *verification* (hashing).

	// Generated from a previous run (simulated here for the test to be self-contained)
	// Fields: inst=i1, prev=p1, state=s1, pol=pv1, ctx=c1, actor=a1, ts=2023-01-01T00:00:00Z
	// We need to know what the hash SHOULD be.

	art := models.NewCommitmentArtifact("i1", "wf", "p1", "s1", "pv1", "c1", "a1", "j")
	if err := art.CalculateHash(); err != nil {
		t.Fatal(err)
	}
	expectedHash := art.ArtifactHash

	// 2. Re-run verification logic
	recalcArt := models.NewCommitmentArtifact("i1", "wf", "p1", "s1", "pv1", "c1", "a1", "j")
	recalcArt.ArtifactID = art.ArtifactID

	if err := recalcArt.CalculateHash(); err != nil {
		t.Fatal(err)
	}

	// 3. ASSERT valid
	if recalcArt.ArtifactHash != expectedHash {
		t.Errorf("Replay Determinism FAILED: \nExpected: %s\nGot:      %s", expectedHash, recalcArt.ArtifactHash)
	}
}
