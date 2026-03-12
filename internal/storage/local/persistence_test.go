package local

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
)

// Storage Persistence Tests

func setupTestStore(t *testing.T) (*Store, string) {
	tmpDir, err := os.MkdirTemp("", "gantral-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	store, err := NewStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	return store, tmpDir
}

func Test_Adversarial_Overwrite(t *testing.T) {
	store, _ := setupTestStore(t)
	ctx := context.Background()

	// 1. Write Artifact A (ID: "123")
	artA := models.NewCommitmentArtifact("inst-1", "wf-1", models.GenesisHash, "APPROVED", "v1", "ctx-A", "user", "j")
	artA.ArtifactHash = "aeebf84acd3705fc5508a55338aae686d17ef026d805c9e6d428835a37387731" // Valid hash format
	
	if err := store.WriteArtifact(ctx, "team-1", artA); err != nil {
		t.Fatalf("Failed first write: %v", err)
	}

	// 2. Create a *different* Artifact B with the same hash
	artB := models.NewCommitmentArtifact("inst-1", "wf-1", models.GenesisHash, "REJECTED", "v1", "ctx-B", "attacker", "j")
	artB.ArtifactHash = artA.ArtifactHash

	// 3. Attempt to Write Artifact B
	err := store.WriteArtifact(ctx, "team-1", artB)

	// 4. ASSERT: Error is ErrImmutableViolation
	if err != artifact.ErrImmutableViolation {
		t.Errorf("Expected ErrImmutableViolation, got: %v", err)
	}

	// 5. ASSERT: Read from disk -> It must still match Artifact A (First write wins)
	readArt, err := store.GetArtifact(ctx, "team-1", artA.ArtifactHash)
	if err != nil {
		t.Fatalf("Failed to read artifact: %v", err)
	}
	if readArt.AuthorityState != "APPROVED" { // A was APPROVED, B was REJECTED
		t.Errorf("Immutability Violation: Artifact content changed! Got state: %s", readArt.AuthorityState)
	}
}

func Test_Survival_After_DB_Wipe(t *testing.T) {
	// 1. Set up storage location (independent of "DB")
	tmpDir, err := os.MkdirTemp("", "gantral-persistence-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Phase 1: Operational
	store1, _ := NewStore(tmpDir)
	art := models.NewCommitmentArtifact("inst-db-wipe", "w", "prev", "APPROVED", "v1", "ctx", "actor", "j")
	if err := art.CalculateHash(); err != nil {
		t.Fatal(err)
	}

	if err := store1.WriteArtifact(context.Background(), "team-1", art); err != nil {
		t.Fatal(err)
	}

	// 2. Simulate DB Loss / Crash
	store1 = nil // discard memory reference
	// In a real system, we'd delete the SQL DB here. Since we only use FS, we just reconnect.

	// 3. Create a *fresh* Store instance pointing to the same disk path
	store2, err := NewStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// 4. Call Get(hash)
	retrievedArt, err := store2.GetArtifact(context.Background(), "team-1", art.ArtifactHash)
	if err != nil {
		t.Fatalf("Failed to retrieve artifact after 'wipe': %v", err)
	}

	// 5. ASSERT: Artifact is retrieved and Hash verifies
	if retrievedArt.ArtifactHash != art.ArtifactHash {
		t.Error("Hash mismatch on retrieval")
	}
}

func Test_Atomic_Failure(t *testing.T) {
	// 1. Use a read-only directory to simulate write failure
	// Note: In some containers/OS, chmod might not strictly prevent root writes,
	// but standard user permissions should fail.
	// NOTE: GitLab CI runners might run as root, which ignores these permissions.
	if os.Getuid() == 0 {
		t.Skip("Skipping atomicity failure test: Running as root ignores directory permissions")
	}

	tmpDir, err := os.MkdirTemp("", "gantral-readonly-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a subdirectory that is read-only
	readOnlyDir := filepath.Join(tmpDir, "locked")
	if err := os.Mkdir(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	} // Read + Execute only, No Write

	store, _ := NewStore(readOnlyDir)
	art := models.NewCommitmentArtifact("inst-fail", "w", "0000", "APPROVED", "v1", "ctx", "act", "j")
	if err := art.CalculateHash(); err != nil {
		t.Fatal(err)
	}

	// 2. Attempt write
	// Expected to fail because we can't create the underlying team directory or file
	err = store.WriteArtifact(context.Background(), "team-1", art)
	if err == nil {
		// If it succeeded, check if the file actually exists (maybe running as root?)
		targetPath := filepath.Join(readOnlyDir, "team-1", art.ArtifactHash+".json")
		if _, statErr := os.Stat(targetPath); statErr == nil {
			t.Skip("Skipping atomicity failure test: Directory permissions were bypassed (running as root or special capabilities)")
		}
		t.Errorf("Expected write to fail in read-only directory, but it succeeded")
	} else {
		t.Logf("Got expected write error: %v", err)
	}

	// 3. ASSERT: No partial/empty file exists at the target path "locked/team-1/HASH.json"
	targetPath := filepath.Join(readOnlyDir, "team-1", art.ArtifactHash+".json")
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("Atomicity Failure: Partial or empty file exists after failed write")
	}
}
