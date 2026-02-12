package policy

import (
	"context"
	"errors"
	"testing"

	corepolicy "github.com/Rainminds/gantral/core/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockEvaluator struct {
	mock.Mock
}

func (m *MockEvaluator) Evaluate(ctx context.Context, p corepolicy.Policy) (corepolicy.EvaluationResult, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(corepolicy.EvaluationResult), args.Error(1)
}

// -- Tests --

func Test_Policy_Down_Fails_Closed(t *testing.T) {
	mockEval := new(MockEvaluator)
	safeEngine := NewFailClosedEngine(mockEval)
	ctx := context.Background()
	p := corepolicy.Policy{ID: "test-policy"}

	// Case: Underlying Engine fails (Network/Timeout)
	mockEval.On("Evaluate", ctx, p).Return(corepolicy.EvaluationResult{}, errors.New("connection timeout"))

	result, err := safeEngine.Evaluate(ctx, p)

	// Assert: No error returned (error handled), but SAFE state forced.
	assert.NoError(t, err)
	assert.True(t, result.ShouldPause)
	assert.Equal(t, "WAITING_FOR_HUMAN", result.NextState)
	assert.Contains(t, result.Reason, "Fail-Closed: Policy Error")
}

func Test_Ambiguous_Result_Fails_Closed(t *testing.T) {
	mockEval := new(MockEvaluator)
	safeEngine := NewFailClosedEngine(mockEval)
	ctx := context.Background()
	p := corepolicy.Policy{ID: "test-policy"}

	// Case: Underlying Engine returns success but empty state (Invalid Logic)
	ambiguousResult := corepolicy.EvaluationResult{
		ShouldPause: false,
		NextState:   "", // BAD
		Reason:      "Who knows?",
	}
	mockEval.On("Evaluate", ctx, p).Return(ambiguousResult, nil)

	result, err := safeEngine.Evaluate(ctx, p)

	// Assert: Fail Closed
	assert.NoError(t, err)
	assert.True(t, result.ShouldPause)
	assert.Equal(t, "WAITING_FOR_HUMAN", result.NextState)
	assert.Contains(t, result.Reason, "Ambiguous Result")
}

func Test_Valid_PassThrough(t *testing.T) {
	mockEval := new(MockEvaluator)
	safeEngine := NewFailClosedEngine(mockEval)
	ctx := context.Background()
	p := corepolicy.Policy{ID: "test-policy"}

	// Case: Valid
	validResult := corepolicy.EvaluationResult{
		ShouldPause: false,
		NextState:   "RUNNING",
		Reason:      "OK",
	}
	mockEval.On("Evaluate", ctx, p).Return(validResult, nil)

	result, err := safeEngine.Evaluate(ctx, p)

	assert.NoError(t, err)
	assert.Equal(t, "RUNNING", result.NextState)
}
