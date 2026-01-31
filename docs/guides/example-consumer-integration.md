---
title: Example Consumer Guide
---

# Example Consumer Guide

⚠️ **This document is not a system specification.**

It is the **normative usage guide** for external consumers (agents, tools, workflows)
that integrate with Gantral.

---

## How External Agents and Systems Should Use Gantral

**Status:** Consumption guide (normative for consumers)  
**Audience:**
- Agent developers
- Platform engineers integrating Gantral
- Teams building tools, automations, or workflows that require human authority

---

## Purpose

This guide explains **when and how external systems must consume Gantral**.

It answers **one question only**:

> **When and how should my agent or system ask Gantral before executing an action?**

This guide is intentionally narrow.  
It does **not** explain Gantral’s internal architecture or enforcement mechanics.

---

## 1. What Gantral Is (From a Consumer’s Perspective)

From the outside, Gantral is an **execution authority gate**.

Gantral:

- Receives a request to perform a potentially material action
- Evaluates whether execution may proceed
- Pauses execution when authority is required
- Resumes or terminates execution based on an explicit decision

Gantral is **not**:

- An agent framework
- A task runner
- A policy engine
- A workflow orchestrator

As a consumer, you treat Gantral as **the authority**, not as infrastructure you control.

---

## 2. When You Must Call Gantral

You must call Gantral **before executing any action** that:

- Has irreversible side effects
- Affects production systems
- Has regulatory, financial, or security impact
- Requires human accountability

### Examples

- Deleting or mutating production infrastructure
- Executing high-cost or high-risk AI operations
- Triggering payments, refunds, or financial transfers
- Sending external communications
- Making customer-impacting changes

If an action is fully reversible, sandboxed, or advisory-only,
Gantral is usually **not required**.

---

## 3. Canonical Consumer Interaction Flow

From a consumer’s perspective, the flow is always:

1. Prepare execution context
2. Ask Gantral for authority
3. Observe Gantral’s execution state
4. Either wait, resume, or abort

> **A consumer must never execute an action until Gantral explicitly allows it.**

---

## 4. Minimal SDK Interaction Pattern

Exact SDK names vary.  
**The semantics do not.**

### Step 1: Create or Identify an Execution Instance

```python
instance = gantral.create_instance(
    workflow_id,
    workflow_version,
    execution_context
)
````

This represents **one immutable execution attempt**.

---

### Step 2: Request Authority

```python
state = gantral.request_authority(
    instance_id = instance.id,
    actor_id = "agent-or-service-id",
    action_context
)
```

Gantral evaluates policy and execution state
and returns the **current authoritative execution state**.

---

### Step 3: Handle the Execution State

Gantral will place the instance in one of the following states.

---

#### APPROVED / RUNNING

* You may proceed with execution
* You **must** use the context returned by Gantral

```python
execute(action)
```

---

#### WAITING_FOR_HUMAN

* **Do not** execute the action
* Pause execution immediately
* Follow an approved pause pattern (see Section 5)

```python
wait_for_update(instance.id)
```

---

#### REJECTED or TERMINATED

* Do not execute the action
* Abort the workflow
* Do not retry

```python
abort("Execution not authorized")
```

---

## 5. Handling Long-Running Approvals (Critical)

Approvals may take **minutes, hours, or days**.

❌ Do **not** rely on local RAM, in-process state, or sleeping threads.

Gantral supports long-running waits via **two valid consumer patterns only**.

---

### 5.1 Pattern A — Persisted Pause (Preferred)

Use this pattern **only if your agent framework supports native persistence / checkpointing**.

When Gantral enters `WAITING_FOR_HUMAN`:

1. Persist agent state using the framework’s native mechanism
2. Exit the process cleanly (zero CPU usage)
3. Wait for Gantral to signal approval
4. Start a new process
5. Restore agent state from persistence
6. Resume execution

Gantral manages **authority**.
Your framework manages **agent memory**.

#### Framework Examples

* CrewAI Flows → native persistence (SQLite / Postgres)
* LangGraph → checkpointing (SQLite / Redis / S3)

> **Rule:**
> If your agent cannot restart from a checkpoint, you must not use this pattern.

---

### 5.2 Pattern B — Split-Agent Pattern (Mandatory Without Persistence)

If your agent framework **cannot** resume from persisted state,
you must **split execution into stages**.

When Gantral enters `WAITING_FOR_HUMAN`:

1. Terminate the current agent execution
2. Persist only minimal handoff context (IDs, references, inputs)
3. Wait for Gantral approval
4. Start a new agent or workflow
5. Resume execution using the approved context

Each stage is a **separate execution**, not a resumed process.

#### Example Split Flow

```
Agent A (pre-approval)
 ├─ gather information
 ├─ prepare action context
 ├─ call Gantral
 └─ exit on WAITING_FOR_HUMAN

