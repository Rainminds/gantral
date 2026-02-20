---
title: Technical Reference Document
---

# **Technical Reference & Architecture Document (TRD)**

## **Version: v6.1**

Status: Authoritative Technical Constitution  
Audience: Core contributors, platform engineers, security reviewers, enterprise architects, design partners

---

# **1\. Purpose**

This document defines the authoritative technical architecture and execution semantics of Gantral.

Gantral is an AI Execution Control Plane that enforces deterministic authority transitions in agentic workflows.

This document defines:

* Architectural invariants  
* Canonical execution model  
* Authority semantics  
* Commitment artifact structure  
* Hash integrity properties  
* Replay guarantees  
* Responsibility boundaries  
* Enterprise adoption characteristics

If an implementation conflicts with this document, the implementation is incorrect.

---

# **1.1 Structural Challenges Addressed**

Enterprise adoption of agentic AI introduces three recurring structural failures that are not solved by agent frameworks or orchestration tools alone.

## **Operational Inefficiency**

Governance thresholds are frequently embedded directly within workflow code:

* Monetary limits hardcoded in agents  
* Risk thresholds implemented inside orchestration logic  
* Team-specific approval rules implemented as workflow forks  
* Environment-specific deployments differing only by governance parameters

This creates:

* Policy–code duplication  
* Redeployment risk for policy changes  
* Governance drift between documented and deployed behavior  
* Slower adaptation to regulatory or risk updates

Gantral eliminates policy–code duplication by separating policy thresholds from workflow implementation and binding authority decisions to versioned policy bundles.

---

## **Operational Fragmentation**

Authority decisions are often evaluated in one system but enforced in another:

* Policy engines detached from execution runtime  
* Human approvals recorded outside canonical workflow state  
* Logs capturing events without structural binding to execution progression  
* Multiple agent frameworks lacking shared authority semantics

In such architectures, authority exists as documentation — not as deterministic execution state.

Gantral represents authority directly as canonical workflow state transitions.

---

## **Broken Chain of Custody**

At AI execution boundaries:

1. An agent proposes an action  
2. A human approves  
3. Execution resumes  
4. Logs attempt reconstruction

Without cryptographic binding, reconstruction depends on:

* Log integrity  
* Runtime availability  
* Policy version drift  
* Operator testimony

Gantral emits cryptographically chained commitment artifacts at authority boundaries, enabling deterministic replay independent of logs or runtime access.

---

# **2\. Scope and Non-Goals**

## **2.1 What Gantral Is**

Gantral is:

* An execution control plane for AI-assisted workflows  
* A deterministic authority state machine  
* An atomic commitment artifact emitter  
* A replay-verifiable execution authority ledger  
* A policy-advisory integration layer  
* Infrastructure positioned above agents and below enterprise accountability systems

Gantral represents authority as canonical workflow state.

---

## **2.2 Explicit Non-Goals**

Gantral does not:

* Build, host, or orchestrate AI agents  
* Persist agent memory or reasoning traces  
* Encode domain-specific business logic  
* Replace workflow runtimes  
* Replace CI/CD, ITSM, or enterprise ticketing systems  
* Act as an identity provider  
* Store secrets or credentials  
* Provide autonomous self-approval mechanisms  
* Guarantee regulatory compliance

Gantral enforces authority transitions only.

---

# **3\. Architectural Invariants**

The following invariants are non-negotiable.

## **3.1 Instance-First Semantics**

All authority, audit, replay, and cost semantics attach to immutable execution instances.

---

## **3.2 Authority Is State**

Authority exists only as canonical workflow state transitions.

---

## **3.3 HITL Is a Blocking State**

Human-in-the-loop is modeled as `WAITING_FOR_HUMAN`, not as an external notification.

---

## **3.4 Atomic Authority Commitment**

Authority transition and commitment artifact emission occur atomically.

If artifact persistence fails:

* The authority state transition MUST NOT be observable.  
* Execution MUST remain in `WAITING_FOR_HUMAN`.  
* Partial transitions are forbidden.

