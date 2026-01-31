---
title: Security Model
---

# Security Model

This document defines the **security model and boundaries** under which Gantral is designed
to operate.

It is written to support **security, risk, and compliance review** in environments where
execution authority, accountability, and post-incident verification are required.

This document does **not** assert compliance with any standard or certification.
It does **not** describe guarantees outside the scope of Gantral’s execution authority.

Authoritative technical guarantees are defined in the Architecture and Implementation
documentation. If a discrepancy exists, those documents are authoritative.

---

## Purpose of the Security Model

Gantral’s security posture is derived from **explicit boundaries**, not implicit trust.

This document exists to answer:

- What Gantral assumes about its environment
- What Gantral explicitly does and does not protect
- What failures are expected and acceptable
- Where responsibility ends

If a risk or protection is not explicitly described here, it must be assumed **out of scope**.

---

## Security Boundary Overview

Gantral is designed as **execution authority infrastructure**.

It enforces *whether* execution may proceed.
It does not participate in *how* execution is reasoned about or carried out.

As a result:

- compromising Gantral must not expose agent internals
- compromising agents must not bypass execution authority
- removing Gantral must not silently allow execution to proceed

Security is enforced through **separation of responsibilities**, **data minimization**, and
**deterministic control semantics**.

---

## Core Security Principles

### 1. Separation of Authority and Intelligence

Gantral enforces execution authority but does not perform reasoning,
planning, or decision-making.

- AI systems may propose actions
- Policy engines may evaluate conditions
- Humans may authorize or deny execution
- Gantral enforces the resulting authority state

Gantral does not interpret intent or correctness.
It enforces authority transitions as state.

---

### 2. Data Minimization by Design

Gantral is intentionally scoped to avoid handling sensitive payloads.

Gantral does **not** store:
- prompts or model inputs
- agent memory or reasoning traces
- raw tool inputs or outputs
- credentials or secrets

Gantral stores only:
- execution state transitions
- authority decisions
- timestamps and identity references
- references (not contents) to external context where required

This minimizes exposure and simplifies data classification.

---

### 3. Deterministic Execution Control

Execution authority decisions in Gantral are:

- instance-scoped
- append-only
- timestamped
- deterministic and replayable

Replay reconstructs **authority and execution state only**.
It does not depend on agent memory, logs, or reconstructed narratives.

Determinism is a security property: it prevents post-hoc reinterpretation.

---

## Identity and Access Boundaries

Gantral integrates with existing identity infrastructure rather than introducing new identity systems.

### Human Identity
- Federated via upstream identity providers (OIDC)
- Identity derived from token claims
- No Gantral-managed user directory

### Machine Identity
- Workload or service identity only
- No static API keys
- No long-lived credentials stored by Gantral

Identity is used for **attribution and authorization**, not trust.

---

## Secrets and Credential Boundaries

Gantral is designed to avoid direct interaction with secrets.

- Secrets are referenced, not persisted
- Resolution occurs at execution edges (e.g., runners)
- External secret managers remain authoritative

Gantral never:
- stores raw credentials
- inspects secret contents
- mediates secret access decisions

---

## Deployment and Network Assumptions

Gantral supports deployment in:

- self-hosted environments
- private cloud or on-prem systems
- restricted or air-gapped networks

Gantral does not require outbound connectivity to function.

Network isolation and data residency controls are enforced by the deployment environment,
not by Gantral itself.

---

## Auditability and Evidence Integrity

Gantral produces:

- immutable execution records
- explicit authority decisions
- deterministic execution histories

Gantral may reference external evidence, but:

- does not inspect payload contents
- does not interpret evidence semantics
- does not authorize based on evidence meaning

Authority remains bound to recorded state transitions, not external artifacts.

---

## Explicit Non-Protections

Gantral does **not** claim to:

- prevent malicious human actions
- secure compromised agent runtimes
- validate correctness of business logic
- enforce safe or ethical AI behavior
- replace security reviews or change management
- guarantee regulatory or legal outcomes

Gantral governs **whether execution may proceed**, not **whether execution is safe or correct**.

---

## Relationship to Verifiability

This security model exists to support **verifiability**, not to replace it.

Security failures are expected and tolerated **if they remain detectable** through:

- commitment artifacts
- deterministic replay
- explicit failure semantics

When proof cannot be produced, Gantral must fail closed or return inconclusive results.

---

## Relationship to Technical Documentation

This document provides a **boundary-level security model**.

Detailed guarantees, invariants, and implementation requirements are defined in:
- Architecture documentation
- Implementation Guide
- Verifiability documentation

If a discrepancy exists, technical documents are authoritative.

---

## Final Boundary Statement

Gantral’s security posture is defined by **what it refuses to do**.

It does not expand its scope to appear safer.
It enforces narrow, explicit authority boundaries and produces evidence when those boundaries
are crossed.

This restraint is intentional.
