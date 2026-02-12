//go:build integration

package integration_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Rainminds/gantral/adapters/secondary/postgres"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"github.com/Rainminds/gantral/infra"
	"github.com/Rainminds/gantral/internal/auth"
	"github.com/Rainminds/gantral/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Helper to generate a dev token
func generateDevToken(secret []byte, sub string, iType string, roles ...string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"type":  iType,
		"roles": roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute).Unix(),
	})
	s, _ := token.SignedString(secret)
	return s
}

// Ensure the original integration test remains active
func TestIntegration_HITL_Flow(t *testing.T) {
	// 0. Setup
	_ = godotenv.Load("../.env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:changeme@127.0.0.1:5432/gantral?sslmode=disable"
		t.Logf("DATABASE_URL not set, using default: %s", dbURL)
	}

	// Run Migrations ensure schema exists
	if err := infra.RunMigrations(dbURL); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	ctx := context.Background()
	store, err := postgres.NewStore(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer store.Close()

	// 1. Initialize Engine with Real Store
	eng := engine.NewEngine(store)

	// 2. Create High Materiality Instance (expect PAUSE)
	pol := policy.Policy{ID: "demo-policy", Materiality: policy.MaterialityHigh}

	inst, err := eng.CreateInstance(ctx, "integration-wf", map[string]interface{}{"key": "val"}, pol)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	if inst.State != engine.StateWaitingForHuman {
		t.Errorf("expected WAITING_FOR_HUMAN, got %s", inst.State)
	}

	// 4. Record Decision (Approve)
	cmd := engine.RecordDecisionCmd{
		InstanceID:    inst.ID,
		Type:          engine.DecisionApprove,
		ActorID:       "integration-tester",
		Justification: "Integration Test Approval",
	}

	updated, err := eng.RecordDecision(ctx, cmd)
	if err != nil {
		t.Fatalf("RecordDecision failed: %v", err)
	}

	if updated.State != engine.StateApproved {
		t.Errorf("expected APPROVED, got %s", updated.State)
	}
}

func TestIntegration_RBAC_Enforcement(t *testing.T) {
	// This test sets up the HTTP middleware stack to verify Role checks works isolation

	// Logger needed for AuthMiddleware
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	secret := []byte("test-secret")
	// Use DevVerifier configured as "Machine Verifier"
	machineVerifier := &auth.DevVerifier{Secret: secret, IdentityType: auth.IdentityTypeMachine}

	authMw := middleware.AuthMiddleware(machineVerifier, logger)

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Chain: Auth -> RBAC -> Handler
	// SCENARIO 1: Machine Token accessing User Endpoint (Should Fail)
	t.Run("Machine accessing User Endpoint", func(t *testing.T) {
		rbacMw := middleware.RequireRole("user", "admin")
		handler := authMw(rbacMw(baseHandler))

		machineToken := generateDevToken(secret, "runner-001", "machine", "runner")

		req := httptest.NewRequest("POST", "/decisions", nil)
		req.Header.Set("Authorization", "Bearer "+machineToken)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
