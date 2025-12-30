package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
)

func TestRoutes(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)
	mux := srv.Routes()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/healthz", http.StatusOK},
		{"POST", "/instances", http.StatusCreated},
		// GET /instances/123 will be 404 because 123 doesn't exist in new engine
		{"GET", "/instances/123", http.StatusNotFound},
		{"GET", "/notfound", http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		if tt.method == "POST" {
			req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(`{}`))
		}
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != tt.status {
			t.Errorf("%s %s expected %d, got %d", tt.method, tt.path, tt.status, w.Code)
		}
	}
}

func TestLoggerMiddleware(t *testing.T) {
	// Mock handler that returns 200
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	middleware := LoggerMiddleware(next)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
