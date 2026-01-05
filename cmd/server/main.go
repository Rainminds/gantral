package main

import (
	"context"
	"log/slog"
	stdhttp "net/http" // Alias standard library
	"os"

	gantralhttp "github.com/Rainminds/gantral/adapters/primary/http" // Alias primary adapter
	"github.com/Rainminds/gantral/adapters/secondary/postgres"
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

	// 5. Start HTTP Server
	// Note: API talks to Temporal for Writes, Postgres for Reads (CQRS).
	srv := gantralhttp.NewServer(port, c, taskQueue, store)

	// Create standard server using the mux from Routes()
	httpServer := &stdhttp.Server{
		Addr:    ":" + port,
		Handler: gantralhttp.LoggingMiddleware(srv.Routes()),
	}

	logger.Info("Server listening...", "addr", ":"+port)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
