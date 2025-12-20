# Core Domain Model

## Workflow Template
The definition of a process. Contains steps, triggers, and references to policies.
- **ID:** Unique identifier.
- **Steps:** Sequence of actions.
- **Policy Refs:** Links to governance policies.

## Instance
A concrete, immutable execution of a Workflow Template.
- **ID:** Unique execution ID (UUID).
- **TeamID:** Tenant/Team ownership.
- **State:** Current state in the state machine.
- **Context:** Snapshot of data at creation (Trigger context).
- **CostMetadata:** Token usage and cost tracking.
- **AuditLog:** Complete history of events for this instance.

## Decision
A record of human intervention.
- **Type:** `APPROVE`, `REJECT`, `OVERRIDE`.
- **User (ActorID):** Identity of the human actor.
- **Role:** The role authorized to make the decision.
- **Justification:** Human notes or reason for override.
- **ContextSnapshot:** What the human saw when deciding.
- **Timestamp:** Exact time of decision.

## Policy
Declarative rules governing execution.
- **Materiality:** Definition of "High Risk" vs "Low Risk".
- **Approvers:** Who is authorized to approve.
- **Timeouts:** How long to wait for human input.
