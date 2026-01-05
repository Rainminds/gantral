# Comprehensive Implementation Guide

**Version:** v4.0 (Authority-First, Federated Execution, Agent-Native Persistence)

**Status:** Reference implementation guide (non-normative)

**Audience:**

* Core contributors  
* Platform engineers  
* Enterprise architects (regulated environments)  
* Security, Risk, and Compliance reviewers  
* AI coding assistants used for implementation

---

## **Purpose**

This document translates the Gantral **PRD** and **TRD v4** into a **buildable, implementation-accurate guide**.

This version incorporates the finalized architectural clarification:

* Gantral owns **execution authority and control state**  
* Agent frameworks own **agent internal state and memory**  
* Long-running approvals require **agent-native checkpointing**  
* Deterministic execution and replay are guaranteed by a workflow runtime (Temporal)

**Rule:** If this guide conflicts with the TRD, the TRD is correct.

---

## **0\. How This Guide Should Be Used**

### **0.1 Documentation Hierarchy (Non-Negotiable)**

**PRD → TRD → Implementation Guide (this document)**

* PRD defines *why* Gantral exists  
* TRD defines *what must be true*  
* This guide defines *one correct way to build it*

This guide is replaceable. The TRD is not.

---

## **1\. System Responsibilities (Explicit)**

### **1.1 What Gantral Owns**

Gantral is an **execution authority layer**. It owns:

* Execution lifecycle and immutable instance boundaries  
* Canonical execution state machine  
* Human-in-the-Loop (HITL) enforcement  
* Authority transitions (approve / reject / override)  
* Policy gating as **transition guards**  
* Time, retries, escalation, and timeout semantics  
* Instance-level audit correctness  
* Deterministic replay of execution decisions (via workflow runtime)

### **1.2 What Gantral Explicitly Does NOT Own**

Gantral does **not**:

* Perform agent reasoning, planning, or cognition  
* Manage prompts, tools, or agent memory  
* Persist or serialize agent internal state  
* Decide outcomes autonomously  
* Store secrets or credentials  
* Act as an identity provider

Gantral governs **whether execution may proceed**, not **how agents reason**.

---

## **2\. Separation of Duties (Enterprise Requirement)**

| Function | Owner |
| ----- | ----- |
| Agent reasoning & memory | Agent frameworks (CrewAI, LangGraph, etc.) |
| Agent execution | Runners |
| Agent checkpointing | Agent frameworks (native persistence) |
| Policy authoring | Platform / Risk / Compliance |
| Policy evaluation | Pluggable policy engine |
| Execution control | Gantral |
| Approval decisions | Humans |
| Audit & review | Independent audit teams |

No component may cross these boundaries.

---

## **3\. Execution Engine (Phase 1\)**

### **Responsibility**

* Enforce canonical execution states  
* Persist immutable execution events  
* Reject invalid or out-of-order transitions

### **Canonical States**

CREATED → RUNNING → WAITING\_FOR\_HUMAN → APPROVED | REJECTED | OVERRIDDEN  
                     ↓                    ↓  
                  RESUMED             TERMINATED → COMPLETED

### **Implementation Guidance**

* Event-sourced execution model  
* Append-only execution history  
* Idempotent transition handlers  
* No implicit or hidden state mutation

**AI assistant instruction:** Implement the state machine exactly as specified. Do not add states.

---

## **4\. HITL State Machine (Phase 1\)**

### **Responsibility**

* Encode human authority as an execution state  
* Prevent agents from resuming execution autonomously

### **Required Guarantees**

Each HITL decision must capture:

* `human_actor_id`  
* `role`  
* `timestamp`  
* `decision` (APPROVE / REJECT / OVERRIDE)  
* `context_snapshot`  
* `justification`

### **Override Semantics (Critical)**

For OVERRIDDEN (and optionally APPROVED) decisions:

* `context_delta` or equivalent override payload is mandatory  
* Execution resumes only with modified context

### **Timeout Handling (Critical)**

* HITL waits are subject to policy-defined timeouts  
* On timeout expiry:  
  * Emit a TIMEOUT execution event  
  * Apply policy-defined behavior (fail-closed or escalate)

Timeout handling is execution control, not policy evaluation.

---

## **5\. Policy Evaluation Layer (Phase 2\)**

### **Core Principle**

Policy evaluation is a **transition guard**, not an execution state and not an authority.

### **When It Runs**

