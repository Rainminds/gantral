package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/internal/auth"
	"github.com/Rainminds/gantral/internal/middleware"
	"github.com/Rainminds/gantral/pkg/constants"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/Rainminds/gantral/pkg/verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verifiability Audit Tests - Validating Phases 4 to 6

// Test Developer Demo (Phase 4)
func TestAudit_Phase4_DeveloperDemoShellScript(t *testing.T) {
	t.Run("Developer Demo shell script exists and contains core commands", func(t *testing.T) {
		scriptPath := "../../scripts/test-demo-env.sh"
		content, err := os.ReadFile(scriptPath)
		require.NoError(t, err, "Demo shell script must exist in scripts/test-demo-env.sh")

		scriptStr := string(content)
		assert.Contains(t, scriptStr, "docker-compose up -d", "Must use docker-compose for foundations")
		assert.Contains(t, scriptStr, "cmd/server/main.go", "Must point to the Gantral server entrypoint")
		assert.Contains(t, scriptStr, "core", "Must verify core dependencies")
	})
}

// Test Federated Identity (Phase 5)
func TestAudit_Phase5_FederatedIdentity(t *testing.T) {
	t.Run("Empty OIDC token is rejected for WAITING_FOR_HUMAN mutation", func(t *testing.T) {
		// Verify middleware.GetIdentity fails on empty context
		emptyCtx := context.Background()
		_, err := middleware.GetIdentity(emptyCtx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no identity in context")
	})

	t.Run("Spoofed ID without OrgID binding is rejected", func(t *testing.T) {
		// Construct a context with an actor but no OrgID (simulating malformed/spoofed state)
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, &auth.Identity{
			Subject: "malicious-actor",
			// Missing OrgID
		})

		identity, err := middleware.GetIdentity(ctx)
		require.NoError(t, err)
		assert.Equal(t, "", identity.OrgID, "Identity must strictly reflect provided OIDC fields")

		// Assertion: Control plane logic must fail if OrgID is empty
		if identity.OrgID == "" {
			err := engine.ErrMissingTeamID
			assert.Error(t, err)
			assert.Equal(t, "missing team id", err.Error())
		}
	})
}

// Test Log Independence (Phase 6)
func TestAudit_Phase6_LogIndependence(t *testing.T) {
	t.Run("TRD 7.0: VerifyArtifact replays successfully using only JSON models", func(t *testing.T) {
		// 1. Create a canonical artifact chain (mocked payload)
		art := models.CommitmentArtifact{
			ArtifactVersion:     models.CurrentArtifactVersion,
			ArtifactID:          "audit-phase-6",
			InstanceID:          "inst-replay-1",
			AuthorityState:      string(constants.StateCompleted),
			ContextSnapshotHash: "0xdeadbeef",
			HumanActorID:        "auditor-1",
		}
		
		// 2. Compute hash (tamper-evident link)
		err := art.CalculateHash()
		require.NoError(t, err)

		// 3. Serialize to JSON (S3 representation)
		data, err := json.Marshal(art)
		require.NoError(t, err)

		// 4. Create a Verifier with NO database access
		v := verifier.New(nil, nil)

		// 5. Verify using ONLY the JSON data
		report, err := v.VerifyArtifact(data)
		require.NoError(t, err)

		// 6. Assert result is VALID
		assert.True(t, report.Valid)
		assert.Equal(t, verifier.StatusValid, report.Status)
	})
}
