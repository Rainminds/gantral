package engine

import "time"

// State represents the current execution state of an Instance.
type State string

// Canonical States as defined in specs/03-state-machine.md
const (
	StateCreated         State = "CREATED"
	StateRunning         State = "RUNNING"
	StateWaitingForHuman State = "WAITING_FOR_HUMAN"
	StateApproved        State = "APPROVED"
	StateRejected        State = "REJECTED"
	StateOverridden      State = "OVERRIDDEN"
	StateResumed         State = "RESUMED"
	StateCompleted       State = "COMPLETED"
	StateTerminated      State = "TERMINATED"
)

// Instance represents a concrete execution of a workflow.
type Instance struct {
	ID        string
	TeamID    string
	State     State
	Context   map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}
