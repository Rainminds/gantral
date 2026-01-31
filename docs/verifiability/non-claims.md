---
title: Explicit Non-Claims
---

# Explicit Non-Claims

This document enumerates what **Gantral explicitly does *not* claim**.

These non-claims are intentional.
They exist to prevent over-interpretation, scope creep, and misplaced trust.

If a property or guarantee is not explicitly claimed elsewhere in the
Verifiability documentation, it must be assumed **out of scope**.

---

## Purpose of Non-Claims

Verifiability depends not only on what a system can prove,
but on what it **refuses to imply**.

Gantral is designed to produce **verifiable execution evidence**.
It is not designed to resolve every governance, security, or legal question.

This document defines the **boundaries** of Gantral’s responsibility.

---

## What Gantral Does Not Claim

Gantral explicitly does **not** claim to:

### Guarantee Legal Admissibility

- Gantral does not guarantee that any artifact or evidence
  will be admissible in court or regulatory proceedings.
- Legal admissibility is determined by external authorities
  based on jurisdiction, procedure, and context.

Gantral provides **technical verifiability**, not legal certification.

---

### Prevent Malicious Actors

- Gantral does not claim to prevent malicious operators.
- Gantral does not claim to stop intentional misuse.
- Gantral does not claim to secure compromised infrastructure.

Gantral’s role is to make misuse **detectable**, not impossible.

---

### Guarantee Correctness of Decisions

- Gantral does not judge whether a decision was correct.
- Gantral does not assess business appropriateness.
- Gantral does not evaluate ethical or normative correctness.

Gantral records **authority and execution**, not judgment quality.

---

### Replace Legal, Risk, or Compliance Functions

- Gantral is not a compliance framework.
- Gantral is not a certification authority.
- Gantral does not replace legal counsel, auditors, or regulators.

Gantral supplies evidence that those functions may examine.

---

### Secure the Entire System Stack

Gantral does not claim to secure:

- operating systems
- hardware
- networks
- identity providers
- application logic outside the execution boundary

Gantral assumes failures may occur elsewhere and designs for **post-incident verification**.

---

### Infer Missing Context

- Gantral does not infer intent.
- Gantral does not reconstruct narratives.
- Gantral does not fill gaps in evidence.

When evidence is insufficient, verification must return **INCONCLUSIVE**.

---

## What Gantral *Does* Claim (For Contrast)

Gantral claims only that:

- execution authority is enforced at execution time
- a commitment artifact is emitted when authority is enforced
- that artifact can be independently replayed and verified
- failure and ambiguity are surfaced, not hidden

Any claim beyond this set is invalid.

---

## Why These Non-Claims Matter

Explicit non-claims ensure that:

- absence of proof is not mistaken for proof
- pressure cannot force positive assertions
- auditors can trust negative or inconclusive outcomes
- the system remains falsifiable under scrutiny

A system that never says “no” cannot be trusted.

---

## Relationship to Other Verifiability Documents

This document constrains:

- **Threat & Adversary Model**  
  by limiting assumed protections

- **Commitment Artifact**  
  by bounding what the artifact can prove

- **Replay Protocol**  
  by defining when verification must refuse

- **Failure Semantics**  
  by justifying inconclusive and invalid outcomes

All verifiability guarantees must be interpreted through these limits.

---

## Guiding Principle

Gantral would rather fail, refuse, or return inconclusive
than imply guarantees it cannot prove.

These non-claims are not omissions.
They are how verifiability remains honest.
