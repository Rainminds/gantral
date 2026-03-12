---
title: Public Product Requirements Document (PRD)
---

**Version:** v13.0
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

Gantral enforces **authority correctness per execution instance**.

Gantral intentionally does **not**:

* provide dashboards  
* manage governance lifecycle  
* orchestrate autonomy tiers

Those capabilities belong to higher-layer platforms.

Gantral remains a **minimal infrastructure kernel**.

---

# **2\. Document Authority Hierarchy**

Gantral documentation follows a strict normative hierarchy.

**1\. Technical Reference Document (TRD)**  
Defines the authoritative architecture, execution semantics, artifact specification, and replay guarantees.

**2\. Implementation Guide**  
Describes how to build compliant implementations of the TRD.

**3\. Product Requirements Document (PRD)**  
Describes product intent, system scope, and adoption rationale.

If any statement in this PRD conflicts with the TRD, **the TRD prevails**.

The **Commitment Artifact schema and cryptographic model are defined exclusively in the TRD Appendix** and must not be redefined elsewhere.

---

# **3\. Strategic Context**

AI adoption does not stall because models fail.

It stalls when organizations attempt to move from:

**Low-risk experimentation**

to

**consequential operational workflows**.

At this boundary:

AI can act.  
Humans remain accountable.

Most enterprises stall between:

**Tier 1 — Human-supervised automation**  
and  
**Tier 2 — Conditional automation**

because execution authority becomes fragmented.

Gantral addresses this boundary directly by converting authority from **organizational convention into infrastructure enforcement**.

---

### **3.1. Infrastructure Layer Thesis**

Complex software ecosystems evolve by separating concerns into independent infrastructure layers.

Examples include:

* container orchestration separating compute scheduling from applications

* identity infrastructure separating authentication from services

* observability infrastructure separating telemetry from execution systems

Agentic AI systems introduce a similar structural separation.

Intelligence systems determine **what should be done**.  
 Execution systems perform **what is done**.

Between them, enterprises require infrastructure that determines **what is allowed to execute**.

Gantral provides this layer.

As agentic systems scale, an **execution authority layer becomes a structural requirement** in the architecture stack.

---

# **4\. Structural Barriers to Scaling AI**

Gantral addresses structural barriers preventing enterprises from scaling AI into consequential workflows.

---

## **4.1 Policy-in-Code Duplication**

Authority logic is often embedded inside:

* workflow definitions  
* orchestrators  
* agent frameworks  
* BPMN engines  
* conditional execution logic

Every governance change requires:

* code modification  
* redeployment  
* environment synchronization  
* cross-team coordination

This results in:

* deployment duplication  
* version drift  
* slower governance adaptation  
* increased operational cost

Gantral externalizes policy evaluation and binds **policy versions to authority decisions**.

---

## **4.2 Environment Fragmentation**

Without version-bound authority semantics:

* approval thresholds diverge across environments  
* dev/staging/prod behave differently  
* production-only inconsistencies emerge

Gantral binds:

* workflow\_version\_id  
* policy\_version\_id  
* context\_snapshot\_hash

to each authority decision.

---

## **4.3 Operational Fragmentation**

As AI adoption expands:

* approval logic is reimplemented repeatedly  
* teams build separate automation stacks  
* governance semantics diverge

Gantral introduces a **uniform authority model independent of runtime**.

---

## **4.4 Broken Chain of Custody**

Typical AI workflows:

1. AI proposes action  
2. Human approves  
3. Execution resumes  
4. Logs attempt reconstruction

This creates ambiguity:

* model version unknown  
* policy version unclear  
* context snapshot missing  
* identity linkage ambiguous

Gantral binds authority decisions to **cryptographically verifiable execution artifacts**.

---

## **4.5 Fragmented Authority Semantics**

Enterprises operate multiple systems:

* workflow orchestrators  
* agent frameworks  
* internal automation tools

Approval semantics differ across systems.

Gantral provides a **canonical authority state machine**.

---

## **4.6 Non-Defendability**

Organizations often cannot deterministically answer:

* Which model version ran?  
* Which workflow version executed?  
* Which policy version governed the decision?  
* Who approved the action?  
* What context existed at approval time?

Logs reconstruct events.

High-impact AI requires **replayable authority evidence**.

---

# **5\. System Scope**

Gantral is the **Execution Authority Kernel**.

It enforces:

* canonical authority state machine  
* explicit transition relations  
* atomic authority transitions  
* identity validation via OIDC  
* policy advisory integration  
* workflow version binding  
* policy version binding  
* context snapshot binding  
* tamper-evident artifact chains  
* offline replay verification  
* fail-closed execution guarantees

Gantral operates **above workflow runtimes and below enterprise governance platforms**.

---

# **6\. System Boundary**

Gantral sits within the execution stack as follows:

Agent Framework  
        ↓  
Workflow Orchestration  
        ↓  
