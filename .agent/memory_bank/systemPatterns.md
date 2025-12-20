# System Architecture & Patterns

> For the normative state machine transitions and domain model, refer to `specs/03-state-machine.md`.


## Architectural Invariants (Non-Negotiable)
1.  **Instance-First:** Audit/cost/auth attaches to instances.
2.  **HITL as State:** Human intervention is an explicit state in the execution graph.
3.  **Human Authority is Final:** AI is advisory; policy-enforced human decisions override AI.
4.  **Determinism > Latency:** Immutable logs and replayability take precedence over speed.
5.  **Declarative Policies:** Rules are YAML/JSON, not code.
6.  **Adapters have NO Logic:** They only emit events and accept decisions.

## High-Level Architecture
- **Layer 1: Enterprise Systems:** (GitHub, Jira, Slack) -> Sources of events.
- **Layer 2: Gantral Control Plane:**
    - `Execution Engine`: Deterministic state machine.
    - `HITL State Machine`: Manages human interaction states.
    - `Policy Engine`: Evaluates materiality and access.
    - `Audit/Replay`: Immutable log storage.
- **Layer 3: Agent Frameworks:** (LangChain, Vellum) -> Executors of tasks.

## Core Domain Model
- **Workflow Template:** The definition (steps, triggers, policy refs).
- **Instance:** A concrete, immutable execution of a workflow (tied to a Team ID).
- **Decision:** A record of human intervention (Approve, Reject, Override).
- **Policy:** Declarative rules for timeouts, approvers, and escalation.

## State Machine States
`CREATED` -> `RUNNING` -> (`POLICY_BREACH`?) -> `WAITING_FOR_HUMAN` -> `APPROVED` | `REJECTED` | `OVERRIDDEN` -> `RESUMED` -> `COMPLETED` | `TERMINATED`
