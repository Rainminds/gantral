# Core Domain Model

## Workflow Template
The definition of a process. Contains steps, triggers, and references to policies.
- **ID:** Unique identifier.
- **Steps:** Sequence of actions.
- **Policy Refs:** Links to governance policies.

## ExecutionInstance
A concrete, immutable execution of a Workflow Template.
- **instance_id:** (UUID, immutable)
- **workflow_id:** String
- **workflow_version:** String
- **owning_team_id:** String
- **current_state:** String (Enum)
- **created_at:** Timestamp
- **terminated_at:** Timestamp

## AuthorityDecision
A record of human intervention.
- **decision_id:** Unique ID.
- **instance_id:** Link to execution.
- **decision_type:** `APPROVE`, `REJECT`, `OVERRIDE`.
- **human_actor_id:** Identity of the human actor.
- **role:** The role authorized to make the decision.
- **justification:** Reason for decision.
- **context_snapshot_hash:** Hash of what the human saw.
- **timestamp:** Exact time of decision.

## CommitmentArtifact
The cryptographic proof of a transition.
- **artifact_version:** Schema version.
- **artifact_id:** Unique ID.
- **instance_id:** Execution instance.
- **prev_artifact_hash:** Hash of the previous artifact (Merkle Chain).
- **authority_state:** State being transitioned to.
- **policy_version_id:** Policy version used.
- **context_hash:** Hash of execution context.
- **human_actor_id:** Signer identity.
- **timestamp:** Emission time.
- **artifact_hash:** Self-hash.

## Policy
Declarative rules governing execution.
- **Materiality:** Definition of "High Risk" vs "Low Risk".
- **Approvers:** Who is authorized to approve.
- **Timeouts:** How long to wait for human input.
