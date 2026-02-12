package policy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/pkg/constants"
)

// Engine is responsible for evaluating policies against execution requests.
type Engine struct{}

// NewEngine creates a new instance of the Policy Engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Evaluate checks the policy rules and determines the next execution state.
func (e *Engine) Evaluate(ctx context.Context, p Policy) (EvaluationResult, error) {
	// Delegate to pure function for core logic
	result := EvaluatePure(p)

	// Observability (Side Effect)
	slog.Info("policy evaluated",
		"policy_id", p.ID,
		"decision", result.NextState,
		"reason", result.Reason,
	)

	return result, nil
}

// EvaluatePure contains the deterministic logic for policy evaluation.
// It must have NO side effects (no logging, no I/O) to be safe for Replay.
func EvaluatePure(p Policy) EvaluationResult {
	result := EvaluationResult{
		ShouldPause: false,
		NextState:   constants.StateRunning,
		Reason:      "Policy allows automatic execution",
	}

	// Rule: HIGH materiality OR explicit human approval requirement pauses execution.
	if p.Materiality == MaterialityHigh || p.RequiresHumanApproval {
		result.ShouldPause = true
		result.NextState = constants.StateWaitingForHuman
		result.Reason = fmt.Sprintf("Execution paused: Materiality=%s, RequiresApproval=%v", p.Materiality, p.RequiresHumanApproval)
	}

	return result
}
