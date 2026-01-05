# 7. Policy Evaluation as Transition Guard

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
We need to enforce governance rules (e.g., "Cost > $50 requires VP approval") dynamically. However, encoding these rules as hard-coded "Execution States" creates a combinatorial explosion of states and makes the system rigid.

## Decision
We will implement Policy Evaluation as a **Synchronous Transition Guard**.

*   **Guard Logic:** Before any state transition (e.g., `RUNNING` -> `COMPLETED`), the engine invokes the Policy Evaluator.
*   **Outcome:** The policy returns a decision (`ALLOW`, `DENY`, `REQUIRE_HUMAN`).
*   **State Impact:**
    *   `ALLOW`: The transition proceeds.
    *   `REQUIRE_HUMAN`: The transition is diverted to `WAITING_FOR_HUMAN`.
    *   `DENY`: The transition is diverted to `TERMINATED`.

## Consequences
*   **Positive:**
    *   **Canonical State Machine:** The core state machine remains simple and fixed, regardless of policy complexity.
    *   **Hot Swapping:** Policies can be updated without redeploying the application code.
    *   **Separation of Duties:** Policy logic (Rego) is completely decoupled from Execution logic (Go/Temporal).
*   **Negative:**
    *   **Latency:** Every transition incurs a policy evaluation overhead (mitigated by caching and local evaluation).
