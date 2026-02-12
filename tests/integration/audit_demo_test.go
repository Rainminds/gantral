//go:build integration

package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuditDemo(t *testing.T) {
	if os.Getenv("CI") != "" && os.Getenv("SKIP_DOCKER_TESTS") != "" {
		t.Skip("Skipping docker-based demo in CI environment without docker socket")
	}

	// Logic to find path
	projectRoot, err := filepath.Abs("../../") // Moving up from tests/e2e
	assert.NoError(t, err)

	scriptRelPath := "examples/audit-demo/run-demo.sh"

	cmd := exec.Command("/bin/bash", scriptRelPath)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Force simulation mode for CI/Test env without Docker
	cmd.Env = append(os.Environ(), "DEMO_SIMULATION=true")

	t.Logf("Running Audit Demo Script from %s...", projectRoot)
	err = cmd.Run()

	assert.NoError(t, err, "Demo script failed")
}
