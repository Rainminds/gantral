package activities

import (
	"context"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

// Mocks

type MockInstanceStore struct {
	mock.Mock
}

func (m *MockInstanceStore) CreateInstance(ctx context.Context, inst *engine.Instance) error {
	args := m.Called(ctx, inst)
	return args.Error(0)
}

func (m *MockInstanceStore) GetInstance(ctx context.Context, id string) (*engine.Instance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*engine.Instance), args.Error(1)
}

func (m *MockInstanceStore) ListInstances(ctx context.Context) ([]*engine.Instance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*engine.Instance), args.Error(1)
}

func (m *MockInstanceStore) GetAuditEvents(ctx context.Context, instanceID string) ([]engine.AuditEvent, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).([]engine.AuditEvent), args.Error(1)
}

func (m *MockInstanceStore) RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error) {
	args := m.Called(ctx, cmd, nextState)
	return args.Get(0).(*engine.Instance), args.Error(1)
}

type MockArtifactEmitter struct {
	mock.Mock
}

func (m *MockArtifactEmitter) EmitArtifact(ctx context.Context, instanceID, prevHash, state, policyVer, contextHash, actorID string) (*models.CommitmentArtifact, error) {
	args := m.Called(ctx, instanceID, prevHash, state, policyVer, contextHash, actorID)
	return args.Get(0).(*models.CommitmentArtifact), args.Error(1)
}

func TestRecordDecision_Chaining(t *testing.T) {
	// Setup
	mockDB := new(MockInstanceStore)
	mockEmitter := new(MockArtifactEmitter)
	activities := &ExecutionActivities{
		DB:              mockDB,
		ArtifactEmitter: mockEmitter,
	}

	// Temporal Test Suite
	s := &testsuite.WorkflowTestSuite{}
	env := s.NewTestActivityEnvironment()
	env.RegisterActivity(activities)

	instanceID := "inst-123"
	prevHash := "hash-abc-123"
	policyVer := "v1.0"
	contextSnapshot := map[string]interface{}{"foo": "bar"}

	// Deterministic hash of {"foo":"bar"}
	expectedContextHash, _ := artifact.HashContext(contextSnapshot)

	// 1. Expect GetInstance to return instance with LastArtifactHash
	mockDB.On("GetInstance", mock.Anything, instanceID).Return(&engine.Instance{
		ID:               instanceID,
		State:            engine.StateWaitingForHuman,
		LastArtifactHash: prevHash,
	}, nil)

	// 2. Expect EmitArtifact with correct prevHash and contextHash
	expectedArtifact := &models.CommitmentArtifact{
		ArtifactID:     "art-new",
		AuthorityState: "APPROVED",
	}
	mockEmitter.On("EmitArtifact", mock.Anything, instanceID, prevHash, "APPROVED", policyVer, expectedContextHash, "user-1").Return(expectedArtifact, nil)

	// 3. Expect RecordDecision (DB) with NewArtifactHash
	expectedCmd := engine.RecordDecisionCmd{
		InstanceID:      instanceID,
		Type:            engine.DecisionApprove,
		ActorID:         "user-1",
		Justification:   "Approved via test",
		Role:            "admin",
		ContextSnapshot: contextSnapshot,
		PolicyVersionID: policyVer,
		NewArtifactHash: "art-new", // Critical Check: Matches ID
	}
	mockDB.On("RecordDecision", mock.Anything, expectedCmd, engine.StateApproved).Return(&engine.Instance{}, nil)

	// Execute via Env
	input := RecordDecisionInput{
		InstanceID:      instanceID,
		DecisionType:    engine.DecisionApprove,
		ActorID:         "user-1",
		Justification:   "Approved via test",
		Role:            "admin",
		ContextSnapshot: contextSnapshot,
		PolicyVersionID: policyVer,
	}

	future, err := env.ExecuteActivity(activities.RecordDecision, input)
	assert.NoError(t, err)

	var art *models.CommitmentArtifact
	err = future.Get(&art)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, expectedArtifact, art)
	mockDB.AssertExpectations(t)
	mockEmitter.AssertExpectations(t)
}
