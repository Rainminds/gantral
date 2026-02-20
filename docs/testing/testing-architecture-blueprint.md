# Testing Architecture Blueprint

Version: v2.0 (Baseline Enforcement Model)

---

# **1\. Purpose**

This document defines the authoritative testing architecture for Gantral.

It ensures that:

* All architectural invariants are mechanically enforced.  
* Authority exists exclusively as canonical workflow state.  
* Commitment artifacts are cryptographically and structurally enforced.  
* Replay is deterministic and independent of runtime, database, and logs.  
* Policy remains advisory and cannot assume authority.  
* Storage, identity, and federation cannot weaken admissibility.  
* Auditor guarantees are executable, not narrative.

Testing in Gantral is not QA.  
It is constitutional enforcement.

---

# **2\. Constitutional Guarantees**

Gantral testing must enforce the following guarantees:

1. Authority is represented exclusively as canonical workflow state.  
2. Only explicitly enumerated transitions are valid.  
3. Authority transition and artifact emission are atomic.  
4. Artifact hash binding is recursive and tamper-evident.  
5. Replay validates:  
   * Hash-chain integrity  
   * Transition validity  
   * workflow\_version\_id consistency  
   * policy\_version\_id consistency  
6. Replay requires no runtime, database, or logs.  
7. Policy evaluation is advisory only.  
8. Artifact store is authoritative for replay.  
9. Identity is verified at decision time.  
10. Unified visibility into running and paused workflows exists.  
11. Agent memory is never part of authority evidence.

All test layers must enforce these guarantees.

---

# **3\. Testing Philosophy (Non-Negotiable)**

## **3.1 Tests Enforce Invariants, Not Features**

Tests must verify:

* Authority is state.  
* HITL is blocking.  
* Policy is advisory only.  
* Artifact emission is atomic.  
* Replay depends only on artifact chain.  
* Fail-closed behavior everywhere.  
* No authority may exist outside canonical state transitions.  
* No reconstruction depends on trust.

## **3.2 Tests Must Be:**

* Deterministic  
* Replay-safe  
* Network-independent (for replay tests)  
* Hash-chain validated  
* Storage-order independent  
* Fail-closed by design  
* Adversarial

Correctness overrides convenience.

---

# **4\. Testing Layers**

Gantral testing is structured in nine enforcement layers.

---

## **Layer 1 – Pure Deterministic Unit Tests**

Location: `/tests/unit`

Covers:

* State transition validation  
* Transition guard correctness  
* Artifact payload construction  
* Recursive hash-chain construction  
* Canonicalization logic  
* Policy output interpretation mapping  
* Input schema validation  
* Error classification determinism

Properties:

* No network  
* No workflow runtime  
* No database dependency  
* Sub-100ms execution

---

## **Layer 2 – Canonical State Machine Matrix Tests**

Location: `/tests/state_machine`

Purpose:

Exhaustively validate the canonical transition relation.

Must enforce:

* All allowed transitions succeed.  
* All non-enumerated transitions terminate execution.  
* Duplicate transitions are rejected.  
* Skipped-state transitions are rejected.  
* Authority states (APPROVED / REJECTED / OVERRIDDEN) only reachable from WAITING\_FOR\_HUMAN.

### **Authority Exclusivity Enforcement**

Additional required tests:

* No artifact emission without a valid state transition.  
* No state transition to authority state without artifact emission.  
* Manual DB authority insert → replay INVALID.  
* Manual artifact injection without canonical state progression → replay INVALID.  
* Authority records cannot exist outside canonical workflow state.

---

## **Layer 3 – Artifact Integrity & Chain Enforcement**

Location: `/tests/artifact`

Must enforce:

* Valid artifact chain → VALID.  
* Hash mismatch → INVALID.  
* Tampered authority\_state → INVALID.  
* Tampered policy\_version\_id → INVALID.  
* Tampered workflow\_version\_id → INVALID.  
* Tampered justification → INVALID.  
* Tampered timestamp → INVALID.  
* Missing prev\_artifact\_hash → INVALID.  
* Duplicate artifact\_id → INVALID.  
* Circular chain reference → INVALID.  
* Chain truncation → INVALID.  
* Chain gap detection → INVALID.

### **Storage-Order Independence**

Replay must:

* Accept artifacts in arbitrary file order.  
* Reconstruct the chain deterministically.  
* Ignore extraneous artifacts not linked to the chain.  
* Reject reconstruction dependent on storage listing order.

---

## **Layer 4 – Replay & Verifier Tests**

Location: `/tests/replay`

Replay must validate:

* Hash-chain integrity.  
* Transition validity.  
* workflow\_version\_id consistency.  
* policy\_version\_id consistency.

Replay must:

* Require no database.  
* Require no runtime.  
* Ignore logs.  
* Ignore telemetry.  
* Ignore agent memory.

### **CLI Contract Locking**

Tests must validate:

* `gantral verify` outputs exactly: VALID / INVALID / INCONCLUSIVE.  
* Exit codes are deterministic.  
* No nondeterministic output ordering.  
* Replay emits no artifacts.  
* Replay produces no side effects.  
* Error classification identical across runs.

