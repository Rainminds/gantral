package api

import (
	"encoding/json"
	"net/http"
)

// handleCreateInstance handles the creation of a new execution instance.
func (s *Server) handleCreateInstance(w http.ResponseWriter, r *http.Request) {
	// 1. Decode generic JSON body (for now just validating JSON structure)
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(response)
}

// handleGetInstance handles retrieving an execution instance by ID.
func (s *Server) handleGetInstance(w http.ResponseWriter, r *http.Request) {
	// 1. Extract ID using Go 1.22 PathValue
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing instance ID", http.StatusBadRequest)
		return
	}

	// 2. Call Engine (Stubbed)
	// In future: instance, err := s.engine.GetInstance(id)

	// 3. Response
	response := map[string]interface{}{
		"instance_id": id,
		"state":       "RUNNING", // stub state
		"created_at":  "2023-10-27T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
