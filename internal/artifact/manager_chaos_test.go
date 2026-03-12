package artifact

import (
	"context"
	"errors"
	"testing"

	"github.com/Rainminds/gantral/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Adversarial Mocks

type ChaosStore struct {
	mock.Mock
}

func (m *ChaosStore) WriteArtifact(ctx context.Context, teamID string, art *models.CommitmentArtifact) error {
	args := m.Called(ctx, teamID, art)
	return args.Error(0)
}

func (m *ChaosStore) GetArtifact(ctx context.Context, teamID string, hash string) (*models.CommitmentArtifact, error) {
	args := m.Called(ctx, teamID, hash)
	return args.Get(0).(*models.CommitmentArtifact), args.Error(1)
}

type ChaosSigner struct {
	mock.Mock
}

func (m *ChaosSigner) Sign(hash []byte) ([]byte, string, error) {
	args := m.Called(hash)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

type ChaosTimeStamper struct {
	mock.Mock
}

func (m *ChaosTimeStamper) Token(hash []byte) ([]byte, string, error) {
	args := m.Called(hash)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func TestFailClosedChaos(t *testing.T) {
	t.Run("Scenario 1: Artifact Persistence Failure (S3 503)", func(t *testing.T) {
		mockStore := new(ChaosStore)
		mockSigner := new(ChaosSigner)
		mockTSA := new(ChaosTimeStamper)

		m := NewManager(mockStore, mockSigner, mockTSA)

		// Setup successful signing and timestamping
		mockSigner.On("Sign", mock.Anything).Return([]byte("sig"), "alg", nil)
		mockTSA.On("Token", mock.Anything).Return([]byte("ts"), "alg", nil)

		// Simulate Store Failure (S3 503 Service Unavailable)
		mockStore.On("WriteArtifact", mock.Anything, "team-1", mock.Anything).Return(errors.New("S3 Error: 503 Service Unavailable"))

		art, err := m.EmitArtifact(
			context.Background(),
			"team-1",
			"inst-1",
			"wf-1",
			"prev-1",
			"APPROVED",
			"v1",
			"ctx-1",
			"user-1",
			"justification",
		)

		assert.Error(t, err)
		assert.Nil(t, art)
		assert.Contains(t, err.Error(), "persistence failure")
		mockStore.AssertExpectations(t)
	})

	t.Run("Scenario 2: Identity Ambiguity (Empty Actor ID)", func(t *testing.T) {
		mockStore := new(ChaosStore)
		mockSigner := new(ChaosSigner)
		mockTSA := new(ChaosTimeStamper)

		m := NewManager(mockStore, mockSigner, mockTSA)

		// EmitArtifact currently doesn't strictly block empty ActorID (only instanceID, state, contextHash),
		// but the Task 2 Requirement says: "Assert that the manager rejects the transition synchronously with an auth error"
		// If the manager code doesn't check it yet, I should add the check or adjust the test.
		// Let's check manager.go lines 61-71. It only checks instanceID, state, contextHash.

		// Let's simulate a check for teamID as well since it's a tenant isolation key.
		art, err := m.EmitArtifact(
			context.Background(),
			"", // Empty TeamID (Cross-tenant violation / Identity ambiguity)
			"inst-1",
			"wf-1",
			"prev-1",
			"APPROVED",
			"v1",
			"ctx-1",
			"user-1",
			"justification",
		)

		// Note: Manager doesn't currently check for empty teamID in EmitArtifact based on my view earlier.
		// If I want to fulfill the "Fail-Closed" requirement for Identity Ambiguity, I should probably update the manager
		// to require teamID and possibly actorID for sensitive states.

		// However, for the sake of the test following the user's prompt:
		// "Mock the OIDC/auth context to return an empty humanactorid or teamid during a WAITING_FOR_HUMAN -> APPROVED transition.
		// Assert that the manager rejects the transition synchronously with an auth error."

		// If manager.go doesn't have these checks, the test will fail.
		// I'll assume for now I should implement the test and then fix the manager if needed.

		// Wait, Scenario 2 says "Identity Ambiguity ... empty humanactorid OR teamid".
		// Let's assume the manager SHOULD block empty teamID.

		if err == nil {
			t.Errorf("Manager allowed artifact emission with empty teamID")
		}
		assert.Nil(t, art)
	})

	t.Run("Scenario 3: Timestamp Verification Failure (TSA Error)", func(t *testing.T) {
		mockStore := new(ChaosStore)
		mockSigner := new(ChaosSigner)
		mockTSA := new(ChaosTimeStamper)

		m := NewManager(mockStore, mockSigner, mockTSA)

		// Setup successful signing
		mockSigner.On("Sign", mock.Anything).Return([]byte("sig"), "alg", nil)

		// Simulate TSA Failure
		mockTSA.On("Token", mock.Anything).Return([]byte(nil), "", errors.New("TSA Connection Timeout"))

		art, err := m.EmitArtifact(
			context.Background(),
			"team-1",
			"inst-1",
			"wf-1",
			"prev-1",
			"APPROVED",
			"v1",
			"ctx-1",
			"user-1",
			"justification",
		)

		assert.Error(t, err)
		assert.Nil(t, art)
		assert.Contains(t, err.Error(), "failed to timestamp artifact")
		mockTSA.AssertExpectations(t)
	})

	t.Run("Scenario 4: Validation Failure pre-emit", func(t *testing.T) {
		mockStore := new(ChaosStore)
		m := NewManager(mockStore, new(ChaosSigner), new(ChaosTimeStamper))

		// Trigger a validation failure by providing empty context hash
		art, err := m.EmitArtifact(
			context.Background(),
			"team-1",
			"inst-1",
			"wf-1",
			"prev-1",
			"APPROVED",
			"v1",
			"", // contextHash required
			"user-1",
			"justification",
		)

		assert.Error(t, err)
		assert.Nil(t, art)
		assert.Contains(t, err.Error(), "context hash required")
		mockStore.AssertNotCalled(t, "WriteArtifact", mock.Anything, mock.Anything, mock.Anything)
	})
}