---

## **3.5 Determinism Over Performance**

Replay correctness supersedes latency optimization.

---

## **3.6 Policy Is Advisory**

Policy engines provide advisory signals only.

Final authority is represented exclusively as workflow state transitions.

---

## **3.7 Version Binding**

Authority decisions must bind:

* `workflow_version_id`  
* `policy_version_id`

Version mismatch during replay yields `INCONCLUSIVE`.

---

## **3.8 Agent State Separation**

Gantral never persists:

* Agent memory  
* Internal plans  
* Tool traces

Agent persistence is the responsibility of agent frameworks.

---

## **3.9 Determinism & Evidence**

Human authority must produce structured, attributable reasoning.

The commitment artifact includes:

* `human_actor_id`  
* `justification`

Deployments MAY configure:

* Empty justification as invalid  
* Minimum reasoning requirements

Approval without attributable reasoning undermines admissibility and may be rejected.

---

# **4\. Formal Execution Model**

## **4.1 Canonical State Set**

The canonical state set S consists of:

* CREATED  
* RUNNING  
* WAITING\_FOR\_HUMAN  
* APPROVED  
* REJECTED  
* OVERRIDDEN  
* RESUMED  
* TERMINATED  
* COMPLETED

---

## **4.2 Transition Relation**

Allowed transitions:

* (CREATED, RUNNING)  
* (RUNNING, WAITING\_FOR\_HUMAN)  
* (WAITING\_FOR\_HUMAN, APPROVED)  
* (WAITING\_FOR\_HUMAN, REJECTED)  
* (WAITING\_FOR\_HUMAN, OVERRIDDEN)  
* (APPROVED, RESUMED)  
* (OVERRIDDEN, RESUMED)  
* (RESUMED, RUNNING)  
* (RUNNING, COMPLETED)  
* (RUNNING, TERMINATED)

Transitions not enumerated above are illegal.

---

## **4.3 Transition Validity**

Let the state sequence for execution instance E be:

σE \= (s₀, s₁, …, sₙ)

For all i in \[0, n−1\]:

(sᵢ, sᵢ₊₁) must belong to the allowed transition relation.

Otherwise execution is invalid and must terminate.

---

# **5\. Commitment Artifact**

## **5.1 Artifact Fields**

Each authority transition emits a commitment artifact containing:

* artifact\_version  
* artifact\_id  
* instance\_id  
* workflow\_version\_id  
* prev\_artifact\_hash  
* authority\_state  
* policy\_version\_id  
* context\_snapshot\_hash  
* human\_actor\_id  
* justification  
* timestamp  
* artifact\_hash

---

## **5.2 Identity Provenance**

Identity is validated using federated OIDC.

Gantral verifies the token and records the subject identifier as `human_actor_id`.

Artifact emission occurs only after identity validation.

Gantral does not operate its own identity directory.

---

## **5.3 Context Snapshot**

`context_snapshot_hash` binds:

* Workflow parameters  
* Policy evaluation inputs  
* Authority-relevant request payload

It does not persist agent memory.

---

# **6\. Hash Integrity Model**

Let H be a collision-resistant cryptographic hash function.

For artifact i:

If i \= 0:  
artifact\_hash\_i \= H(payload\_i)

If i \> 0:  
artifact\_hash\_i \= H(payload\_i concatenated with artifact\_hash\_\{i−1\})

Integrity Property:

For any i \> 0, modification of:

* payload\_i  
* artifact\_hash\_\{i−1\}

causes artifact\_hash\_i to change and invalidates all subsequent artifacts during replay.

Artifacts form a recursive tamper-evident chain.

---

# **7\. Replay Determinism**

Replay validation MUST verify:

1. Hash-chain integrity  
2. Valid state transitions  
3. workflow\_version\_id consistency  
4. policy\_version\_id consistency

Replay reconstructs the authority-state projection of execution.

Replay outcomes:

* VALID  
* INVALID  
* INCONCLUSIVE

