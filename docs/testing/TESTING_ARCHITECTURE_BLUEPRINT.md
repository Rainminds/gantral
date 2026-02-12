# Testing Architecture Blueprint

Version: v1.0 (Baseline Enforcement Model)

---

## **1\. Purpose**

This document defines the **authoritative testing architecture** for Gantral.

It ensures that:

* All TRD invariants are enforced by executable tests  
* All PRD claims are mechanically verifiable  

Testing in Gantral is not QA.  
It is **constitutional enforcement**.

---

## **2\. Testing Philosophy (Non-Negotiable)**

### **2.1 Tests Enforce Invariants, Not Features**

Tests must verify:

* Authority is state  
* HITL is blocking  
* Policy is advisory only  
* Artifact emission is atomic  
* Replay depends only on artifact chain  
* Fail-closed behavior everywhere

---

### **2.2 Tests Must Be:**

* Deterministic  
* Replay-safe  
* Network-independent (for replay tests)  
* Hash-chain validated  
* Fail-closed by design  
* Adversarial

---

## **3\. Testing Layers**

Gantral testing is structured in 7 layers:

---

### **Layer 1 – Pure Unit Tests (Deterministic Logic)**

Location:

`/tests/unit`

Covers:

* State machine transitions  
* Artifact hash construction  
* Artifact hash chain validation  
* Policy interpretation mapping  
* Input validation

Properties:

* No network  
* No Temporal runtime  
* No DB dependency  
* \<100ms execution

---

### **Layer 2 – Authority State Machine Matrix Tests**

Location:

`/tests/state_machine`

Purpose:  
Exhaustively validate canonical transitions defined in TRD


.

All allowed transitions must succeed.  
All invalid transitions must panic or terminate.

Table-driven.

---

### **Layer 3 – Artifact Integrity & Tamper Tests**

Location:

`/tests/artifact`

Covers:

* Valid artifact chain  
* Hash mismatch  
* Corrupted artifact payload  
* Missing previous hash  
* Timestamp mutation  
* Authority state mutation  
* Policy version mutation  
* Actor mutation

Replay must return:

* VALID  
* INVALID  
* INCONCLUSIVE (where defined)

Aligned to Phase 6 requirements


.

---

### **Layer 4 – Replay & Verifier Tests**

Location:

`/tests/replay`

Tests:

* Replay with correct chain  
* Replay without Gantral service  
* Replay without DB  
* Replay after DB deletion  
* Replay after workflow code change  
* Replay under modified agent code

Replay must depend only on artifact chain.

---

### **Layer 5 – Policy Interface Tests**

Location:

`/tests/policy`

Validates:

* Input schema correctness  
* ALLOW behavior  
* REQUIRE\_HUMAN behavior  
* DENY behavior  
* Timeout scheduling  
* Escalation update  
* Policy engine unavailability behavior (materiality-based fail-open/closed)

Policy must never:

* Pause execution directly  
* Approve execution directly  
* Emit artifacts

---

### **Layer 6 – Integration Tests (Temporal \+ OPA)**

Location:

`/tests/integration`

Tests full execution lifecycle:

* Trigger instance  
* Policy guard  
* WAITING\_FOR\_HUMAN  
* Approval  
* Artifact emission  
* Resume  
* Completion

Verifies:

* Atomic artifact emission  
* No execution without artifact  
* No hidden retries

---

### **Layer 7 – Adversarial & Failure Tests**

Location:

`/tests/adversarial`

Tests:

* Artifact substitution  
* Log substitution  
* Missing artifact  
* Partial artifact  
* Authority ambiguity  
* Policy ambiguity  
* Concurrent approval attempts  
* Escalation owner unavailable  
* Timeout race conditions

Must prove:

Gantral refuses ambiguity.

---

## **4\. Test Naming Convention**

`Test_StateMachine_InvalidTransition_Panics`  
`Test_Artifact_TamperedHash_Invalid`  
`Test_Replay_WithoutDB_Succeeds`  
`Test_Policy_DENY_TransitionsToTerminated`  
`Test_HITL_NoDecision_BlocksExecution`  
`Test_Execution_NoArtifact_FailsClosed`

No clever names. Explicit intent only.

---

## **5\. Coverage Categories**

Gantral test coverage must include:

* State coverage  
* Transition coverage  
* Failure coverage  
* Artifact mutation coverage  
* Replay isolation coverage  
* Timeout coverage  
* Escalation coverage  
* Identity coverage  
* Schema backward compatibility  
* Hash chain continuity

---

## **6\. Required Test Enforcement Gates**

CI must fail if:

* State transition coverage \< 100%  
* Artifact schema changes without version bump  
* Replay tests fail  
* Fail-closed tests fail  
* Hash mutation tests pass unexpectedly  
* Any nondeterminism detected

## **7\. Test Organization

### Co-located Unit Tests (`*_test.go`)
Location: Next to production code (e.g., `core/engine/machine_test.go`)
Purpose: Test individual functions/methods in isolation
Scope: Single package, internal API, fast (<10ms), no dependencies
Package: Same as production code (e.g., `package engine`)
Coverage: Edge cases, validation, error handling, business logic

### Organized Tests (`/tests/*`)
Location: Centralized test directory
Purpose: Test system-level behavior, invariants, integration
Scope: Multi-package, external API, requires dependencies
Package: External test package (e.g., `package statemachine_test`, `package integration_test`)
Coverage: State machine matrix, artifact integrity, replay, E2E

### Decision Rule
- If testing **internal function** → co-located `*_test.go`
- If testing **system invariant** → `/tests/` organized
- If testing **cross-package** → `/tests/` organized
- If requiring **DB/Temporal/OPA** → `/tests/integration/`
