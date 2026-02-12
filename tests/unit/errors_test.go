package unit

import (
	"fmt"
	"testing"

	"github.com/Rainminds/gantral/core/engine"
)

// Test_Deterministic_Errors_SectionZ8 ensures error types are stable.
func Test_Deterministic_Errors_SectionZ8(t *testing.T) {
	t.Parallel()

	// 1. Invalid Transition Error Stability
	// ErrInvalidTransition is a struct, allowing type assertion
	err := engine.Transition(&engine.Instance{State: engine.StateCompleted}, engine.StateRunning)

	if err == nil {
		t.Fatal("Expected error")
	}

	// Check Type
	invalidErr, ok := err.(engine.ErrInvalidTransition)
	if !ok {
		t.Errorf("Expected ErrInvalidTransition type, got %T", err)
	}

	// Check Value Stability
	if invalidErr.From != engine.StateCompleted || invalidErr.To != engine.StateRunning {
		t.Errorf("Error content mismatch. Got From=%s To=%s", invalidErr.From, invalidErr.To)
	}

	// Check String Stability
	expectedMsg := fmt.Sprintf("invalid transition from %s to %s", engine.StateCompleted, engine.StateRunning)
	if err.Error() != expectedMsg {
		t.Errorf("Error message mismatch.\nExp: %s\nGot: %s", expectedMsg, err.Error())
	}
}
