package http

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	// We can use nil dependencies for testing route existence if we don't execute them deeply,
	// or assume Handler methods handle nil gracefully during simple wiring checks?
	// Better to use empty struct defaults or mocks if needed.
	// For routing check, we just verify paths are registered.

	// Since NewServer constructs a Handler which is used in Routes, we need to pass strict types?
	// NewServer(port, client, queue, store)

	srv := NewServer("8080", nil, "queue", nil)
	// Note: Calling methods on srv.Routes() might panic if handlers access nil client/store immediately.
	// But Routes() just registers them.
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
		w.Write([]byte("ok"))
	})

	middleware := LoggingMiddleware(next) // Name changed in implementation to LoggingMiddleware or similar? handlers.go has LoggingMiddleware.

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
