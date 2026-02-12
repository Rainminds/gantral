---
title: Public Product Requirements Document (PRD)
---

**Version:** v5.1 (Verifiability & Authority Clarified)  
**Status:** Public, canonical product document  
**Audience:** Enterprises, platform engineers, contributors, regulators, partners  

---

## 1. Purpose of This Document

This document defines the **public product boundaries, responsibilities, and intent** of:

- **Gantral** — the open-core AI Execution Control Plane  
- **Gantrio** — the commercial enterprise control experience built on top of Gantral  

It exists to:

- Prevent category confusion (e.g. “agent platform”, “workflow engine”)
- Make the **Gantral vs Gantrio split unambiguous**
- Serve as a guardrail against scope creep
- Provide a regulator- and auditor-aware explanation of **execution-time authority**
- Define what Gantral can and cannot **prove after execution**

If an implementation or product claim contradicts this document, it is incorrect.

---

## 2. The Core Problem (Non-Negotiable)

Large organizations are adopting AI across SDLC, operations, finance, compliance, support, and internal workflows.

What breaks is **not model capability**.

What breaks is:

- Execution control  
- Human authority  
- Accountability  
- Auditability at organizational scale  

### What Happens in Practice

- AI usage is fragmented and team-driven  
- Human review exists, but only as behavior (Slack, email, checklists)  
- Approvals are informal, inconsistent, and non-enforceable  
- No system of record exists for *who allowed what to run, and when*

Organizations do not lack AI tools.  
They lack an **AI execution control plane**.

---

## 3. The Mental Model

The system is composed of **two strictly separated layers**:

- **Gantral** decides **whether execution may proceed**
- **Gantrio** helps humans **see, manage, and operate those decisions**

> **Gantral decides. Gantrio shows and manages.**

This separation is **structural, not cosmetic** and is required for verifiability.

---

## 4. What Gantral Is (Open Core)

**Gantral is an AI Execution Control Plane.**

It is infrastructure that enforces how AI-assisted and agentic workflows:

- start  
- pause for required authority  
- resume, override, or terminate  
- emit verifiable execution evidence  

Gantral sits **above agent frameworks** and **below enterprise systems**.

### Gantral Owns

- Canonical **authority state machine**
- Immutable execution instances
- Human-in-the-Loop (HITL) as a blocking execution state
- Authority transitions (approve / reject / override)
- Policy **evaluation interfaces** (advisory only)
- **Commitment artifacts** emitted at authority boundaries
- Deterministic replay of authority decisions
- Control APIs and SDKs

### Gantral Explicitly Does NOT Own

- User interfaces
- Org or team modeling
- RBAC UX
- Policy authoring tools
- Workflow builders or editors
- Approval inboxes or notifications
- Cost dashboards or optimization logic
- Integrations UX

Gantral **enforces authority**.  
It does not visualize, explain, or manage it for humans.

---

## 5. Determinism, Evidence, and Failure

When execution authority is exercised, Gantral emits a **commitment artifact**.

Human authority transitions must capture attributable reasoning sufficient to demonstrate active judgement, not merely procedural approval.

This artifact:

- is emitted **atomically** with the authority state transition
- binds execution context references, authority decision, and time
- is replayable without access to Gantral systems
- does not depend on agent memory, logs, or narrative reconstruction

If authority cannot be enforced **and** recorded:

> **Execution must not proceed.**

Commitment artifacts are structured to contain decision-grade evidence only. Telemetry, logs, and agent reasoning traces are out-of-scope and non-authoritative.

Gantral is designed to **fail closed**.
Ambiguity results in refusal or inconclusive outcomes, not best-effort execution.

---

## 6. Policy vs Authority

Policy engines (e.g. OPA):

- evaluate conditions
- advise escalation or denial
- do **not** grant execution authority

Authority:

- exists only as execution state
- is enforced by Gantral
- is the sole mechanism by which execution proceeds

Policy advises.  
Authority enforces.

This separation is required for replay and third-party verification.

---

## 7. What Gantrio Is (Commercial Platform)

**Gantrio is the enterprise control experience for Gantral.**

Gantrio exists because organizations cannot operate execution authority
using APIs and logs alone.

Gantrio provides:

- Enterprise UI
- Approval inboxes and escalation views
- Org, team, and role modeling
- RBAC and delegation
- Policy authoring and lifecycle UX
- Compliance reporting and exports
- Cost attribution and usage visibility
- Integrations (Jira, GitHub, ServiceNow, Slack, etc.)
- Managed hosting and enterprise support

Gantrio **never enforces execution authority**.

If Gantrio is unavailable, Gantral must still behave correctly.

---

## 8. Gantral vs Gantrio — Feature Boundary Table

| Capability | Gantral (Open Core) | Gantrio (Commercial) |
|---------|------------------|-------------------|
| Authority state machine | ✅ | ❌ |
| HITL enforcement | ✅ | ❌ |
| Authority transitions | ✅ | ❌ |
| Commitment artifacts | ✅ | ❌ |
| Deterministic replay | ✅ | ❌ |
| Policy evaluation interface | ✅ | ❌ |
| Policy authoring UX | ❌ | ✅ |
| Approval inboxes | ❌ | ✅ |
| Org / team modeling | ❌ | ✅ |
| RBAC & delegation | ❌ | ✅ |
| Compliance reports | ❌ | ✅ |
| Cost dashboards | ❌ | ✅ |
| Integrations UI | ❌ | ✅ |
| Managed hosting | ❌ | ✅ |

This boundary is **intentional and enforced**.

---

## 9. What Gantral Will Explicitly Never Become

To avoid erosion of guarantees, Gantral will never:

- build or host agents
- provide workflow builders
- optimize or route models
- make autonomous decisions
- store agent memory or prompts
- act as an identity provider

These are **structural non-goals**.

---

## 10. What Changes After Adoption

**Before:**

- approvals outside execution
- governance by discipline
- audit reconstructed manually

**After:**

- execution pauses for authority
- decisions enforced technically
- audit trails are native and replayable

Governance becomes a **by-product of execution**.

---

## 11. Final Reminder

Gantral is not about what AI *can* do.

It is about what organizations are willing to  
**allow AI to do — and how they prove it**.

Gantrio exists to make that control usable at enterprise scale.