---
title: Public Product Requirements Document (PRD)
---

**Version:** v6.1
**Status:** Public, canonical product document  
**Audience:** Enterprises, platform engineers, contributors, regulators, partners  

# **Product Requirements Document**

---

# **1\. Executive Summary**

Gantral is **Execution Infrastructure for Deterministic Agentic AI**.

It enables enterprises to:

* Separate policy thresholds from workflow code  
* Eliminate governance drift across teams  
* Enforce human authority as canonical workflow state  
* Bind execution decisions to versioned policy and workflow definitions  
* Replay authority independently under adversarial conditions

Gantral addresses three structural enterprise failures:

* Operational inefficiency from policy–code duplication  
* Governance fragmentation across systems  
* Broken chain-of-custody at execution boundaries

Gantral is not a governance dashboard.  
It is execution authority infrastructure.

---

# **2\. The Enterprise Problem**

As agentic AI systems execute material actions — financial transfers, production changes, access grants — governance must move from convention to infrastructure.

## **2.1 Operational Inefficiency**

In many AI deployments:

* Approval thresholds are hardcoded in agent logic  
* Risk limits are embedded in orchestration code  
* Teams fork workflows to represent different limits  
* Policy updates require redeployment

This results in:

* Redeployment risk  
* Configuration drift  
* Slower regulatory adaptation  
* Increased operational surface area

Gantral externalizes policy from workflow implementation.

Policy changes no longer require agent modification or workflow redeploy.

---

## **2.2 Operational Fragmentation**

Authority decisions are frequently:

* Evaluated in policy engines  
* Enforced elsewhere  
* Recorded in logs  
* Approved via chat or ticketing systems

Authority becomes narrative, not structural.

Gantral represents authority as canonical workflow state transitions.

Execution cannot progress without structurally recorded authority.

---

## **2.3 Broken Chain of Custody**

Without cryptographic binding:

* Logs can be modified  
* Policy versions can drift  
* Approval context can be lost  
* Reconstruction depends on runtime access

Gantral replaces reconstruction with deterministic replay.

Authority becomes cryptographically bound to execution state.

---

# **3\. Deterministic Authority Model**

Gantral defines a canonical state machine:

CREATED → RUNNING → WAITING\_FOR\_HUMAN  
→ APPROVED / REJECTED / OVERRIDDEN  
→ RESUMED → COMPLETED / TERMINATED

Only explicitly enumerated transitions are valid.

If artifact persistence fails:

Execution does not proceed.

Fail-closed semantics are mandatory.

No inferred authority.  
No silent fallback.  
No post-hoc reconstruction.

---

# **4\. Commitment Artifacts (Admissible Evidence)**

At each authority boundary, Gantral emits a commitment artifact binding:

* instance\_id  
* workflow\_version\_id  
* policy\_version\_id  
* authority\_state  
* human\_actor\_id (OIDC validated)  
* justification  
* context\_snapshot\_hash  
* timestamp  
* artifact\_hash

Artifacts form a recursive, tamper-evident hash chain.

Modification invalidates downstream history.

Artifacts are:

* Log-independent  
* Runtime-independent  
* Database-independent  
* Replayable offline

This ensures admissible chain-of-custody.

---

# **5\. Replay Under Adversarial Conditions**

Replay verifies:

1. Hash-chain integrity  
2. Valid state transitions  
3. workflow\_version\_id consistency  
4. policy\_version\_id consistency

Replay requires no:

* Gantral services  
* Workflow runtime  
* Logs  
* Agent memory

Replay outcomes:

* VALID  
* INVALID  
* INCONCLUSIVE

Authority decisions remain inspectable under hostile review.

---

# **6\. Policy Separation (OPA Integration)**

Gantral separates policy evaluation from authority enforcement.

In the reference implementation:

* Policy evaluation uses Open Policy Agent (OPA)  
* Policies are authored in Rego  
* Policy bundles are versioned independently  
* Policy evaluation is advisory only

OPA does not grant authority.

Gantral enforces authority as canonical workflow state.

The evaluated policy\_version\_id is recorded in each artifact.

Policy updates do not require workflow redeploy.

This reduces:

* Operational duplication  
* Redeployment risk  
* Governance drift

---

# **7\. Storage Model**

Gantral distinguishes between:

* Non-authoritative execution indices (PostgreSQL)  
* Authoritative commitment artifacts (append-only object storage)

Artifact storage is write-once.

Replay depends exclusively on artifact chain integrity.

Databases and logs are not authoritative for admissibility.

---

# **8\. Reference Implementation Status**

Gantral is implemented in:

* Go  
* Temporal (deterministic workflow runtime)

It includes:

* Canonical authority state machine  
* OPA integration  
* Artifact emission  
* Standalone replay CLI

Gantral is an active open-source initiative.

---

# **9\. Performance Characteristics**

Authority transitions require:

* Hash computation  
* Artifact persistence

In human-in-the-loop workflows, approval latency dominates.

Cryptographic overhead is negligible relative to workflow duration.

Gantral is designed for correctness over micro-optimization.

---

# **10\. Enterprise Impact**

After adoption:

* Policy updates occur without redeploy  
* Governance is consistent across teams  
* Approval reasoning is attributable  
* Authority is version-bound  
* Audit becomes deterministic  
* Redeployment risk decreases  
* Execution authority becomes standardized infrastructure

Gantral centralizes authority semantics without centralizing agents.

---

# **11\. Design Partner Adoption Model**

Gantral is designed for incremental introduction.

Organizations can:

* Start at high-materiality boundaries  
* Integrate without rewriting agents  
* Retain existing runtimes  
* Gradually externalize governance

Gantral introduces deterministic authority above existing systems.

---

# **12\. Explicit Non-Claims**

Gantral does not:

* Guarantee regulatory compliance  
* Replace legal review  
* Provide certification  
* Interpret domain-specific business logic  
* Make autonomous decisions

Gantral provides verifiable execution-time authority evidence.

Organizations remain responsible for compliance evaluation.

---

