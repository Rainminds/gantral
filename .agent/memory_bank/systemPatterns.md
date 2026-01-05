# System Architecture & Patterns

> For the normative state machine transitions and domain model, refer to `core/engine/types.go` and `specs/03-state-machine.md`.

## Architectural Invariants (Non-Negotiable)
1.  **Instance-First:** Audit/cost/auth attaches to instances (single executions).
2.  **HITL as State:** Human intervention is an explicit state in the execution graph.
3.  **Human Authority is Final:** AI is advisory; policy-enforced human decisions override AI.
4.  **Determinism > Latency:** Immutable logs and replayability take precedence over speed.
5.  **Federated Execution:** Gantral never executes user code. Runners pull tasks.
6.  **Stateless Waiting:** Processes MUST exit during long human pauses (resume from checkpoint).

## High-Level Architecture
- **Layer 1: Enterprise Systems:** (IdP, Policy Repo) -> Sources of Truth.
- **Layer 2: Gantral Control Plane:**
    - `API Server`: HTTP Gatekeeper & Policy Guard.
    - `Postgres`: Event Store (Instances, Decisions, Audit).
- **Layer 3: Execution Runtime:**
    - `Temporal Cluster`: Manages timers, retries, and queues.
- **Layer 4: Federated Runners:**
    - `Worker`: Go process polling Temporal. Executes logic.
    - `Agent`: Python/Node process started by Worker.

## Core Domain Model
- **Workflow Template:** The definition (steps, triggers, policy refs).
- **Instance:** A concrete, immutable execution of a workflow.
- **Decision:** A record of human intervention (Approve, Reject, Override).
- **Policy:** Declarative rules acting as Transition Guards.
- **Runner:** The trusted agent executing code in User Space.

## State Machine States (Canonical)
`CREATED` -> `RUNNING` -> (`POLICY_CHECK`) -> `WAITING_FOR_HUMAN` -> `APPROVED` | `REJECTED` | `OVERRIDDEN` -> `RESUMED` -> `COMPLETED` | `TERMINATED`
