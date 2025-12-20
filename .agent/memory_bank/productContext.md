# Product Context: Gantral & Gantrio

## The Core Problem
Enterprises adopt AI tools (agents, copilots) but lack a unified "AI Execution Control Plane" to manage governance, auditability, and Human-in-the-Loop (HITL) workflows at scale. Gantral solves the question: "Who authorized this AI action?"

## What Gantral Is (Open Core)
- **Infrastructure, not an Agent Builder:** It sits *above* frameworks like LangChain/CrewAI and *below* enterprise processes.
- **Execution Control Plane:** Standardizes how workflows execute, pause for humans, and resume.
- **System of Record:** Provides immutable audit logs for every AI decision.

## Key Differentiation
- **HITL is a State, not a Feature:** Human intervention is modeled as a first-class state transition (`WAITING_FOR_HUMAN` -> `APPROVED/REJECTED`).
- **Instance-First Semantics:** All governance attaches to a specific *execution instance*, not the general agent definition.
- **Materiality Assessment:** Declarative policies determine if a workflow is "high risk" (requires human) or "low risk" (auto-approve).

## Strategic Non-Goals
- NO model training or fine-tuning.
- NO building or hosting agents.
- NO prompt optimization.
