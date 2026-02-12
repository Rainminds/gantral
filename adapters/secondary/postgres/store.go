package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/ports"
	"github.com/Rainminds/gantral/infra/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*db.Queries
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Store{
		Queries: db.New(pool),
		pool:    pool,
	}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// Ensure Store implements InstanceStore
var _ ports.InstanceStore = (*Store)(nil)

func (s *Store) CreateInstance(ctx context.Context, inst *engine.Instance) error {
	triggerBytes, _ := json.Marshal(inst.TriggerContext)
	policyBytes, _ := json.Marshal(inst.PolicyContext)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Correctly using the embedded promoted method directly
	qtx := s.WithTx(tx)

	// 1. Create Instance
	_, err = qtx.CreateInstance(ctx, db.CreateInstanceParams{
		ID:               inst.ID,
		WorkflowID:       inst.WorkflowID,
		State:            string(inst.State),
		TriggerContext:   triggerBytes,
		PolicyContext:    policyBytes,
		PolicyVersionID:  inst.PolicyVersionID,
		LastArtifactHash: inst.LastArtifactHash,
	})
	if err != nil {
		return err
	}

	// 2. Create Audit Event (INSTANCE_CREATED)
	eventPayload := map[string]interface{}{
		"workflow_id": inst.WorkflowID,
		"state":       inst.State,
	}
	payloadBytes, _ := json.Marshal(eventPayload)

	_, err = qtx.CreateAuditEvent(ctx, db.CreateAuditEventParams{
		ID:         fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		InstanceID: inst.ID,
		EventType:  "INSTANCE_CREATED",
		Payload:    payloadBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create audit event: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetInstance(ctx context.Context, id string) (*engine.Instance, error) {
	row, err := s.Queries.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting instance: %w", err)
	}
	return mapDBInstance(row), nil
}

func (s *Store) ListInstances(ctx context.Context) ([]*engine.Instance, error) {
	rows, err := s.Queries.ListInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing instances: %w", err)
	}

	result := make([]*engine.Instance, len(rows))
	for i, r := range rows {
		result[i] = mapDBInstance(r)
	}
	return result, nil
}

func (s *Store) RecordDecision(ctx context.Context, cmd engine.RecordDecisionCmd, nextState engine.State) (*engine.Instance, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Correctly using the embedded promoted method directly
	qtx := s.WithTx(tx)

	// 1. Create Decision Record
	decisionID := fmt.Sprintf("dec-%d", time.Now().UnixNano())
	snapshotBytes, _ := json.Marshal(cmd.ContextSnapshot)
	deltaBytes, _ := json.Marshal(cmd.ContextDelta)

	_, err = qtx.CreateDecision(ctx, db.CreateDecisionParams{
		ID:              decisionID,
		InstanceID:      cmd.InstanceID,
		Type:            string(cmd.Type),
		ActorID:         cmd.ActorID,
		Justification:   cmd.Justification,
		Role:            cmd.Role,
		ContextSnapshot: snapshotBytes,
		ContextDelta:    deltaBytes,
		PolicyVersionID: cmd.PolicyVersionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create decision: %w", err)
	}

	// 1.5 Fetch current state for audit log (to handle Overrides correct)
	currentInst, err := qtx.GetInstance(ctx, cmd.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current instance state: %w", err)
	}

	// 2. Update Instance State
	err = qtx.UpdateInstanceState(ctx, db.UpdateInstanceStateParams{
		ID:               cmd.InstanceID,
		State:            string(nextState),
		LastArtifactHash: cmd.NewArtifactHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update instance state: %w", err)
	}

	// 3. Create Audit Event (DECISION_RECORDED)
	eventPayload := map[string]interface{}{
		"decision_id":   decisionID,
		"decision_type": cmd.Type,
		"actor_id":      cmd.ActorID,
		"from_state":    currentInst.State,
		"to_state":      nextState,
	}
	payloadBytes, _ := json.Marshal(eventPayload)

	_, err = qtx.CreateAuditEvent(ctx, db.CreateAuditEventParams{
		ID:         fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		InstanceID: cmd.InstanceID,
		EventType:  "DECISION_RECORDED",
		Payload:    payloadBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit event: %w", err)
	}

	// 4. Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 5. Return updated instance
	return s.GetInstance(ctx, cmd.InstanceID)
}

func (s *Store) GetAuditEvents(ctx context.Context, instanceID string) ([]engine.AuditEvent, error) {
	// Implements ports.InstanceStore.GetAuditEvents using generated SQLC code.
	rows, err := s.Queries.GetAuditEvents(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("query audit events: %w", err)
	}

	var events []engine.AuditEvent
	for _, r := range rows {
		var payload map[string]interface{}
		_ = json.Unmarshal(r.Payload, &payload)

		events = append(events, engine.AuditEvent{
			ID:         r.ID,
			InstanceID: r.InstanceID,
			EventType:  r.EventType,
			Payload:    payload,
			Timestamp:  r.Timestamp.Time,
		})
	}
	return events, nil
}

func mapDBInstance(row db.Instance) *engine.Instance {
	var trigger map[string]interface{}
	var policy map[string]interface{}
	_ = json.Unmarshal(row.TriggerContext, &trigger)
	_ = json.Unmarshal(row.PolicyContext, &policy)

	return &engine.Instance{
		ID:               row.ID,
		WorkflowID:       row.WorkflowID,
		State:            engine.State(row.State),
		TriggerContext:   trigger,
		PolicyContext:    policy,
		PolicyVersionID:  row.PolicyVersionID,
		LastArtifactHash: row.LastArtifactHash,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}
}
