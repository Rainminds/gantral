---
sidebar_position: 2
title: Example Consumer Guide
---

# Example Consumer Guide

⚠️ **This document is not a system specification.**

It is the **normative usage guide** for external consumers (agents, tools, workflows) integrating with Gantral.

---

## How External Agents and Systems Should Use Gantral

**Status:** Consumption guide (normative for consumers)  
**Audience:**  
- Agent developers  
- Platform engineers integrating Gantral  
- Teams building tools, automations, or workflows that require human approval  

---

## Purpose

This guide explains **when and how external agents and systems should consume Gantral**.

It answers **one question only**:

> **When and how should my agent or system call Gantral?**

This guide is intentionally narrow.  
It does **not** explain how Gantral works internally.

---

## 1. What Gantral Is (From a Consumer’s Perspective)

Gantral is an **execution control plane**.

From the outside, Gantral:

- Receives a request to perform a sensitive action  
- Decides whether execution may continue  
- Pauses execution if human approval is required  
- Resumes or terminates execution based on that decision  

Gantral is **not**:

- An agent framework  
- A task runner  
- A policy engine  
- A workflow orchestrator  

As a consumer, you treat Gantral as a **decision gate**, not as a runtime.

---

## 2. When You Should Call Gantral

You should call Gantral **before executing any action** that:

- Has irreversible side effects  
- Affects production systems  
- Has regulatory, financial, or security impact  
- Requires human accountability  

### Examples

- Deleting production infrastructure  
- Executing high-cost AI operations  
- Triggering payments or refunds  
- Sending external communications  
- Making customer-impacting changes  

If an action can be safely retried, undone, or sandboxed, Gantral is usually **not required**.

---

## 3. High-Level Integration Flow

From a consumer’s point of view, the flow is always:

1. Prepare action context  
2. Ask Gantral for a decision  
3. Handle Gantral’s response  
4. Either wait, resume, or abort  

> **Your system must not execute the action until Gantral explicitly allows it.**

---

## 4. Minimal SDK Interaction Pattern

Exact function names may vary by SDK.  
The **semantics do not**.

### Step 1: Create or Identify an Execution Instance

```python
instance = gantral.create_instance(
    workflow_id,
    workflow_version,
    context
)
````

This represents a **single execution attempt**.

---

### Step 2: Request a Decision

```python
response = gantral.request_decision(
    instance_id = instance.id,
    actor_id = "agent-or-service-id",
    action_context
)
```

Gantral evaluates policy and execution state.

---

### Step 3: Handle the Response

Gantral will return **one of the following outcomes**.

#### APPROVED

* You may proceed with execution
* Use the **exact context returned by Gantral**

```python
execute(action)
```

---

#### WAITING_FOR_HUMAN

* **Do not** execute the action
* Pause the agent or workflow
* Follow the appropriate pause pattern (see Section 5)

```python
wait_for_update(instance_id)
```

---

#### REJECTED

* Do not execute the action
* Abort the workflow

```python
abort("Rejected by policy or human")
```

---

#### TERMINATED

* Treat as a hard stop
* Do not retry

---

## 5. Handling Long-Running Approvals (Critical)

Many approvals take **minutes, hours, or days**.

❌ Do **not** rely on local RAM or in-process state for these waits.

Gantral supports long-running approvals via **two valid consumer patterns**.

---

### 5.1 Pattern A – Persisted Pause (Preferred)

Use this pattern **only if your agent framework supports native persistence / checkpointing**.

When Gantral returns `WAITING_FOR_HUMAN`:

1. Persist agent state using the framework’s native checkpointing
2. Exit the process cleanly (zero CPU usage)
3. Wait for Gantral to signal approval
4. Start a new process
5. Restore agent state from persistence
6. Resume execution

Gantral manages **execution authority**.
Your framework manages **agent memory**.

#### Framework Examples

* CrewAI Flows → `@persist` (SQLite / Postgres)
* LangGraph → Checkpointers (SQLite / Redis / S3)

> **Rule:** If your agent cannot be restarted from a checkpoint, you must not use this pattern.

---

### 5.2 Pattern B – Split-Agent Pattern

*(Mandatory if No Persistence)*

Some agent frameworks do **not** support native checkpointing or resume.

If your agent **cannot** be restarted from persisted state, you must **not** pause execution in memory.

You must split execution into **multiple stages**.

#### The Split-Agent Pattern

When Gantral returns `WAITING_FOR_HUMAN`:

1. Terminate the current agent execution
2. Persist **only minimal handoff context**
   (IDs, references, inputs — not internal memory)
3. Wait for Gantral approval
4. Start a new agent or workflow
5. Resume execution using the approved context

Each stage is a **separate execution**, not a resumed process.

#### Example Split Flow

```
Agent A (pre-approval)
 ├─ gathers information
 ├─ prepares action context
 ├─ calls Gantral
 └─ exits on WAITING_FOR_HUMAN

