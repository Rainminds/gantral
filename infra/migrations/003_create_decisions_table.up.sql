CREATE TABLE decisions (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES instances(id),
    type TEXT NOT NULL, -- ENUM: APPROVE, REJECT, OVERRIDE
    actor_id TEXT NOT NULL,
    justification TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
