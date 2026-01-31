---
title: Failure Semantics
---

# Failure Semantics

This document defines **how and when verifiability intentionally fails** in Gantral.

Failure is not an error condition.
Failure is a **designed outcome** when execution authority cannot be proven
under hostile, post-incident scrutiny.

Gantral fails **closed by construction**.

---

## Purpose of Failure Semantics

Verifiability is meaningful only if the system can say:

> “Authority cannot be proven here.”

This document answers:

- When does verification **refuse** to produce proof?
- When does verification return **INVALID**?
- When does verification return **INCONCLUSIVE**?
- Why these outcomes are correct behavior

If a system always produces a positive answer, it is not verifiable.

---

## Failure Categories

Gantral defines three distinct failure categories:

| Category | Meaning |
|--------|--------|
| **INVALID** | Evidence is present but demonstrably incorrect or tampered |
| **INCONCLUSIVE** | Evidence is missing, incomplete, or insufficient |
| **REFUSAL** | Execution is blocked because proof cannot be produced |

Each category is intentional and non-overlapping.

---

## INVALID Outcomes

An **INVALID** outcome means that verification **actively detects a violation**.

### Conditions That Produce INVALID

Verification must return **INVALID** when:

- The commitment artifact integrity check fails
- Artifact fields have been altered or substituted
- Authority state contradicts execution transition
- Execution transition violates published semantics
- Artifact was emitted outside the execution commit path
- Conflicting artifacts exist for the same execution instance
- Replay requires prohibited inputs (logs, testimony, platform access)

### Interpretation

An INVALID result means:

> “Execution authority cannot be proven because the evidence is compromised or inconsistent.”

INVALID is a **hard failure**.

---

## INCONCLUSIVE Outcomes

An **INCONCLUSIVE** outcome means that verification **cannot establish proof**.

### Conditions That Produce INCONCLUSIVE

Verification must return **INCONCLUSIVE** when:

- The commitment artifact is missing
- Required artifact fields are absent
- Execution semantics version is unavailable
- Timing or sequencing cannot be established
- Context fingerprint is insufficient
- Only partial evidence survives
- Verification would require inference or narrative

### Interpretation

An INCONCLUSIVE result means:

> “Execution authority cannot be determined from available evidence.”

This is not an error.
It is an explicit refusal to speculate.

---

## Execution-Time Refusal (Prevention)

Failure semantics apply **before execution**, not only after.

Gantral must refuse execution when:

- Authority state is ambiguous
- Required human authority is unavailable
- Escalation cannot be resolved
- Policy evaluation returns conflicting signals
- Commitment artifact emission fails

### Rule

> **No artifact → no execution**

Execution must not proceed if verifiable evidence cannot be produced.

---

## Negative Proof Examples

The following are **expected and correct outcomes**.

### Example 1: Missing Artifact

- Execution occurred
- No commitment artifact exists

**Outcome:** INCONCLUSIVE  
**Reason:** Authority cannot be proven.

---

### Example 2: Log-Based Reconstruction

- Logs indicate approval
- Artifact is missing or incomplete

**Outcome:** INCONCLUSIVE  
**Reason:** Logs are non-authoritative.

---

### Example 3: Altered Artifact

- Artifact exists
- Integrity binding fails

**Outcome:** INVALID  
**Reason:** Evidence is compromised.

---

### Example 4: Policy Disagreement

- Policy evaluation logic changed after execution
- Artifact references unavailable version

**Outcome:** INCONCLUSIVE  
**Reason:** Evaluation context cannot be reconstructed.

---

### Example 5: Authority Revoked Mid-Execution

- Execution started under valid authority
- Authority revoked before commit
- Artifact reflects ambiguity

**Outcome:** REFUSAL or INCONCLUSIVE  
**Reason:** Authority was not provably active at commit time.

---

## Why Failure Is a Feature

Gantral prefers:

- refusal over ambiguity
- inconclusive over false certainty
- invalid over narrative reconstruction

This ensures that:

- absence of proof is never mistaken for proof
- pressure cannot force a positive result
- auditors and regulators can trust negative outcomes

---

## Relationship to Replay Protocol

Failure semantics are enforced during replay:

- **INVALID** is returned when violations are detected
- **INCONCLUSIVE** is returned when proof is insufficient
- Replay must never “best-guess” or infer

Replay that does not honor failure semantics is incorrect.

---

## Relationship to Threat Model

Failure semantics exist because Gantral assumes:

- hostile reviewers
- uncooperative operators
- substituted logs
- missing testimony

If Gantral attempted to “fill in the gaps,”
it would violate its own threat model.

---

## Non-Goals of Failure Semantics

Failure semantics do **not** attempt to:

- determine blame
- assess intent
- judge correctness
- satisfy legal standards
- maximize coverage

They exist solely to preserve **verifiability integrity**.

---

## Guiding Principle

If execution authority cannot be proven **without trust or inference**,
Gantral must fail.

Failure is not weakness.
Failure is how verifiability remains honest.
