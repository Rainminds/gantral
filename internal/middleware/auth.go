package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Rainminds/gantral/internal/auth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type contextKey string

const (
	UserContextKey contextKey = "user_context"
)

var (
	authFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gantral_auth_failures_total",
		Help: "The total number of authentication failures by reason",
	}, []string{"reason", "identity_type"})
)

// AuthMiddleware creates a middleware that verifies Bearer tokens.
func AuthMiddleware(verifier auth.TokenVerifier, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handleAuthError(w, logger, "missing_header", "unknown", "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handleAuthError(w, logger, "invalid_format", "unknown", "invalid authorization header format")
				return
			}

			tokenString := parts[1]
			identity, err := verifier.Verify(r.Context(), tokenString)
			if err != nil {
				// Log with redaction (never log the token)
				logger.Info("auth_failed", "reason", "verify_failed", "error", err.Error())
				handleAuthError(w, logger, "verify_failed", "unknown", "invalid token")
				return
			}

			// Success - inject identity into context
			ctx := context.WithValue(r.Context(), UserContextKey, identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole creates a middleware that enforces RBAC.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := GetIdentity(r.Context())
			if err != nil {
				// Should have been caught by AuthMiddleware, but fail closed just in case
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, role := range identity.Roles {
				for _, allowed := range allowedRoles {
					if role == allowed {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			// Also allow checks by identity type if needed, e.g. "machine" role mapping
			// For now assuming roles are populated in the token or the verifier maps logic to roles.

			w.WriteHeader(http.StatusForbidden)
			// Check error to satisfy errcheck, use assignment to satisfy staticcheck (no empty branch)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "insufficient_permissions"})
		})
	}
}

func handleAuthError(w http.ResponseWriter, logger *slog.Logger, reason string, identityType string, clientMsg string) {
	authFailures.WithLabelValues(reason, identityType).Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	resp := map[string]string{
		"error": clientMsg,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed to write auth error response", "error", err)
	}
}

// GetIdentity helper to retrieve identity from context
func GetIdentity(ctx context.Context) (*auth.Identity, error) {
	identity, ok := ctx.Value(UserContextKey).(*auth.Identity)
	if !ok {
		return nil, errors.New("no identity in context")
	}
	return identity, nil
}
