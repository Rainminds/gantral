package engine

import (
	"time"

	"github.com/Rainminds/gantral/pkg/constants"
)

// State represents the current execution state of an Instance.
type State string

// Canonical States as defined in specs/03-state-machine.md
const (
	StateCreated         State = constants.StateCreated
	StateRunning         State = constants.StateRunning
	StateWaitingForHuman State = constants.StateWaitingForHuman
	StateApproved        State = constants.StateApproved
	StateRejected        State = constants.StateRejected
	StateOverridden      State = constants.StateOverridden
	StateResumed         State = constants.StateResumed
	StateCompleted       State = constants.StateCompleted
	StateTerminated      State = constants.StateTerminated
)

// Instance represents a concrete execution of a workflow.
type Instance struct {
	ID               string                 `json:"id"`
	WorkflowID       string                 `json:"workflow_id"`
	State            State                  `json:"state"`
	TriggerContext   map[string]interface{} `json:"trigger_context"`
	PolicyContext    map[string]interface{} `json:"policy_context"`
	PolicyVersionID  string                 `json:"policy_version_id"`
	LastArtifactHash string                 `json:"last_artifact_hash"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// AuditEvent represents an immutable record of a state change or decision.
type AuditEvent struct {
	ID         string                 `json:"id"`
	InstanceID string                 `json:"instance_id"`
	EventType  string                 `json:"event_type"`
	Payload    map[string]interface{} `json:"payload"`
	Timestamp  time.Time              `json:"timestamp"`
}
