---
sidebar_position: 1
title: What is Gantral?
---

# What Is Gantral?

Gantral is an **open-source AI Execution Control Plane**.

It standardizes how AI-enabled workflows are executed, paused, escalated, approved, overridden, and audited across teams and systems.

Gantral exists to solve a specific problem faced by large organizations:

> **AI adoption breaks execution control and accountability—not just model quality.**

As AI tools spread across the software development lifecycle (SDLC) and operational workflows, organizations lose a consistent way to answer fundamental questions:

- What ran?
- Under whose authority?
- With what configuration?
- What human approved or overrode the outcome?
- Can this decision be replayed and audited?

Gantral provides infrastructure-level mechanisms to record and surface answers to these questions.

---

## The Core Idea: Authority vs. Intelligence

Gantral introduces a **shared execution plane** that separates **Authority** from **Reasoning**.

*   **Agents (CrewAI, LangGraph, etc.)** provide the *Intelligence*. They plan, reason, and generate code.
*   **Gantral** provides the *Authority*. It decides if the agent is allowed to proceed, pauses for human input, and records the outcome.

This separation prevents AI-driven execution from advancing past governed states without explicit authorization.

---

## How It Works (The Runner Pattern)

Gantral does not "host" your agents like a PaaS. It orchestrates them via a **Federated Runner** model (similar to GitHub Actions runners).

1.  **Agents** run in your own infrastructure (Kubernetes, Lambda, etc.).
2.  **Runners** pull tasks from Gantral.
3.  **Gantral** enforces policy barriers (e.g., "High Materiality = Pause").
4.  If a human is required, Gantral suspends the workflow. The agent process can exit (Zero CPU).
5.  On approval, Gantral reschedules the task, and the Runner launches a new agent process to complete the work.

---

## What Gantral Owns

Gantral owns **execution semantics**, not agent memory.

Specifically, Gantral provides:
- A deterministic execution state machine.
- Human-in-the-Loop (HITL) as a blocking state transition.
- Instance-level isolation for audit, cost, and accountability.
- Declarative control policies (Rego-based).
- Immutable execution records with replay capability.

Gantral is intentionally boring, predictable, and auditable.

---

## Mental Models

Gantral can be understood as:

- **"Sudo" for AI** — An agent tries to execute a command, but Gantral intercepts it and asks, "Are you authorized?"
- **Kubernetes for Semantics** — It manages the *lifecycle state* of AI processes, not just the containers.
- **Terraform for Process** — It defines the "Infrastructure of Decision Making" as code.

---

## Next Steps

*   See how the architecture ensures safety in the **[Technical Reference (TRD)](../architecture/trd.md)**.
*   See it in action with the **[Persistent Agent Demo](../guides/demo.md)**.
