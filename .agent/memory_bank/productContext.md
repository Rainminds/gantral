# Product Context: Gantral & Gantrio

## The Core Problem
Enterprises adopt AI tools (agents, copilots) but lack a unified "AI Execution Control Plane" to manage governance, auditability, and Human-in-the-Loop (HITL) workflows at scale. Gantral solves the question: "Who authorized this AI action?"

## What Gantral Is (Open Core)
- **Infrastructure, not an Agent Builder:** It sits *above* frameworks like LangChain/CrewAI and *below* enterprise processes.
- **Execution Control Plane:** Standardizes how workflows execute, pause for humans, and resume.
- **System of Record:** Provides immutable audit logs for every AI decision.
- **Federated Authority:** Manages the *decision* to execute, while delegating the *actual execution* to runners in user infrastructure.

## Key Differentiation
- **HITL is a State, not a Feature:** Human intervention is modeled as a first-class state transition (`WAITING_FOR_HUMAN` -> `APPROVED/REJECTED`).
- **Federated Execution:** Code and data never leave the customer's VPC. Gantral orchestrates via "Runners" (Pull-based).
- **Agent-Native Persistence:** Supports long-term pauses (days/weeks) by relying on Agent Framework checkpointing (scale-to-zero waiting).
- **Policy as Guard:** Policies are evaluated *synchronously* at state transitions, making it impossible to bypass governance.

## Strategic Non-Goals
- NO model training or fine-tuning.
- NO building or hosting agents.
- NO prompt optimization.
- NO managing secrets (References only).
- NO user database (Federated Identity only).
