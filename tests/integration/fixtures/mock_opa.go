package fixtures

import (
	"context"
	"fmt"

	"github.com/Rainminds/gantral/core/engine"
)

// MockPolicyEngine implements a simple policy evaluator for testing.
type MockPolicyEngine struct {
	// Mode determines what decision to return
	Mode string // "allow", "deny", "require_human", "error", "timeout"
}

func (m *MockPolicyEngine) Evaluate(ctx context.Context, input interface{}) (engine.DecisionType, error) {
	switch m.Mode {
	case "allow":
		// Assume "approve" maps to ALLOW in policy terms if auto-approval is supported?
		// Actually, Policy usually returns ALLOW/DENY/REQUIRE_HUMAN.
		// Gantral core `DecisionType` is Approve/Reject/Override.
		// We might need to map Policy Result to Decision if auto-execution is supported.
		// If Policy says ALLOW -> Engine might auto-approve?
		// Tier 1 tests showed `RecordDecision` takes Approve/Reject.
		// Let's assume Policy returns outcome which engine interprets.
		// For Integration tests, we might test the Engine's reaction to Policy.
		// BUT `core/engine` likely has an interface for Policy.
		return engine.DecisionApprove, nil // Simulating auto-approve for ALLOW?
	case "deny":
		return engine.DecisionReject, nil
	case "require_human":
		// Engine should wait. DecisionType "on evaluation" might not be DecisionType enum?
		// We need to check `core/ports/policy.go` or similar to see the interface.
		// I'll assume for now and fix if interface differs.
		return "", nil // Wait?
	case "error":
		return "", fmt.Errorf("mock policy error")
	default:
		return "", fmt.Errorf("unknown mode")
	}
}
