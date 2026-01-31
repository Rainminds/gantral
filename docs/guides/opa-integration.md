---
title: OPA Integration Reference
sidebar_label: OPA Integration
---

# OPA Integration Reference

## Policy Evaluation with Open Policy Agent (OPA)

**Status:** Reference implementation guide (non-normative)  

**Audience:**
- Platform engineers  
- Core contributors  
- Enterprise implementers  
- Solution architects  
- Security, Risk, and Compliance reviewers  

---

## Purpose

This document explains **how to implement Gantral’s policy evaluation layer using Open Policy Agent (OPA)** as a *reference implementation*.

OPA is used here to demonstrate **one correct way** to implement Gantral’s **pluggable policy evaluation interface**, while remaining fully compliant with the Gantral **PRD v5.1** and **TRD v4.1**.

This guide is intentionally **non-normative**.

> **Rule:** If this guide conflicts with the TRD, this guide is wrong.

---

## 1. How This Guide Fits into Gantral

### What This Guide IS

- A concrete implementation reference  
- A bridge between PRD intent and TRD guarantees  
- A practical example for contributors and adopters  
- Suitable for enterprise and regulated environments  

### What This Guide IS NOT

- A definition of Gantral’s execution semantics  
- A requirement to use OPA  
- A policy authoring standard  
- A replacement for the TRD  

OPA is an **implementation choice**, not a dependency.

---

## 2. Responsibility Split (Non-Negotiable)

Gantral enforces a strict separation of concerns.

### Gantral Owns

- Authority state machine  
- Execution state transitions  
- Human-in-the-Loop (HITL) enforcement  
- Authority capture and commitment artifacts  
- Deterministic replay of authority decisions  

### Policy Evaluator (OPA) Provides

- Declarative policy evaluation  
- Advisory signals during execution transitions  
- Approver eligibility and escalation hints  

### Policy Evaluators Never

- Pause execution  
- Approve actions  
- Override decisions  
- Resume workflows  
- Write audit or authority records  

**Policy advises. Authority enforces.**

---

## 3. Where Policy Evaluation Occurs

Policy evaluation occurs:

- **Synchronously** during execution transitions  
- As a **transition guard**, not a state  
- **Before** entering `WAITING_FOR_HUMAN`

Canonical flow:

```

RUNNING
└─ policy evaluated
├─ ALLOW → continue execution
├─ REQUIRE_HUMAN → WAITING_FOR_HUMAN
└─ DENY → TERMINATED

````

There is **no `CHECK_POLICY` execution state**.

---

## 4. Policy Evaluation Interface (Contract)

Gantral exposes an **evaluator-agnostic policy interface**.
OPA is one implementation.

### 4.1 Input Schema (Structured, Immutable)

Policy evaluators receive **structured context**, never raw prompts or tool payloads.

```json
{
  "instance": {
    "instance_id": "uuid",
    "workflow_id": "string",
    "workflow_version": "string",
    "materiality": "LOW | MEDIUM | HIGH",
    "owning_team_id": "string"
  },
  "execution": {
    "current_state": "RUNNING",
    "step": "string",
    "cost_estimate": 123.45
  },
  "context": {
    "domain": "sdlc | finance | operations | support",
    "trigger": "PR | INCIDENT | PAYMENT | OTHER",
    "attributes": {}
  },
  "identity": {
    "actor_id": "service_or_human_id",
    "roles": ["AGENT_EXECUTOR"]
  },
  "policy": {
    "policy_version_id": "string"
  }
}
````

**Critical guarantees:**

* Inputs are immutable snapshots
* `policy_version_id` is always explicit
* No evaluator sees agent memory or tool payloads

---

### 4.2 Output Schema (Advisory Only)

```json
{
  "decision": "ALLOW | REQUIRE_HUMAN | DENY",
  "approver_roles": ["TECH_LEAD", "RISK_OFFICER"],
  "timeout": "30m",
  "escalation_roles": ["PLATFORM_ADMIN"]
}
```

OPA emits **signals only**.
Gantral interprets and enforces them.

---

## 5. Using OPA as the Evaluator

OPA evaluates Rego policies using the schema above.

### 5.1 Deployment Options

OPA may be deployed as:

* Sidecar to the Gantral execution engine (recommended initially)
* Centralized policy service
* Embedded library (advanced / optional)

Deployment choice does **not** change semantics.

---

## 6. Example Rego Policy (Illustrative)

```rego
package gantral.policy

default decision = "ALLOW"

decision = "REQUIRE_HUMAN" {
  input.instance.materiality == "HIGH"
}

decision = "DENY" {
  input.execution.cost_estimate > 100000
  not input.context.attributes.override_allowed
}

approver_roles = ["TECH_LEAD"] {
  input.instance.materiality == "MEDIUM"
}

approver_roles = ["RISK_OFFICER"] {
  input.instance.materiality == "HIGH"
}
```

This example is illustrative only.
Gantral does not interpret Rego semantics.

---

## 7. Mapping Policy Signals to Execution Transitions

| Policy Signal     | Gantral Behavior                   |
| ----------------- | ---------------------------------- |
| `ALLOW`           | Continue execution                 |
| `REQUIRE_HUMAN`   | Transition to `WAITING_FOR_HUMAN`  |
| `DENY`            | Transition to `TERMINATED`         |
| `timeout` present | Schedule `TIMEOUT` execution event |

### Timeout Semantics

If a `timeout` is provided:

1. Gantral schedules a `TIMEOUT` event
2. If still in `WAITING_FOR_HUMAN` at expiry:

   * **Fail-closed** (default), or
   * **Escalate** using `escalation_roles`

Timeout handling is **execution control**, not policy evaluation.

---

## 8. HITL Handling (Gantral-Owned)

Once in `WAITING_FOR_HUMAN`:

* Gantral determines eligible approvers
* Gantral captures decision, role, justification, and context
* Gantral enforces resume, override, or termination
* Gantral emits the commitment artifact

OPA is **not invoked** during human decision capture.

---

## 9. CI/CD and Offline Policy Validation (Optional)

OPA policies may be:

* Unit tested
* Validated in CI
* Reviewed independently of Gantral code

This is recommended but **not required** at runtime.

---

## 10. Failure Modes and Guarantees

### Evaluator Unavailable

* HIGH materiality → fail closed
* MEDIUM materiality → configurable
* LOW materiality → fail open allowed

All failures emit explicit execution events.

### Determinism & Replay

* Policy inputs are immutable
* `policy_version_id` is recorded with each evaluation
* Replay uses the **same policy version**
* Latest policy is never substituted during replay

Replay divergence is therefore detectable.

---

## 11. Beyond OPA

Because Gantral depends only on the policy interface:

* Enterprises may use internal policy engines
* Multiple evaluators may coexist
* Regulator-mandated evaluators can be substituted

Gantral remains unchanged.

---

## 12. When to Use This Guide

This guide is most useful:

* After reading the TRD
* Before implementing the policy layer
* During enterprise architecture review
* During audits explaining policy vs authority separation

---

## Final Reminder

OPA is an **advisor**.

Gantral is the **authority**.

Humans remain accountable.

This guide exists to help teams implement policy evaluation **correctly**,
not to redefine how Gantral works.

