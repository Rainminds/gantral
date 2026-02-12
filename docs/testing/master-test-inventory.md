# Exhaustive Master Test Inventory

Version: v1.0 (Baseline Enforcement Model)

---

# **Purpose**

This document defines the complete, adversarially hardened, operationally mature test inventory for Gantral.

It ensures:

* Authority invariants are mechanically enforced  
* Replay is deterministic and cryptographically stable  
* Canonicalization is explicit and versioned  
* Artifact integrity survives hostile reconstruction  
* Identity drift does not weaken admissibility  
* Storage and concurrency failures are fail-closed  
* Verifier cannot be exploited  
* Schema evolution does not silently weaken guarantees  
* Temporal runtime behavior remains deterministic  
* Federated runners operate safely  
* Agent integrations preserve execution isolation  
* Performance and load do not compromise correctness  
* Observability is correct and non-leaking  
* Disaster recovery preserves admissibility  
* Compatibility and upgrade paths do not introduce drift  
* Chaos conditions do not violate authority guarantees

If a test class is not defined here, it is considered out of scope.

---

# **SECTION A – State Machine Tests**

## **Valid Transitions**

* CREATED → RUNNING  
* RUNNING → WAITING\_FOR\_HUMAN  
* WAITING\_FOR\_HUMAN → APPROVED  
* WAITING\_FOR\_HUMAN → REJECTED  
* WAITING\_FOR\_HUMAN → OVERRIDDEN  
* APPROVED → RESUMED  
* OVERRIDDEN → RESUMED  
* RESUMED → RUNNING  
* RUNNING → COMPLETED  
* RUNNING → TERMINATED

## **Invalid Transitions (Must Panic or Terminate)**

* RUNNING → APPROVED  
* RUNNING → REJECTED  
* COMPLETED → RUNNING  
* TERMINATED → RUNNING  
* WAITING\_FOR\_HUMAN → RUNNING  
* CREATED → COMPLETED  
* Any transition from COMPLETED  
* Any transition from TERMINATED  
* Duplicate identical transition  
* Transition skipping intermediate state

---

# **SECTION B – HITL Tests**

* No execution without valid decision  
* Approval requires non-empty justification  
* Override requires context\_delta  
* Concurrent approvals → only first valid  
* Concurrent override \+ reject → deterministic resolution  
* Missing identity → reject  
* Role mismatch → reject  
* Escalation updates eligible roles correctly  
* Timeout before approval → follow policy behavior  
* Rejection → TERMINATED  
* Override → RESUMED with context\_delta applied  
* Approval after timeout → rejected  
* Approval after termination → rejected  
* Approval by non-eligible role → rejected  
* Decision captures identity snapshot  
* Decision captures role snapshot  
* HITL signal during replay handled deterministically  
* HITL query during WAITING\_FOR\_HUMAN consistent snapshot  
* SLA breach metric emitted correctly  
* Approval after workflow restart deterministic  
* Whitespace-only justification rejected  
* Empty string rejected  
* Justification must exceed configurable minimum length (if enforced)  
* Override requires both context\_delta and justification

---

# **SECTION C – Policy Engine Integration Tests**

* ALLOW → continue  
* REQUIRE\_HUMAN → WAITING\_FOR\_HUMAN  
* DENY → TERMINATED  
* Timeout present → schedule TIMEOUT event  
* Escalation roles applied correctly  
* Policy version recorded in decision  
* Policy evaluator unavailable (high materiality) → fail closed  
* Policy evaluator unavailable (low materiality) → configurable  
* Policy cannot emit artifacts  
* Policy cannot mutate execution state  
* Policy cannot override authority  
* OPA Rego syntax error → fail closed  
* OPA timeout → fail closed  
* OPA service restart mid-evaluation → fail closed  
* Policy schema mismatch → reject  
* Unknown policy decision → fail closed  
* Policy output strictly validated  
* Multiple policy conflicts resolved deterministically  
* Policy dry\_run does not alter authority  
* Policy deterministic across repeated evaluation

---

# **SECTION D – Artifact Tests**

