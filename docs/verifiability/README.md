---
title: Verifiability Overview
---

# Verifiability

This section documents how **Gantral execution decisions can be independently verified**
without trusting Gantral operators, infrastructure, user interfaces, or narratives.

The purpose of this section is not to explain *why* a decision was made,  
but to make it possible to **prove whether execution authority was functioning at the exact
moment a consequential action occurred**.

These documents are written for readers who are:
- skeptical
- adversarial
- operating post-incident
- unwilling to rely on testimony, dashboards, or internal assurances

If you are attempting to **disprove** Gantral’s execution-time authority claims, start here.

---

## What Verifiability Means in Gantral

In Gantral, **verifiability** means:

> A third party can independently reconstruct and evaluate an execution decision using only  
> (a) a recorded execution artifact and  
> (b) publicly documented execution semantics.

Verification must be possible:
- without access to Gantral services
- without access to internal databases
- without operator credentials
- without human explanation or testimony

Verification outcomes are strictly defined and limited to:
- **VALID** — authority and execution state are consistent and untampered
- **INVALID** — tampering, substitution, or inconsistency detected
- **INCONCLUSIVE** — insufficient or incomplete evidence

---

## What Verifiability Does *Not* Mean

Verifiability **does not imply**:
- legal admissibility in any jurisdiction
- regulatory approval or certification
- correctness of human judgment
- ethical or normative evaluation
- compliance classification of a system

Those determinations are made by **external authorities**, not by Gantral.

Gantral’s responsibility is narrower and structural:
to ensure that **execution authority leaves behind evidence that can be independently inspected**.

---

## How This Section Is Structured

This section is intentionally organized as a **chain of custody**, not a developer guide.

Each document addresses a specific requirement for hostile, third-party verification:

### Threat & Adversary Model
Defines the adversaries assumed, their capabilities, and the explicit limits of the threat model.

→ `threat-model.md`

---

### Commitment Artifact
Defines the execution-time artifact emitted by Gantral, including:
- what it contains
- when it is emitted
- what is bound into it
- what cannot be altered after emission

This artifact is the **unit of verification**.

→ `commitment-artifact.md`

---

### Replay Protocol
Defines how an independent third party can verify execution authority using:
- the artifact
- public Gantral execution semantics

This protocol explicitly prohibits reliance on logs, credentials, or platform access.

→ `replay-protocol.md`

---

### Failure Semantics
Documents when and how verification intentionally fails, including:
- missing artifacts
- altered artifacts
- ambiguous or incomplete authority

Gantral fails closed by design.

→ `failure-semantics.md`

---

### Explicit Non-Claims
Enumerates what Gantral explicitly does **not** attempt to prove or guarantee.

This document exists to prevent over-interpretation and scope creep.

→ `non-claims.md`

---

## Relationship to Admissibility

Verifiability is a **technical property**.

Legal or regulatory **admissibility** is an outcome determined by courts, regulators,
or other external authorities.

Gantral is designed so that execution decisions are *technically verifiable under hostile replay*.
Whether such evidence is deemed admissible is outside Gantral’s control.

Verifiability is a **necessary condition** for admissibility — not a guarantee of it.

---

## Intended Audience

This section is written for:
- auditors
- regulators
- litigators
- internal security review teams
- platform engineers performing post-incident analysis

It is **not** an onboarding guide and does not attempt to be user-friendly.

Clarity, falsifiability, and explicit limits take precedence over approachability.

---

## Guiding Principle

If a claim cannot be independently replayed by a hostile third party,  
it is not considered verifiable.

Everything in this section exists to satisfy that principle.
