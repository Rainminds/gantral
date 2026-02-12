package workflow

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Rainminds/gantral/internal/replay"
	"github.com/Rainminds/gantral/pkg/models"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

// ReplayInterceptor enforces that replayed activities have valid artifacts.
// It matches the specific requirement: "If the runtime proposes a history... execution must Fail-Closed."
type ReplayInterceptor struct {
	interceptor.WorkerInterceptorBase
	Guard *replay.ReplayGuard
}

// NewReplayInterceptor creates the interceptor.
func NewReplayInterceptor(guard *replay.ReplayGuard) *ReplayInterceptor {
	return &ReplayInterceptor{Guard: guard}
}

// InterceptWorkflow wraps the workflow execution to hook into outbound calls (Activities).
func (r *ReplayInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	return &replayWorkflowInbound{
		WorkflowInboundInterceptorBase: interceptor.WorkflowInboundInterceptorBase{Next: next},
		guard:                          r.Guard,
	}
}

// replayWorkflowInbound intercepts inbound calls to init outbound interceptors.
type replayWorkflowInbound struct {
	interceptor.WorkflowInboundInterceptorBase
	guard *replay.ReplayGuard
}

func (w *replayWorkflowInbound) Init(outbound interceptor.WorkflowOutboundInterceptor) error {
	i := &replayWorkflowOutbound{
		WorkflowOutboundInterceptorBase: interceptor.WorkflowOutboundInterceptorBase{Next: outbound},
		guard:                           w.guard,
	}
	return w.Next.Init(i)
}

// replayWorkflowOutbound intercepts outbound calls (ExecuteActivity).
type replayWorkflowOutbound struct {
	interceptor.WorkflowOutboundInterceptorBase
	guard *replay.ReplayGuard
}

// ExecuteActivity intercepts activity calls to valid results on replay.
func (w *replayWorkflowOutbound) ExecuteActivity(
	ctx workflow.Context,
	activityType string,
	args ...interface{},
) workflow.Future {
	// Call the next interceptor (or SDK core)
	f := w.Next.ExecuteActivity(ctx, activityType, args...)

	// Wrap the future to validate the result upon completion
	return &replayFuture{
		Future: f,
		ctx:    ctx,
		guard:  w.guard,
		name:   activityType,
	}
}

// replayFuture wraps workflow.Future to inspect results.
type replayFuture struct {
	workflow.Future
	ctx   workflow.Context
	guard *replay.ReplayGuard
	name  string
}

// Get intercepts the result retrieval.
func (f *replayFuture) Get(ctx workflow.Context, valuePtr interface{}) error {
	// 1. Get the result from the SDK (which fetches from History on replay)
	err := f.Future.Get(ctx, valuePtr)
	if err != nil {
		return err // Activity failed, nothing to validate (or we assume failure is safe)
	}

	// 2. Continuous Validation (Defense in Depth)
	// We Validate Artifacts on BOTH Live Execution and Replay.
	// - On Replay: Prevents History Injection (Security Requirement).
	// - On Live: Catches Hash/State bugs immediately (Fail Fast).

	// 3. Inspect if the result is an Artifact.
	// valuePtr is a pointer to the destination (e.g. **CommitmentArtifact or *interface{}).
	// We check if valuePtr points to a CommitmentArtifact.
	// Since we know the specific activity types or can check type, we try type assertion.

	// Use a helper to extract the artifact safely.
	if artifact, ok := extractArtifact(valuePtr); ok {
		slog.Debug("Validating replayed artifact", "activity", f.name, "id", artifact.ArtifactID)
		// Use a detached context for the store lookup, as workflow.Context is not compatible.
		if vErr := f.guard.ValidateReplay(context.Background(), artifact); vErr != nil {
			// FAIL-CLOSED: Panic or Token Cancellation.
			// Panicking inside a workflow on replay is the standard way to reject deterministic violation.
			// It crashes the replay, preventing the worker from proceeding on this bad history.
			panic(fmt.Sprintf("REPLAY GUARD VIOLATION: %v", vErr))
		}
	}

	return nil
}

// extractArtifact attempts to cast the valuePtr to a *models.CommitmentArtifact.
func extractArtifact(valuePtr interface{}) (*models.CommitmentArtifact, bool) {
	if valuePtr == nil {
		return nil, false
	}

	// 1. Direct Typed Pointer: **CommitmentArtifact
	if ptr, ok := valuePtr.(**models.CommitmentArtifact); ok && ptr != nil && *ptr != nil {
		return *ptr, true
	}

	// 2. Interface Wrapper: *interface{} -> *CommitmentArtifact OR map[string]interface{}
	if ptr, ok := valuePtr.(*interface{}); ok && ptr != nil {
		val := *ptr
		if val == nil {
			return nil, false
		}

		// 2a. Already a struct (unlikely from history unless custom converter)
		if art, ok := val.(*models.CommitmentArtifact); ok {
			return art, true
		}

		// 2b. Map (JSON unmarshalled)
		if m, ok := val.(map[string]interface{}); ok {
			// Check for required fields to identify it as an Artifact
			id, hasID := m["artifact_id"].(string)
			state, hasState := m["authority_state"].(string)
			// hash field is removed in v1, ID is the hash.

			if hasID && hasState {
				// Reconstruct for validation
				// extracting other fields for integrity hash check
				// extracting other fields for integrity hash check
				instID, _ := m["instance_id"].(string)
				prevHash, _ := m["prev_artifact_hash"].(string)
				policyVer, _ := m["policy_version_id"].(string)
				ctxHash, _ := m["context_hash"].(string)
				actorID, _ := m["human_actor_id"].(string)
				timestamp, _ := m["timestamp"].(string)
				version, _ := m["artifact_version"].(string)

				return &models.CommitmentArtifact{
					ArtifactID:       id,
					AuthorityState:   state,
					InstanceID:       instID,
					PrevArtifactHash: prevHash,
					PolicyVersionID:  policyVer,
					ContextHash:      ctxHash,
					HumanActorID:     actorID,
					Timestamp:        timestamp,
					ArtifactVersion:  version,
				}, true
			}
		}
	}
	return nil, false
}