* Artifact emitted on APPROVED  
* Artifact emitted on REJECTED  
* Artifact emitted on OVERRIDDEN  
* No artifact on ALLOW  
* Artifact hash correct  
* prev\_artifact\_hash correct  
* Artifact version embedded  
* Canonicalization version embedded  
* Hash algorithm version embedded  
* Tampered authority\_state → INVALID  
* Tampered policy\_version\_id → INVALID  
* Tampered timestamp → INVALID  
* Missing artifact → execution fails  
* Artifact emission failure → execution aborts  
* Artifact write failure → execution aborts  
* No retry on artifact failure  
* Duplicate artifact emission rejected  
* Artifact emission atomic with state transition  
* Artifact immutability enforced  
* Artifact hash computed before persistence  
* No partial artifact commit  
* Justification included in context\_hash  
* Tampered justification → INVALID replay  
* Justification canonicalization deterministic

---

# **SECTION E – Replay Tests**

* Replay valid chain → VALID  
* Replay modified artifact → INVALID  
* Replay missing artifact → INVALID  
* Replay without DB → VALID  
* Replay without workflow runtime → VALID  
* Replay after workflow refactor → VALID  
* Replay ignores agent memory  
* Replay ignores logs  
* Replay ignores telemetry  
* Replay deterministic across runs  
* Replay fails if canonicalization version mismatched  
* Replay fails if hash algorithm mismatched  
* Replay fails if artifact schema unknown  
* Replay after runtime version upgrade consistent  
* Replay under alternate JSON parser consistent  
* Replay under load deterministic  
* Replay performance bounded

---

# **SECTION F – Fail-Closed Tests**

* Ambiguous authority → TERMINATED  
* Missing policy decision → TERMINATED  
* Missing artifact → TERMINATED  
* Hash mismatch → INVALID  
* Partial artifact chain → INVALID  
* Corrupt artifact JSON → INVALID  
* Unknown state in artifact → INVALID  
* Unknown enum value → INVALID  
* Incomplete artifact → INVALID  
* Nil pointer in transition → terminate  
* Missing canonicalization version → terminate

---

# **SECTION G – Identity & Security Tests**

* OIDC identity required  
* No local user accepted  
* Service identity accepted  
* Team isolation enforced  
* Role mismatch denied  
* OIDC token expiry during WAITING\_FOR\_HUMAN handled  
* JWT signature manipulation rejected  
* JWT algorithm downgrade rejected  
* Privilege escalation attempt rejected  
* Service identity impersonation rejected  
* Approval with expired token rejected  
* Replay valid without IdP availability  
* Identity deletion does not invalidate artifact  
* Role removal does not invalidate artifact  
* Identity snapshot bound to artifact

---

# **SECTION H – Temporal Determinism & Workflow Runtime Tests**

* No time.Now in workflow logic  
* No random()  
* No nondeterministic map iteration  
* No external I/O inside deterministic workflow  
* Replay produces identical workflow history  
* Workflow panic on nondeterminism  
* Activity retry deterministic  
* Activity failure propagation correct  
* Signal handling deterministic  
* Query handling deterministic  
* Workflow version upgrade preserves replay  
* ContinueAsNew deterministic  
* Workflow cancellation deterministic  
* Timer handling deterministic  
* Workflow history corruption detected

---

# **SECTION I – Adversarial Tests**

* Replace artifact in chain  
* Inject fake artifact  
* Reorder artifact chain  
* Substitute logs  
* Delete DB records  
* Remove first artifact  
* Modify only last artifact  
* Truncate chain  
* Double approval attempt  
* Race between timeout and approval  
* Authority revoked during WAITING\_FOR\_HUMAN  
* Network partition simulation  
* Partial artifact write interruption  
* Duplicate artifact chain insertion  
* Forked artifact chain detection  
* Artifact replay across instances rejected

---

# **SECTION J – End-to-End Tests**

* WAITING\_FOR\_HUMAN visible  
* Approval required  
* Artifact generated  
* Offline verify works  
* System shutdown → verify works  
* DB deletion → verify works  
* Rolling restart during WAITING\_FOR\_HUMAN  
* Restart during artifact emission  
* Backup restore → replay valid