[human authority decision]

Agent B (post-approval)
 ├─ load approved context
 ├─ execute the action
 └─ complete
```

#### Hard Rules for Non-Persistent Frameworks

If your framework does not support persistence:

* ❌ Do NOT keep state in memory
* ❌ Do NOT sleep or poll locally
* ❌ Do NOT attempt to serialize internal agent memory
* ✅ Split execution cleanly

Violating these rules breaks auditability and is unsupported.

---

### 5.3 Choosing the Correct Pattern

| Agent Capability     | Required Pattern    |
| -------------------- | ------------------- |
| Native checkpointing | Persisted Pause     |
| No checkpointing     | Split-Agent Pattern |

---

## 6. Required Consumer Rules (Non-Negotiable)

Consumers **must**:

* Treat Gantral execution state as authoritative
* Pause immediately when instructed
* Persist or split execution correctly
* Resume only when explicitly allowed
* Use the execution context returned by Gantral

Consumers **must not**:

* Execute speculatively
* Retry locally while waiting
* Cache approvals
* Bypass Gantral on failure
* Store agent memory inside Gantral context

Violations **invalidate auditability**.

---

## 7. Timeouts and Escalation

Gantral may:

* Apply approval timeouts
* Escalate to different approvers
* Terminate execution on timeout

As a consumer:

* You do **not** manage timeouts
* You react only to Gantral’s execution state

Timeout behavior is centrally enforced.

---

## 8. Error Handling Expectations

If Gantral returns an error or ambiguous state:

* Treat it as **no authority granted**
* Do not execute the action
* Surface the failure to operators or logs

Fail-open behavior is **never permitted**.

---

## 9. Example Reference Implementation

A working example is provided at:

```
/examples/agent-proxy
```

The example demonstrates:

* Persisted pause on `WAITING_FOR_HUMAN`
* Split-agent execution where persistence is unavailable
* Process exit and restart semantics

Use it as a **reference**, not as a dependency.

---

## 10. What This Guide Intentionally Does Not Cover

This guide does **not** explain:

* Gantral’s internal state machine
* Policy authoring or Rego syntax
* Audit or replay internals
* Deployment or infrastructure setup

Refer to the **TRD** and **Implementation Guide** for those topics.

---

## 11. Agent Logs vs Execution Evidence

Agent-side logs are useful for debugging.

They are **not sufficient** for execution authority or accountability.

For workflows requiring defensible approvals:

* Tool inputs/outputs may be captured externally
* Evidence must be immutable and independently attested
* Gantral consumes **references**, not raw payloads

Gantral never treats agent logs as authoritative evidence.

### Optional Pattern: Evidence-Aware Approval

* Agent invokes tool
* Runner captures raw I/O
* Evidence stored externally
* Evidence reference attached to execution context
* Human reviews → approves → execution resumes

#### When External Evidence Capture Is Recommended

* Regulatory or financial decisions
* Irreversible actions
* Approvals that must be defensible months later

#### When Agent Logs Are Sufficient

* Advisory workflows
* Non-production environments
* Low-impact actions

---

## 12. Final Reminder

Gantral exists to ensure:

> **Actions occur only when an organization is willing to be accountable for them.**

As a consumer, your responsibility is simple:

**Ask first.
Pause correctly.
Act only when allowed.**

