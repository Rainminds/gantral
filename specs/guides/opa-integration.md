# **Gantral – Implementation Guide**

## **Policy Evaluation with OPA (Reference Implementation)**

**Status:** Reference implementation guide (non-normative)

**Audience:**

* Platform engineers  
* Core contributors  
* Enterprise implementers  
* Solution architects  
* Security, Risk, and Compliance reviewers

**Purpose:**  
This guide explains **how to implement Gantral’s policy evaluation layer using Open Policy Agent (OPA)** as a *reference implementation*, while remaining fully compliant with the Gantral **PRD** and **TRD**.

This guide is intentionally **non-normative**:

* It does **not** redefine Gantral’s architecture or execution semantics  
* It demonstrates *one correct way* to implement the pluggable policy evaluation interface

**Rule:** If this guide conflicts with the TRD, this guide is wrong.

---

## **1\. How This Guide Fits into Gantral**

### **What This Guide IS**

* A concrete implementation guide  
* A reference for contributors and adopters  
* A practical bridge between the PRD (intent) and TRD (constitution)  
* A starting point for real enterprise deployments

### **What This Guide IS NOT**

* A replacement for the TRD  
* A binding architectural dependency  
* A definition of Gantral’s execution semantics

OPA is used here as a **reference evaluator**, not as a core dependency.

---

## **2\. Recap: Correct Responsibility Split**

Gantral enforces a strict separation of concerns.

### **Gantral Owns**

* Execution state machine  
* HITL transitions  
* Authority capture  
* Audit and deterministic replay  
* Instance lifecycle

### **Policy Evaluator (OPA in this guide) Provides**

* Declarative policy evaluation  
* Advisory signals influencing execution transitions  
* Approver eligibility and escalation hints

### **Policy Evaluators Never**

* Pause execution  
* Approve actions  
* Override decisions  
* Write audit logs

They **only advise**. Gantral enforces.

---

## **3\. Where Policy Evaluation Happens (Critical)**

Policy evaluation occurs:

* **Synchronously** during execution transitions  
* As a **transition guard**, not a durable execution state  
* **Before** entering `WAITING_FOR_HUMAN`

Canonical flow:

RUNNING  
  └── policy evaluated  
        ├── ALLOW → continue execution  
        └── REQUIRE\_HUMAN → WAITING\_FOR\_HUMAN

There is **no `CHECK_POLICY` execution state**.

---

## **4\. Policy Evaluation Interface (Contract)**

Gantral exposes an internal policy evaluation interface. The interface is evaluator-agnostic; OPA is one implementation.

### **4.1 Input (Structured & TRD-Aligned Schema)**

Policy evaluators receive **structured context**, never raw prompts.

```json
{  
  "instance": {  
    "instance\_id": "uuid",  
    "workflow\_id": "string",  
    "workflow\_version": "string",  
    "materiality": "LOW | MEDIUM | HIGH",  
    "owning\_team\_id": "string"  
  },  
  "execution": {  
    "current\_state": "RUNNING",  
    "step": "string",  
    "cost\_estimate": 123.45  
  },  
  "context": {  
    "domain": "sdlc | finance | operations | support",  
    "trigger": "PR | INCIDENT | PAYMENT | OTHER",  
    "attributes": { }  
  },  
  "identity": {  
    "actor\_id": "service\_account\_or\_human\_id",  
    "roles": \["AGENT\_EXECUTOR"\]  
  },  
  "policy": {  
    "policy\_version\_id": "string"  
  }  
}
```

**Notes (Critical):**

* `actor_id` naming is aligned exactly with TRD terminology  
* `policy_version_id` is **explicitly included** to guarantee deterministic replay
* **Hashing Requirement:** The full JSON input payload must be hashed (SHA-256) and bound to the resulting `CommitmentArtifact`.

---

### **4.2 Output (Advisory Decision Signals)**

```json
{  
  "decision": "ALLOW | REQUIRE\_HUMAN | DENY",  
  "approver\_roles": \["TECH\_LEAD", "RISK\_OFFICER"\],  
  "timeout": "30m",  
  "escalation\_roles": \["PLATFORM\_ADMIN"\]  
}
```

