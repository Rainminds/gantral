# Architecture & Invariants

## Scope
Gantral is an **Execution Control Plane** that sits between enterprise systems and AI agent frameworks. It strictly manages governance, auditability, and Human-in-the-Loop flows.

## Architectural Invariants
1.  **Instance-First:** Audit, cost, and auth context attach to specific execution instances, not abstract definitions.
2.  **HITL as State:** Human intervention is an explicit, first-class state (`WAITING_FOR_HUMAN`) in the execution graph, not a side-effect.
3.  **Human Authority is Final:** Policy-enforced human decisions always override AI outputs.
4.  **Determinism > Latency:** Immutable logs and deterministic replayability take precedence over execution speed.
5.  **Declarative Policies:** Governance rules must be defined in data (YAML/JSON), not code.
6.  **Adapters have NO Logic:** Integration adapters only emit events and accept decisions; they contain no business logic.

## High-Level Architecture
### Layer 1: Enterprise Systems (Sources)
- GitHub, Jira, Slack, PagerDuty.
- Sources of events that trigger workflows.

### Layer 2: Gantral Control Plane
- **Execution Engine:** Deterministic state machine managing flow.
- **HITL State Machine:** Manages human interaction states & approvals.
- **Policy Engine:** Evaluates materiality and access control.
- **Audit/Replay:** Immutable log storage for all actions.

### Layer 3: Agent Frameworks (Executors)
- LangChain, Vellum, CrewAI.
- External systems that perform the actual AI tasks.
