---
title: Implementation Guide
---

# Comprehensive Implementation Guide

**Version:** v4.1 (Authority-First, Federated Execution, Atomic Commit)  
**Status:** Reference implementation guide (non-normative)  
**Audience:** Core contributors, platform engineers, enterprise architects (regulated environments), security / risk / compliance reviewers, AI coding assistants  

---

## Purpose

This document translates the **Gantral PRD v5.1** and **TRD v4.1** into a **buildable, implementation-accurate guide**.

It incorporates finalized architectural clarifications:

- Gantral owns **execution authority and control state**
- Agent frameworks own **agent internal state and memory**
- Long-running approvals require **agent-native checkpointing**
- Deterministic execution and replay are guaranteed by a **workflow runtime (e.g. Temporal)**
- Authority enforcement is **atomic with evidence emission**
- Execution instances are **explicitly owned and governed** (team + policy binding)

> **Rule:** If this guide conflicts with the TRD, the TRD is correct.

---

## 0. How This Guide Should Be Used

### 0.1 Documentation Hierarchy (Non-Negotiable)

```

PRD → TRD → Implementation Guide (this document)

```

- PRD defines **why** Gantral exists  
- TRD defines **what must be true**  
- This guide defines **one correct way to build it**

This guide is replaceable.  
The TRD is not.

---

## 1. System Responsibilities (Explicit)

### 1.1 What Gantral Owns

Gantral is an **execution authority layer**. It owns:

- Execution lifecycle and immutable execution instance boundaries
- Canonical **authority state machine**
- Human-in-the-Loop (HITL) enforcement
- Authority transitions (APPROVE / REJECT / OVERRIDE)
- Policy gating as **transition guards**
- Time, retries, escalation, and timeout semantics
- Instance-level audit correctness
- Deterministic replay of authority decisions (via workflow runtime)
- **Commitment artifact emission** at authority boundaries
- Binding of execution instances to **owning team and applicable policy set**

### 1.2 What Gantral Explicitly Does NOT Own

Gantral does **not**:

- Perform agent reasoning, planning, or cognition
- Manage prompts, tools, or agent memory
- Persist or serialize agent internal state
- Decide outcomes autonomously
- Store secrets or credentials
- Act as an identity provider

Gantral governs **whether execution may proceed**, not **how agents reason**.

---

## 2. Separation of Duties (Enterprise Requirement)

| Function                 | Owner                                      |
| ------------------------ | ------------------------------------------ |
| Agent reasoning & memory | Agent frameworks (CrewAI, LangGraph, etc.) |
| Agent execution          | Runners                                    |
| Agent checkpointing      | Agent frameworks (native persistence)      |
| Policy authoring         | Platform / Risk / Compliance               |
| Policy evaluation        | Pluggable policy engine                    |
| Execution control        | Gantral                                    |
| Approval decisions       | Humans                                     |
| Audit & review           | Independent audit teams                    |

No component may cross these boundaries.

---

## 3. Execution Engine (Phase 1)

### Responsibility

- Enforce canonical authority states
- Persist immutable execution events
- Reject invalid or out-of-order transitions

### Canonical States

```

CREATED → RUNNING → WAITING_FOR_HUMAN → APPROVED | REJECTED | OVERRIDDEN
↓                ↓
RESUMED         TERMINATED → COMPLETED

```

### Implementation Guidance

- Event-sourced execution model
- Append-only execution history
- Idempotent transition handlers
- No implicit or hidden state mutation

> **AI assistant instruction:** Implement the state machine exactly as specified. Do not add states.

---

## 4. HITL State Machine (Phase 1)

### Responsibility

- Encode human authority as an execution state
- Prevent agents from resuming execution autonomously

### Required Guarantees

Each HITL decision must capture:

- `human_actor_id`
- `role`
- `timestamp`
- `decision` (APPROVE / REJECT / OVERRIDE)
- `context_snapshot`
- `justification`

### Override Semantics (Critical)

For **OVERRIDDEN** (and optionally APPROVED) decisions:

- `context_delta` (or equivalent override payload) is mandatory
- Execution resumes **only** with modified context

### Timeout Handling (Critical)

- HITL waits are subject to policy-defined timeouts
- On timeout expiry:
  - Emit a `TIMEOUT` execution event
  - Apply policy-defined behavior (fail-closed or escalate)

Timeout handling is **execution control**, not policy evaluation.

---

## 5. Policy Evaluation Layer (Phase 2)

### Core Principle

Policy evaluation is a **transition guard**, not an execution state and not an authority.

### When It Runs

- Synchronously during execution transitions
- Before entering `WAITING_FOR_HUMAN`

### What It Produces (Advisory Signals)

- `ALLOW`
- `REQUIRE_HUMAN`
- `DENY`

Policy results are **never persisted** as durable execution state.

---

## 6. Failure Modes and Atomicity (Critical)

### Purpose

This section defines **when execution must stop**, **when evidence must not be emitted**, and **how atomicity is enforced**.

