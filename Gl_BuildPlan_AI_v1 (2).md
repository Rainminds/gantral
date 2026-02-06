# **Gantral – Phase-wise Build Plan for AI Assistants**

**Version:** v10.0 (Authority-First, Federated Execution, Agent-Native Persistence)

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

## **Phase 4 – Developer Experience & Executable Proof**

**Goal:** Prove Gantral is usable, enforceable, and understandable via runnable, end-to-end demos — without altering core semantics.

**Phase 4 is considered incomplete unless an external engineer can run the demos end-to-end in \<30 minutes.**

---

### **Phase 4.1 – Executable Demo Environment (New – Required)**

**Objective:**  
 Provide a docker-compose-based environment that brings up Gantral, runtime, policy evaluator, runner, and agents.

**Acceptance Criteria:**

* `docker compose up` starts all required components

* No Kubernetes, Helm, or cloud dependencies

* Deterministic startup with no manual steps

---

### **Phase 4.2 – Persistent Agent Reference Implementation**

**Objective:**  
 Demonstrate Gantral integration with an agent framework that supports native checkpointing.

**Acceptance Criteria:**

* Agent runs until WAITING\_FOR\_HUMAN

* Agent persists internal state via framework-native persistence

* Agent process exits cleanly

* Upon approval, a **new process** resumes execution

* No agent internal state is stored in Gantral or Temporal

---

### **Phase 4.3 – Split-Agent Reference Implementation**

**Objective:**  
 Demonstrate Gantral integration where the agent framework does **not** support persistence.

**Acceptance Criteria:**

* Pre-approval agent exits on WAITING\_FOR\_HUMAN

* Minimal handoff context is persisted externally

* Post-approval agent starts as a new process

* Execution resumes using approved context only

---

### **Phase 4.4 – Scripted External Interaction (New – Critical)**

**Objective:**  
 Ensure the demo is runnable by outsiders with no prior context.

**Required Commands:**

* `trigger` – create execution instance

* `status` – observe execution state

* `approve` – capture human decision

**Acceptance Criteria:**

* WAITING\_FOR\_HUMAN is visibly enforced

* Approval is explicit and attributable

* RESUMED state occurs only after approval

---

### **Phase 4.5 – Acceptance Test for External Adoption**

**Phase 4 is complete only if:**

`git clone`

`cd examples/persistent-agent`

`docker compose up`

`./trigger.sh`

`./status.sh   # shows WAITING_FOR_HUMAN`

`./approve.sh`

`./status.sh   # shows RESUMED / COMPLETED`

If this flow cannot be executed by an external engineer, **Phase 4 is not complete**.

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

### **Phase 5.5 – Evidence Capture & Tool Mediation (Non-Authoritative)**

**Acceptance Criteria**

* Runs at runner boundary  
* Stores evidence externally  
* Gantral stores only references  
* No tool-level blocking  
* No policy inspection of payloads

**Constraint:**

* Phase 5.5 must not be started until Phase 4 acceptance tests pass.  
* Evidence capture must not be introduced into Phase 4 demos.

---

## **Stop Conditions (Critical)**

Do **NOT** proceed if:

* Acceptance criteria fail  
* Execution states or schemas drift from the TRD  
* Determinism or authority guarantees are violated  
* Agent internal state leaks into Gantral or Temporal history

---

## **Phase 6 – Verifiability & Admissibility (Edward-Readiness)**

Goal:  
Transform documented guarantees into mechanically verifiable artifacts that survive hostile audit, infrastructure compromise, and time.

This phase MUST NOT alter execution semantics, state machines, or authority rules.

### **Phase 6.1 – Commitment Artifact Implementation**

Objective:  
Implement a concrete, inspectable commitment artifact emitted atomically with authority transitions.

Acceptance Criteria:  
• Artifact schema versioned and frozen  
• Artifact emitted only on APPROVED / REJECTED / OVERRIDDEN  
• Artifact emission is atomic with execution state transition  
• No artifact → no execution continuation  
• Artifact includes hash chain to prior artifact

### **Phase 6.2 – Artifact Storage & Log Independence**

Objective:  
Ensure artifacts are independent of operational logs and databases.

Acceptance Criteria:  
• Artifacts written to append-only storage  
• Artifacts not reconstructible from logs or DB  
• Deleting Gantral DB does not invalidate artifacts  
• DB records are non-authoritative indices only

### **Phase 6.3 – Offline Verification Tooling**

Objective:  
Enable third-party verification without Gantral access.

Acceptance Criteria:  
• CLI or library verifier implemented (gantral-verify)  
• No network access required  
• Outputs: VALID / INVALID / INCONCLUSIVE  
• Verifier depends only on artifact(s) and schema

### **Phase 6.4 – Authority-Only Replay Enforcement**

Objective:  
Guarantee replay depends solely on authority artifacts.

Acceptance Criteria:  
• Replay consumes only artifact chain  
• Agent memory, logs, and tool outputs excluded  
• Changing agent code does not affect replay outcome

### **Phase 6.5 – Fail-Closed Guarantees**

Objective:  
Eliminate ambiguous or best-effort execution paths.

Acceptance Criteria:  
• Execution terminates on missing or partial artifacts  
• Hash mismatches invalidate execution  
• No retries or inferred authority  
• Failures surface explicitly

### **Phase 6.6 – Auditor Verification Demo**

Objective:  
Demonstrate offline verification from an auditor’s perspective.

Acceptance Criteria:  
• Execution produces artifact(s)  
• Gantral services can be shut down  
• Offline verifier validates artifact(s)  
• Demo runnable in under 30 minutes

## **Stop Conditions (Non-Negotiable)**

Do NOT proceed if:  
• Artifact emission is non-atomic  
• Replay requires Gantral services  
• Logs or dashboards are treated as evidence  
• Execution continues under ambiguity

---

## **Final Reminder**

Gantral is an **execution authority layer**, not a workflow engine and not an agent framework.

Determinism is guaranteed by the runtime.  
Authority is enforced by Gantral.  
Human accountability is final.

