ALTER TABLE instances ADD COLUMN policy_context JSONB DEFAULT '{}' NOT NULL;