Gantral prefers **failure over ambiguity**.  
Execution must never proceed if authority cannot be **enforced and proven atomically**.

---

### 6.1 Atomic Authority Commit

An authority transition is valid **only if** all of the following occur atomically:

1. Authority decision is evaluated
2. Execution state transition is committed
3. Commitment artifact is emitted and durably persisted

These steps are inseparable.

> **Rule:**  
> If a commitment artifact cannot be emitted, the execution transition **must not occur**.

No artifact → no execution.

---

### 6.2 Execution-Time Failure Modes

Gantral must refuse or terminate execution when:

- Authority state is ambiguous
- Required human authority is unavailable
- Policy evaluation produces conflicting signals
- HITL timeout expires without resolution
- Execution transition is invalid or out of order
- Commitment artifact emission fails
- Persistence guarantees cannot be met

Failure to proceed is **correct behavior**.

---

### 6.3 Artifact Emission Failures

If artifact emission fails due to:

- Storage unavailability
- Integrity binding failure
- Partial write or timeout
- Workflow runtime error

Gantral must:

- Abort the execution transition
- Mark the instance as `TERMINATED`
- Surface the failure explicitly

Execution must never continue optimistically.

---

### 6.4 Replay-Oriented Failure Semantics

Failure modes are designed to support **hostile replay**:

| Condition             | Replay Outcome |
|----------------------|----------------|
| Missing artifact     | INCONCLUSIVE   |
| Altered artifact     | INVALID        |
| Ambiguous authority  | INCONCLUSIVE   |
| Invalid transition   | INVALID        |

Failures must be **detectable during replay**.

---

### 6.5 No Best-Effort Execution

Gantral must not:

- Retry authority transitions silently
- Infer missing authority
- Reconstruct approvals from logs
- Rely on operator narratives

Best-effort execution is incompatible with verifiability.

---

## 7. Policy Evaluation Interface (Contract)

### Input (TRD-Aligned)

- `instance_id`
- `workflow_id`
- `workflow_version`
- `materiality`
- `current_state`
- `actor_id`
- `roles`
- `execution_context`
- `policy_version_id`
- `dry_run`

### Output (Advisory Only)

- `decision` (ALLOW / REQUIRE_HUMAN / DENY)
- `eligible_approver_roles`
- `escalation_roles`
- `timeout`

Gantral remains the **sole enforcement authority**.

---

## 8. Execution History, Audit & Replay

### Responsibility

- Store authoritative execution records
- Enable deterministic replay of authority decisions
- Support regulatory audits

### Determinism Model

- Gantral workflows execute inside a deterministic workflow runtime
- Replay re-executes workflow logic from recorded authority history
- Replay does **not** rehydrate agent internal memory

Agent memory replay is optional and **non-authoritative**.

---

## 9. APIs & SDKs

### API Principles

- Explicit state transitions
- Versioned endpoints
- No implicit side effects

### SDKs

- Thin wrappers only
- No business logic
- Server remains the source of truth

---

## 10. Adapters & Framework-Native Integrations

### Responsibility

- Translate external events into Gantral triggers
- Forward execution decisions back to external systems

Adapters contain **no business logic**, **no policy checks**, and **no execution authority**.

### 10.1 Framework-Native Agent Integrations (Critical)

Runners and SDKs must support **agent-native suspension and resume semantics**.

| Agent Outcome | Meaning                       | Gantral Action      |
|--------------|-------------------------------|---------------------|
| Completed    | Agent finished normally       | `COMPLETED`         |
| Failed       | Agent crashed or errored      | `TERMINATED`        |
| Suspended    | Agent checkpointed and exited | `WAITING_FOR_HUMAN` |

Gantral never loads, inspects, or persists agent internal state.

---

## 11. Security, Identity & Federation

### 11.1 Human Identity

- OIDC federation only
- Identity derived from token claims
- No Gantral-managed users

### 11.2 Machine Identity

- Workload identity only
- No static API keys

---

## 12. Runner Pattern (Federated Execution)

- Runners execute inside team-owned infrastructure
- Pull-based task queues
- Data locality preserved
- Network isolation respected

**Gantral is the authority. Runners are the executors.**

---

## 13. Deployment Model

- Kubernetes (primary)
- Self-hosted
- Air-gapped supported

---

## 14. Implementation Phases

- Phase 1 – Control Foundations
- Phase 2 – Governance Hardening
- Phase 3 – Enterprise Integration
- Phase 4 – Federated Execution

Phases must be completed **strictly in order**.

---

## 15. What This Guide Protects Against

- Accidental autonomy
- Agent memory leakage into audit logs
- Non-deterministic approvals
- Optimistic execution
- Regulator rejection

---

## 16. Final Reminder

Gantral is not about what AI can do.

It is about what organizations are willing to  
**allow AI to do — and how they prove it**.

This guide exists to ensure Gantral is built **correctly, safely, and credibly**.
