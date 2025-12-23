# What Is Gantral

Gantral is an **open-source AI Execution Control Plane**.

It standardizes how AI-enabled workflows are executed, paused, escalated, approved, overridden, and audited across teams and systems.

Gantral exists to solve a specific problem faced by large organizations:

> AI adoption breaks execution, control, and accountability — not model quality.

As AI tools spread across the software development lifecycle (SDLC) and operational workflows, organizations lose a consistent way to answer fundamental questions:

- What ran?
- Under whose authority?
- With what configuration?
- What human approved or overrode the outcome?
- Can this decision be replayed and audited?

Gantral provides infrastructure-level mechanisms to record and surface answers to these questions.

---

## The Core Idea

Gantral introduces a **shared execution plane** that sits:

- **Above** AI agent frameworks (LangChain, CrewAI, Vellum, custom agents)
- **Below** enterprise processes (SDLC, incident management, governance)

Agents perform computation.  
Gantral governs execution.

This separation is designed to prevent AI-driven execution from advancing past governed states without explicit authorization.

---

## What Gantral Owns

Gantral owns **execution semantics**, not agent intelligence.

Specifically, Gantral provides:

- A deterministic execution state machine
- Human-in-the-Loop (HITL) as a first-class state transition
- Instance-level isolation for audit, cost, and accountability
- Declarative control policies (materiality, escalation, authority)
- Immutable execution records with replay capability

Gantral is intentionally boring, predictable, and auditable.

---

## What Gantral Enables

With Gantral, organizations can:

- Standardize HITL across AI workflows
- Scale AI usage across hundreds of teams without duplicating agents
- Enforce governance policies without modifying agent code
- Produce audit-ready execution records by default
- Separate experimentation (agents) from accountability (execution)

Gantral does not make AI more powerful.

It makes AI **safe, governable, and operable at scale**.

---

## Mental Models

Gantral can be understood as:

- **Kubernetes** — but for AI execution semantics, not containers
- **Terraform** — but for AI process control, not infrastructure
- **ServiceNow** — but for AI execution governance, not ITSM

These analogies are conceptual and do not imply feature parity, compatibility, or equivalence.

Gantral defines a new control layer specific to AI-enabled workflows.
