package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Rainminds/gantral/core/errors"
)

// handleCreateInstance handles the creation of a new execution instance.
func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	// 1. Decode generic JSON body (for now just validating JSON structure)
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, errors.Wrap(errors.ErrInvalidInput, "invalid json body"))
		return
	}

	// 2. Call Engine (Stubbed for now)
	// In future: instance, err := s.engine.CreateInstance(...)

	// 3. Response
	response := map[string]string{
		"status":      "created",
		"instance_id": "inst-12345-stub",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
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

	// 2. Call Engine (Stubbed)
	// In future: instance, err := s.engine.GetInstance(id)

	// Stub logic for 404
	if id == "notfound" {
		s.writeError(w, errors.ErrNotFound)
		return
	}

	// 3. Response
	response := map[string]interface{}{
		"instance_id": id,
		"state":       "RUNNING", // stub state
		"created_at":  "2023-10-27T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// handleHealthz returns a simple 200 OK status.
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
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
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
