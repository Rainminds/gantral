---
title: Replay Protocol
---

# Replay Protocol

This document defines the **third-party replay and verification protocol**
for Gantral execution decisions.

It answers the question:

> Given only a commitment artifact and public Gantral execution semantics,
> how can an independent party verify whether execution authority was
> functioning at the exact moment a consequential action occurred?

This protocol is designed for **post-incident, adversarial verification**.

---

## Purpose of the Replay Protocol

The replay protocol exists to ensure that verification:

- does not depend on Gantral infrastructure
- does not require operator cooperation
- does not rely on logs, dashboards, or testimony
- produces deterministic, reproducible outcomes

If verification requires **trust**, **access**, or **interpretation**,
then replay has failed.

---

## Inputs Required

The replay protocol requires **only** the following inputs:

1. **Commitment Artifact**
   - The immutable execution-time artifact emitted by Gantral
   - Obtained from any source (archive, evidence store, cold storage)

2. **Public Execution Semantics**
   - The documented rules governing:
     - authority states
     - execution state transitions
     - escalation and pause behavior

No other inputs are permitted.

---

## Inputs Explicitly Prohibited

Replay **must not** rely on:

- Gantral services or APIs
- internal databases
- application or workflow logs
- dashboards or UIs
- operator testimony or explanation
- credentials, secrets, or access tokens

If any of the above are required,
verification must return **INCONCLUSIVE**.

---

## Replay Preconditions

Before replay begins, the verifier must establish:

- the artifact is intact and readable
- the artifact has not been altered
- the execution semantics version referenced by the artifact is available

If these conditions are not met,
replay must not proceed.

---

## Deterministic Replay Procedure

The replay procedure is **deterministic**.
Given identical inputs, it must always produce identical outputs.

### Step 1: Artifact Integrity Check

- Verify the artifact’s integrity binding
- Validate cryptographic hashes
- If integrity validation fails:
  → **INVALID**

---

### Step 2: Artifact Completeness Check

- Confirm all required fields are present
- Confirm no ambiguous or missing authority data
- If required fields are missing:
  → **INCONCLUSIVE**

---

### Step 3: Context Reconstruction

Using the artifact:
- reconstruct the execution instance
- identify the execution state transition
- identify the authority state at commit time
- identify the policy evaluation logic version

No external context may be inferred.

---

### Step 4: Authority Evaluation

Using public execution semantics:
- evaluate whether the recorded authority state
  was sufficient to permit the recorded execution transition

If authority was insufficient or contradictory:
→ **INVALID**

---

### Step 5: State Transition Validation

- Validate that the execution state transition
  is permitted under Gantral semantics
- Validate escalation and pause rules

If the transition violates semantics:
→ **INVALID**

---

### Step 6: Temporal Validation

- Validate the timestamp boundary
- Confirm ordering and sequencing constraints

If timing cannot be established:
→ **INCONCLUSIVE**

---

### Step 7: Final Determination

Based on the above steps, the verifier must return **exactly one** outcome.

---

## Replay Outcomes

Replay produces one of three outcomes:

### VALID

Returned when:
- artifact integrity is intact
- authority was sufficient at execution time
- execution state transition is valid
- no prohibited dependencies were required

This means:
> Execution authority is proven for this instance.

---

### INVALID

Returned when:
- artifact integrity is broken
- authority was insufficient or contradictory
- execution transition violates semantics
- substitution or tampering is detected

This means:
> Execution authority cannot be proven and evidence is invalid.

---

### INCONCLUSIVE

Returned when:
- artifact is missing or incomplete
- required semantics are unavailable
- timing or sequencing cannot be established
- verification would require prohibited inputs

This means:
> Execution authority cannot be determined.

Failing closed is intentional.

---

## Negative Proof and Refusal

The replay protocol is explicitly designed to **refuse** verification
when proof is insufficient.

Examples:
- “Authority cannot be proven here.”
- “Replay fails due to missing artifact.”
- “Verification requires external trust.”

Such outcomes are **correct behavior**, not errors.

---

## Independence Guarantee

Replay and verification must be possible:

- without Gantral’s cooperation
- without contacting Gantral operators
- without runtime system access
- without privileged credentials

If cooperation is required, replay is invalid.

---

## Relationship to Other Verifiability Documents

This protocol depends on:

- **Threat & Adversary Model**  
  Defines the hostile conditions replay must survive.

- **Commitment Artifact**  
  Defines the sole unit of verification.

- **Failure Semantics**  
  Defines expected refusal and inconclusive cases.

All three must be read together.

---

## Non-Goals of the Replay Protocol

This protocol does **not** attempt to:

- explain intent or rationale
- judge correctness of decisions
- assess business appropriateness
- determine legal admissibility
- replace regulatory or judicial review

Its sole purpose is **technical verification of execution authority**.

---

## Guiding Principle

If a third party cannot independently replay and verify execution authority
using only the commitment artifact and public semantics,
then Gantral must not claim that authority can be proven.

The replay protocol defines how that claim is tested.
