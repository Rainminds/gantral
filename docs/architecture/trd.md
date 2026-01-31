---
title: Technical Reference Document
---

# Gantral – Technical Reference & Architecture Document

**Version:** v4.1 (Authority-First, Verifiable Execution, Atomic Evidence)  
**Status:** Living technical reference (authoritative)  
**Audience:** Core contributors, platform engineers, security reviewers, enterprise architects  

---

## Purpose

This document is the **technical constitution** of the Gantral open-source core.

It defines the **architectural invariants, execution semantics, authority boundaries,
and verifiability guarantees** that must hold across all implementations.

Gantral deliberately separates **execution authority** from **agent reasoning and memory**.

If an implementation conflicts with this document, the implementation is incorrect.

---

## 1. Scope & Non-Goals

### 1.1 What Gantral Is (Technical Definition)

Gantral is an **AI Execution Control Plane**.

It enforces how AI-assisted and agentic workflows:

- execute
- pause at governed boundaries
- require and record authority
- resume, override, or terminate
- emit verifiable execution evidence

Technically, Gantral provides:

- A **deterministic, instance-first authority state machine**
- Human-in-the-Loop (HITL) as a **blocking execution state**
- Instance-level isolation for audit, cost, and accountability
- Declarative control over materiality, authority, and escalation
- A **pluggable policy evaluation interface (advisory only)**
- Commitment artifacts emitted at authority boundaries
- Deterministic replay of **authority decisions only**
- Control APIs and SDKs that sit above agent frameworks and below enterprise systems

Gantral is domain-agnostic by design.  
SDLC workflows are an initial wedge, not the boundary.

---

### 1.2 Explicit Non-Goals (Hard Exclusions)

Gantral will not:

- build, host, or orchestrate AI agents
- encode domain-specific business logic
- manage agent prompts, plans, or internal memory
- serialize or persist agent internal state
- optimize, train, or fine-tune models
- provide autonomous decision loops
- make probabilistic decisions without human authority
- replace CI/CD, ITSM, ticketing, or enterprise systems
- store raw secrets or credentials
- act as an identity provider

---

## 2. Architectural Invariants (Non-Negotiable)

The following invariants must hold across all implementations:

- **Instance-first semantics**  
  All audit, cost, authority, and replay semantics attach to immutable execution instances.

- **Authority is state**  
  Execution authority exists only as execution state transitions.

- **HITL is a state transition**  
  Human intervention is enforced inside the execution graph, not externally.

- **Human authority is final**  
  AI output and policy evaluation are advisory only.

- **Determinism over performance**  
  Replayability and auditability take precedence over latency.

- **Declarative control**  
  Materiality, authority, and escalation rules are configuration, not embedded logic.

- **Adapters contain no business logic**  
  Integrations emit events and receive decisions only.

- **Identity federation required**  
  Gantral must not maintain a standalone user directory.

- **No secret persistence**  
  Gantral stores only secret references; resolution occurs at execution edges.

Violating these invariants invalidates the system.

---

### New Invariant: Agent-Native Persistence (Critical)

- **Execution State vs Agent State Separation**
  - Gantral owns execution and authority state
  - Agent frameworks own agent internal state (memory, plans, tool context)

- **Checkpointability Requirement**
  - Agents must be restartable from externally persisted checkpoints
  - Agent state persistence is the agent framework’s responsibility

- **No Agent State in Gantral**
  - Agent internal state must never be serialized into Gantral execution history

---

## 3. High-Level Architecture

### 3.1 Logical Architecture Layers

```

+--------------------------------------------------+
| Enterprise Systems                               |
| (GitHub, Jira, Slack, ServiceNow, Custom Apps)   |
+------------------------▲-------------------------+
| Events / Decisions
+------------------------|-------------------------+
| Gantral Control Plane                          |
| - Execution Engine                              |
| - Authority State Machine                       |
| - Policy Interface (Advisory)                   |
| - Instance Registry                             |
| - Commitment Artifact Emission                  |
+------------------------▲-------------------------+
| SDK / API
+------------------------|-------------------------+
| Agent Frameworks & Runners (Distributed)        |
+--------------------------------------------------+

```

---

### 3.2 Runtime Model (Authoritative)

1. Trigger received (event, schedule, external signal)  
2. Workflow template selected  
3. Immutable execution instance created  
4. Execution proceeds in deterministic workflow runtime  
5. Policy evaluated as a **transition guard**  
6. Execution continues or enters `WAITING_FOR_HUMAN`  
7. Agent framework checkpoints and suspends  
8. Human authority decision captured (if required)  
9. Authority state transition committed  
10. **Commitment artifact emitted atomically**  
11. Execution resumes or terminates via new process  

