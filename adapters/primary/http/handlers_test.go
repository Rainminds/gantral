package http

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// --- Mocks ---

type MockTemporalClient struct {
	mock.Mock
	client.Client // Embed to satisfy interface (will panic if unused methods called)
}

func (m *MockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	called := m.Called(ctx, options, workflow, args)
	// Return a simplified run stub
	if run, ok := called.Get(0).(client.WorkflowRun); ok {
		return run, called.Error(1)
	}
	return nil, called.Error(1)
}

func (m *MockTemporalClient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error {
	args := m.Called(ctx, workflowID, runID, signalName, arg)
	return args.Error(0)
}

type MockWorkflowRun struct {
	mock.Mock
}

func (m *MockWorkflowRun) GetID() string    { return "test-wf-id" }
func (m *MockWorkflowRun) GetRunID() string { return "test-run-id" }

// Needed to satisfy interface
func (m *MockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error { return nil }
func (m *MockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, options client.WorkflowRunGetOptions) error {
	return nil
}

type MockReadStore struct {
	mock.Mock
}

func (m *MockReadStore) CreateInstance(ctx context.Context, inst *engine.Instance) error {
	return m.Called(ctx, inst).Error(0)
}
func (m *MockReadStore) GetInstance(ctx context.Context, id string) (*engine.Instance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*engine.Instance), args.Error(1)
}
func (m *MockReadStore) ListInstances(ctx context.Context) ([]*engine.Instance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*engine.Instance), args.Error(1)
}
func (m *MockReadStore) RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error) {
	args := m.Called(ctx, cmd, nextState)
	return args.Get(0).(*engine.Instance), args.Error(1)
}
func (m *MockReadStore) GetAuditEvents(ctx context.Context, instanceID string) ([]engine.AuditEvent, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).([]engine.AuditEvent), args.Error(1)
}

// --- Tests ---

func TestCreateInstance(t *testing.T) {
	mockTemporal := new(MockTemporalClient)
	mockRun := new(MockWorkflowRun)
	handler := &Handler{
		TemporalClient: mockTemporal,
		TaskQueue:      "test-queue",
	}

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"workflow_id": "test-wf", "policy": {"id": "p1"}}`
		req := httptest.NewRequest("POST", "/instances", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		// Expect ExecuteWorkflow
		mockTemporal.On("ExecuteWorkflow",
			mock.Anything,
			mock.MatchedBy(func(opts client.StartWorkflowOptions) bool {
				return opts.TaskQueue == "test-queue" && strings.HasPrefix(opts.ID, "inst-")
			}),
			mock.Anything, // Workflow func
			mock.Anything, // Args
		).Return(mockRun, nil)

		handler.CreateInstance(w, req)

		if w.Code != stdhttp.StatusAccepted {
			t.Errorf("expected 202, got %d", w.Code)
		}

		var resp CreateInstanceResponse
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if resp.Status != "PENDING" {
			t.Errorf("expected PENDING, got %s", resp.Status)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/instances", strings.NewReader("{invalid"))
		w := httptest.NewRecorder()
		handler.CreateInstance(w, req)
		if w.Code != stdhttp.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestRecordDecision(t *testing.T) {
	mockTemporal := new(MockTemporalClient)
	handler := &Handler{
		TemporalClient: mockTemporal,
	}

	t.Run("Success Approve", func(t *testing.T) {
		reqBody := `{"type": "APPROVE", "actor_id": "user1", "justification": "LGTM"}`
		req := httptest.NewRequest("POST", "/instances/inst-1/decisions", strings.NewReader(reqBody))
		req.SetPathValue("id", "inst-1")
		w := httptest.NewRecorder()

		mockTemporal.On("SignalWorkflow",
			mock.Anything,
			"inst-1",
			"",
			"HumanDecision",
			mock.MatchedBy(func(arg interface{}) bool {
				// We can't easily cast to RecordDecisionInput here due to internal type,
				// but we verify the call happens.
				return true
			}),
		).Return(nil)

		handler.RecordDecision(w, req)

		if w.Code != stdhttp.StatusAccepted {
			t.Errorf("expected 202, got %d", w.Code)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		reqBody := `{"type": "APPROVE"}`
		req := httptest.NewRequest("POST", "/instances/missing/decisions", strings.NewReader(reqBody))
		req.SetPathValue("id", "missing")
		w := httptest.NewRecorder()

		// Use the correct type for NotFound error from Temporal SDK
		// serviceerror.NotFound is usually what is returned.
		notFoundErr := &serviceerror.NotFound{Message: "Workflow not found"}

		mockTemporal.On("SignalWorkflow", mock.Anything, "missing", "", "HumanDecision", mock.Anything).Return(notFoundErr)

		handler.RecordDecision(w, req)

		if w.Code != stdhttp.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestGetAuditLogs(t *testing.T) {
	mockStore := new(MockReadStore)
	handler := &Handler{
		ReadStore: mockStore,
	}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/instances/inst-1/audit", nil)
		req.SetPathValue("id", "inst-1")
		w := httptest.NewRecorder()

		events := []engine.AuditEvent{
			{ID: "evt-1", EventType: "CREATED", Timestamp: time.Now()},
		}
		mockStore.On("GetAuditEvents", mock.Anything, "inst-1").Return(events, nil)

		handler.HandleGetAuditLogs(w, req)

		if w.Code != stdhttp.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		_ = json.NewDecoder(w.Body).Decode(&resp)
		if _, ok := resp["events"]; !ok {
			t.Error("expected events in response")
		}
	})
}
