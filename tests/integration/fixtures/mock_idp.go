package fixtures

import (
	"context"
	"fmt"
)

// MockIDP implements a simple OIDC provider for testing.
type MockIDP struct {
	ValidTokens map[string]UserInfo
}

type UserInfo struct {
	ID    string
	Role  string
	Email string
}

func (m *MockIDP) VerifyToken(ctx context.Context, token string) (*UserInfo, error) {
	if info, ok := m.ValidTokens[token]; ok {
		return &info, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func NewMockIDP() *MockIDP {
	return &MockIDP{
		ValidTokens: map[string]UserInfo{
			"valid-user-1": {ID: "user-1", Role: "approver", Email: "user1@example.com"},
			"valid-admin":  {ID: "admin-1", Role: "admin", Email: "admin@example.com"},
		},
	}
}
