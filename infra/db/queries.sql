-- name: CreateInstance :one
INSERT INTO instances (
    id,
    workflow_id,
    state,
    trigger_context,
    policy_context
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetInstance :one
SELECT * FROM instances
WHERE id = $1 LIMIT 1;

UPDATE instances
SET state = $2, updated_at = NOW()
WHERE id = $1;

-- name: CreateDecision :one
INSERT INTO decisions (
    id, instance_id, type, actor_id, justification
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetDecisionsByInstance :many
SELECT * FROM decisions
WHERE instance_id = $1
ORDER BY created_at DESC;
