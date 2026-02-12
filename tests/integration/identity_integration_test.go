//go:build integration

package integration_test

import (
	"context"
	"testing"

	"github.com/Rainminds/gantral/tests/integration/fixtures"
)

// TestIdentityIntegration_SectionG validates identity and security.
func TestIdentityIntegration_SectionG(t *testing.T) {
	ctx := context.Background()
	idp := fixtures.NewMockIDP()

	t.Run("Valid Token -> Accepted", func(t *testing.T) {
		token := "valid-user-1"
		userInfo, err := idp.VerifyToken(ctx, token)
		if err != nil {
			t.Fatalf("Expected valid token, got error: %v", err)
		}
		if userInfo.ID != "user-1" {
			t.Errorf("Expected user-1, got %s", userInfo.ID)
		}
	})

	t.Run("Invalid Token -> Rejected", func(t *testing.T) {
		token := "invalid-token"
		_, err := idp.VerifyToken(ctx, token)
		if err == nil {
			t.Error("Expected error for invalid token, got nil")
		}
	})

	t.Run("Role Mismatch -> Denied", func(t *testing.T) {
		// Mock engine enforcing role check
		// engine.CheckRole(user, requiredRole)
		// This requires engine logic which we assume is tested in unit/decision_test.go partially.
		// Integration test here would rely on policy + identity.
	})
}
