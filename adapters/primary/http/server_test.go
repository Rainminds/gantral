package http

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	// Use nil dependencies for route registration check.
	// NewServer constructs the Handler; we verifies Routes() registers paths correctly.

	srv := NewServer("8080", nil, "queue", nil)
	// Routes() registers handlers but doesn't execute them, so nil dependencies are safe here.
	mux := srv.Routes()

	// Check Healthz (does not require client/store)
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusOK {
		t.Errorf("/healthz expected 200, got %d", w.Code)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	next := stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	middleware := LoggingMiddleware(next) // Name changed in implementation to LoggingMiddleware or similar? handlers.go has LoggingMiddleware.

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