---

### 3.3 Execution Plane Responsibility Boundary

**Gantral owns**
- execution and authority state
- authority transitions
- time, retries, escalation, replay
- audit correctness
- artifact emission

**Agent frameworks own**
- reasoning, planning, tool execution
- agent memory and checkpoints

**Runners**
- execute agent processes
- translate lifecycle signals into Gantral events

Agents cannot advance execution past governed states independently.

---

### 3.4 Execution State vs Execution Readiness

Execution **state** is canonical and auditable.  
Execution **readiness** is transient and derived.

Formally:

```

runnable(instance) =
execution_state == RUNNING
AND no pending HITL decisions
AND no policy blocks
AND execution plane capacity available

```

Readiness is never persisted.

---

## 4. Core Domain Model

### 4.1 Core Entities

**Workflow (Template)**  
- workflow_id, version  
- step graph  
- trigger definitions  
- policy references  
- materiality level  

**Execution Instance**  
- immutable instance_id  
- workflow_id + version  
- owning_team_id  
- trigger context  
- execution state  
- timestamps  
- cost metadata  

**HITL Decision**  
- decision_id  
- instance_id  
- decision_type (APPROVE / REJECT / OVERRIDE)  
- human_actor_id  
- role  
- justification  
- context_snapshot  
- context_delta (required for OVERRIDE)  

---

## 5. Authority State Machine

### 5.1 Canonical States

```

CREATED
↓
RUNNING
↓ (policy requires authority)
WAITING_FOR_HUMAN
↙      ↓      ↘
OVERRIDDEN APPROVED REJECTED
↓        ↓        ↓
RESUMED   RESUMED  TERMINATED
↓
COMPLETED

```

---

### 5.2 State Guarantees

- append-only transitions
- no in-place mutation
- timestamped and attributable
- authority transitions always produce artifacts

---

### 5.3 Policy Evaluation Semantics (Critical)

- Policy evaluation introduces **no new execution state**
- Evaluation occurs synchronously during transitions
- Results determine **whether authority is required**
- Policy engines never hold authority
- Policy evaluation must not depend on tool payloads

Policy checks are **transition guards**, not durable state.

---

## 6. Policy Interface & Evaluation Layer

Policy evaluators provide:
- ALLOW / REQUIRE_HUMAN / DENY
- approver eligibility
- escalation and timeout signals

Gantral enforces:
- authority transitions
- execution control
- audit semantics

---

## 7. Determinism & Replay Model

Gantral guarantees deterministic replay of **authority decisions**.

### 7.1 Replay Semantics

- Replay re-executes workflow logic from recorded authority history
- Identical inputs and versions yield identical authority outcomes
- Replay **never** rehydrates agent internal memory
- Replay requires no access to Gantral services

### 7.2 Responsibility Boundary

- Gantral defines *what* must be replayable
- Workflow runtime defines *how* replay occurs
- Agent replay is optional and non-authoritative

---

## 8. APIs & SDKs

- REST and gRPC APIs (OpenAPI 3.1)
- Core API groups:
  - /workflows
  - /instances
  - /decisions
  - /policies
  - /audit

SDKs are thin wrappers. All authority logic is server-side.

---

## 9. Data & Storage Model

| Purpose | Approach |
|------|--------|
| Execution state | Event-sourced workflow history |
| Authority artifacts | Immutable append-only store |
| Snapshots | Object storage |
| Caching | Optional (non-authoritative) |

Agent internal state stores are explicitly out of scope.

---

## 10. Security Architecture

- Identity: OAuth 2.0 / OIDC (federated)
- Authorization: policy-driven
- Secrets: external secret managers only
- Audit: 100% authority capture, tamper-evident artifacts

---

## 11. Privacy & Compliance

- data minimization by default
- configurable redaction
- self-hosted and air-gapped supported

Compliance is a **structural property of the execution model**, not a checklist.

---

## 12. Standards & Interoperability

- OpenAPI 3.1
- OpenTelemetry
- OAuth 2.0 / OIDC
- CNCF-aligned design principles

---

## 13. Execution Evidence References (Non-Authoritative)

Gantral may store **references** to external evidence.

- evidence captured outside Gantral
- Gantral stores hashes or pointers only
- Gantral does not inspect or interpret evidence
- authority remains governed solely by state transitions

---

## 14. Final Principle

Gantral is not about what AI can do.

It is about what organizations are willing to **allow AI to do — and how they prove it**.

This document is authoritative.
