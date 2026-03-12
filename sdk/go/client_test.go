package sdk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/Rainminds/gantral/sdk/go"
	"github.com/stretchr/testify/assert"
)

// mockAuthorityCheckpointJSON represents the exact JSON schema returned by the Gantral backend
// during an instance creation or decision recording.
const mockAuthorityCheckpointJSON = `{
  "id": "inst-12345",
  "workflow_id": "wf-finance-approval",
  "owning_team_id": "team-alpha",
  "state": "WAITING_FOR_HUMAN",
  "trigger_context": {
    "amount": 5000,
    "currency": "USD"
  },
  "policy": {
    "id": "pol-v1",
    "materiality": "HIGH",
    "requires_human_approval": true,
    "approval_timeout_seconds": 3600
  },
  "artifact_hash": "a1b2c3d4e5f6g7h8",
  "artifact_version": 1,
  "created_at": "2026-03-10T10:00:00Z",
  "updated_at": "2026-03-10T10:05:00Z"
}`

func TestClient_AuthorityCheckpoint_Deserialization(t *testing.T) {
	// Setup a mock HTTP server to return the hardcoded JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(mockAuthorityCheckpointJSON))
	}))
	defer ts.Close()

	// Initialize the SDK Client
	client := sdk.NewClient(ts.URL)

	// Invoke an SDK method that deserializes the response (e.g., CreateInstance)
	// We pass empty maps/structs for the request because we are strictly testing the response JSON deserialization.
	inst, err := client.CreateInstance(
		context.Background(),
		"wf-finance-approval",
		map[string]interface{}{"amount": 5000, "currency": "USD"},
		sdk.Policy{},
	)

	// Assert no error in HTTP call or deserialization
	assert.NoError(t, err)
	assert.NotNil(t, inst)

	// Prevent field loss: strictly assert that all fields mapped properly into the SDK struct.
	// This proves the SDK models are functionally parallel to the backend API schema.
	assert.Equal(t, "inst-12345", inst.ID)
	assert.Equal(t, "wf-finance-approval", inst.WorkflowID)
	assert.Equal(t, "team-alpha", inst.OwningTeamID)
	assert.Equal(t, sdk.StateWaitingForHuman, inst.State)

	// Assert Context Map translation
	assert.NotNil(t, inst.TriggerContext)
	assert.Equal(t, float64(5000), inst.TriggerContext["amount"]) // JSON numbers decode to float64
	assert.Equal(t, "USD", inst.TriggerContext["currency"])

	// Assert Policy Sub-Struct translation
	assert.Equal(t, "pol-v1", inst.Policy.ID)
	assert.Equal(t, sdk.MaterialityHigh, inst.Policy.Materiality)
	assert.True(t, inst.Policy.RequiresHumanApproval)
	assert.Equal(t, int64(3600), inst.Policy.ApprovalTimeoutSeconds)

	// Assert Metadata translation
	assert.Equal(t, "a1b2c3d4e5f6g7h8", inst.ArtifactHash)
	assert.Equal(t, 1, inst.ArtifactVersion)
	assert.Equal(t, "2026-03-10T10:00:00Z", inst.CreatedAt)
	assert.Equal(t, "2026-03-10T10:05:00Z", inst.UpdatedAt)
}
