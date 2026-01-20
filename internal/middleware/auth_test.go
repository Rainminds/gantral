package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Rainminds/gantral/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVerifier
type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) Verify(ctx context.Context, tokenString string) (*auth.Identity, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Identity), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("Valid Token", func(t *testing.T) {
		mockVerifier := new(MockVerifier)
		mw := AuthMiddleware(mockVerifier, logger)

		expectedIdentity := &auth.Identity{Subject: "user123", Type: auth.IdentityTypeHuman}
		mockVerifier.On("Verify", mock.Anything, "valid.token").Return(expectedIdentity, nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := GetIdentity(r.Context())
			assert.NoError(t, err)
			assert.Equal(t, "user123", identity.Subject)
			assert.Equal(t, auth.IdentityTypeHuman, identity.Type)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer valid.token")
		rec := httptest.NewRecorder()

		mw(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockVerifier.AssertExpectations(t)
	})

	t.Run("No Header", func(t *testing.T) {
		mockVerifier := new(MockVerifier)
		mw := AuthMiddleware(mockVerifier, logger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/", nil)
		// No Authorization header
		rec := httptest.NewRecorder()

		mw(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "missing authorization header")
	})

	t.Run("Invalid Format", func(t *testing.T) {
		mockVerifier := new(MockVerifier)
		mw := AuthMiddleware(mockVerifier, logger)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		rec := httptest.NewRecorder()

		mw(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid authorization header format")
	})

	t.Run("Invalid Token (Verification Failed)", func(t *testing.T) {
		mockVerifier := new(MockVerifier)
		mw := AuthMiddleware(mockVerifier, logger)

		mockVerifier.On("Verify", mock.Anything, "bad.token").Return(nil, errors.New("signature invalid"))

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Next handler should not be called")
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer bad.token")
		rec := httptest.NewRecorder()

		mw(nextHandler).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("RequireRole Success", func(t *testing.T) {
		// Test the RBAC middleware specifically
		identity := &auth.Identity{
			Subject: "runner1",
			Roles:   []string{"runner"},
		}

		handler := RequireRole("runner")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest("POST", "/poll", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, identity)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req.WithContext(ctx))

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireRole Forbidden", func(t *testing.T) {
		identity := &auth.Identity{
			Subject: "user1",
			Roles:   []string{"user"},
		}

		handler := RequireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Should not run")
		}))

		req := httptest.NewRequest("POST", "/admin", nil)
		ctx := context.WithValue(req.Context(), UserContextKey, identity)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req.WithContext(ctx))

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
