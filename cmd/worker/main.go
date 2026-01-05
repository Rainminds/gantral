package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Rainminds/gantral/adapters/secondary/postgres"
	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/workflows"
	"github.com/Rainminds/gantral/pkg/config"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// 1. Initialize Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file if it exists
	_ = godotenv.Load()

	// 2. Configuration
	// Fail fast for structural dependencies
	temporalHost := config.MustGetEnv("TEMPORAL_HOST_PORT")
	taskQueue := config.MustGetEnv("TASK_QUEUE")
	// Fail fast if critical config is missing
	dbURL := config.MustGetEnv("DATABASE_URL")

	logger.Info("Starting Worker",
		"temporal_host", temporalHost,
		"task_queue", taskQueue,
	)

	// 3. Initialize DB Adapter (for Activities)
	ctx := context.Background()
	store, err := postgres.NewStore(ctx, dbURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	// 4. Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort: temporalHost,
	})
	if err != nil {
		logger.Error("Unable to create Temporal client", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// 5. Register Worker
	w := worker.New(c, taskQueue, worker.Options{})

	// Register Workflows
	w.RegisterWorkflow(workflows.GantralExecutionWorkflow)

	// Register Activities
	activityImpl := &activities.ExecutionActivities{DB: store}
	w.RegisterActivity(activityImpl)

	// 6. Run with Graceful Shutdown
	// InterruptCh() captures SIGINT and SIGTERM
	logger.Info("Worker started successfully")

	// Create a channel to listen for OS signals specifically if worker.Run doesn't handle everything or if we want custom logging before exit
	// worker.Run blocks until interrupt.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("Worker stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("Worker stopped gracefully")
}