OPA returns signals only. Gantral interprets and enforces them.

---

## **5\. Using OPA as the Policy Evaluator**

OPA evaluates policies written in Rego using the input schema above.

### **5.1 Deployment Models**

OPA may be deployed as:

* Sidecar to the Gantral execution engine  
* Centralized policy service  
* Embedded library (advanced / optional)

Sidecar deployment is recommended initially for simplicity and isolation.

---

## **6\. Example Rego Policy (Illustrative)**

```rego
package gantral.policy

default decision \= "ALLOW"

decision \= "REQUIRE\_HUMAN" {  
  input.instance.materiality \== "HIGH"  
}

decision \= "DENY" {  
  input.execution.cost\_estimate \> 100000  
  not input.context.attributes.override\_allowed  
}

approver\_roles \= \["TECH\_LEAD"\] {  
  input.instance.materiality \== "MEDIUM"  
}

approver\_roles \= \["RISK\_OFFICER"\] {  
  input.instance.materiality \== "HIGH"  
}
```

---

## **7\. Mapping OPA Output to Execution Transitions (Updated)**

| OPA Signal | Gantral Behavior |
| ----- | ----- |
| `ALLOW` | Continue execution |
| `REQUIRE_HUMAN` | Transition to `WAITING_FOR_HUMAN` |
| `DENY` | Transition to `TERMINATED` |
| `timeout` present | Schedule a `TIMEOUT` execution event |

### **Timeout Semantics (Critical)**

If a `timeout` value is present:

1. Gantral schedules a `TIMEOUT` event at the specified duration  
2. Upon expiry, **if the instance is still in `WAITING_FOR_HUMAN`**:  
   * **Fail-closed:** transition to `TERMINATED`, or  
   * **Escalate:** emit an `ESCALATION` event and update eligible approver roles using `escalation_roles`

The exact behavior is determined by policy configuration.

Timeout handling is **execution control**, not policy evaluation.

---

## **8\. HITL Handling (Gantral-Owned)**

Once in `WAITING_FOR_HUMAN`:

* Gantral determines eligible approvers  
* Gantral captures decision, role, justification, and context  
* Gantral enforces resume, override, or termination  
* Gantral seals the audit record

OPA is **not invoked** during human decision capture.

---

## **9\. CI/CD and Offline Validation (Optional)**

OPA policies may be:

* Unit tested  
* Validated in CI pipelines  
* Reviewed independently of Gantral code

This is recommended but **non-blocking** to runtime execution.

---

## **10\. Failure Modes & Guarantees (Clarified)**

### **Policy Evaluator Unavailable**

* HIGH materiality → fail closed  
* MEDIUM materiality → configurable  
* LOW materiality → fail open allowed

All failures emit explicit execution events.

### **Determinism Guarantee**

* Policy inputs are immutable snapshots  
* `policy_version_id` is recorded with every evaluation  
* Replay uses the **same policy version**, never the latest

This prevents replay drift under policy evolution.

---

## **11\. Extending Beyond OPA**

Because Gantral depends only on the policy interface:

* Enterprises may replace OPA with internal policy engines  
* Regulators may mandate specific evaluators  
* Multiple evaluators may coexist

Gantral remains unchanged.

---

## **12\. How and Where to Use This Guide**

### **Who Should Read This**

* Contributors implementing the policy layer  
* Enterprises deploying self-hosted Gantral  
* Security and compliance reviewers  
* Partners building adapters or extensions

### **Where This Guide Lives**

Recommended locations:

* `/docs/implementation/policy-evaluation-opa.md`  
* Public documentation site (advanced section)  
* Reference repository README

### **When to Use It**

* After reading the TRD  
* Before implementing the policy layer  
* During enterprise architecture reviews  
* During audits to explain policy–execution separation

---

## **13\. Final Reminder**

OPA is an **implementation choice**.

Gantral is the **execution authority**.

Humans remain accountable.

This guide exists to help teams build correctly — not to redefine the system.
