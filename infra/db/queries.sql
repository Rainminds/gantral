-- name: CreateInstance :exec
INSERT INTO instances (
    id,
    workflow_id,
    state,
    trigger_context
) VALUES (
    $1, $2, $3, $4
);

-- name: GetInstance :one
SELECT * FROM instances
WHERE id = $1 LIMIT 1;

-- name: UpdateState :exec
UPDATE instances
SET state = $2, updated_at = NOW()
WHERE id = $1;
