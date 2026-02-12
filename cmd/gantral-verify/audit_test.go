package main_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Rainminds/gantral/pkg/models"
)

// Hostile Auditor Tests
// These tests execute the compiled binary or main entry point to verify behavior
// in a "Clean Room" environment (no DB access, pure file logic).

func buildVerifier(t *testing.T) string {
	binPath := filepath.Join(os.TempDir(), "gantral-verify")
	// Build the CLI tool
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "." // Current directory (cmd/gantral-verify)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build verifier: %v\nOutput: %s", err, output)
	}
	return binPath
}

func Test_Offline_Verification(t *testing.T) {
	binPath := buildVerifier(t)
	defer os.Remove(binPath)

	// 1. Create a valid artifact JSON
	art := models.NewCommitmentArtifact("inst-valid", "prev", "APPROVED", "v1", "ctx", "auditor")
	if err := art.CalculateHashAndSetID(); err != nil {
		t.Fatal(err)
	}

	data, _ := json.Marshal(art)
	tmpFile, err := os.CreateTemp("", "audit-valid-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.Write(data)
	tmpFile.Close()

	// 2. Run the `gantral-verify` binary against it
	cmd := exec.Command(binPath, "file", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	outStr := string(output)

	// 3. ASSERT: Output contains "VALID" and Exit Code is 0
	if err != nil {
		t.Fatalf("Verification failed unexpectedly: %v\nOutput: %s", err, outStr)
	}
	if !strings.Contains(outStr, "VALID") {
		t.Errorf("Expected output to contain 'VALID', got:\n%s", outStr)
	}
	if strings.Contains(outStr, "INVALID") {
		t.Error("Output contained 'INVALID' for a valid artifact")
	}
}

func Test_Tampered_File(t *testing.T) {
	binPath := buildVerifier(t)
	defer os.Remove(binPath)

	// 1. Create a valid artifact
	art := models.NewCommitmentArtifact("inst-tampered", "prev", "APPROVED", "v1", "ctx", "attacker")
	if err := art.CalculateHashAndSetID(); err != nil {
		t.Fatal(err)
	}

	validJSON, _ := json.Marshal(art)

	// 2. Modify the file on disk (change 1 byte in the content, keep ID same)
	// Replace "attacker" with "hacker___" (same length to avoid JSON error, or just change content)
	tamperedJSON := strings.Replace(string(validJSON), "attacker", "hacker__", 1)

	tmpFile, err := os.CreateTemp("", "audit-tampered-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.WriteString(tamperedJSON)
	tmpFile.Close()

	// 3. Run `gantral-verify`
	cmd := exec.Command(binPath, "file", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	outStr := string(output)

	// 4. ASSERT: Output contains "INVALID" / "Hash Mismatch" and Exit Code is 1
	if err == nil {
		t.Fatal("Verifier should have failed (exit code 1), but succeeded (0)")
	}
	if !strings.Contains(outStr, "INVALID") {
		t.Errorf("Expected 'INVALID' in output, got:\n%s", outStr)
	}
}

func Test_Chain_Verification(t *testing.T) {
	binPath := buildVerifier(t)
	defer os.Remove(binPath)

	tmpDir, _ := os.MkdirTemp("", "audit-chain-*")
	defer os.RemoveAll(tmpDir)

	// Create Chain: A -> B -> C
	// Note: We need timestamps to allow sorting if the tool relies on it.
	// Models sets timestamp to Now(). We sleep slightly to ensure order.

	// A
	artA := models.NewCommitmentArtifact("inst", models.GenesisHash, "RUNNING", "v1", "ctxA", "sys")
	_ = artA.CalculateHashAndSetID()
	_ = os.WriteFile(filepath.Join(tmpDir, "1_A.json"), mustMarshal(artA), 0644)
	time.Sleep(10 * time.Millisecond)

	// B
	artB := models.NewCommitmentArtifact("inst", artA.ArtifactID, "APPROVED", "v1", "ctxB", "sys")
	_ = artB.CalculateHashAndSetID()
	_ = os.WriteFile(filepath.Join(tmpDir, "2_B.json"), mustMarshal(artB), 0644)
	time.Sleep(10 * time.Millisecond)

	// C
	artC := models.NewCommitmentArtifact("inst", artB.ArtifactID, "COMPLETED", "v1", "ctxC", "sys")
	_ = artC.CalculateHashAndSetID()
	_ = os.WriteFile(filepath.Join(tmpDir, "3_C.json"), mustMarshal(artC), 0644)

	// Run Verify Chain
	cmd := exec.Command(binPath, "chain", tmpDir)
	output, err := cmd.CombinedOutput()
	outStr := string(output)

	if err != nil {
		t.Fatalf("Chain verification failed: %v\nOutput: %s", err, outStr)
	}
	if !strings.Contains(outStr, "CHAIN VALID") {
		t.Errorf("Expected 'CHAIN VALID', got:\n%s", outStr)
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
