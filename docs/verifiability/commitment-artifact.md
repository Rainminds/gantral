---
title: Commitment Artifact
---

# Commitment Artifact

This document defines the **commitment artifact** emitted by Gantral at execution time.

The commitment artifact is the **unit of verification** in Gantral.
It is the only object considered authoritative during post-incident,
third-party verification.

This document answers, without hand-waving:
- what the artifact is
- when it is emitted
- what it contains
- what it binds
- what cannot be altered after emission

If a claim about execution authority cannot be derived from this artifact,
that claim is **out of scope**.

---

## What the Commitment Artifact Is

A commitment artifact is an **immutable execution-time evidence object**
emitted by Gantral at the moment execution authority is enforced.

It binds:
- authorization state
- execution state transition
- policy evaluation context
- execution context fingerprint
- time boundary

into a single, verifiable record.

The artifact is **not** a log entry.
It is **not** a narrative.
It is **not** reconstructed after execution.

It is emitted **as part of the execution commit path**.

---

## When the Artifact Is Emitted

The commitment artifact is emitted:

- **at execution time**
- **at the authority boundary**
- **at the moment a consequential action is committed**

Emission occurs **atomically with (or as close as computably possible to)**:
- the authorization decision
- the execution state transition

There is no valid execution path in which:
- execution proceeds
- authority is enforced
- and no commitment artifact is emitted

**No artifact → no execution.**

---

## What the Artifact Contains

The following fields are **normative**.
Names are illustrative; semantics are required.

### Core Fields

- **execution_instance_id**  
  A unique identifier for the execution instance.

- **execution_state_transition**  
  The state transition being committed  
  (e.g. `WAITING_FOR_HUMAN → RUNNING`, `RUNNING → TERMINATED`).

- **authority_state_hash**  
  A cryptographic commitment to the authority state at commit time,
  including whether authority was automated, human, or escalated.

- **policy_version_reference**  
  A reference to the **policy evaluation logic version** used.  
  This is **not** a policy decision and **not** an approval.

- **execution_context_fingerprint**  
  A deterministic fingerprint of the execution context
  sufficient to prevent out-of-context replay.

- **timestamp_boundary**  
  A cryptographically bound time boundary indicating
  when the commit occurred.

---

### Optional but Recommended Fields

- **previous_artifact_hash**  
  Links this artifact to a prior execution artifact
  to form a verifiable chain.

- **human_decision_reference**  
  Present only when human authority was exercised.
  References the decision without embedding narrative.

- **integrity_binding**  
  Cryptographic binding over all artifact fields.

---

## What the Artifact Explicitly Does *Not* Contain

The commitment artifact deliberately excludes:

- application logs
- UI state
- dashboards or workflow metadata
- operator explanations
- intent narratives
- reconstructed timelines

Logs may explain.  
Narratives may persuade.  
The artifact is what can be verified.

---

## Artifact Properties (Required)

A valid commitment artifact must be:

- **Immutable**  
  Any modification invalidates verification.

- **Deterministically Replayable**  
  A third party can replay authority and execution semantics
  without access to Gantral systems.

- **Context-Complete**  
  No external inference or missing data is required
  to evaluate authority at commit time.

- **Independent of Trust**  
  Verification does not rely on:
  - operator testimony
  - database integrity
  - dashboards
  - platform access

If any of these properties do not hold,
verification must return **INVALID** or **INCONCLUSIVE**.

---

## What the Artifact Is Sufficient to Prove

Given the artifact and public execution semantics,
a verifier can determine:

- whether execution was allowed, paused, escalated, or terminated
- whether human authority was required
- whether a human decision was captured
- which policy evaluation logic was applied
- when authority was exercised
- whether execution state transitions were consistent

These are **execution facts**, not interpretations.

---

## What the Artifact Cannot Prove

The commitment artifact does **not** prove:

- correctness of human judgment
- intent or motivation
- business appropriateness
- compliance classification
- ethical evaluation

Absence of proof in these areas is intentional.

---

## Failure Conditions

The following conditions are explicitly defined:

- **Missing artifact**  
  → verification is **INCONCLUSIVE**

- **Altered artifact**  
  → verification is **INVALID**

- **Artifact emitted outside commit path**  
  → verification is **INVALID**

- **Context incomplete or ambiguous**  
  → verification is **INCONCLUSIVE**

Gantral fails closed by design.

---

## Relationship to Replay Protocol

The commitment artifact is the **only required input**
to the replay protocol, alongside public execution semantics.

The replay protocol must not:
- consult logs
- query Gantral services
- rely on credentials
- infer missing context

If verification requires any of the above,
the artifact is insufficient.

---

## Design Rationale

The commitment artifact exists to ensure that:

- authorization and execution are inseparable
- authority is bound at execution time
- post-incident verification does not depend on trust
- substitution and narrative reconstruction are detectable

This design favors **falsifiability over convenience**.

---

## Guiding Principle

If execution authority cannot be proven using only the commitment artifact
and public execution semantics,
then Gantral must not claim that it can be proven.

The commitment artifact defines the boundary of that claim.

---

## Implementation Specification

Gantral’s commitment artifact is implemented according to a
separate, normative implementation specification that defines
exact field structure, hash construction, atomic emission rules,
and offline verification logic.

That specification is used by Gantral implementations and
independent verifiers, but is not required to understand
the verifiability claims made in this document.
