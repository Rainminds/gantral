-- name: CreateInstance :one
INSERT INTO instances (
    id,
    workflow_id,
    state,
    trigger_context,
    policy_context,
    policy_version_id,
    last_artifact_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetInstance :one
SELECT * FROM instances
WHERE id = $1 LIMIT 1;

-- name: UpdateInstanceState :exec
UPDATE instances
SET state = $2, last_artifact_hash = $3, updated_at = NOW()
WHERE id = $1;

-- name: CreateDecision :one
INSERT INTO decisions (
    id, instance_id, type, actor_id, justification, role, context_snapshot, context_delta, policy_version_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetDecisionsByInstance :many
SELECT * FROM decisions
WHERE instance_id = $1;

-- name: ListInstances :many
SELECT * FROM instances
ORDER BY created_at DESC;

-- name: CreateAuditEvent :one
INSERT INTO audit_events (
    id, instance_id, event_type, payload
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetAuditEvents :many
SELECT * FROM audit_events
WHERE instance_id = $1
ORDER BY timestamp ASC;
