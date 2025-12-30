package api

import (
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/web"
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
	mux.HandleFunc("GET /instances", s.handleListInstances)
	mux.HandleFunc("GET /instances/{id}", s.handleGetInstance)
	mux.HandleFunc("POST /instances/{id}/decisions", s.handleRecordDecision)
	mux.HandleFunc("GET /healthz", s.handleHealthz)

	// Serve Static Files
	// We use fs.Sub to root the file server at "static" directory
	staticFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		// Should not happen if build is correct
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(staticFS)))
	// Also handle /dashboard specifically if desired, but user asked for /
	// Note: API routes registered above take precedence because they are more specific (in Go 1.22+ if using method spec)
	// but purely path-based matching depends on length.
	// Actually, "POST /instances" is more specific than "/".
	// "GET /" matches everything else.

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
		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", duration,
		)
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
