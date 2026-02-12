package policy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	corepolicy "github.com/Rainminds/gantral/core/policy"
	"github.com/Rainminds/gantral/pkg/constants"
)

var (
	// ErrPolicyUnreachable indicates the policy engine could not be contacted.
	ErrPolicyUnreachable = errors.New("policy engine unreachable")
)

// Evaluator defines the policy evaluation interface.
type Evaluator interface {
	Evaluate(ctx context.Context, p corepolicy.Policy) (corepolicy.EvaluationResult, error)
}

// FailClosedEngine wraps a Policy Engine to enforce strict fail-closed behavior.
// If the underlying engine fails or returns an ambiguous result, this wrapper
// ensures the decision is "WAITING_FOR_HUMAN" (safe state), never "RUNNING".
type FailClosedEngine struct {
	inner Evaluator
}

// NewFailClosedEngine creates a new safety wrapper.
func NewFailClosedEngine(inner Evaluator) *FailClosedEngine {
	return &FailClosedEngine{inner: inner}
}

// Evaluate acts as a firewall for policy decisions.
func (e *FailClosedEngine) Evaluate(ctx context.Context, p corepolicy.Policy) (corepolicy.EvaluationResult, error) {
	// 1. Delegated Evaluation
	result, err := e.inner.Evaluate(ctx, p)

	// 2. Fail-Closed Guard (Error Case)
	if err != nil {
		slog.Error("SECURITY ALERT: Policy evaluation failed. Enforcing fail-closed.",
			"error", err,
			"policy_id", p.ID)

		return corepolicy.EvaluationResult{
			ShouldPause: true,
			NextState:   constants.StateWaitingForHuman, // Fail Safe
			Reason:      fmt.Sprintf("Fail-Closed: Policy Error (%v)", err),
		}, nil
		// Swallowing error allows the workflow to "park" in a SAFE state
		// rather than crash-looping. This is intentional Fail-Safe design.
	}

	// 3. Fail-Closed Guard (Invalid/Ambiguous Result)
	if result.NextState == "" {
		slog.Error("SECURITY ALERT: Policy returned empty state. Enforcing fail-closed.", "policy_id", p.ID)
		return corepolicy.EvaluationResult{
			ShouldPause: true,
			NextState:   constants.StateWaitingForHuman,
			Reason:      "Fail-Closed: Ambiguous Result",
		}, nil
	}

	// 4. Pass-through valid result
	return result, nil
}