Replay requires no runtime, database, or log access.

---

# **8\. Reference Architecture**

## **8.1 Layered Model**

Enterprise Systems  
↓  
Gantral Control Plane  
↓  
Deterministic Workflow Runtime  
↓  
Execution Systems

Side integrations:

* Policy Engine (e.g., OPA)  
* Identity Provider (OIDC)  
* Append-only Artifact Store

---

## **8.2 Control Plane Responsibilities**

Gantral:

* Enforces canonical authority state transitions  
* Evaluates advisory policy signals  
* Emits commitment artifacts  
* Maintains append-only authority history  
* Provides unified visibility into running and paused workflows  
* Exposes replay-verification tooling

Gantral is a structural control plane.

It unifies authority semantics across heterogeneous agent frameworks and eliminates forked governance logic embedded in workflow code.

---

## **8.3 Runtime Responsibilities**

Deterministic workflow runtime:

* Provides durable timers  
* Preserves event ordering  
* Supports long-running execution  
* Enables deterministic replay of workflow logic

Gantral governs authority.  
Runtime executes logic.

---

# **9\. Policy Integration & Separation**

Policy evaluation may be implemented using Open Policy Agent (OPA).

Policies:

* Are authored in Rego  
* Are versioned independently from workflow code  
* May be updated without redeploying agent workflows

Gantral records `policy_version_id` within each commitment artifact.

Policy updates do not require modifying workflow code.

Policy evaluation is advisory.  
Authority remains represented exclusively as state transitions.

This separation reduces operational duplication and redeployment risk.

---

# **10\. Storage Model**

Execution indices: PostgreSQL (non-authoritative)

Commitment artifacts: Append-only object storage (authoritative)

Artifact store is authoritative for replay.

Write-once configuration is recommended.

---

# **11\. Security Architecture**

Identity: OAuth 2.0 / OIDC federation  
Authorization: Policy-driven  
Secrets: External secret managers only  
Audit: 100% authority capture via artifacts

Gantral does not store raw secrets.

---

# **12\. Unified Authority Visibility**

Gantral provides unified visibility into execution instances across teams and domains, including:

* Running workflows  
* Paused workflows awaiting authority  
* Authority progression history

This reduces governance ambiguity and operational blind spots.

Visibility refers strictly to authority state and execution metadata. Gantral does not provide analytics, dashboards, or cost intelligence.

---

# **13\. Non-Functional Characteristics**

Authority transitions require:

* Hash computation  
* Artifact persistence

In human-in-the-loop workflows, human latency dominates pause duration.

Deterministic correctness takes precedence over micro-optimizations.

---

# **14\. Enterprise Adoption Considerations**

Gantral is designed for incremental and reversible adoption.

Organizations may:

* Introduce Gantral at high-materiality authority boundaries  
* Retain existing workflow runtimes  
* Gradually externalize policy thresholds from workflow code  
* Maintain agent frameworks without modifying internal memory models

Gantral does not require rewriting agents.  
It introduces deterministic authority semantics above them.

---

# **15\. Testing & Constitutional Enforcement**

Testing enforces:

* Transition correctness  
* Atomic artifact emission  
* Hash integrity  
* Version consistency  
* Replay determinism  
* Fail-closed behavior

Testing is mechanical enforcement of invariants.

---

# **16\. Auditor Considerations**

Gantral provides:

* Cryptographically bound authority transitions  
* Version-bound policy evaluation  
* Identity validation at decision time  
* Deterministic replay independent of logs

Gantral does not guarantee regulatory compliance.

It provides verifiable execution-time authority evidence.

---

# **17\. Legal Disclaimer**

Gantral is open-source infrastructure provided as-is.

This document does not constitute legal advice or regulatory certification.

Organizations deploying Gantral are responsible for independent security and compliance evaluation.

---

# **18\. Foundational Principle**

Gantral is not about what AI can do.

It is about what organizations allow AI to do —  
and how that authority is structurally enforced and provably replayable.

This document is authoritative.

---

