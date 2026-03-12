package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/constants"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Hardening Audit Tests - Validating Phase 7

// 1. Artifact Schema Lock (7.0 & 7.2)
func TestAudit_Phase7_ArtifactSchemaLock(t *testing.T) {
	t.Run("Artifact version strictly v1, extra fields fail hashing", func(t *testing.T) {
		// TRD Invariant: Spec version MUST be embedded
		assert.Equal(t, "v1", models.CurrentArtifactVersion)

		// Test unsupported version failure
		art := models.CommitmentArtifact{
			ArtifactVersion: "v9.9", // Illegal version
			InstanceID:      "inst-1",
			AuthorityState:  string(constants.StateApproved),
		}

		_, err := art.PayloadForHashing()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported artifact version")
	})
}

// 2. Tenant Isolation (7.3)
func TestAudit_Phase7_TenantIsolation(t *testing.T) {
	t.Run("Cross-tenant access globally blocked at Engine level", func(t *testing.T) {
		mockStore := &mockStoreForHardening{}
		e := engine.NewEngine(mockStore)

		// Simulate Team B attempting to access Team A's instance
		instanceID := "team-a-instance"
		mockStore.instance = &engine.Instance{
			ID:           instanceID,
			OwningTeamID: "team-a",
		}

		_, err := e.GetInstance(context.Background(), "team-b", instanceID)
		require.Error(t, err)
		assert.ErrorIs(t, err, engine.ErrCrossTenantViolation)
	})
}

// 3. Fail-Closed Semantics (7.5)
func TestAudit_Phase7_FailClosedSemantics(t *testing.T) {
	t.Run("S3 persistence failure blocks authority decision (Atomicity)", func(t *testing.T) {
		mockStore := new(ChaosStoreMock)
		m := artifact.NewManager(mockStore, new(ChaosSignerMock), new(ChaosTSAMock))

		// Simulate S3 write failure
		mockStore.On("WriteArtifact", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("S3_STORAGE_UNAVAILABLE"))

		_, err := m.EmitArtifact(
			context.Background(),
			"team-1",
			"inst-1",
			"wf-v1",
			"prev-hash",
			string(constants.StateApproved),
			"policy-v1",
			"context-hash",
			"actor-1",
			"justification",
		)

		// Fail-Closed Assertion
		require.Error(t, err)
		assert.Contains(t, err.Error(), "persistence failure")
	})
}

// Mocks for Phase 7 Verification

type mockStoreForHardening struct {
	instance *engine.Instance
}

var _ engine.InstanceStore = (*mockStoreForHardening)(nil)

func (m *mockStoreForHardening) CreateInstance(ctx context.Context, inst *engine.Instance) error { return nil }
func (m *mockStoreForHardening) GetInstance(ctx context.Context, id string) (*engine.Instance, error) {
	return m.instance, nil
}
func (m *mockStoreForHardening) ListInstances(ctx context.Context, teamID string) ([]*engine.Instance, error) {
	return nil, nil
}
func (m *mockStoreForHardening) RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error) {
	return nil, nil
}
func (m *mockStoreForHardening) GetAuditEvents(ctx context.Context, id string) ([]engine.AuditEvent, error) {
	return nil, nil
}

type ChaosStoreMock struct {
	mock.Mock
}

func (m *ChaosStoreMock) WriteArtifact(ctx context.Context, teamID string, art *models.CommitmentArtifact) error {
	args := m.Called(ctx, teamID, art)
	return args.Error(0)
}
func (m *ChaosStoreMock) GetArtifact(ctx context.Context, teamID string, hash string) (*models.CommitmentArtifact, error) {
	return nil, nil
}

type ChaosSignerMock struct{ mock.Mock }

func (m *ChaosSignerMock) Sign(hash []byte) ([]byte, string, error) {
	return []byte("sig"), "alg", nil
}

type ChaosTSAMock struct{ mock.Mock }

func (m *ChaosTSAMock) Token(hash []byte) ([]byte, string, error) {
	return []byte("ts"), "alg", nil
}
