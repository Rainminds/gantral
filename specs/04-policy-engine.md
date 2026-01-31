# Policy Engine

Gantral uses a policy engine to evaluate "Materiality" and enforcement rules. Policies are declarative content.

## Schema Example

## Interface Contract (Code-Level)

### Input
- `instance_id`
- `workflow_id`, `workflow_version`
- `materiality`
- `current_state`
- `actor_id`, `roles`
- `policy_version_id`
- `dry_run`

### Output (Advisory Only)
- `decision` (ALLOW / REQUIRE_HUMAN / DENY)
- `approver_roles`
- `timeout`
- `escalation_roles`

## Concepts
- **Materiality:** The assessment of risk. High/Low/Critical.
- **Rules:** Condition -> Action mappings.
- **Approvers:** Roles or users required to sign off.

## Integrity & Hashing
To ensure adversarial auditability:
- **Policy Input (Context) Hash:** The exact JSON input provided to the policy engine MUST be hashed (SHA-256).
- **Policy Output (Decision) Hash:** The exact decision returned MUST be hashed.
- **Binding:** These hashes MUST be included in the `CommitmentArtifact` to prove *why* a transition occurred.
