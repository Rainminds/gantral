package authority

import (
	"context"
	"testing"

	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

// -- Mocks --

type MockEmitter struct {
	mock.Mock
}

func (m *MockEmitter) EmitArtifact(
	ctx context.Context,
	teamID string,
	instanceID string,
	workflowVersionID string,
	prevHash string,
	state string,
	policyVersionID string,
	contextHash string,
	actorID string,
	justification string,
) (*models.CommitmentArtifact, error) {
	args := m.Called(ctx, teamID, instanceID, workflowVersionID, prevHash, state, policyVersionID, contextHash, actorID, justification)
	return args.Get(0).(*models.CommitmentArtifact), args.Error(1)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetInstance(ctx context.Context, id string) (*engine.Instance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*engine.Instance), args.Error(1)
}

func (m *MockDB) RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error) {
	args := m.Called(ctx, cmd, nextState)
	return args.Get(0).(*engine.Instance), args.Error(1)
}
func (m *MockDB) CreateInstance(ctx context.Context, inst *engine.Instance) error { return nil }
func (m *MockDB) ListInstances(ctx context.Context, teamID string) ([]*engine.Instance, error) {
	return nil, nil
}
func (m *MockDB) GetAuditEvents(ctx context.Context, instanceID string) ([]engine.AuditEvent, error) {
	return nil, nil
}

type MockStore struct {
	WriteFn func(ctx context.Context, art *models.CommitmentArtifact) error
}

func (m *MockStore) WriteArtifact(ctx context.Context, art *models.CommitmentArtifact) error {
	return m.WriteFn(ctx, art)
}

func (m *MockStore) GetArtifact(ctx context.Context, artifactID string) (*models.CommitmentArtifact, error) {
	return nil, nil // Not used in these tests yet
}

func Test_Opaque_Handling(t *testing.T) {
	// The "Privacy Wall" Test
	// Objective: Verify Control Plane accepts an Evidence Hash and mints an artifact
	// WITHOUT requiring the raw payload (ContextSnapshot can be nil/empty).

	mockEmitter := new(MockEmitter)
	mockDB := new(MockDB)

	executor := &activities.ExecutionActivities{
		DB:              mockDB,
		ArtifactEmitter: mockEmitter,
	}

	// Temporal Test Suite setup
	s := &testsuite.WorkflowTestSuite{}
	env := s.NewTestActivityEnvironment()
	env.RegisterActivity(executor)

	evidenceHash := "sha256:valid-hash-of-opaque-payload"
	instanceID := "inst-privacy"

	// 1. Setup DB to return instance (for chaining)
	mockDB.On("GetInstance", mock.Anything, instanceID).Return(&engine.Instance{
		ID:                instanceID,
		LastArtifactHash:  "prev-hash",
		WorkflowVersionID: "wf-1",
	}, nil)

	// 2. Expect EmitArtifact to use the EvidenceHash directly
	expectedArtifact := &models.CommitmentArtifact{
		ArtifactID: "new-art-id",
	}
	mockEmitter.On("EmitArtifact",
		mock.Anything,
		"team-1",
		instanceID,
		"wf-1",
		"prev-hash",
		"APPROVED",
		"v1",
		evidenceHash, // <--- CRITICAL: Must match input hash
		"tool-runner",
		"",
	).Return(expectedArtifact, nil)

	// 3. Expect DB Record
	mockDB.On("RecordDecision", mock.Anything, mock.Anything, engine.StateApproved).Return(&engine.Instance{}, nil)

	// 4. Exec via Env
	input := activities.RecordDecisionInput{
		TeamID:          "team-1",
		InstanceID:      instanceID,
		DecisionType:    engine.DecisionApprove,
		ActorID:         "tool-runner",
		Role:            "machine",
		PolicyVersionID: "v1",
		EvidenceHash:    evidenceHash, // Providing Hash ONLY
		ContextSnapshot: nil,          // Raw payload is MISSING/NIL
	}

	future, err := env.ExecuteActivity(executor.RecordDecision, input)
	assert.NoError(t, err)

	var art *models.CommitmentArtifact
	err = future.Get(&art)
	assert.NoError(t, err)

	// 5. Assertions
	assert.Equal(t, expectedArtifact, art)

	// Assert Validated
	mockEmitter.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
