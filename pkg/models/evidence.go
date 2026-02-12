package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ExecutionEvidence represents the raw data of a tool execution.
// The Control Plane never sees this payload; it only sees the hash.
type ExecutionEvidence struct {
	EvidenceID    string          `json:"evidence_id"`
	InstanceID    string          `json:"instance_id"`
	ToolName      string          `json:"tool_name"`
	InputPayload  json.RawMessage `json:"input_payload"`
	OutputPayload json.RawMessage `json:"output_payload"`
	Timestamp     string          `json:"timestamp"`
}

// NewExecutionEvidence creates a new evidence container.
func NewExecutionEvidence(instanceID, toolName string, input, output []byte) *ExecutionEvidence {
	return &ExecutionEvidence{
		// EvidenceID will be calculated deterministically from content hash.
		InstanceID:    instanceID,
		ToolName:      toolName,
		InputPayload:  input,
		OutputPayload: output,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}
}

// ByteArray returns the canonical JSON representation of the evidence.
// This is what gets hashed.
func (e *ExecutionEvidence) CanonicalPayload() ([]byte, error) {
	if e.InstanceID == "" {
		return nil, errors.New("missing instance_id")
	}
	if e.ToolName == "" {
		return nil, errors.New("missing tool_name")
	}

	// We create a map to ensure consistent key ordering (standard lib guarantees sorted keys)
	// EvidenceID is EXCLUDED from the hash because it is derived FROM the hash.
	payload := map[string]interface{}{
		"instance_id":    e.InstanceID,
		"tool_name":      e.ToolName,
		"input_payload":  e.InputPayload,
		"output_payload": e.OutputPayload,
		"timestamp":      e.Timestamp,
	}

	return json.Marshal(payload)
}

// CalculateHashAndSetID computes the SHA-256 hash and sets it as the EvidenceID.
func (e *ExecutionEvidence) CalculateHashAndSetID() error {
	hash, err := e.CalculateHash()
	if err != nil {
		return err
	}
	e.EvidenceID = hash
	return nil
}

// CalculateHash returns the SHA-256 hash of the canonical payload.
func (e *ExecutionEvidence) CalculateHash() (string, error) {
	bytes, err := e.CanonicalPayload()
	if err != nil {
		return "", fmt.Errorf("failed to canonicalize evidence: %w", err)
	}

	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}