---

## **Layer 5 – Policy Advisory Boundary Enforcement**

Location: `/tests/policy`

Must enforce:

* ALLOW → continue.  
* REQUIRE\_HUMAN → WAITING\_FOR\_HUMAN.  
* DENY → TERMINATED.

### **Policy Boundary Guarantees**

Policy must never:

* Emit artifacts.  
* Approve execution directly.  
* Mutate workflow state.  
* Inject authority decisions.

Required tests:

* Malicious policy output containing artifact-like fields → rejected.  
* Policy returning approval decision directly → ignored.  
* Policy unavailability at high materiality → fail-closed.  
* REQUIRE\_HUMAN returned but no WAITING\_FOR\_HUMAN transition → execution fails.

### **Context Snapshot Binding Enforcement**

* Snapshot captured before policy evaluation.  
* Snapshot used for policy evaluation must equal snapshot hashed in artifact.  
* Mutation between evaluation and emission → execution aborts.  
* Snapshot must exclude agent internal memory.

---

## **Layer 6 – Integration Tests (Runtime \+ Policy \+ Storage)**

Location: `/tests/integration`

Covers:

* Full lifecycle execution.  
* Atomic artifact emission.  
* No hidden retries.  
* No execution past authority boundary without artifact persistence.

### **Storage Authority Boundary Enforcement**

* DB shows APPROVED but artifact missing → replay INVALID.  
* Artifact chain valid but DB missing → replay VALID.  
* Cross-region replication incomplete → replay INVALID.  
* Partial artifact persistence → execution aborts.  
* Replay cannot succeed if chain incomplete due to replication lag.

---

## **Layer 7 – Auditor Conformance Tests**

Location: `/tests/auditor`

Must prove:

* Authority decisions are cryptographically bound.  
* Policy version recorded in artifact.  
* Identity verified at decision time.  
* Context bound via hash.  
* Replay independent of logs and runtime.

Tests include:

* IdP offline → replay VALID.  
* Logs modified → replay unaffected.  
* Policy bundle unavailable → replay unaffected.  
* Organization rename → replay VALID.  
* Role deletion → replay VALID.  
* Database deleted → replay VALID.  
* Replay depends only on artifact chain.

---

## **Layer 8 – Unified Visibility Tests**

Location: `/tests/visibility`

Must enforce:

* WAITING\_FOR\_HUMAN instances visible.  
* Authority progression visible historically.  
* Visibility reflects canonical workflow state only.  
* Visibility cannot contradict artifact chain.  
* Visibility failure does not affect replay validity.

---

## **Layer 9 – Agent Isolation & Memory Boundary Tests**

Location: `/tests/agent`

Must enforce:

* Agent memory never included in artifact payload.  
* Attempt to include agent memory in snapshot → rejected.  
* Post-approval agent mutation does not affect replay.  
* Agent restart does not alter artifact chain.  
* Agent cannot bypass authority through retries.

---

# **5\. Coverage Categories**

Gantral test coverage must include:

* State coverage  
* Transition coverage  
* Authority exclusivity coverage  
* Artifact mutation coverage  
* Replay isolation coverage  
* Context snapshot binding coverage  
* Policy advisory boundary coverage  
* Storage-order independence coverage  
* Database contradiction coverage  
* Identity drift coverage  
* CLI contract coverage  
* Auditor independence coverage  
* Unified visibility coverage  
* Agent isolation coverage  
* Cross-region storage coverage  
* Hash-chain continuity  
* Canonicalization determinism

---

# **6\. CI Enforcement Gates**

CI must fail if:

* State transition coverage \< 100%.  
* Any authority state reachable without artifact emission.  
* Artifact schema changes without version bump.  
* Replay output differs from golden baseline.  
* Canonical bytes differ unexpectedly.  
* Policy advisory boundary violated.  
* CLI contract changes.  
* Storage-order replay differs.  
* Any nondeterminism detected.  
* Any ambiguity bypasses execution.  
* Any authority exists outside canonical state.

---

# **7\. Release Gate Requirements**

Before release:

* Full cross-language canonical equivalence validated.  
* Replay stable across previous workflow versions.  
* Golden artifact regression suite passes.  
* Cross-region storage replay validated.  
* Auditor conformance suite passes.  
* No schema drift without compatibility matrix update.  
* No canonicalization change without cross-language verification.  
* No workflow version binding drift.

---

# **8\. Constitutional Enforcement Statement**

With this blueprint:

* Authority cannot exist outside canonical workflow state.  
* Policy cannot assume authority.  
* Artifact chain cannot be reordered, truncated, or substituted.  
* Replay cannot depend on logs, database, or runtime.  
* Storage ordering cannot influence verification.  
* Identity drift cannot invalidate admissibility.  
* Agent memory cannot influence authority evidence.  
* Unified visibility reflects canonical state only.  
* CLI verification behavior is deterministic and stable.  
* Auditor guarantees are mechanically provable.

Gantral’s admissibility claim is therefore:

Not a design assertion.  
Not a documentation claim.  
But an executable invariant.

---

