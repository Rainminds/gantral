# Architectural Invariants

This document defines architectural invariants for Gantral.

Invariants are constraints that guide design and implementation decisions.

They are intended to preserve execution control and auditability.

---

## Invariant 1: Instance-First Semantics

Execution metadata is associated with instances, not agents.

Audit, cost, and accountability attach to a specific execution instance.

---

## Invariant 2: Human Authority Is Explicit

Where human involvement is configured, execution requires an explicit human decision.

Human decisions are modeled within execution, not inferred externally.

---

## Invariant 3: Determinism Over Convenience

Execution behavior should be predictable and reconstructable.

Features that introduce hidden non-determinism are discouraged.

---

## Invariant 4: Declarative Control

Policies, escalation rules, and authority constraints should be expressed declaratively.

Embedding control logic directly in agent code is discouraged.

---

## Invariant 5: Adapters Contain No Business Logic

Integrations should emit events and consume decisions.

Business or governance logic should not be embedded in adapters.

---

## Interpretation

These invariants describe architectural intent.

They are not contractual guarantees and may evolve through documented governance processes.