Gantral (Execution Authority Kernel)  
        ↓  
Execution Systems

External integrations include:

* Policy Engines (OPA or equivalent)  
* Identity Providers (OIDC)  
* Append-only Artifact Storage

Gantral does **not control agent reasoning or workflow orchestration**.

It governs **authority transitions only**.

---

# **7\. Deterministic Authority Model**

Authority is represented as **canonical execution state**, not metadata.

Canonical state progression:

CREATED  
→ RUNNING  
→ WAITING\_FOR\_HUMAN  
→ APPROVED / REJECTED / OVERRIDDEN  
→ RESUMED  
→ COMPLETED / TERMINATED

Rules:

* only enumerated transitions are valid  
* illegal transitions terminate execution  
* transitions are atomic  
* artifact persistence is mandatory

Authority becomes **structural rather than interpretive**.

---

# **8\. Commitment Artifacts**

When authority decisions occur, Gantral emits **Commitment Artifacts**.

Artifacts bind:

* execution instance identity  
* workflow version  
* policy version  
* authority decision  
* human identity  
* context snapshot  
* cryptographic ordering

Artifacts form a **tamper-evident hash chain**.

### **Canonical Artifact Specification**

The artifact schema, hashing model, and signature structure are defined in:

**TRD — Appendix A: Artifact Specification v1**

This PRD intentionally **does not duplicate the artifact structure** to prevent specification drift.

---

# **9\. Reference Integration Flow**

A typical execution flow using Gantral proceeds as follows:

1. A workflow runtime begins executing an automation task.  
2. The workflow reaches an **authority checkpoint** and invokes Gantral.  
3. Gantral evaluates policy via an external policy engine.  
4. If policy returns **ALLOW**, execution continues.  
5. If policy returns **REQUIRE\_HUMAN**, the execution instance transitions to **WAITING\_FOR\_HUMAN**.  
6. A human authority decision is captured and bound to the execution instance.  
7. Gantral emits a commitment artifact and execution resumes.

Authority decisions are therefore **structurally bound to workflow progression** rather than reconstructed from logs.

---

# **10\. Replay**

Gantral provides deterministic **offline replay verification**.

Replay validates:

* artifact hash chain integrity  
* authority transition correctness  
* workflow version consistency  
* policy version consistency

Replay outputs:

* VALID  
* INVALID  
* INCONCLUSIVE

Replay requires only the artifact chain.

Replay does **not require**:

* runtime systems  
* logs  
* databases

---

# **11\. Policy Separation**

Gantral integrates with external policy engines such as **Open Policy Agent**.

Policy engines return advisory signals:

* ALLOW  
* REQUIRE\_HUMAN  
* DENY

Gantral interprets these signals and enforces authority transitions.

Policy remains **advisory**.

Authority remains **structural execution state**.

---

# **12\. Explicit Non-Claims**

Gantral intentionally does **not** claim to:

* prevent malicious operators  
* secure compromised infrastructure  
* guarantee regulatory compliance  
* replace organizational governance processes  
* determine the legality or ethics of execution decisions

Gantral provides **structural authority enforcement and verifiable evidence**.

Organizational governance remains the responsibility of the deploying organization.

---

# **13\. Measurable Outcomes**

Gantral enables measurable improvements in enterprise AI adoption.

### **Policy Redeploy Cycle Reduction**

* reduce governance change lead time from weeks to hours  
* eliminate workflow forks caused solely by governance variation

### **Environment Drift Reduction**

* eliminate hidden governance differences across environments  
* reduce production-only inconsistencies

### **Audit Preparation Reduction**

* reduce audit preparation cycles by 30–70%  
* eliminate cross-system log stitching

### **Chain-of-Custody Assurance**

* reduce incident reconstruction time from weeks to minutes  
* enable deterministic authority replay

---

# **14\. Non-Goals**

Gantral intentionally excludes:

* policy lifecycle management  
* enterprise dashboards  
* approval queue optimization  
* cross-workflow analytics  
* autonomy orchestration  
* managed hosting services

These capabilities belong in higher-level platforms built on top of the kernel.

Gantral remains:

* deterministic  
* minimal  
* open infrastructure  
* vendor-neutral

---

# **15\. Ecosystem & Extensibility**

Gantral is designed as **neutral infrastructure**, not an end-user platform.

The kernel intentionally leaves higher-level capabilities outside its scope, including:

* approval user interfaces

* governance lifecycle systems

* autonomy optimization tooling

* analytics and authority intelligence

* compliance and audit packaging

These capabilities can be implemented as independent systems built on top of the Gantral kernel.

This separation enables an ecosystem of interoperable tools while preserving a **minimal deterministic authority substrate**.

---

# **16\. Foundational Principle**

Gantral is not about what AI can do.

It is about **what organizations allow AI to do — and how that authority is structurally enforced and provably replayable**.

Authority becomes infrastructure.

---

