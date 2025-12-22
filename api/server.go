package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Rainminds/gantral/core/engine"
)

// Server holds the dependencies for the HTTP API.
type Server struct {
	engine *engine.Engine
}

// NewServer creates a new API server with the given engine.
func NewServer(e *engine.Engine) *Server {
	return &Server{
		engine: e,
	}
}

// Routes returns the http.ServeMux with all registered routes.
// It uses Go 1.22+ pattern matching.
func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register routes using Go 1.22 method + path pattern
	mux.HandleFunc("POST /instances", s.handleCreateInstance)
	mux.HandleFunc("GET /instances/{id}", s.handleGetInstance)

	// In a real implementation, we might wrap this mux in more middleware
	return mux
}

// LoggerMiddleware wraps an http.Handler to log request details.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.status, duration)
	})
}

// responseWriter is a wrapper to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
