---
title: Auditor Verification Guide
---

# Auditor Verification Guide

This guide describes how to **independently verify a Gantral execution decision**
using a commitment artifact obtained long after execution occurred.

It is written for auditors, regulators, investigators, and internal review teams
operating **post-incident**, under **adversarial assumptions**.

This guide assumes:
- no access to Gantral systems
- no access to operator testimony
- no access to dashboards or logs
- no cooperation from engineering teams

If verification depends on any of the above, it has failed.

---

## Scenario

> You are given a commitment artifact from an AI-driven execution that occurred **three years ago**.  
> The execution resulted in a consequential real-world action.  
> You are asked to determine whether execution authority was functioning correctly at the time.

Your task is **not** to judge whether the decision was correct.  
Your task is to determine **whether authority can be proven**.

---

## What You Are Given

At minimum, you should be provided:

- a **commitment artifact** (file, record, or object)
- a reference to **public Gantral execution semantics**
- the claimed execution outcome (e.g. “this action was authorized”)

You should **not** require:
- application logs
- workflow dashboards
- human explanations
- access to Gantral APIs
- credentials or secrets

If any of these are required, verification must stop.

---

## Step 1: Establish Preconditions

Before attempting verification, confirm:

- the artifact is intact and readable
- the artifact has not been modified
- the execution semantics version referenced by the artifact is available

If any of these are not true:

**Outcome:** INCONCLUSIVE  
**Reason:** Verification cannot proceed.

---

## Step 2: Verify Artifact Integrity

Check that:

- cryptographic bindings are intact
- hashes validate correctly
- the artifact has not been altered or partially substituted

If integrity validation fails:

**Outcome:** INVALID  
**Reason:** Evidence is compromised.

If integrity validation succeeds, continue.

---

## Step 3: Verify Artifact Completeness

Confirm the artifact contains all required fields, including:

- execution instance identifier
- execution state transition
- authority state hash
- policy evaluation logic reference
- execution context fingerprint
- timestamp boundary

If required fields are missing or ambiguous:

**Outcome:** INCONCLUSIVE  
**Reason:** Evidence is insufficient.

---

## Step 4: Reconstruct Authority State

Using the artifact alone:

- identify the execution state transition being committed
- determine whether authority was:
  - automated
  - human-granted
  - overridden
  - revoked

Do **not** infer intent or reasoning.
Do **not** consult external narratives.

If authority state cannot be reconstructed deterministically:

**Outcome:** INCONCLUSIVE

---

## Step 5: Validate Authority Against Semantics

Using public Gantral execution semantics:

- determine whether the recorded authority state
  was sufficient to permit the recorded execution transition
- confirm that no forbidden transition occurred

If the transition violates semantics:

**Outcome:** INVALID

---

## Step 6: Temporal Validation

Evaluate the timestamp boundary:

- confirm ordering of authority and execution
- confirm authority was active **at commit time**
- confirm no retroactive approval is implied

If timing or ordering cannot be established:

**Outcome:** INCONCLUSIVE

---

## Step 7: Determine Final Outcome

Based on the steps above, return **exactly one** outcome.

---

## Interpretation of Outcomes

### VALID

**Meaning:**
- Execution authority is proven for this instance
- Authority was present at the moment execution occurred
- No trust in narratives or systems was required

**What This Does Not Mean:**
- The decision was correct
- The outcome was appropriate
- The system is compliant

---

### INVALID

**Meaning:**
- Evidence is present but contradictory or compromised
- Authority cannot be proven
- Execution should be treated as unauthorized

**Common Causes:**
- artifact tampering
- invalid state transitions
- retroactive approval attempts

---

### INCONCLUSIVE

**Meaning:**
- Authority cannot be determined from available evidence
- Absence of proof is not proof of absence

**Common Causes:**
- missing artifacts
- incomplete records
- unavailable semantics
- reliance on logs or testimony

INCONCLUSIVE is a **correct and expected outcome**.

---

## What You Should *Not* Do

During verification, you must not:

- infer intent from context
- rely on operator explanations
- reconstruct events from logs
- assume approval “must have happened”
- accept screenshots or dashboards as evidence

If such steps feel necessary, verification has failed.

---

## Relationship to Other Documents

This guide depends on:

- **Threat & Adversary Model**  
  Defines the hostile conditions assumed.

- **Commitment Artifact**  
  Defines the unit of verification.

- **Replay Protocol**  
  Defines the deterministic verification procedure.

- **Failure Semantics**  
  Defines how and why verification fails.

These documents must be read together.

---

## Final Reminder

This guide does not determine blame.
It does not determine correctness.
It does not determine compliance.

It determines one thing only:

> **Can execution authority be proven without trust?**

If the answer is anything other than **YES**,  
the correct outcome is **INVALID** or **INCONCLUSIVE**.

That result is not a failure of the audit.

It is the system behaving honestly.
