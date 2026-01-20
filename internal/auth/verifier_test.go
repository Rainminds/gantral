package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVerifier for testing
type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) Verify(ctx context.Context, tokenString string) (*Identity, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Identity), args.Error(1)
}

func TestMultiVerifier(t *testing.T) {
	t.Run("First Verifier Succeeds", func(t *testing.T) {
		v1 := new(MockVerifier)
		v2 := new(MockVerifier)

		expected := &Identity{Subject: "user1", Type: IdentityTypeHuman}
		v1.On("Verify", mock.Anything, "token").Return(expected, nil)

		mv := NewMultiVerifier(v1, v2)
		identity, err := mv.Verify(context.Background(), "token")

		assert.NoError(t, err)
		assert.Equal(t, expected, identity)
		v1.AssertExpectations(t)
		v2.AssertNotCalled(t, "Verify")
	})

	t.Run("Second Verifier Succeeds", func(t *testing.T) {
		v1 := new(MockVerifier)
		v2 := new(MockVerifier)

		v1.On("Verify", mock.Anything, "token").Return(nil, errors.New("invalid signature"))

		expected := &Identity{Subject: "machine1", Type: IdentityTypeMachine}
		v2.On("Verify", mock.Anything, "token").Return(expected, nil)

		mv := NewMultiVerifier(v1, v2)
		identity, err := mv.Verify(context.Background(), "token")

		assert.NoError(t, err)
		assert.Equal(t, expected, identity)
		v1.AssertExpectations(t)
		v2.AssertExpectations(t)
	})

	t.Run("All Verifiers Fail", func(t *testing.T) {
		v1 := new(MockVerifier)
		v2 := new(MockVerifier)

		v1.On("Verify", mock.Anything, "token").Return(nil, errors.New("fail 1"))
		v2.On("Verify", mock.Anything, "token").Return(nil, errors.New("fail 2"))

		mv := NewMultiVerifier(v1, v2)
		identity, err := mv.Verify(context.Background(), "token")

		assert.Error(t, err)
		assert.Nil(t, identity)
		assert.Contains(t, err.Error(), "fail 1")
		assert.Contains(t, err.Error(), "fail 2")
	})
}
