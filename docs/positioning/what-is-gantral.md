---
title: What is Gantral?
---

# What Is Gantral?

Gantral is an **open-source Execution Authority Kernel**.

It provides **Deterministic Authority Infrastructure for Scaling AI**.

Gantral enforces and records **execution-time authority** in AI-assisted and agentic workflows operating in consequential domains.

Gantral does not evaluate model quality, reasoning correctness, or business logic.

Gantral governs **whether execution is admissible**.

---

## Why Gantral Exists

AI adoption does not stall because models fail.

It stalls when organizations attempt to move from:

- Low-risk experimentation  
to  
- Consequential, high-impact workflows  

When AI systems begin to:

- Move money  
- Modify infrastructure  
- Change access control  
- Influence regulatory posture  
- Trigger irreversible operational outcomes  

At that boundary:

AI can act.  
Humans remain accountable.

Most enterprise stacks today can:

- Orchestrate workflows  
- Evaluate policy  
- Apply guardrails  
- Log events  
- Produce dashboards  

Very few structurally enforce execution authority.

Authority is reconstructed from logs instead of enforced at runtime.

Gantral exists to ensure that when an AI-enabled system performs a consequential action,
an organization can later determine — **without relying on trust, narrative, or log stitching**:

- what executed  
- under whose authority  
- with which execution context  
- whether human approval was required  
- whether authority can be independently verified  

Gantral makes authority:

**explicit · enforced · verifiable**

---

## The Execution Authority Gap

Most governance failures in AI systems are not caused by models.

They are caused by how execution authority is implicitly assumed.

Human approval and accountability often exist only as convention:

- pull-request reviews  
- chat approvals  
- informal runbooks  
- screenshots and logs  

These mechanisms:

- are not technically enforceable  
- are not structurally bound to execution  
- do not produce durable, replayable evidence  

At scale, this leads to:

- fragmented execution control across teams  
- approvals detached from execution context  
- post-incident reconstruction based on inference  

Gantral converts authority from convention into infrastructure.

---

## Before and After: Execution Authority

### Before

Authority is implied by process and reconstructed after the fact.

```mermaid
flowchart TB
   A[AI Agent or Model
   • Policy logic embedded in prompts or code
   • Proposes or initiates actions autonomously]

   A --> M[Execution Management
   • Automation workflows
   • Task runners or orchestration
   • Executes without explicit authority checks]

   A -.-> H[Human Reviewer
   • Informal notification or review
   • No explicit authority boundary]

   H --> T[Human Action Tools
   • Dashboards
   • Ticketing systems
   • Admin consoles]

   M --> R[Action in Real World
   • Data changes
   • Customer impact
   • External side effects]

   T --> R

   R -.-> L[Post-hoc Evidence
   • Logs
   • Tickets
   • Chat records
   • Human memory]
````

Authority is reconstructed from logs and memory.

---

### After

Authority is enforced as canonical execution state.

```mermaid
flowchart TB
   P[Policy & Governance
   • Intent & constraints
   • Risk thresholds
   • Advisory only]

   subgraph AI["AI & Agent Systems"]
       LG[LangGraph]
       VE[Vellum]
       CC[Claude Cowork]
   end

   H[Human Operator
   • Reviews context
   • Exercises judgment
   • Accountable actor]

   M[Workflow Orchestration
   • Awaits authority]

   subgraph Gantral["Gantral — Execution-Time Authority"]
       G[Authority Boundary
       • Pause execution
       • Require explicit decision
       • Approve / Reject / Override
       • Record authority deterministically]
   end

   R[Real-World Actions]

   D[Deterministic Authority Record
   • Who authorized
   • What was approved
   • When & under what conditions]

   LG --> M
   VE --> M
   CC --> M

   H --> M
   P -.-> G

   M --> G
   G -->|Authorized| R
   G -.-> D
   G -->|Pause / Escalate| H
