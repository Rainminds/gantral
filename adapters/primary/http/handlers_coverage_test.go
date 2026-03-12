package http

import (
	"context"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/stretchr/testify/mock"
)

func TestHandlers_ErrorPaths(t *testing.T) {
	mockTemporal := new(MockTemporalClient)
	mockStore := new(MockReadStore)
	handler := &Handler{
		TemporalClient: mockTemporal,
		ReadStore:      mockStore,
		TaskQueue:      "test-queue",
	}

	t.Run("CreateInstance_NoTeamID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/instances", strings.NewReader(`{}`))
		w := httptest.NewRecorder()
		handler.CreateInstance(w, req)
		if w.Code != stdhttp.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("RecordDecision_NoTeamID", func(t *testing.T) {
		// No path value set, returns 400
		req := httptest.NewRequest("POST", "/instances/1/decisions", strings.NewReader(`{}`))
		w := httptest.NewRecorder()
		handler.RecordDecision(w, req)
		if w.Code != stdhttp.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("GetInstance_StoreError", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/instances/1", nil)
		req.SetPathValue("id", "1")
		ctx := context.WithValue(req.Context(), TeamIDKey, "team-1")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockStore.On("GetInstance", mock.Anything, "1").Return((*engine.Instance)(nil), errors.New("db fail")).Once()

		handler.HandleGetInstance(w, req)
		if w.Code != stdhttp.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("ListInstances_NoTeamID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/instances", nil)
		w := httptest.NewRecorder()
		handler.HandleListInstances(w, req)
		if w.Code != stdhttp.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})

	t.Run("GetAuditLogs_StoreError", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/instances/1/audit", nil)
		req.SetPathValue("id", "1")
		ctx := context.WithValue(req.Context(), TeamIDKey, "team-1")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		inst := &engine.Instance{ID: "1", OwningTeamID: "team-1"}
		mockStore.On("GetInstance", mock.Anything, "1").Return(inst, nil).Once()
		mockStore.On("GetAuditEvents", mock.Anything, "1").Return([]engine.AuditEvent(nil), errors.New("fail")).Once()

		handler.HandleGetAuditLogs(w, req)
		if w.Code != stdhttp.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}
