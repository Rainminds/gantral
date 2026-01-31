---
title: Threat & Adversary Model
---

# Threat & Adversary Model

This document defines the **explicit threat and adversary assumptions** under which
Gantral’s verifiability guarantees are designed to hold.

The purpose of this document is **scope clarity**, not reassurance.

If an adversary or capability is not explicitly listed here, it is **out of scope**.
Gantral makes no claims outside these bounds.

---

## Purpose of the Threat Model

Gantral is designed to produce **execution-time evidence** that survives
**post-incident, adversarial reconstruction**.

This threat model exists to answer one question:

> Under what adversarial conditions should a third party still be able to
> independently verify that execution authority was functioning at the moment
> a consequential action occurred?

This document deliberately assumes **uncooperative and hostile conditions**.

---

## Adversaries Assumed

Gantral assumes the presence of the following adversaries:

### Post-Incident Investigator
- Reviewing execution after an incident or failure
- Unwilling to rely on operator explanations
- Treats all narratives as potentially biased

### Auditor or Regulator
- Evaluating execution controls retrospectively
- Requires reproducible, third-party verification
- Does not trust platform operators or vendors

### Litigant or External Reviewer
- Attempting to disprove that authority was active
- Incentivized to find ambiguity or substitution
- Treats absence of proof as failure

### Internal Audit Acting Adversarially
- Operating independently from engineering teams
- Assumes internal incentives may bias reporting
- Challenges sequence, timing, and completeness

### Platform Operator Acting Defensively
- May withhold testimony
- May present selective logs
- May attempt to reconstruct events favorably

Gantral assumes **no cooperation** from any of the above.

---

## Adversary Capabilities Assumed

The following adversary actions are explicitly assumed to be possible:

- Substitution, omission, or reordering of logs  
- Withholding or refusal of operator testimony  
- Loss of dashboards, UIs, or workflow metadata  
- Replay of events out of original execution context  
- Challenges to sequence, timing, and causality  
- Claims that approval occurred “out of band”  

Gantral’s verifiability guarantees are designed to hold **despite** these actions.

---

## Adversaries and Capabilities *Not* Assumed

Gantral explicitly does **not** claim resistance against:

- Nation-state level attackers  
- Kernel-level or hypervisor compromise  
- Hardware trust violations (TPM, enclave failure)  
- Undetectable cryptographic primitive breakage  
- Complete loss of the commitment artifact  

If these occur, verification **may be impossible**.

Such scenarios are **out of scope**.

---

## Trust Assumptions (Minimal and Explicit)

Gantral makes the following minimal assumptions:

- Cryptographic primitives behave as specified
- Public execution semantics are available
- At least one copy of the commitment artifact survives

Gantral does **not** assume:
- honest operators
- intact databases
- truthful narratives
- reliable dashboards
- continuous system availability

---

## Implications for Verifiability

Given this threat model:

- Logs are treated as **untrusted**
- Testimony is treated as **irrelevant**
- Dashboards are treated as **non-authoritative**
- Only execution-time artifacts are considered **verifiable evidence**

If authority cannot be proven under these conditions,
verification must return **INCONCLUSIVE** or **INVALID**.

Failing closed is intentional.

---

## Relationship to Other Verifiability Documents

This threat model constrains and informs:

- **Commitment Artifact**  
  What must be bound at execution time to survive adversarial review.

- **Replay Protocol**  
  Why replay must function without platform access or cooperation.

- **Failure Semantics**  
  When and why verification intentionally fails.

These documents must be read together.

---

## Non-Goals of This Threat Model

This document does **not** attempt to:

- Prevent malicious behavior  
- Secure compromised infrastructure  
- Guarantee correctness or intent  
- Assert legal admissibility  
- Replace legal or regulatory judgment  

Gantral’s responsibility is limited to **producing verifiable execution evidence**
within the scope defined above.

---

## Guiding Principle

If execution authority cannot be independently verified
under the adversarial conditions defined here,
then Gantral must not claim that it can be proven.

This threat model defines the boundary of those claims.