```

Authority is not inferred.

It is enforced.

---

## The Core Distinction: Authority vs Intelligence

Gantral introduces a strict separation between:

### Intelligence

* Planning
* Reasoning
* Tool selection
* Action proposals

### Authority

* Whether execution may proceed
* Whether human approval is required
* Whether execution must terminate
* Whether authority can be proven later

Agents provide intelligence.

Gantral enforces authority.

This separation is structural, not conceptual.

---

## How Execution Is Governed

Gantral operates as an **execution-time authority layer**.

A typical flow:

1. An agent proposes an action.
2. Execution reaches a governed boundary.
3. External policy (e.g., OPA) may evaluate conditions (advisory only).
4. Gantral enforces one of:

   * Continue
   * Pause for human authority
   * Reject / Terminate
5. Authority transitions are enforced via a deterministic state machine.
6. A commitment artifact is emitted atomically.

If authority cannot be enforced and recorded,
execution must not proceed.

Gantral fails closed.

---

## What Gantral Owns

Gantral owns **execution authority invariants**, not lifecycle governance.

It provides:

* A deterministic authority state machine
* Explicit `WAITING_FOR_HUMAN` blocking semantics
* Explicit `APPROVED / REJECTED / OVERRIDDEN` transitions
* Identity binding at authority boundaries
* Policy version binding
* Workflow version binding
* Context snapshot binding
* Atomic authority transition + artifact emission
* Tamper-evident artifact chains
* Log-independent replay

Gantral enforces authority correctness **per execution instance**.

---

## Policy Advisory Integration (OPA)

Gantral supports external policy engines in an advisory role.

Policy may be authored in:

* Open Policy Agent (OPA) / Rego
* Custom enterprise policy services

At authority checkpoints:

* Policy evaluates conditions.
* Advisory output is returned.
* `policy_version_id` is bound to the execution instance.
* Gantral enforces the resulting authority transition structurally.

Policy remains advisory.
Authority enforcement remains internal and deterministic.

---

## Federated Runner Model

Gantral uses a federated execution model:

* Agents execute in team-owned infrastructure.
* Gantral does not inspect agent memory or tool payloads.
* Execution pauses at authority boundaries.
* Resume signals inject fresh execution context.
* Long waits allow agent processes to exit cleanly.

Gantral governs **permission to execute**, not execution mechanics.

---

## What Gantral Does Not Do

Gantral explicitly does not:

* Manage policy lifecycle
* Provide dashboards
* Provide cross-workflow analytics
* Orchestrate autonomy tiers
* Replace orchestration engines
* Inspect model prompts or memory
* Guarantee regulatory compliance

Gantral records and enforces **authority**, not intent, correctness, or compliance status.

---

## When Gantral Is Appropriate

Gantral is designed for workflows that:

* Affect production systems
* Have regulatory, financial, or security impact
* Require explicit human accountability
* Must be auditable months or years later

Gantral is not necessary for:

* Advisory-only agents
* Exploratory or sandbox workflows
* Low-impact, reversible actions

Gantral becomes rational infrastructure when:

“Could we independently prove authority correctness without relying on logs?”

must be answered with confidence.

---

## How to Think About Gantral

Useful mental models:

* **“sudo for AI”**
  Execution is intercepted and requires authority before proceeding.

* **“Execution-time constitution”**
  Authority is enforced as state, not reconstructed from metadata.

* **“Chain of custody for automation”**
  Authority is bound to execution at the moment it occurs.

---

## Where to Go Next

* **Authority & Enforcement**
  See the **[Authority State Machine](../architecture/authority-state-machine.md)**

* **Proof & Audit**
  Start with **[Verifiability Overview](../verifiability/README.md)**

* **Technical Semantics**
  Read the **[Technical Reference (TRD)](../architecture/trd.md)**

---

Gantral is intentionally narrow.

It does not attempt to make AI correct, ethical, or intelligent.

It makes **execution authority deterministic, enforceable, and provable**.
