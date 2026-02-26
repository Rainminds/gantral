---
title: Public Product Requirements Document (PRD)
---

**Version:** v11.0
**Status:** Public, canonical product document  
**Audience:** Enterprises, platform engineers, contributors, regulators, partners  

---

# **1\. Executive Summary**

Gantral is **Deterministic Authority Infrastructure for Scaling AI**.

It is an open-source **Execution Authority Kernel** that enforces execution-time authority in AI-assisted and agentic enterprise workflows operating in consequential domains.

Gantral ensures that when AI systems influence:

* Financial movement  
* Infrastructure changes  
* Access control  
* Regulatory posture  
* Customer-impacting operations

execution authority is:

* Deterministic  
* Version-bound  
* Identity-bound  
* Context-bound  
* Cryptographically committed  
* Log-independent  
* Replayable  
* Fail-closed

Gantral does not provide dashboards.  
It does not manage lifecycle governance.  
It does not orchestrate autonomy tiers.

It enforces execution authority correctness per execution instance.

---

# **2\. Strategic Context**

AI adoption does not stall because models fail.

It stalls when organizations attempt to move from:

* Low-risk experimentation  
  to  
* Consequential, high-impact workflows

At this boundary:

AI can act.  
Humans remain accountable.

Most enterprises stall between Tier 1 (human-supervised automation) and Tier 2 (conditional automation) because authority becomes fragmented.

Gantral addresses the authority boundary directly.

It converts authority from convention into infrastructure.

---

# **3\. Structural Barriers to Scaling AI (Gantral-Relevant)**

Gantral directly addresses the following structural barriers identified in enterprise AI adoption.

---

## **3.1 Operational Inefficiency — Policy-in-Code Duplication**

Authority logic is frequently embedded inside:

* Workflow definitions  
* Orchestrators  
* Agent frameworks  
* BPMN engines  
* Conditional execution logic

Every change requires:

* Code modification  
* Redeployment  
* Dev → staging → prod synchronization  
* Cross-team coordination

Consequences:

* Deployment duplication  
* Version drift  
* Environment inconsistencies  
* Slower release cycles  
* Increased operational cost

Gantral externalizes policy evaluation from workflow implementation and binds policy version at authority checkpoints.

---

## **3.2 Environment Fragmentation (Dev / Staging / Prod)**

Without version-bound authority semantics:

* Thresholds differ across environments  
* Approval rules diverge  
* Production-only inconsistencies appear  
* Rollback events increase

Gantral binds:

* workflow\_version\_id  
* policy\_version\_id  
* context\_snapshot\_hash

to each authority decision, eliminating hidden authority drift.

---

## **3.3 Operational Fragmentation Across Teams**

As AI expands:

* Agent code is duplicated  
* Approval logic is reimplemented  
* Governance semantics drift

Authority becomes runtime-specific instead of enterprise-coherent.

Gantral introduces a uniform, deterministic authority state machine independent of orchestration or agent framework.

---

## **3.4 Broken Chain of Custody (AI ↔ Human Handoff)**

Common failure patterns:

* AI recommendation in one system  
* Human approval in another  
* Execution in a third  
* Logs reconstruct what happened

Without structural binding:

* Model version may be unknown  
* Policy version unclear  
* Context snapshot missing  
* Identity linkage ambiguous

Gantral binds authority \+ identity \+ policy \+ context \+ workflow version into commitment artifacts at execution time.

---

## **3.5 Fragmented Authority Semantics Across Runtimes**

Multiple:

* Orchestrators  
* Agent frameworks  
* Internal services

Without shared authority semantics:

* Approval behavior diverges  
* Escalation logic differs  
* Human-in-the-loop rules drift

Gantral provides a canonical authority state machine that is orchestrator-agnostic and agent-agnostic.

---

## **3.6 Non-Defendability**

Enterprises often cannot deterministically answer:

* Which model version ran?  
* Which workflow version?  
* Which policy version?  
* Who approved?  
* What context existed at decision time?

Logs reconstruct.  
High-impact AI requires proof.

Gantral enables deterministic replay independent of logs and runtime.

---

## **3.7 Fragmented or Missing Audit Logs**

Audit evidence is often:

* Distributed  
* Inconsistent  
* Environment-specific  
* Not version-bound

Audit becomes investigative instead of replayable.

Gantral replaces log stitching with cryptographically verifiable artifact chains.

---

## **3.8 Authority–Intelligence Boundary Confusion**

Enterprises assume they have governance because they have:

* Policy engines  
* Logs  
* Observability  
* Orchestration frameworks

But:

* Policy evaluation is not execution authority  
* Logs are not admissible proof  
* Orchestration is not accountability  
* Guardrails are not structural enforcement

Gantral enforces authority as canonical execution state.

---

## **3.9 Black-Box Infrastructure Risk**

Closed, opaque enforcement layers create:

* Limited transparency  
* Vendor lock-in  
* Security review challenges  
* Durability risk

Gantral provides:

* Open-source deterministic kernel  
* Transparent authority semantics  
* Vendor-neutral substrate  
* Log-independent replay

Authority becomes infrastructure — not a black box.

---

# **4\. Product Scope**

Gantral is the **Execution Authority Kernel**.

It enforces:

