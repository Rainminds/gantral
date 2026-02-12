package authority

import (
	"context"
	"errors"
	"testing"

	"github.com/Rainminds/gantral/internal/artifact"
	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockArtifactStore struct {
	mock.Mock
}

func (m *MockArtifactStore) Write(ctx context.Context, art *models.CommitmentArtifact) error {
	args := m.Called(ctx, art)
	return args.Error(0)
}

func (m *MockArtifactStore) Get(ctx context.Context, id string) (*models.CommitmentArtifact, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CommitmentArtifact), args.Error(1)
}

// Since Engine is a struct, we can't mock it directly easily unless we interface it.
// Phase 6.6 requirement was to wrap *Engine.
// But for testing the wrapper, we need to inject errors.
// The wrapper wraps *policy.Engine. *policy.Engine logic is hardcoded in core/policy/engine.go unless we interface it.
// Wait, core/policy/engine.go: type Engine struct{}.
// Evaluating it always succeeds currently.
// TO TEST FAIL-CLOSED, we need to simulate failure.
// We can't easily simulate failure in the *real* engine struct without mocking internal behavior or dependency.
// For the purpose of this test, we might need to modify core/policy to use an interface or allow error injection.
// OR, we assume the FailClosedEngine logic (if err != nil) is correct, but we can't trigger it with current Engine.
//
// Let's Verify Authority Safety first (ConsistencyGuard).
// For Policy Safety, I will check if I can modify Engine to interface or add a "SimulateError" flag for testing.
// Or better, define an interface in internal/policy/safety.go that the wrapper uses, allowing mocking.

// -- Tests --

func Test_Phantom_Artifact_Panic(t *testing.T) {
	mockStore := new(MockArtifactStore)
	guard := NewConsistencyGuard(mockStore)
	ctx := context.Background()

	// Case: Artifact Missing
	mockStore.On("Get", ctx, "phantom-123").Return(nil, artifact.ErrArtifactNotFound)

	err := guard.EnsureStateConsistency(ctx, "inst-1", "phantom-123")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrStateAmbiguous))
	assert.Contains(t, err.Error(), "artifact phantom-123 not found")
}

func Test_Cross_Instance_Contamination(t *testing.T) {
	mockStore := new(MockArtifactStore)
	guard := NewConsistencyGuard(mockStore)
	ctx := context.Background()

	// Case: Artifact Exists but belongs to wrong instance
	art := &models.CommitmentArtifact{
		ArtifactID: "leak-123",
		InstanceID: "other-inst",
	}
	mockStore.On("Get", ctx, "leak-123").Return(art, nil)

	err := guard.EnsureStateConsistency(ctx, "inst-1", "leak-123")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrStateAmbiguous))
	assert.Contains(t, err.Error(), "instance mismatch")
}

func Test_Valid_Consistency(t *testing.T) {
	mockStore := new(MockArtifactStore)
	guard := NewConsistencyGuard(mockStore)
	ctx := context.Background()

	// Case: Valid
	art := &models.CommitmentArtifact{
		ArtifactID: "valid-123",
		InstanceID: "inst-1",
	}
	mockStore.On("Get", ctx, "valid-123").Return(art, nil)

	err := guard.EnsureStateConsistency(ctx, "inst-1", "valid-123")

	assert.NoError(t, err)
}

// -- Policy Safety Tests --
// To test Policy fail-closed, we need an interface. I'll define one locally for the wrapper to use/test,
// or check if I can modify the wrapper to use an interface.
// `internal/policy/safety.go` defined `FailClosedEngine` wrapping `*Engine`.
// I should probably change `FailClosedEngine` to wrap an interface `Evaluator`.
