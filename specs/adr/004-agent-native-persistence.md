# 4. Agent-Native Persistence for Long-Running Pauses

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
Human-in-the-Loop (HITL) approvals often have **undefined latency** (minutes to weeks).

Keeping an Agent process (container/pod) alive and blocked on I/O during this wait is financially wasteful and operationally brittle. If the underlying node patches or restarts, the in-memory state of the Agent is lost.

## Decision
We will enforce a **Stateless / Resume-able Execution Contract**.

Gantral **will not** serialize or store the internal memory of the Agent (e.g., LLM context window). Instead, we rely on the **Agent Framework's Native Persistence**.

*   **On Pause:** The Agent process MUST persist its internal state (via framework checkpoints) and **exit**.
*   **On Resume:** Gantral triggers a **new** process. The Agent MUST reload its state from its own persistence layer and continue execution.

## Consequences
*   **Positive:**
    *   **Zero-Cost Waiting:** No compute resources are consumed during approval waits.
    *   **Framework Agnostic:** We do not need to build custom serializers for every Agent framework (CrewAI, LangChain, etc.).
    *   **Resilience:** Agents survive system reboots implicitly.
*   **Negative:**
    *   **Developer Constraint:** Developers cannot write simple scripts; they must use frameworks that support checkpointing (or manually split their logic).
