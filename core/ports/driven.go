package ports

import (
	"context"

	"github.com/Rainminds/gantral/core/engine"
)

// InstanceStore defines the secondary port for persistence.
// This is implemented by Adapters (Postgres, Memory).
type InstanceStore interface {
	CreateInstance(ctx context.Context, inst *engine.Instance) error
	GetInstance(ctx context.Context, id string) (*engine.Instance, error)
	// ListInstances retrieves all instances.
	ListInstances(ctx context.Context) ([]*engine.Instance, error)

	// GetAuditEvents retrieves the immutable event log for an instance.
	GetAuditEvents(ctx context.Context, instanceID string) ([]engine.AuditEvent, error)
	RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error)
}
