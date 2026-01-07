---

sidebar_position: 3
title: Implementation Guide
-------------------------------------

# Comprehensive Implementation Guide

**Version:** v4.0 (Authority-First, Federated Execution, Agent-Native Persistence)
**Status:** Reference implementation guide (non-normative)
**Audience:** Core contributors, platform engineers, enterprise architects (regulated environments), security / risk / compliance reviewers, AI coding assistants

---

## Purpose

This document translates the **Gantral PRD and TRD v4** into a **buildable, implementation-accurate guide**.

It incorporates the finalized architectural clarifications:

* Gantral owns **execution authority and control state**
* Agent frameworks own **agent internal state and memory**
* Long-running approvals require **agent-native checkpointing**
* Deterministic execution and replay are guaranteed by a **workflow runtime (e.g. Temporal)**

> **Rule:** If this guide conflicts with the TRD, the TRD is correct.

---

## 0. How This Guide Should Be Used

### 0.1 Documentation Hierarchy (Non-Negotiable)

```
PRD → TRD → Implementation Guide (this document)
```

* PRD defines **why** Gantral exists
* TRD defines **what must be true**
* This guide defines **one correct way to build it**

This guide is replaceable. The TRD is not.

---

## 1. System Responsibilities (Explicit)

### 1.1 What Gantral Owns

Gantral is an **execution authority layer**. It owns:

* Execution lifecycle and immutable instance boundaries
* Canonical execution state machine
* Human-in-the-Loop (HITL) enforcement
* Authority transitions (APPROVE / REJECT / OVERRIDE)
* Policy gating as **transition guards**
* Time, retries, escalation, and timeout semantics
* Instance-level audit correctness
* Deterministic replay of execution decisions (via workflow runtime)

### 1.2 What Gantral Explicitly Does NOT Own

Gantral does **not**:

* Perform agent reasoning, planning, or cognition
* Manage prompts, tools, or agent memory
* Persist or serialize agent internal state
* Decide outcomes autonomously
* Store secrets or credentials
* Act as an identity provider

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

* Enforce canonical execution states
* Persist immutable execution events
* Reject invalid or out-of-order transitions

### Canonical States

```
CREATED → RUNNING → WAITING_FOR_HUMAN → APPROVED | REJECTED | OVERRIDDEN
                     ↓                ↓
                  RESUMED         TERMINATED → COMPLETED
```

### Implementation Guidance

* Event-sourced execution model
* Append-only execution history
* Idempotent transition handlers
* No implicit or hidden state mutation

> **AI assistant instruction:** Implement the state machine exactly as specified. Do not add states.

---

## 4. HITL State Machine (Phase 1)

### Responsibility

* Encode human authority as an execution state
* Prevent agents from resuming execution autonomously

### Required Guarantees

Each HITL decision must capture:

* `human_actor_id`
* `role`
* `timestamp`
* `decision` (APPROVE / REJECT / OVERRIDE)
* `context_snapshot`
* `justification`

### Override Semantics (Critical)

For **OVERRIDDEN** (and optionally APPROVED) decisions:

* `context_delta` (or equivalent override payload) is mandatory
* Execution resumes **only** with modified context

### Timeout Handling (Critical)

* HITL waits are subject to policy-defined timeouts
* On timeout expiry:

  * Emit a `TIMEOUT` execution event
  * Apply policy-defined behavior (fail-closed or escalate)

Timeout handling is **execution control**, not policy evaluation.

---

## 5. Policy Evaluation Layer (Phase 2)

### Core Principle

Policy evaluation is a **transition guard**, not an execution state and not an authority.

### When It Runs

* Synchronously during execution transitions
* Before entering `WAITING_FOR_HUMAN`

### What It Produces (Advisory Signals)

* `ALLOW`
* `REQUIRE_HUMAN`
* `DENY`

Policy results are **never persisted** as durable execution state.

---

## 6. Policy Evaluation Interface (Contract)

### Input (TRD-Aligned)

* `instance_id`
* `workflow_id`
* `workflow_version`
* `materiality`
* `current_state`
* `actor_id`
* `roles`
* `execution_context`
* `policy_version_id`
* `dry_run`

