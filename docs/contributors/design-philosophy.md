---
title: Design Philosophy
---

# Design Philosophy

Gantral is infrastructure.

Infrastructure succeeds when it is predictable, boring, and trustworthy.

This philosophy guides all design and contribution decisions.

---

## Core Design Principles

### 1. Execution Control Over Intelligence

Gantral governs *how* AI participates in workflows, not *what* AI decides.

Improvements to model quality, prompts, or agent behavior are explicitly out of scope.

---

### 2. Human Authority Is Explicit

Human involvement is modeled as part of execution, not as an external review step.

Where workflows are material, human decisions are treated as authoritative within the execution model.

---

### 3. Determinism Over Convenience

Gantral prioritizes:

- Deterministic execution
- Immutable records
- Replayability
- Clear state transitions

Convenience features that weaken these properties are avoided.

---

### 4. Declarative Over Imperative

Policies, escalation rules, and control logic should be declarative.

Embedding governance logic directly in code is discouraged.

---

### 5. Instance-First Semantics

Audit, cost, and accountability attach to execution instances, not agents or workflows.

Designs that blur instance boundaries are rejected.

---

## Design Trade-offs

Gantral intentionally trades:

- Speed for correctness
- Flexibility for clarity
- Autonomy for accountability

These trade-offs are deliberate and should not be “optimized away.”

---

## When in Doubt

If a design decision raises ambiguity about:

- Who is accountable
- Whether a human can intervene
- How an auditor would reconstruct events

The design is likely incorrect.

---

## Scope Discipline

Adding features is easy. Removing them later is not.

Design proposals should err on the side of exclusion unless there is a clear execution-control justification.
