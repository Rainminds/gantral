package evidence

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Rainminds/gantral/pkg/models"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/interceptor"
)

// EvidenceInterceptor captures tool inputs/outputs and produces a hash.
type EvidenceInterceptor struct {
	interceptor.WorkerInterceptorBase
}

// NewEvidenceInterceptor creates a new interceptor.
func NewEvidenceInterceptor() *EvidenceInterceptor {
	return &EvidenceInterceptor{}
}

// InterceptActivity wraps activity execution to capture evidence.
func (e *EvidenceInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	return &evidenceActivityInboundInterceptor{root: next}
}

type evidenceActivityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
	root interceptor.ActivityInboundInterceptor
}

// ExecuteActivity intercepts the activity execution.
func (a *evidenceActivityInboundInterceptor) ExecuteActivity(
	ctx context.Context,
	in *interceptor.ExecuteActivityInput,
) (interface{}, error) {
	// 1. Capture Input
	var inputBytes []byte
	if len(in.Args) > 0 {
		var err error
		inputBytes, err = json.Marshal(in.Args[0])
		if err != nil {
			slog.Warn("Failed to marshal activity input for evidence", "error", err)
			inputBytes = []byte("{}")
		}
	}

	// 2. Execute Activity (The Tool)
	result, err := a.root.ExecuteActivity(ctx, in)

	// 3. Capture Output
	var outputBytes []byte
	if err == nil {
		var marshalErr error
		outputBytes, marshalErr = json.Marshal(result)
		if marshalErr != nil {
			slog.Warn("Failed to marshal activity output for evidence", "error", marshalErr)
			outputBytes = []byte("{}")
		}
	} else {
		outputBytes = []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	// 4. Create Evidence Object
	info := activity.GetInfo(ctx)
	toolName := info.ActivityType.Name
	if toolName == "" {
		toolName = "tool-execution"
	}

	// Use WorkflowID as InstanceID for correlation
	instanceID := info.WorkflowExecution.ID
	if instanceID == "" {
		instanceID = "unknown-workflow-id"
	}

	ev := models.NewExecutionEvidence(
		instanceID,
		toolName,
		inputBytes,
		outputBytes,
	)

	// 5. Hash Locally
	// We use the model to calculate the ID (hash) deterministically
	if err := ev.CalculateHashAndSetID(); err != nil {
		slog.Error("CRITICAL: Failed to calculate evidence hash", "error", err)
		return nil, fmt.Errorf("evidence calculation failed: %w", err)
	}
	hash := ev.EvidenceID

	// 6. Generic Upload (S3/File/Mem)
	bucketURL := os.Getenv("EVIDENCE_BUCKET_URL")
	if bucketURL == "" {
		// Default to local temp dir if not set (safe fallback for dev)
		tempDir := os.TempDir()
		bucketURL = fmt.Sprintf("file://%s", tempDir)
	}

	bucket, err := blob.OpenBucket(ctx, bucketURL)
	// Fail-Closed: We MUST persist evidence before returning.
	// If we cannot prove what happened, we treat the execution as failed.
	if err != nil {
		slog.Error("Failed to open evidence bucket", "url", bucketURL, "error", err)
		return nil, fmt.Errorf("security violation: failed to open evidence storage: %w", err)
	}
	defer bucket.Close()

	evidenceJSON, _ := json.Marshal(ev) // Canonical payload is internal, this is for storage

	// Key: evidence/{hash}.json
	key := fmt.Sprintf("evidence/%s.json", hash)

	if err := bucket.WriteAll(ctx, key, evidenceJSON, nil); err != nil {
		slog.Error("Failed to upload evidence", "key", key, "error", err)
		return nil, fmt.Errorf("failed to persist evidence: %w", err)
	}

	slog.Info("Evidence Persisted", "hash", hash, "storage", bucketURL, "key", key)

	return result, err
}
