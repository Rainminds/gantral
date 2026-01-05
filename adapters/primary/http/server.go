package http

import (
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/Rainminds/gantral/core/ports"
	"github.com/Rainminds/gantral/web"
	"go.temporal.io/sdk/client"
)

// Server holds the dependencies for the HTTP API.
type Server struct {
	handler *Handler
}

// NewServer creates a new API server.
func NewServer(port string, temporalClient client.Client, taskQueue string, readStore ports.InstanceStore) *Server {
	return &Server{
		handler: &Handler{
			TemporalClient: temporalClient,
			TaskQueue:      taskQueue,
			ReadStore:      readStore,
		},
	}
}

// Routes returns the http.ServeMux with all registered routes.
func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Register routes using Go 1.22 method + path pattern
	mux.HandleFunc("POST /instances", s.handler.CreateInstance)
	mux.HandleFunc("POST /instances/{id}/decisions", s.handler.RecordDecision)
	mux.HandleFunc("GET /instances/{id}/audit", s.handler.HandleGetAuditLogs)
	mux.HandleFunc("GET /instances/{id}", s.handler.HandleGetInstance)
	mux.HandleFunc("GET /instances", s.handler.HandleListInstances)
	mux.HandleFunc("GET /healthz", s.handler.HealthCheck)

	// Serve Static Files
	staticFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(staticFS)))

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
