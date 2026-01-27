package main

import (
	"context"
	"log/slog"
	stdhttp "net/http" // Alias standard library
	"os"

	"time"

	gantralhttp "github.com/Rainminds/gantral/adapters/primary/http" // Alias primary adapter
	"github.com/Rainminds/gantral/adapters/secondary/postgres"
	"github.com/Rainminds/gantral/internal/auth"
	"github.com/Rainminds/gantral/internal/middleware"
	"github.com/Rainminds/gantral/pkg/config"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

func main() {
	// 1. Initialize Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file if it exists (silence error if missing, for cloud envs)
	_ = godotenv.Load()

	// 2. Configuration
	port := config.GetEnv("PORT", "8080") // Port defaults are standard convention

	// Fail fast for structural dependencies
	temporalHost := config.MustGetEnv("TEMPORAL_HOST_PORT")
	taskQueue := config.MustGetEnv("TASK_QUEUE")
	dbURL := config.MustGetEnv("DATABASE_URL")

	logger.Info("Starting Gantral API Server",
		"port", port,
		"temporal_host", temporalHost,
	)

	// 3. Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort: temporalHost,
	})
	if err != nil {
		logger.Error("Unable to create Temporal client", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// 4. Connect to Postgres (Read Path)
	store, err := postgres.NewStore(context.Background(), dbURL)
	if err != nil {
		logger.Error("Unable to create Postgres store", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	// 5. Setup Authentication
	var verifiers []auth.TokenVerifier
	devMode := config.GetEnv("DEV_MODE", "false") == "true"

	if devMode {
		logger.Info("Auth: Dev Mode Enabled (HS256)")
		secretKey := config.GetEnv("DEV_AUTH_SECRET", "dev-secret-key")

		// In dev mode, we trust the token to declare if it's human or machine.
		// For simplicity, we create one verifier that we'll assume is "Human" by default type,
		// but the verification logic in DevVerifier preserves roles.
		devVerifier := &auth.DevVerifier{
			Secret:       []byte(secretKey),
			IdentityType: auth.IdentityTypeHuman, // Default
		}
		verifiers = append(verifiers, devVerifier)

	} else {
		// Production: Multi-Issuer

		// 1. Human Identity (e.g. Okta/Auth0)
		if oidcIssuer := config.GetEnv("OIDC_ISSUER", ""); oidcIssuer != "" {
			oidcClientID := config.MustGetEnv("OIDC_CLIENT_ID")
			v, err := auth.NewOIDCVerifier(context.Background(), oidcIssuer, oidcClientID, auth.IdentityTypeHuman)
			if err != nil {
				logger.Error("Failed to initialize User OIDC verifier", "error", err)
				os.Exit(1)
			}
			verifiers = append(verifiers, v)
			logger.Info("Auth: Added User Issuer", "issuer", oidcIssuer)
		}

		// 2. Machine Identity (e.g. AWS/K8s)
		if machineIssuer := config.GetEnv("MACHINE_ISSUER", ""); machineIssuer != "" {
			machineAud := config.GetEnv("MACHINE_AUDIENCE", "gantral-core")
			v, err := auth.NewOIDCVerifier(context.Background(), machineIssuer, machineAud, auth.IdentityTypeMachine)
			if err != nil {
				logger.Error("Failed to initialize Machine OIDC verifier", "error", err)
				os.Exit(1)
			}
			verifiers = append(verifiers, v)
			logger.Info("Auth: Added Machine Issuer", "issuer", machineIssuer)
		}
	}

	if len(verifiers) == 0 {
		logger.Error("Auth: No verifiers configured! Set DEV_MODE=true or configure OIDC_ISSUER/MACHINE_ISSUER")
		os.Exit(1)
	}

	multiVerifier := auth.NewMultiVerifier(verifiers...)
	authMiddleware := middleware.AuthMiddleware(multiVerifier, logger)

	// 6. Start HTTP Server
	// Note: API talks to Temporal for Writes, Postgres for Reads (CQRS).
	srv := gantralhttp.NewServer(port, c, taskQueue, store)
	mux := srv.Routes()

	// 7. Manual RBAC implementation since we can't easily inject into the mux returned by adapters logic
	// We wrap the entire mux with a handler that checks specific paths.
	rbacHandler := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		path := r.URL.Path
		method := r.Method

		// Rule 1: /tasks/poll (POST) -> Machine only (Role: runner)
		if path == "/tasks/poll" && method == "POST" {
			middleware.RequireRole("runner")(mux).ServeHTTP(w, r)
			return
		}

		// Rule 2: /decisions (POST) -> User/Admin only
		if path == "/decisions" && method == "POST" {
			middleware.RequireRole("admin", "user")(mux).ServeHTTP(w, r)
			return
		}

		// Default: Pass through (AuthMiddleware already validated identity exists)
		mux.ServeHTTP(w, r)
	})

	// Chain: Logging -> Auth -> RBAC -> Routes
	finalHandler := gantralhttp.LoggingMiddleware(authMiddleware(rbacHandler))

	httpServer := &stdhttp.Server{
		Addr:              ":" + port,
		Handler:           finalHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("Server listening...", "addr", ":"+port)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