---

# **SECTION K – Canonicalization & Deterministic Serialization Tests**

* Same map different key order → identical canonical bytes  
* Nested map ordering deterministic  
* Deep nesting deterministic  
* Struct reordering deterministic  
* Canonical output stable across runtimes  
* Cross-language canonical equivalence (Go, Python, Rust, Node)  
* Numeric normalization deterministic  
* NaN/Infinity rejected  
* \-0 normalized or rejected  
* UTF-8 normalization stable  
* Escaped vs unescaped equivalent  
* Duplicate keys rejected  
* Canonical hash stable across reserialization  
* Canonicalization does not rely on runtime map iteration  
* Deterministic byte stream guaranteed

---

# **SECTION L – Schema Versioning & Migration Tests**

* Artifact schema version embedded  
* Unknown schema rejected  
* Backward-compatible additions supported  
* Field removal rejected  
* Multi-version coexistence supported  
* Upgrade/downgrade compatibility matrix validated  
* Canonicalization version migration safe  
* Hash algorithm migration safe  
* DB migration rollback safe  
* Forward compatibility validated

---

# **SECTION M – Cryptographic Integrity Tests**

* SHA-256 enforced  
* Wrong algorithm rejected  
* Truncated hash rejected  
* Mixed encoding rejected  
* Case sensitivity enforced  
* Empty hash rejected  
* Duplicate hash field rejected  
* Hash normalization strict

---

# **SECTION N – Temporal Consistency Tests**

* Artifact timestamp monotonic  
* Timestamp \>= instance creation  
* Timestamp \<= termination  
* Duplicate timestamps deterministic  
* UTC normalization enforced  
* Clock skew simulated  
* Leap second normalized

---

# **SECTION O – Concurrency & Atomicity Tests**

* Atomic commit across DB \+ storage  
* Concurrent artifact emission safe  
* Concurrent resume prevented  
* Lock contention safe  
* Parallel instance creation safe  
* Replay cannot occur mid-commit  
* Distributed runner race resolution deterministic

---

# **SECTION P – Storage Integrity & Disaster Recovery Tests**

* Write-once storage enforced  
* Overwrite attempt rejected  
* Partial write detected  
* Corrupt blob detected  
* Artifact deletion detected  
* Cross-region replication consistent  
* Backup/restore preserves replay validity  
* Storage permission changes detected

---

# **SECTION Q – Identity Drift & Historical Verification Tests**

* Replay valid without live IdP  
* Organization rename safe  
* Team deletion safe  
* Role mapping change safe  
* Identity federation provider migration safe  
* Historical verification independent of identity system

---

# **SECTION R – Verifier Robustness Tests**

* Maximum artifact size enforced  
* Maximum nesting depth enforced  
* JSON bomb attack rejected  
* Malformed UTF-8 rejected  
* Control characters rejected  
* Schema spoofing rejected  
* Continuous fuzz corpus integration  
* Malicious payload injection rejected

---

# **SECTION S – Numeric & Semantic Ambiguity Tests**

* \-0 deterministic  
* Large integers preserved exactly  
* Scientific notation deterministic  
* "1" vs 1 distinguished  
* Null vs absent field distinguished  
* Boolean vs string distinguished  
* Deep map\[string\]interface{} determinism enforced  
* Mixed-type arrays deterministic

---

# **SECTION T – Runner & Federated Execution Tests**

* Runner registration validated  
* Runner heartbeat enforced  
* Task queue distribution correct  
* Runner failure → task reassigned  
* Secret resolution at runner boundary only  
* No secret persistence in control plane  
* Evidence capture reference only  
* Runner impersonation rejected  
* Cross-team isolation enforced  
* Multiple runners same instance prevented

---

# **SECTION U – Agent Integration Tests**

Persistent agent:

* Checkpoint saved  
* Process exits on WAITING\_FOR\_HUMAN  
* Resume starts new process  
* Agent state not stored in Gantral

