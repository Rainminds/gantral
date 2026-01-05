package ports

import (
	"context"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
)

// InstanceService defines the primary port for interacting with Instances.
// This is implemented by the Core (Engine).
type InstanceService interface {
	CreateInstance(ctx context.Context, workflowID string, triggerContext map[string]interface{}, pol policy.Policy) (*engine.Instance, error)
	GetInstance(ctx context.Context, id string) (*engine.Instance, error)
	ListInstances(ctx context.Context) ([]*engine.Instance, error)
	RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd) (*engine.Instance, error)
}