[human approval happens]

Agent B (post-approval)
 ├─ loads approved context
 ├─ performs the action
 └─ completes execution
```

#### Hard Rules for Non-Persistent Frameworks

If your framework does not support persistence:

* ❌ Do NOT keep agent state in memory
* ❌ Do NOT sleep, poll, or block locally
* ❌ Do NOT attempt to serialize internal agent state yourself
* ✅ Do split execution into pre-approval and post-approval stages

Failure to follow this pattern breaks auditability and is **not supported**.

---

### 5.3 Choosing the Correct Pattern

| Agent Framework Capability | Required Pattern    |
| -------------------------- | ------------------- |
| Native checkpointing       | Persisted Pause     |
| No checkpointing           | Split-Agent Pattern |

---

## 6. Required Consumer Behavior (Rules)

Consumers **must**:

* Treat Gantral decisions as authoritative
* Pause execution when instructed
* Persist state or split execution correctly
* Resume only when explicitly allowed
* Use the context returned by Gantral

Consumers **must not**:

* Retry locally while waiting for approval
* Execute actions speculatively
* Cache approval decisions
* Bypass Gantral on failure
* Store agent memory inside Gantral context

Violating these rules **breaks auditability**.

---

## 7. Timeouts and Escalation

Gantral may:

* Apply timeouts while waiting for approval
* Escalate to different approvers
* Terminate execution on timeout

As a consumer:

* You do **not** manage timeouts yourself
* You must react to Gantral’s final state

Timeout behavior is controlled **centrally**.

---

## 8. Error Handling Expectations

If Gantral returns an error:

* Treat it as **no approval granted**
* Do not execute the action
* Surface the error to operators or logs

Failing open is **never permitted**.

---

## 9. Example Reference

A complete working example is provided at:

```
/examples/agent-proxy
```

The example demonstrates:

* Persisted pause on `WAITING_FOR_HUMAN`
* Split-agent execution where persistence is unavailable
* Process exit and restart

Use it as a **reference**, not as a library.

---

## 10. What This Guide Intentionally Does Not Cover

This guide does **not** explain:

* Gantral’s internal state machine
* Policy authoring or Rego syntax
* Audit and replay internals
* Deployment or infrastructure setup

Refer to the **TRD** and **Implementation Guide** for those topics.

---

## 11. Agent Logs vs Execution Evidence

Agent-side logging is useful for debugging.

It is **not sufficient** for execution authority, auditability, or human accountability.

For workflows requiring defensible approvals:

* Tool inputs/outputs may be captured outside the agent
* Evidence must be immutable and independently attested
* Gantral consumes **references**, not raw payloads

Gantral never trusts agent logs as authoritative evidence.

### Optional Pattern: Evidence-Aware Approval

* Agent calls tool
* Runner captures raw I/O
* Evidence stored externally
* Evidence reference attached to execution context
* Human reviews → approves → execution resumes

#### When External Evidence Capture Is Recommended

* Regulatory, financial, or safety-critical approvals
* Decisions that must be defensible months later
* Irreversible actions

#### When Agent-Side Logging Is Sufficient

* Advisory or exploratory workflows
* Non-production environments
* Actions with no lasting side effects

---

## 12. Final Reminder

Gantral exists to ensure:

> **Actions happen only when an organization is willing to be accountable for them.**

As a consumer, your responsibility is simple:

**Ask first. Persist or split correctly. Act only when allowed.**
