# Gantral – Phase-wise Build Plan for AI Assistants

**Version:** v9.0 (Authority-First, Federated Execution, Agent-Native Persistence)

**Status:** Authoritative execution plan (current)

**Audience:**

* Human maintainers  
* AI coding assistants (as build agents)  
* Core contributors  
* Enterprise platform teams

---

## **Purpose**

This document defines the **authoritative, phase-wise build plan** for Gantral.

It reflects the finalized architectural decisions:

* Gantral is the **execution authority**  
* Execution is **federated** via distributed Runners  
* Identity and secrets are **federated and externalized**  
* Deterministic execution and replay are guaranteed by a workflow runtime (Temporal)  
* **Agent internal state is owned by agent frameworks, not Gantral**

Gantral:

* Defines and enforces execution authority, policy gating, and HITL semantics  
* Delegates deterministic execution and replay mechanics to the workflow substrate  
* Requires **agent-native persistence OR explicit split-agent execution** for long-running approvals

**Constitutional rule:** If this document conflicts with the TRD, the TRD prevails.

---

## **How to Use This Document with AI Assistants**

For each task:

1. Paste **only** the relevant phase \+ task section into the AI assistant  
2. Explicitly state the phase and task number  
3. Enforce the listed **Non-Goals**  
4. Validate output against **Acceptance Criteria**

**AI assistants must never be asked to work across phases.**

---

## **Global Rules for AI Assistants (Non-Negotiable)**

AI assistants must:

* Implement **only** what is specified in the current task  
* Never add, rename, or bypass execution states  
* Never embed policy logic outside the policy interface  
* Never assume autonomy or self-approval  
* Never mutate immutable execution history  
* Never change schema field names defined by the TRD  
* Never persist or serialize agent internal state inside Gantral  
* Never introduce identity stores, secrets, or execution shortcuts

Violations invalidate the output.

---

## **Phase 1 – Control Foundations (Gantral Core)**

**Goal:** Establish deterministic execution control with strict human authority.

**Status:** Complete / final hardening

### **Phase 1.1 – Canonical Execution State Machine**

**Objective**  
Implement the canonical execution state machine exactly as defined in the TRD.

**Canonical States**

* CREATED  
* RUNNING  
* WAITING\_FOR\_HUMAN  
* APPROVED  
* REJECTED  
* OVERRIDDEN  
* RESUMED  
* TERMINATED  
* COMPLETED

**Acceptance Criteria**

* Invalid transitions are deterministically rejected  
* All valid transitions emit immutable execution events  
* State transitions are idempotent

---

### **Phase 1.2 – Instance Model & Execution History**

**Objective**  
Persist authoritative execution history suitable for deterministic workflow replay.

**Acceptance Criteria**

* Every state transition emits exactly one immutable event  
* Execution history is append-only  
* History contains all information required for replay by the workflow runtime

---

### **Phase 1.3 – HITL State & Decision Capture**

**Objective**  
Encode Human-in-the-Loop (HITL) as a first-class execution state.

**Acceptance Criteria**

* Execution cannot resume without a valid HITL decision  
* Overrides resume execution only with modified context (`context_delta` required)  
* Human identity, role, timestamp, and justification are always captured

---

## **Phase 2 – Governance Hardening**

**Goal:** Introduce policy evaluation, timeouts, escalation, failure handling, and audit semantics without altering execution semantics.

**Status:** Complete / final hardening

### **Phase 2.1 – Policy Evaluation Interface**

**Acceptance Criteria**

* Input schema includes: instance\_id, workflow\_id, workflow\_version, materiality, current\_state, actor\_id, roles, policy\_version\_id, dry\_run  
* Output schema includes: decision (ALLOW / REQUIRE\_HUMAN / DENY), approver\_roles, timeout, escalation\_roles  
* Interface is deterministic and side-effect free

---

### **Phase 2.2 – Policy Evaluation as Transition Guard**

**Acceptance Criteria**

* ALLOW → execution continues  
* REQUIRE\_HUMAN → transition to WAITING\_FOR\_HUMAN  
* DENY → transition to TERMINATED  
* Presence of timeout schedules a TIMEOUT event within the workflow  
* On timeout expiry:  
  * Fail-closed → TERMINATED  
  * Escalate → update approver roles and remain WAITING\_FOR\_HUMAN

---

