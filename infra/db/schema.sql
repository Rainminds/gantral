CREATE TABLE instances (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL,
    state TEXT NOT NULL,
    trigger_context JSONB NOT NULL,
    policy_context JSONB DEFAULT '{}' NOT NULL,
    policy_version_id TEXT NOT NULL DEFAULT '',
    last_artifact_hash TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE decisions (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES instances(id),
    type TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    justification TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT '',
    context_snapshot JSONB NOT NULL DEFAULT '{}',
    context_delta JSONB NOT NULL DEFAULT '{}',
    policy_version_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_events (
    id TEXT PRIMARY KEY,
    instance_id TEXT NOT NULL REFERENCES instances(id),
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
