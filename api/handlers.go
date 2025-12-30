package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/errors"
	"github.com/Rainminds/gantral/core/policy"
)

// handleCreateInstance handles the creation of a new execution instance.
func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	// 1. Decode generic JSON body (for now just validating JSON structure)
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, errors.Wrap(errors.ErrInvalidInput, "invalid json body"))
		return
	}

	// 2. Call Engine
	// Construct policy based on input for testing
	pol := policy.Policy{ID: "default-policy"}
	if mat, ok := body["materiality"].(string); ok && mat == "HIGH" {
		pol.Materiality = policy.MaterialityHigh
	}

	instance, err := s.engine.CreateInstance(r.Context(), "default-wf", body, pol)
	if err != nil {
		s.writeError(w, err)
		return
	}

	// 3. Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(instance); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// handleGetInstance handles retrieving an execution instance by ID.
func (s *Server) handleGetInstance(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ID using Go 1.22 PathValue
	id := r.PathValue("id")
	if id == "" {
		s.writeError(w, errors.Wrap(errors.ErrInvalidInput, "missing instance id"))
		return
	}

	// 2. Call Engine
	instance, err := s.engine.GetInstance(r.Context(), id)
	if err != nil {
		s.writeError(w, errors.ErrNotFound)
		return
	}

	// 3. Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(instance); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// handleHealthz returns a simple 200 OK status.
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

// writeError maps domain errors to HTTP status codes and writes a JSON error response.
func (s *Server) writeError(w http.ResponseWriter, err error) {
	var status int
	switch {
	case errors.Is(err, errors.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, errors.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, errors.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, errors.ErrUnauthorized):
		status = http.StatusUnauthorized
	default:
		status = http.StatusInternalServerError
	}

	if status == http.StatusInternalServerError {
		slog.Error("internal server error", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	}); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

// HandleRecordDecision handles the recording of a human decision.
func (s *Server) handleRecordDecision(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")
	if instanceID == "" {
		s.writeError(w, errors.Wrap(errors.ErrInvalidInput, "instance ID is required"))
		return
	}

	var req struct {
		Type          string `json:"type"`
		ActorID       string `json:"actor_id"`
		Justification string `json:"justification"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, errors.Wrap(errors.ErrInvalidInput, "invalid request body"))
		return
	}

	cmd := engine.RecordDecisionCmd{
		InstanceID:    instanceID,
		Type:          engine.DecisionType(req.Type),
		ActorID:       req.ActorID,
		Justification: req.Justification,
	}

	inst, err := s.engine.RecordDecision(r.Context(), cmd)
	if err != nil {
		slog.Error("failed to record decision", "error", err)
		s.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(inst); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// handleListInstances handles retrieving all execution instances.
func (s *Server) handleListInstances(w http.ResponseWriter, r *http.Request) {
	instances, err := s.engine.ListInstances(r.Context())
	if err != nil {
		s.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(instances); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}