* Synchronously during execution transitions  
* Before entering `WAITING_FOR_HUMAN`

### **What It Produces (Advisory Signals)**

* ALLOW  
* REQUIRE\_HUMAN  
* DENY

Policy results are **never persisted as durable execution state**.

---

## **6\. Policy Evaluation Interface (Contract)**

### **Input (TRD-Aligned)**

* instance\_id  
* workflow\_id  
* workflow\_version  
* materiality  
* current\_state  
* actor\_id  
* roles  
* execution\_context  
* policy\_version\_id  
* dry\_run

### **Output (Advisory Only)**

* decision (ALLOW / REQUIRE\_HUMAN / DENY)  
* eligible\_approver\_roles  
* escalation\_roles  
* timeout

Gantral remains the **sole enforcement authority**.

---

## **7\. Execution History, Audit & Replay**

### **Responsibility**

* Store authoritative execution records  
* Enable deterministic replay of **authority decisions**  
* Support regulatory audits

### **Determinism Model**

* Gantral workflows execute inside a deterministic workflow runtime (Temporal)  
* Replay re-executes workflow logic from recorded execution history  
* Replay does **not** rehydrate agent internal memory

Agent memory replay is optional and non-authoritative.

---

## **8\. APIs & SDKs**

### **API Principles**

* Explicit state transitions  
* Versioned endpoints  
* No implicit side effects

### **SDKs**

* Thin wrappers only  
* No business logic  
* Server remains the source of truth

---

## **9\. Adapters & Framework-Native Integrations (Updated)**

### **Responsibility**

* Translate external events into Gantral triggers  
* Forward execution decisions back to external systems

Adapters contain **no business logic, no policy checks, and no execution authority**.

### **9.1 Framework-Native Agent Integrations (Critical)**

Runners and SDKs **must support agent-native suspension and resume semantics**.

Agent frameworks may deliberately suspend execution ("hibernation") when:

* A HITL approval is required  
* Execution must pause for an extended duration

### **Required Runner / SDK Behavior**

The SDK **must distinguish** between:

| Agent Outcome | Meaning | Gantral Action |
| ----- | ----- | ----- |
| Completed | Agent finished normally | Transition to COMPLETED |
| Failed | Agent crashed or errored | Transition to TERMINATED |
| Suspended | Agent checkpointed and exited | Transition to WAITING\_FOR\_HUMAN |

### **Suspension Handling**

* Agent framework persists internal state using its **native checkpointing**  
* Agent process exits cleanly  
* Runner emits a **SUSPENDED** execution event  
* Gantral workflow waits on approval signal

### **Resume Handling**

* Upon APPROVED / OVERRIDDEN:  
  * Gantral schedules a new execution task  
  * Runner launches a **new agent process**  
  * Agent framework restores state from checkpoint  
  * Execution resumes deterministically

Gantral never loads, inspects, or persists agent internal state.

---

## **10\. Security, Identity & Federation**

### **10.1 Human Identity**

* OIDC federation only (Okta, Azure AD, etc.)  
* Identity derived from token claims  
* No Gantral-managed users

### **10.2 Machine Identity**

* Workload identity (AWS IAM, K8s SA, etc.)  
* No static API keys

### **10.3 Secrets & Connections**

* Gantral stores references only (e.g. `vault://...`)  
* Secrets resolved by Runners at execution edge  
* Raw credentials never enter Gantral

---

## **11\. Runner Pattern (Federated Execution)**

* Runners execute inside team-owned infrastructure  
* Pull-based task queues  
* Data locality preserved  
* Network isolation respected

Gantral is the **authority**. Runners are the **executors**.

---

## **12\. Deployment Model**

* Kubernetes (primary)  
* Self-hosted  
* Air-gapped deployments supported

---

## **13\. Implementation Phases**

* Phase 1 – Control Foundations  
* Phase 2 – Governance Hardening  
* Phase 3 – Enterprise Integration  
* Phase 4 – Federated Execution

Phases must be completed **strictly in order**.

---

## **14\. What This Guide Protects Against**

* Accidental autonomy  
* Agent memory leakage into audit logs  
* Non-deterministic approvals  
* Regulator rejection  
* Framework lock-in

---

## **15\. Final Reminder**

Gantral is not about what AI can do.

It is about what organizations are willing to allow AI to do — and how they prove it.

This guide exists to ensure Gantral is built **correctly, safely, and credibly**.