### **Phase 2.3 – Audit Semantics (Replay-Aware)**

**Acceptance Criteria**

* Every decision is recorded with the applicable `policy_version_id`  
* Execution history is sufficient to deterministically replay the workflow

---

## **Phase 3 – Enterprise Integration (Temporal-backed)**

**Goal:** Integrate Gantral safely into enterprise environments using a deterministic workflow runtime.

**Status:** In progress / finalization

### **Phase 3.1 – Workflow Runtime Integration**

**Acceptance Criteria**

* All execution flows run inside Temporal workflows  
* Workflow code is deterministic and replay-safe  
* Gantral does not implement a separate replay engine

---

### **Phase 3.2 – Adapters Framework**

**Acceptance Criteria**

* Adapters emit events only  
* Adapters receive decisions only  
* No business logic or policy checks in adapters

---

### **Phase 3.3 – SDKs (Thin Wrappers)**

**Acceptance Criteria**

* SDKs map 1:1 to Gantral APIs  
* SDKs distinguish **Agent Failed** vs **Agent Suspended** outcomes  
* No hidden retries, logic, or side effects

---

### **Phase 3.4 – Observability & Compliance Outputs**

**Acceptance Criteria**

* HITL SLA metrics available  
* Exportable audit reports  
* All compliance artifacts derivable from workflow execution history

---

## **Phase 4 – Developer Experience & Framework Examples (Updated)**

**Goal:** Prove Gantral is usable by real agent frameworks without changing core semantics.

### **Phase 4.1 – Reference Agent Proxy**

**Acceptance Criteria**

* Example runs end-to-end against Gantral Core  
* Demonstrates WAITING\_FOR\_HUMAN → APPROVED → RESUMED  
* No new execution semantics introduced

---

### **Phase 4.2 – Policy Examples Library**

**Acceptance Criteria**

* Policies map cleanly to the documented schema  
* Policies do not assume internal implementation details

---

### **Phase 4.3 – Consumer Guide (Updated)**

**Acceptance Criteria**

* Explicit documentation of **Persisted Pause** pattern  
* Explicit documentation of **Split-Agent Pattern**  
* Decision table mapping framework capability → required pattern

---

### **Phase 4.4 – Framework Reference Implementations (New – Critical)**

**Objective**  
Build concrete, end-to-end reference implementations demonstrating both supported execution patterns.

**Deliverables**

* CrewAI Flow example using native persistence (`@persist`)  
* LangGraph example using checkpointers (SQLite / Redis)  
* Split-Agent example for a non-persistent framework

**Acceptance Criteria**

* Agent pauses on `WAITING_FOR_HUMAN`  
* Process is terminated  
* Gantral approval is granted  
* A **new process** resumes execution  
* No agent internal state is stored in Gantral

Failure to complete this phase blocks external adoption.

---

## **Phase 5 – Federated Execution & Enterprise Identity**

**Goal:** Enable secure, federated execution across multiple teams with strict isolation.

### **Phase 5.1 – Identity Federation & Claim Mapping**

**Acceptance Criteria**

* User identity derived from upstream IdP claims  
* No Gantral-managed user database  
* Team isolation enforced via derived `team_id`

---

### **Phase 5.2 – Service Identity & Keyless Authentication**

**Acceptance Criteria**

* ServiceIdentity entity implemented  
* AWS IAM / K8s SA / equivalent supported  
* Services can create instances without API keys

---

### **Phase 5.3 – Runner Protocol & Task Queues**

**Acceptance Criteria**

* Runners subscribe to explicit task queues  
* Work scheduled centrally, executed locally  
* Runner emits **COMPLETED / FAILED / SUSPENDED** signals

---

### **Phase 5.4 – Connection Registry & Secret Resolution**

**Acceptance Criteria**

* Connection entities store references only  
* Secrets resolved at execution edge by Runners  
* Raw credentials never persist in Gantral

---

## **Stop Conditions (Critical)**

Do **NOT** proceed if:

* Acceptance criteria fail  
* Execution states or schemas drift from the TRD  
* Determinism or authority guarantees are violated  
* Agent internal state leaks into Gantral or Temporal history

---

## **Final Reminder**

Gantral is an **execution authority layer**, not a workflow engine and not an agent framework.

Determinism is guaranteed by the runtime.  
Authority is enforced by Gantral.  
Human accountability is final.

