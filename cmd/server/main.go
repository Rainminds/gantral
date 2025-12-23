package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rainminds/gantral/api"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/infra"
	"github.com/Rainminds/gantral/pkg/logger"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Initialize structured logging
	logger.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err := run(ctx); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited properly")
}

func run(ctx context.Context) error {
	// 0. Load .env file if it exists
	_ = godotenv.Load()

	// 1. Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// 1. Database & Migrations
	slog.Info("starting database connection", "url", dbURL)

	// Run migrations
	if err := infra.RunMigrations(dbURL); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Connect to DB
	store, err := infra.NewStore(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer store.Close()
	slog.Info("connected to database")

	// 2. Initialize Core Engine
	eng := engine.NewEngine()

	// 3. Setup API Server using errgroup
	srv := api.NewServer(eng) // Pass store later when Engine supports it
	// Wrap with logger middleware
	handler := api.LoggerMiddleware(srv.Routes())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	g, gCtx := errgroup.WithContext(ctx)

	// Start HTTP Server
	g.Go(func() error {
		slog.Info("starting api server", "addr", ":8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("http server error: %w", err)
		}
		return nil
	})

	// Wait for context cancellation (signal or error) to shutdown
	g.Go(func() error {
		<-gCtx.Done()
		slog.Info("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		return nil
	})

	return g.Wait()
}
