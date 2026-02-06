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
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

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
	artA := models.NewCommitmentArtifact("inst-1", "genesis", "APPROVED", "v1", "ctx-A", "user")
	artA.ArtifactID = "123" // Force ID for collision testing
	artA.ArtifactHash = "hash-A"

	if err := store.Write(ctx, artA); err != nil {
		t.Fatalf("Failed first write: %v", err)
	}

	// 2. Create a *different* Artifact B with the same ID ("123")
	artB := models.NewCommitmentArtifact("inst-1", "genesis", "REJECTED", "v1", "ctx-B", "attacker")
	artB.ArtifactID = "123" // Same ID
	artB.ArtifactHash = "hash-B"

	// 3. Attempt to Write Artifact B
	err := store.Write(ctx, artB)

	// 4. ASSERT: Error is ErrArtifactAlreadyExists
	if err != artifact.ErrArtifactAlreadyExists {
		t.Errorf("Expected ErrArtifactAlreadyExists, got: %v", err)
	}

	// 5. ASSERT: Read "123" from disk -> It must still match Artifact A (First write wins)
	readArt, err := store.Get(ctx, "123")
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
	defer os.RemoveAll(tmpDir)

	// Phase 1: Operational
	store1, _ := NewStore(tmpDir)
	art := models.NewCommitmentArtifact("inst-db-wipe", "prev", "APPROVED", "v1", "ctx", "actor")
	art.CalculateHashAndSetID()

	if err := store1.Write(context.Background(), art); err != nil {
		t.Fatal(err)
	}
	targetID := art.ArtifactID

	// 2. Simulate DB Loss / Crash
	store1 = nil // discard memory reference
	// In a real system, we'd delete the SQL DB here. Since we only use FS, we just reconnect.

	// 3. Create a *fresh* Store instance pointing to the same disk path
	store2, err := NewStore(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// 4. Call Get(id)
	retrievedArt, err := store2.Get(context.Background(), targetID)
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

	tmpDir, err := os.MkdirTemp("", "gantral-readonly-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory that is read-only
	readOnlyDir := filepath.Join(tmpDir, "locked")
	os.Mkdir(readOnlyDir, 0555) // Read + Execute only, No Write

	store, _ := NewStore(readOnlyDir)
	art := models.NewCommitmentArtifact("inst-fail", "0000", "APPROVED", "v1", "ctx", "act")
	art.CalculateHashAndSetID()

	// 2. Attempt write
	// Expected to fail because we can't create the underlying file
	err = store.Write(context.Background(), art)
	if err == nil {
		// If it succeeded, check if the file actually exists (maybe running as root?)
		// Ideally we skip this test if we are root/can write.
		t.Log("Write succeeded (unexpected permissions?), verifying content consistency at least.")
	} else {
		t.Logf("Got expected write error: %v", err)
	}

	// 3. ASSERT: No partial/empty file exists at the target path "locked/ID.json"
	targetPath := filepath.Join(readOnlyDir, art.ArtifactID+".json")
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Error("Atomicity Failure: Partial or empty file exists after failed write")
	}
}
