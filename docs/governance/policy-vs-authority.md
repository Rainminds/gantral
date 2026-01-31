---
title: Policy vs Authority
---

# Policy vs Authority

This document defines the **explicit separation between policy evaluation and execution authority**
in Gantral.

This separation is not an architectural preference.
It is a **requirement for verifiability**.

If policy and authority are conflated, execution decisions cannot be independently proven.

---

## Core Distinction

Gantral distinguishes between:

- **Policy** — advisory evaluation of conditions and constraints  
- **Authority** — enforced permission to execute

Only authority allows execution to proceed.

---

## What Policy Is (Advisory)

Policy engines (such as **OPA**) are used to:

- evaluate intent and constraints
- apply organizational rules and thresholds
- determine whether escalation or human review is required

Policy evaluation produces **signals**, not permissions.

Typical policy outputs:
- `ALLOW`
- `REQUIRE_HUMAN`
- `DENY`

These outputs **do not**:
- pause execution
- resume execution
- terminate execution
- grant authority

Policy engines do not control execution.
They advise Gantral on **whether authority may be required**.

---

## What Authority Is (Enforced)

Authority is the **ability for execution to proceed**.

In Gantral, authority:

- exists only as execution state
- is enforced by the authority state machine
- is exercised explicitly by humans or by pre-approved automated paths
- produces a commitment artifact when exercised

Execution proceeds **only** when authority is present.

---

## Where Authority Lives

Authority lives in the **execution state machine**, not in policy.

Examples:
- `RUNNING` — execution is authorized
- `WAITING_FOR_HUMAN` — authority is absent
- `RESUMED` — authority has been explicitly restored
- `TERMINATED` — authority has been revoked

There is no concept of “implicit” or “soft” authority.

---

## Why This Separation Is Required

### 1. Verifiability Requires Determinism

Policy evaluation may change over time:
- policies evolve
- interpretations shift
- evaluation logic is updated

Authority must remain **stable and replayable**.

By recording authority as state, Gantral ensures:
- replay does not depend on re-evaluating policy
- authority can be proven years later
- verification does not require policy engines to be available

---

### 2. Policy Engines Are Not Evidence

Policy engines:
- do not emit immutable artifacts
- do not bind execution context
- do not record authority decisions

A policy result alone is **not evidence**.

Only an execution-time authority transition, recorded as state and bound into
a commitment artifact, is verifiable.

---

### 3. Preventing Silent Bypass

If policy engines were allowed to grant authority directly:

- execution could proceed without a durable record
- authority could be inferred post-hoc
- replay would require trusting logs or narratives

By separating policy from authority:
- all execution paths are forced through the authority state machine
- all authority decisions are recorded
- bypass is structurally prevented

---

## How Policy and Authority Interact

1. Execution reaches a controlled boundary
2. Policy is evaluated as a **guard**
3. Policy returns an advisory signal
4. Gantral enforces authority by:
   - continuing execution
   - pausing for human review
   - terminating execution

Policy never enforces execution.
Gantral always does.

---

## Relationship to the Authority State Machine

The authority state machine defines:
- when execution may proceed
- when execution must stop
- when human authority is required

Policy evaluation may trigger a transition to `WAITING_FOR_HUMAN`,
but it cannot transition execution out of that state.

This invariant is required for verifiability.

---

## Relationship to OPA Integrations

Gantral commonly integrates with **OPA** as a policy evaluation engine.

OPA:
- evaluates rules
- returns decisions
- remains stateless with respect to execution authority

Gantral:
- enforces authority
- records authority transitions
- emits commitment artifacts

OPA advises.
Gantral enforces.

---

## Failure Behavior

If policy evaluation:
- fails
- times out
- produces conflicting results

Gantral must:
- refuse execution
- fail closed
- surface the ambiguity

Execution must never proceed on ambiguous policy signals.

---

## Guiding Principle

Policy determines **what should happen**.

Authority determines **what is allowed to happen**.

Only authority is enforceable.
Only authority is verifiable.
