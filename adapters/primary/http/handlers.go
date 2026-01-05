package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"github.com/Rainminds/gantral/core/ports"
	"github.com/Rainminds/gantral/core/workflows"
	"github.com/google/uuid"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	TemporalClient client.Client
	TaskQueue      string
	ReadStore      ports.InstanceStore // CQRS Read Path
}

// CreateInstanceRequest defines the payload for creating an instance.
type CreateInstanceRequest struct {
	WorkflowID     string                 `json:"workflow_id"`
	TriggerContext map[string]interface{} `json:"trigger_context"`
	Policy         policy.Policy          `json:"policy"`
}

// CreateInstanceResponse defines the response.
type CreateInstanceResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// CreateInstance handles POST /instances.
// It starts a Temporal Workflow.
func (h *Handler) CreateInstance(w http.ResponseWriter, r *http.Request) {
	var req CreateInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Generate Instance ID (Execution ID)
	instanceID := fmt.Sprintf("inst-%s", uuid.New().String())

	workflowOptions := client.StartWorkflowOptions{
		ID:        instanceID,
		TaskQueue: h.TaskQueue,
	}

	input := workflows.WorkflowInput{
		WorkflowID:     req.WorkflowID,
		TriggerContext: req.TriggerContext,
		Policy:         req.Policy,
	}

	we, err := h.TemporalClient.ExecuteWorkflow(r.Context(), workflowOptions, workflows.GantralExecutionWorkflow, input)
	if err != nil {
		slog.Error("Failed to start workflow", "error", err)
		http.Error(w, "failed to start execution", http.StatusInternalServerError)
		return
	}

	slog.Info("Workflow Started", "instance_id", instanceID, "run_id", we.GetRunID())

	resp := CreateInstanceResponse{
		ID:     instanceID,
		Status: "PENDING",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/instances/%s", instanceID))
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}

// RecordDecisionRequest defines the payload for a human decision.
type RecordDecisionRequest struct {
	Type          string `json:"type"`
	ActorID       string `json:"actor_id"`
	Justification string `json:"justification"`
}

// RecordDecision handles POST /instances/{id}/decisions.
// It sends a Signal to the Temporal Workflow.
func (h *Handler) RecordDecision(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")
	if instanceID == "" {
		http.Error(w, "instance id required", http.StatusBadRequest)
		return
	}

	var req RecordDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Map to Signal Input
	var dType engine.DecisionType
	switch req.Type {
	case "APPROVE":
		dType = engine.DecisionApprove
	case "REJECT":
		dType = engine.DecisionReject
	case "OVERRIDE":
		dType = engine.DecisionOverride
	default:
		http.Error(w, "invalid decision type", http.StatusBadRequest)
		return
	}

	signalArg := activities.RecordDecisionInput{
		InstanceID:    instanceID,
		DecisionType:  dType,
		ActorID:       req.ActorID,
		Justification: req.Justification,
		Role:          "unknown_via_api",
		// PolicyVersionID, ContextSnapshot etc left empty for now as API doesn't provide them
	}

	err := h.TemporalClient.SignalWorkflow(r.Context(), instanceID, "", workflows.SignalHumanDecision, signalArg)
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			http.Error(w, "instance not found or completed", http.StatusNotFound)
			return
		}
		slog.Error("Failed to signal workflow", "error", err)
		http.Error(w, "failed to record decision", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "SIGNAL_SENT"})
}

// HandleGetAuditLogs retrieves audit logs for an instance.
func (h *Handler) HandleGetAuditLogs(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")
	if instanceID == "" {
		http.Error(w, "instance id required", http.StatusBadRequest)
		return
	}

	events, err := h.ReadStore.GetAuditEvents(r.Context(), instanceID)
	if err != nil {
		slog.Error("failed to get audit events", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
	})
}

// HandleGetInstance retrieves a single instance by ID.
func (h *Handler) HandleGetInstance(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")
	if instanceID == "" {
		http.Error(w, "instance id required", http.StatusBadRequest)
		return
	}

	inst, err := h.ReadStore.GetInstance(r.Context(), instanceID)
	if err != nil {
		slog.Error("failed to get instance", "error", err)
		http.Error(w, "instance not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inst)
}

// HandleListInstances retrieves all instances.
func (h *Handler) HandleListInstances(w http.ResponseWriter, r *http.Request) {
	instances, err := h.ReadStore.ListInstances(r.Context())
	if err != nil {
		slog.Error("failed to list instances", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"instances": instances,
	})
}

// HealthCheck handles GET /health.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// LoggingMiddleware wraps an http.Handler to log request details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter (simple version for brevity)
		next.ServeHTTP(w, r)

		slog.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}
