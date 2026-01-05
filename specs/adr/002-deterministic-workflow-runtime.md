# 2. Use Temporal as Deterministic Workflow Runtime

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
Gantral must guarantee deterministic replay and auditability for all execution decisions. An auditor must be able to verify exactly why a transition occurred, even years later.

Additionally, Human-in-the-Loop (HITL) workflows require durable suspensions—the ability to pause execution for days or weeks without consuming compute resources or risking state loss during restarts.

Building a custom distributed event-sourcing engine to guarantee these properties is high-risk and undifferentiated heavy lifting.

## Decision
We will use Temporal as the underlying Workflow Substrate (Execution Runtime).

Gantral retains ownership of the Execution Authority Model (State Machine, Policy Gating, HITL Semantics), while delegating the mechanics of durability, retries, and history replay to Temporal.

## Consequences
*   **Positive:**
    *   **Guaranteed Audit:** We get "Replay-as-a-Service," allowing us to mathematically prove the sequence of events.
    *   **Resource Efficiency:** We can "hibernate" workflows during long human approvals (0 CPU usage) and wake them up reliably.
    *   **Safety:** Eliminates an entire class of distributed systems bugs (race conditions, lost timers).
*   **Negative:**
    *   **Dependency Weight:** Requires a Temporal Server (and its DB) in the deployment stack.
    *   **Strict Determinism:** Code within the workflow logic must be strictly deterministic (no random numbers, no direct API calls), which constrains contribution.