### Output (Advisory Only)

* `decision` (ALLOW / REQUIRE_HUMAN / DENY)
* `eligible_approver_roles`
* `escalation_roles`
* `timeout`

Gantral remains the **sole enforcement authority**.

---

## 7. Execution History, Audit & Replay

### Responsibility

* Store authoritative execution records
* Enable deterministic replay of authority decisions
* Support regulatory audits

### Determinism Model

* Gantral workflows execute inside a deterministic workflow runtime
* Replay re-executes workflow logic from recorded execution history
* Replay does **not** rehydrate agent internal memory

Agent memory replay is optional and **non-authoritative**.

---

## 8. APIs & SDKs

### API Principles

* Explicit state transitions
* Versioned endpoints
* No implicit side effects

### SDKs

* Thin wrappers only
* No business logic
* Server remains the source of truth

---

## 9. Adapters & Framework-Native Integrations

### Responsibility

* Translate external events into Gantral triggers
* Forward execution decisions back to external systems

Adapters contain **no business logic**, **no policy checks**, and **no execution authority**.

### 9.1 Framework-Native Agent Integrations (Critical)

Runners and SDKs must support **agent-native suspension and resume semantics**.

Agent frameworks may deliberately suspend execution when:

* HITL approval is required
* Execution must pause for an extended duration

#### Required Runner / SDK Behavior

| Agent Outcome | Meaning                       | Gantral Action                    |
| ------------- | ----------------------------- | --------------------------------- |
| Completed     | Agent finished normally       | Transition to `COMPLETED`         |
| Failed        | Agent crashed or errored      | Transition to `TERMINATED`        |
| Suspended     | Agent checkpointed and exited | Transition to `WAITING_FOR_HUMAN` |

#### Suspension Handling

* Agent framework persists internal state using native checkpointing
* Agent process exits cleanly
* Runner emits a `SUSPENDED` execution event
* Gantral workflow waits on approval signal

#### Resume Handling

Upon APPROVED / OVERRIDDEN:

* Gantral schedules a new execution task
* Runner launches a new agent process
* Agent framework restores state from checkpoint
* Execution resumes deterministically

Gantral never loads, inspects, or persists agent internal state.

---

## 10. Security, Identity & Federation

### 10.1 Human Identity

* OIDC federation only (Okta, Azure AD, etc.)
* Identity derived from token claims
* No Gantral-managed users

### 10.2 Machine Identity

* Workload identity (AWS IAM, Kubernetes Service Accounts, etc.)
* No static API keys

### 10.3 Secrets & Connections

* Gantral stores references only (e.g. `vault://...`)
* Secrets resolved by Runners at execution edge
* Raw credentials never enter Gantral

---

## 11. Runner Pattern (Federated Execution)

* Runners execute inside team-owned infrastructure
* Pull-based task queues
* Data locality preserved
* Network isolation respected

**Gantral is the authority. Runners are the executors.**

---

## 12. Optional Tool Mediation & Evidence Capture (Runner-Side)

Runners may optionally mediate tool execution to:

* Capture tool inputs and outputs
* Store evidence externally
* Generate immutable evidence references

### Constraints

* Tool mediation must not block tool execution
* Gantral must not inspect payloads
* Policies may require **presence of evidence**, not its contents

### Non-Goals

* Must not authorize or deny tool execution
* Must not introduce new execution states
* Must not pass raw tool payloads into Gantral APIs
* Must not modify agent reasoning behavior

---

## 13. Deployment Model

* Kubernetes (primary)
* Self-hosted
* Air-gapped deployments supported

---

## 14. Implementation Phases

* Phase 1 – Control Foundations
* Phase 2 – Governance Hardening
* Phase 3 – Enterprise Integration
* Phase 4 – Federated Execution

Phases must be completed **strictly in order**.

---

## 15. What This Guide Protects Against

* Accidental autonomy
* Agent memory leakage into audit logs
* Non-deterministic approvals
* Regulator rejection
* Framework lock-in

---

## 16. Final Reminder

Gantral is not about what AI can do.

It is about what organizations are willing to **allow AI to do — and how they prove it**.

This guide exists to ensure Gantral is built **correctly, safely, and credibly**.
