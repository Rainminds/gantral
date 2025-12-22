package engine

import (
	"fmt"
	"time"
)

// ErrInvalidTransition is returned when a state transition is not allowed.
type ErrInvalidTransition struct {
	From State
	To   State
}

func (e ErrInvalidTransition) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s", e.From, e.To)
}

// AllowedTransitions defines the strict map of valid state transitions.
// Based on specs/03-state-machine.md.
var AllowedTransitions = map[State][]State{
	StateCreated: {
		StateRunning,
	},
	StateRunning: {
		StateWaitingForHuman,
		StateCompleted,
		StateTerminated,
	},
	StateWaitingForHuman: {
		StateApproved,
		StateRejected,
		StateOverridden,
	},
	StateApproved: {
		StateResumed,
	},
	StateRejected: {
		StateTerminated, // Or remediation path, but strictly TERMINATED for now based on simple model
	},
	StateOverridden: {
		StateResumed,
	},
	StateResumed: {
		StateRunning,
	},
	StateCompleted:  {}, // Terminal state
	StateTerminated: {}, // Terminal state
}

// Transition attempts to move the instance to the target state.
// It returns an error if the transition is invalid.
func Transition(instance *Instance, target State) error {
	allowed, ok := AllowedTransitions[instance.State]
	if !ok {
		// If current state is not even in the map (e.g. unknown state), it's invalid
		return ErrInvalidTransition{From: instance.State, To: target}
	}

	isValid := false
	for _, s := range allowed {
		if s == target {
			isValid = true
			break
		}
	}

	if !isValid {
		return ErrInvalidTransition{From: instance.State, To: target}
	}

	instance.State = target
	instance.UpdatedAt = time.Now()
	return nil
}