* Canonical authority state machine  
* Explicit transition relations  
* Atomic authority transition \+ artifact emission  
* Identity validation (OIDC)  
* Policy version binding  
* Workflow version binding  
* Context snapshot binding  
* Tamper-evident artifact chains  
* Offline replay verification  
* Fail-closed semantics

Gantral does not:

* Manage policy lifecycle  
* Provide dashboards  
* Provide cross-workflow analytics  
* Orchestrate autonomy tiers  
* Replace orchestration runtimes  
* Provide managed hosting

It enforces authority correctness per execution instance.

---

# **5\. Deterministic Authority Model**

Canonical state machine:

CREATED → RUNNING → WAITING\_FOR\_HUMAN  
→ APPROVED / REJECTED / OVERRIDDEN  
→ RESUMED → COMPLETED / TERMINATED

Rules:

* Only enumerated transitions are valid  
* Illegal transitions terminate execution  
* Authority transitions are atomic  
* Artifact persistence is mandatory  
* Failure to persist \= execution does not proceed

Authority is modeled as execution state, not metadata.

---

# **6\. Commitment Artifacts**

At each authority transition, Gantral emits a commitment artifact binding:

* instance\_id  
* workflow\_version\_id  
* policy\_version\_id  
* authority\_state  
* human\_actor\_id  
* justification  
* context\_snapshot\_hash  
* timestamp  
* artifact\_hash  
* prev\_artifact\_hash

Artifacts form a recursive hash chain.

Integrity properties:

* Modification invalidates downstream history  
* Replay validates transition correctness  
* Replay requires no runtime or logs

Authority history becomes cryptographically verifiable.

---

# **7\. Replay**

Replay verifies:

1. Hash-chain integrity  
2. Valid state transitions  
3. workflow\_version\_id consistency  
4. policy\_version\_id consistency

Replay outputs:

* VALID  
* INVALID  
* INCONCLUSIVE

Replay is:

* Runtime-independent  
* Log-independent  
* Database-independent

Authority becomes inspectable under adversarial conditions.

---

# **8\. Policy Separation (OPA Integration)**

Gantral integrates with external policy engines (e.g., OPA).

Policy:

* Evaluates materiality thresholds  
* Returns advisory decision  
* Is versioned independently

Gantral:

* Records policy\_version\_id  
* Enforces authority transition  
* Treats policy as advisory only

Policy changes do not require workflow redeploy.

This reduces:

* Policy redeploy cycles  
* Governance duplication  
* Configuration drift

---

# **9\. Measurable Outcomes (Gantral-Relevant)**

Gantral enables measurable acceleration in enterprise AI adoption.

---

## **9.1 Policy Redeploy Cycle Reduction**

* Reduce policy change lead time from weeks to hours  
* Eliminate workflow forks created solely for governance variation  
* Lower regression risk from governance edits

---

## **9.2 Environment Drift Reduction**

* Eliminate hidden authority differences across environments  
* Reduce production-only inconsistencies  
* Lower rollback events due to governance mismatch

---

## **9.3 Audit Preparation Time Reduction**

* Reduce audit preparation cycles by 30–70%  
* Eliminate cross-system log stitching  
* Remove dependency on runtime access during review  
* Produce replay-ready artifact bundles

---

## **9.4 Chain-of-Custody Risk Reduction**

* Reduce post-incident authority reconstruction from weeks to deterministic replay in minutes  
* Lower legal discovery preparation cost  
* Eliminate cross-system log reconciliation

---

## **9.5 Cross-Team Governance Duplication Reduction**

* Eliminate reimplementation of approval handlers  
* Reduce duplicated workflow forks  
* Standardize authority semantics across teams

---

## **9.6 Governance Change Velocity**

* Increase governance change speed without increasing redeploy risk  
* Bind policy evolution to versioned execution instances  
* Prevent fear-driven stagnation in AI progression

---

# **10\. Architecture Overview**

Gantral integrates at explicit authority boundaries:

Agents → Orchestration → Gantral → Execution

Gantral:

* Pauses execution at governed checkpoints  
* Validates identity  
* Binds policy version  
* Emits commitment artifact  
* Resumes execution only upon valid authority transition

Orchestration remains external.

---

# **11\. Design Principles**

Gantral is built on five principles:

1. Authority must be structurally enforced, not reconstructed.  
2. Authority must remain separate from intelligence.  
3. Execution authority must be deterministic and replayable.  
4. Replay must be log-independent.  
5. Enforcement must fail closed.

---

# **12\. Enterprise Position**

Gantral is:

* Minimal by design  
* Infrastructure-grade  
* Self-hosted  
* Open source  
* Agent-agnostic  
* Orchestrator-agnostic  
* Vendor-neutral

It is the constitutional enforcement layer beneath AI systems operating in consequential domains.

---

# **13\. Non-Claims**

Gantral does not:

* Guarantee regulatory compliance  
* Replace legal review  
* Provide certification  
* Interpret business logic  
* Make autonomous decisions

It provides verifiable execution-time authority evidence.

---

# **Final Framing**

AI scaling does not fail because intelligence is insufficient.

It fails because authority is fragmented.

Gantral removes that fragmentation — deterministically, transparently, and measurably.

Authority becomes infrastructure.

AI moves from pilot to platform.

