---
title: Admissible Execution
sidebar_label: Admissible Execution
---

> **Canonical Reference**
>
> This document is derived from the position paper  
> **“Admissible Execution: Invariants for AI Execution Authority”**, published on Zenodo.
>
> **Read the authoritative version (DOI):**  
> https://doi.org/10.5281/zenodo.18471561
>
> Version: 1.0 · January 2026

---

## Purpose

This document defines the **minimum admissibility bar** for execution-time authority in AI-driven systems.

It does not describe an implementation, protocol, or product.
It specifies **what must be true** for execution authority to survive audit, incident review, and adversarial scrutiny.

---

## When this matters

You should read this if you are:
- designing execution controls for AI or automation
- responsible for audit, risk, or regulatory defensibility
- evaluating whether “human-in-the-loop” claims are enforceable
- building systems where decisions have irreversible effects

---

## The Admissibility Bar (Summary)

At a minimum, admissible execution requires that:
- authority is explicit and enforced at runtime
- execution halts in the absence of authority
- approvals bind to execution instances
- execution history is immutable and replayable
- evidence does not depend on operator testimony or mutable logs

Anything weaker is auditability.
Admissibility is stronger.

---

## Invariants (Normative)

For the full invariant set, see the canonical paper.

The following invariants are **non-negotiable**:
- Authority is separate from intelligence
- Execution halts without authority
- Authority binds to an execution instance
- Authority is explicit, not ambient

Violation of any invariant renders execution authority inadmissible.

---

## Relationship to Gantral

Gantral is designed to **meet these admissibility requirements** by providing
verifiable execution-time authority guarantees.

This document defines the bar.
Gantral implements the mechanisms.

---

> For the complete definitions, invariants, and failure modes,  
> refer to the canonical Zenodo publication.
