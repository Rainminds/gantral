package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
)

func TestHandleHealthz(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	srv.handleHealthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", w.Body.String())
	}
}

func TestHandleCreateInstance(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)

	t.Run("Valid Request", func(t *testing.T) {
		body := `{"workflow_id": "test-wf"}`
		req := httptest.NewRequest(http.MethodPost, "/instances", strings.NewReader(body))
		w := httptest.NewRecorder()

		srv.handleCreateInstance(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp map[string]string
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp["status"] != "created" {
			t.Errorf("expected status 'created', got %q", resp["status"])
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/instances", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()

		srv.handleCreateInstance(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestHandleGetInstance(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)

	t.Run("Valid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/instances/123", nil)
		req.SetPathValue("id", "123") // Go 1.22 feature

		w := httptest.NewRecorder()
		srv.handleGetInstance(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		var resp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp["instance_id"] != "123" {
			t.Errorf("expected instance_id '123', got %v", resp["instance_id"])
		}
	})

	t.Run("Missing ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/instances/", nil)
		// No path value set leads to empty ID check
		w := httptest.NewRecorder()
		srv.handleGetInstance(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/instances/notfound", nil)
		req.SetPathValue("id", "notfound")
		w := httptest.NewRecorder()
		srv.handleGetInstance(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
