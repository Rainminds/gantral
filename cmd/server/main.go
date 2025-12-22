package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rainminds/gantral/api"
	"github.com/Rainminds/gantral/core/engine"
)

func main() {
	// 1. Initialize Engine
	eng := engine.NewEngine()

	// 2. Initialize API Server
	srv := api.NewServer(eng)

	// 3. Configure HTTP Server with Middleware
	// We wrap the Mux with our LoggerMiddleware
	handler := api.LoggerMiddleware(srv.Routes())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// 4. Start Server in Goroutine
	go func() {
		log.Printf("Starting Gantral Core on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 5. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait for signal

	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
