package policy

import (
	"context"
	"fmt"
	"log/slog"
)

// Engine is responsible for evaluating policies against execution requests.
type Engine struct{}

// NewEngine creates a new instance of the Policy Engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Evaluate checks the policy rules and determines the next execution state.
func (e *Engine) Evaluate(ctx context.Context, p Policy) (EvaluationResult, error) {
	result := EvaluationResult{
		ShouldPause: false,
		NextState:   "RUNNING",
		Reason:      "Policy allows automatic execution",
	}

	// Rule: HIGH materiality OR explicit human approval requirement pauses execution.
	if p.Materiality == MaterialityHigh || p.RequiresHumanApproval {
		result.ShouldPause = true
		result.NextState = "WAITING_FOR_HUMAN"
		result.Reason = fmt.Sprintf("Execution paused: Materiality=%s, RequiresApproval=%v", p.Materiality, p.RequiresHumanApproval)
	}

	// Observability
	slog.Info("policy evaluated",
		"policy_id", p.ID,
		"decision", result.NextState,
		"reason", result.Reason,
	)

	return result, nil
}
