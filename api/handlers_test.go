package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
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

		// var resp map[string]string
		// if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		// 	t.Fatalf("failed to decode response: %v", err)
		// }
		// if resp["status"] != "created" {
		// 	t.Errorf("expected status 'created', got %q", resp["status"])
		// }

		var inst engine.Instance
		if err := json.NewDecoder(w.Body).Decode(&inst); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if inst.ID == "" {
			t.Error("expected instance ID to be set")
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
		// Seed an instance first for in-memory engine
		// Re-using 'eng' from outer scope.
		created, _ := eng.CreateInstance(context.Background(), "test-wf", nil, policy.Policy{ID: "p1"})
		idToFetch := created.ID

		req := httptest.NewRequest(http.MethodGet, "/instances/"+idToFetch, nil)
		req.SetPathValue("id", idToFetch)

		w := httptest.NewRecorder()
		srv.handleGetInstance(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var inst engine.Instance
		if err := json.NewDecoder(w.Body).Decode(&inst); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if inst.ID != idToFetch {
			t.Errorf("expected instance_id '%s', got %v", idToFetch, inst.ID)
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

func TestHandleListInstances(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)

	// Seed one
	eng.CreateInstance(context.Background(), "wf-1", nil, policy.Policy{ID: "p1"})

	req := httptest.NewRequest(http.MethodGet, "/instances", nil)
	w := httptest.NewRecorder()

	srv.handleListInstances(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var list []engine.Instance
	if err := json.NewDecoder(w.Body).Decode(&list); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 instance, got %d", len(list))
	}
}

func TestHandleRecordDecision(t *testing.T) {
	eng := engine.NewEngine()
	srv := NewServer(eng)

	// Seed a WAITING instance
	// Need to force HIGH materiality to get it paused
	hiPol := policy.Policy{ID: "hi", Materiality: policy.MaterialityHigh}
	inst, _ := eng.CreateInstance(context.Background(), "wf-1", nil, hiPol)
	id := inst.ID

	reqBody := `{"type": "APPROVE", "actor_id": "tester", "justification": "ok"}`
	req := httptest.NewRequest(http.MethodPost, "/instances/"+id+"/decisions", strings.NewReader(reqBody))
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	srv.handleRecordDecision(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}
