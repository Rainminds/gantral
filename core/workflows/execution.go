package workflows

import (
	"time"

	"github.com/Rainminds/gantral/core/activities"
	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	// SignalHumanDecision is the signal name for HITL decisions.
	SignalHumanDecision = "HumanDecision"

	// TaskQueue is the default task queue for Gantral.
	TaskQueue = "gantral-core"
)

// WorkflowInput defines strict inputs for the execution workflow.
type WorkflowInput struct {
	WorkflowID     string
	TriggerContext map[string]interface{}
	Policy         policy.Policy
}

// WorkflowResult defines the output of the execution workflow.
type WorkflowResult struct {
	InstanceID string
	FinalState engine.State
}

// GantralExecutionWorkflow orchestrates the lifecycle of a Gantral Instance.
// It is deterministic and handles: Creation -> Policy Eval -> HITL -> Completion.
func GantralExecutionWorkflow(ctx workflow.Context, input WorkflowInput) (WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting Gantral Execution Workflow", "workflow_id", input.WorkflowID)

	// A. Setup Activities
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    100 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// B. Policy Evaluation (Deterministic Logic)
	// We evaluate purely in code since policy dictionary is passed in input.
	// If policy engine required external calls, we'd use Activity.
	// Replicating logic from core/policy/engine.go here for determinism or calling a shared deterministic function.

	// Inline logic to ensure safety:
	shouldPause := false
	nextState := engine.StateRunning
	reason := "Policy allows automatic execution"

	if input.Policy.Materiality == policy.MaterialityHigh || input.Policy.RequiresHumanApproval {
		shouldPause = true
		nextState = engine.StateWaitingForHuman
		reason = "Execution paused by policy"
	}

	policyResult := map[string]interface{}{
		"should_pause": shouldPause,
		"reason":       reason,
		"policy_id":    input.Policy.ID,
	}

	// C. Persist Instance (Create)
	var inst *engine.Instance
	persistInput := activities.PersistInstanceInput{
		InstanceID:      workflow.GetInfo(ctx).WorkflowExecution.ID,
		WorkflowID:      input.WorkflowID,
		TriggerContext:  input.TriggerContext,
		Policy:          nil,             // Metadata
		PolicyVersionID: input.Policy.ID, // Assuming ID implies version for now
		InitialState:    nextState,
		PolicyResult:    policyResult,
	}

	var a *activities.ExecutionActivities // nil struct for name resolution
	err := workflow.ExecuteActivity(ctx, a.PersistInstance, persistInput).Get(ctx, &inst)
	if err != nil {
		logger.Error("Failed to persist instance", "error", err)
		return WorkflowResult{}, err
	}

	// D. HITL Loop
	if shouldPause {
		logger.Info("Blocking for Human Decision", "instance_id", inst.ID)

		// Defaults
		const ApprovalTimeout = 24 * time.Hour
		var decisionInput activities.RecordDecisionInput

		// Setup Selector
		msg := "HITL Decision Received"
		selector := workflow.NewSelector(ctx)
		signalChan := workflow.GetSignalChannel(ctx, SignalHumanDecision)
		timerFuture := workflow.NewTimer(ctx, ApprovalTimeout)

		// 1. Handle Signal
		selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &decisionInput)
		})

		// 2. Handle Timeout
		selector.AddFuture(timerFuture, func(f workflow.Future) {
			msg = "HITL Timeout Exceeded"
			// Construct System Rejection
			decisionInput = activities.RecordDecisionInput{
				InstanceID:    inst.ID,
				DecisionType:  engine.DecisionReject,
				ActorID:       "SYSTEM",
				Justification: "Approval Timeout (24h) Exceeded",
				Role:          "SYSTEM",
			}
		})

		// Wait for one
		selector.Select(ctx)
		logger.Info(msg, "instance_id", inst.ID)

		// Validate (if strictly from signal, if timeout we constructed valid input)
		if decisionInput.InstanceID != inst.ID && decisionInput.ActorID != "SYSTEM" {
			logger.Warn("Received signal for wrong instance", "expected", inst.ID, "got", decisionInput.InstanceID)
			// For robustness, we could loop back but assuming happy path or timeout for now.
		}

		// Ensure InstanceID is set for Timeout case if it wasn't
		if decisionInput.InstanceID == "" {
			decisionInput.InstanceID = inst.ID
		}

		// Record Decision via Activity
		err := workflow.ExecuteActivity(ctx, a.RecordDecision, decisionInput).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to record decision", "error", err)
			return WorkflowResult{}, err
		}

		// Update local state based on decision
		switch decisionInput.DecisionType {
		case engine.DecisionApprove, engine.DecisionOverride:
			inst.State = engine.StateApproved
		case engine.DecisionReject:
			inst.State = engine.StateRejected
		}
	}

	logger.Info("Workflow Completed", "final_state", inst.State)
	return WorkflowResult{
		InstanceID: inst.ID,
		FinalState: inst.State,
	}, nil
}
