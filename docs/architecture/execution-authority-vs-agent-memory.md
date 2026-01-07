# Execution Authority vs Agent Memory vs Runtime

This diagram shows the **strict separation of responsibilities** in Gantral.

It answers one question unambiguously:

**Who owns what state, and why Temporal is still required even when agents manage their own memory.**

---

## **High-level Responsibility Split**

* **Execution Authority** → *Gantral*  
* **Execution Runtime** → *Deterministic Workflow Engine (Temporal)*  
* **Agent Memory & Reasoning** → *Agent Frameworks (CrewAI, LangGraph, etc.)*

---

## **Single-Page Architecture Diagram**

```mermaid
flowchart TB

  %% ===== External Triggers =====
  EXT["External Trigger Event / Schedule / API"]

  %% ===== Gantral =====
  subgraph G["Gantral: Execution Authority"]
    G1["Execution State Machine RUNNING / WAITING_FOR_HUMAN / ..."]
    G2["HITL Enforcement Approve / Reject / Override"]
    G3["Policy Evaluation Interface Transition Guards"]
    G4["Audit & Authority Log Immutable Decisions"]
  end

  %% ===== Runtime =====
  subgraph T["Deterministic Runtime (Temporal)"]
    T1["Workflow Orchestration"]
    T2["Durable Timers & Signals"]
    T3["Deterministic Replay Authority Replay"]
  end

  %% ===== Runners =====
  subgraph R["Distr Runners Team-owned Infra"]
    R1["Agent Process Launcher"]
    R2["Lifecycle Signals COMPLETED / FAILED / SUSPENDED"]
  end

  %% ===== Agent Framework =====
  subgraph A["Agent Framework"]
    A1["Agent Reasoning & Planning"]
    A2["Agent Memory Conversation / Tools / State"]
    A3["Native Persistence Checkpoint DB / S3"]
  end

  %% ===== Flows =====
  EXT --> G1

  G1 --> T1
  G2 --> T2
  G3 --> T1

  T1 --> R1
  R1 --> A1

  A1 --> A2
  A2 --> A3

  %% Suspend / Resume Flow
  A1 -->|WAITING_FOR_HUMAN| R2
  R2 -->|SUSPENDED| T2
  T2 -->|Approval Signal| R1

  %% Completion Flow
  A1 -->|Done / Error| R2
  R2 --> T1

  %% Audit
  G1 --> G4
  G2 --> G4
  T3 --> G4

```
---

## **How to Read This Diagram**

### **1\. Gantral \= Execution Authority**

* Owns execution state  
* Decides **when** execution may proceed  
* Enforces human approval  
* Produces the **authoritative audit log**

Gantral never sees:

* Agent memory  
* Prompts  
* Tool state

---

### **2\. Temporal \= Execution Runtime**

* Owns ordering, time, retries, and durability  
* Guarantees deterministic replay of **authority decisions**  
* Waits safely for hours or days during HITL

Temporal does **not**:

* Run agent reasoning  
* Store agent memory  
* Decide outcomes

---

### **3\. Agent Framework \= Reasoning \+ Memory**

* Owns cognition, planning, and tools  
* Persists internal state using **native persistence**  
* Can safely exit and resume later

If native persistence is **not supported**:

* The agent **must be split** into pre-approval and post-approval stages

---

### **4\. Runners \= Execution Boundary**

* Launch agent processes  
* Detect lifecycle outcomes:  
  * COMPLETED  
  * FAILED  
  * SUSPENDED (hibernation)  
* Translate agent signals into Gantral execution events

Runners never make decisions.

---

## **Why This Separation Matters**

* **Auditability**: Authority can be replayed without agent memory  
* **Cost efficiency**: No compute during long approvals  
* **Framework freedom**: Any agent framework can be used  
* **Regulator trust**: Clear ownership of decisions

---

## **One-line Summary**

**Gantral decides. Temporal remembers. Agents think.**

---

## **Recommended Placement in the Repo**

Add this file as:

/docs/architecture/execution-authority-vs-agent-memory.md

Then:

1. Link it from:  
   * `README.md`  
   * `docs/architecture/README.md`  
2. Reference it from:  
   * PRD (conceptual grounding)  
   * TRD (visual companion)  
   * Consumer Guide (for developers)

This diagram becomes the **canonical mental model** for contributors, users, and reviewers.