Split-agent:

* Pre-approval agent exits cleanly  
* Minimal context persisted  
* Post-approval resumes correctly  
* No agent memory leakage

General:

* Agent failure does not bypass authority  
* Agent upgrade safe  
* Agent restart safe

---

# **SECTION V – Performance & Load Tests**

* High concurrency instance creation  
* Artifact emission under load  
* Policy evaluation under load  
* Replay under load  
* DB connection exhaustion handled  
* Storage latency spike handled  
* Lock contention stress tested  
* SLA under concurrency verified  
* Performance regression detection in CI

---

# **SECTION W – Observability & Telemetry Tests**

* Metrics emitted for state transitions  
* HITL wait metrics accurate  
* Policy decision metrics accurate  
* Artifact emission metrics accurate  
* Trace propagation correct  
* No sensitive data in telemetry  
* Alert on SLA breach triggered correctly  
* Audit log completeness verified

---

# **SECTION X – Compatibility & Upgrade Matrix Tests**

* Go minor version compatibility  
* OS compatibility  
* Architecture compatibility  
* Rolling upgrade safe  
* Version skew tolerated  
* Upgrade during WAITING\_FOR\_HUMAN safe  
* Downgrade handling deterministic  
* Dependency version bump safe

---

# **SECTION Y – Chaos & Fault Injection Tests**

* Random component restart  
* Network partition  
* DB latency spike  
* Storage latency spike  
* Policy engine crash  
* Runner crash mid-execution  
* Clock skew injection  
* Randomized failure injection  
* Recovery preserves correctness

---

# **SECTION Z – Deterministic Boundary & Residual Edge Case Tests**

## **Z1 – Justification Boundary Conditions**

* Whitespace-only justification rejected (already covered; reaffirmed)
* Justification exceeding maximum configured length → rejected
* Extremely large justification input (DoS attempt) → rejected
* Justification truncation must never occur silently
* Justification line ending normalization (CRLF vs LF) deterministic
* Justification included in canonical context before hashing
* Tampered justification → INVALID replay

---

## **Z2 – context\_delta Determinism & Integrity**

* context\_delta included in context\_hash
* context\_delta canonicalization deterministic
* Deeply nested context\_delta deterministic
* context\_delta duplicate keys rejected
* Tampered context\_delta → INVALID replay
* Override without context\_delta → rejected
* context\_delta conflicting key overwrite attempt → rejected

---

## **Z3 – Authority Decision Object Immutability**

* Full authority decision object canonicalized before hashing
* Mutation of decision object after artifact emission → INVALID replay
* Decision object snapshot immutable once artifact committed
* Decision identity \+ role \+ justification bound into artifact hash

---

## **Z4 – Multi-Artifact Chain Edge Cases**

* Skipped prev\_artifact\_hash → INVALID
* Circular artifact reference → INVALID
* Duplicate artifact\_id in chain → INVALID
* Artifact referencing future artifact → INVALID
* Non-sequential chain ordering → INVALID
* Artifact chain gap detection
* Hash referencing incorrect previous artifact detected

---

## **Z5 – Replay Idempotency & Immutability**

* Replay run multiple times → identical result
* Replay produces no side effects
* Replay does not emit artifacts
* Replay does not mutate internal state
* Replay failure deterministic across runs
* Replay error classification consistent

---

## **Z6 – Canonical Depth & Overflow Protection**

* Maximum canonicalization depth enforced
* Nested array depth limit enforced
* Extremely large exponent numeric rejected
* Numeric overflow rejected
* Canonicalization stack overflow prevented
* Canonicalization memory bounded

---

## **Z7 – Deterministic Sorting & Locale Independence**

* Key sorting independent of locale
* Sorting consistent across environments
* Duplicate key conflict rejected before sorting
* Stable lexicographic ordering guaranteed

---

## **Z8 – Deterministic Error Semantics**

* All ambiguous errors terminate execution
* No warning-level bypass permitted
* No best-effort downgrade paths
* Error classification stable across runs
* Error codes deterministic for identical inputs

