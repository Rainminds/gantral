---
title: What Gantral Is Not
---

Clear boundaries are essential for infrastructure.

Gantral is intentionally narrow.
Its value comes from **what it enforces** and, equally, from **what it refuses to do**.

If Gantral attempted to solve adjacent problems in AI development, automation,
or analytics, it would undermine its core purpose: **execution authority**.

This document defines those boundaries explicitly.

---

## Gantral Is NOT an Agent Framework

Gantral does not:

- build AI agents
- define agent architectures
- manage prompts or toolchains
- select, train, or fine-tune models
- persist or inspect agent memory or reasoning state

Agent frameworks are responsible for **intelligence**.
Gantral is responsible for **authority**.

This separation ensures that:
- agent evolution does not destabilize execution control
- auditability does not depend on agent internals
- governance remains independent of reasoning behavior

---

## Gantral Is NOT a Workflow Engine or Orchestrator

Gantral does not:

- replace workflow engines (Temporal, Airflow, etc.)
- own business process definitions
- schedule arbitrary jobs
- optimize execution performance
- define task graphs or retries

Gantral governs **whether execution may proceed**, not **how workflows are designed or optimized**.

Workflow engines handle orchestration.
Gantral enforces authority at execution boundaries.

---

## Gantral Is NOT an End-to-End SDLC Automation Platform

Gantral does not:

- replace CI/CD systems
- replace GitHub, Jira, ServiceNow, or ticketing tools
- claim to automate the full SDLC
- promise autonomous engineering organizations

Gantral integrates *into* existing SDLC systems.
It does not attempt to subsume them.

Its role is to ensure that **execution decisions inside workflows are authorized and attributable**.

---

## Gantral Is NOT an Observability, Monitoring, or Analytics Tool

Gantral does not:

- benchmark model quality
- measure accuracy, precision, or recall
- analyze agent performance
- replace observability or AgentOps platforms
- interpret logs or metrics

Gantral may **consume signals** from observability systems,
but it does not generate insight or analysis from them.

Its concern is **authority**, not measurement.

---

## Gantral Is NOT a Policy Authoring System

Gantral does not:

- define organizational policy
- replace policy engines
- interpret regulatory text
- encode compliance logic directly

Policy engines (e.g. OPA) provide **advisory evaluation**.
Gantral enforces **authority decisions** as execution state.

Policy advises.
Authority enforces.

---

## Gantral Is NOT a UI-First or Convenience-Driven Product

Gantral does not prioritize:

- visual workflow builders
- no-code orchestration
- dashboard-driven execution control
- convenience over determinism

Explicit semantics and enforceable boundaries take precedence over usability shortcuts.

If a choice must be made, Gantral favors:
- correctness over convenience
- determinism over flexibility
- failure over ambiguity

---

## Gantral Is NOT Autonomous AI Infrastructure

Gantral does not:

- allow AI systems to approve themselves
- support self-authorizing execution loops
- remove human accountability from material actions
- “trust” AI judgment for irreversible execution

If an action is material,
Gantral enforces **explicit authority**.

Human involvement is not optional where authority is required.

---

## Gantral Is NOT a Compliance or Legal Guarantee

Gantral does not:

- guarantee regulatory compliance
- certify legal admissibility
- replace auditors, regulators, or legal review
- determine correctness or intent

Gantral provides **technical verifiability**, not legal judgment.

Absence of proof is treated honestly as absence of proof.

---

## Why These Non-Goals Matter

Infrastructure fails in regulated environments not because it lacks features,
but because it lacks **clear limits**.

Gantral’s non-goals ensure that:

- execution authority remains explicit
- audit does not depend on inference
- failure is detectable
- responsibility is not blurred

If a use case requires:
- silent autonomy
- implicit approval
- narrative-based reconstruction

Gantral is not the right tool.

That is not a limitation.
It is the design.
