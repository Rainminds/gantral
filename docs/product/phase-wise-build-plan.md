---
title: Roadmap & Build Status
sidebar_label: Roadmap
---


# Gantral Roadmap & Build Status


**Current Status:** Phase 7 Complete 

**Next Milestone:** Phase 8 (Determinism Assurance)


This document outlines the authoritative build plan for Gantral. We follow a strict "Authority-First" architecture.


---


## ✅ Phase 1: Control Foundations (Core)
*Goal: Establish deterministic execution control with strict human authority.*


- [x] **1.1 Canonical State Machine**: Implemented strict transitions (RUNNING -> WAITING_FOR_HUMAN -> APPROVED).
- [x] **1.2 Immutable History**: All transitions emit append-only events.
- [x] **1.3 HITL Semantics**: Human-in-the-Loop state is first-class and blocking.


## ✅ Phase 2: Governance Hardening
*Goal: Introduce policy evaluation and audit safety.*


- [x] **2.1 Policy Interface**: Pluggable Rego-based policy evaluation.
- [x] **2.2 Transition Guards**: Policies enforce `REQUIRE_HUMAN` or `DENY` logic.
- [x] **2.3 Audit Replay**: Execution history supports deterministic replay.


## ✅ Phase 3: Enterprise Integration
*Goal: Scalable, durable execution via Temporal.*


- [x] **3.1 Workflow Runtime**: Temporal integration for durability.
- [x] **3.2 Adapters**: Event-driven adapter architecture.
- [x] **3.3 SDKs**: Python SDK for agent interaction.


## ✅ Phase 4: Developer Experience & Demos
*Goal: Prove usability with run-ready examples.*


- [x] **4.1 Demo Environment**: Docker Compose stack with no K8s dependencies.
- [x] **4.2 Persistent Agent**: Reference implementation for agents with native checkpointing (`sys.exit(3)`).
- [x] **4.3 Split-Agent**: Reference pattern for stateless agents (Pre/Post split).
- [x] **4.4 Scripted Interaction**: CLI scripts for `trigger`, `status`, `approve`.
- [x] **4.5 Verification**: Validated "Stranger Test" (clone -> run -> works).


---


## ✅ Phase 5: Federated Execution 
*Goal: Enable secure, multi-team execution with zero trust.*


- [x] **5.1 Identity Federation**: Implemented OIDC/JWT Middleware with Dev Mode (HS256) and Production (RS256) support. Verified via Docker.
- [x] **5.2 Service Identity**: Implemented Multi-Verifier (Chain) and RBAC Middleware. Runners restricted to Polling; Users restricted to Decisions.
- [x] **5.3 Runner Protocol**: Pull-based task queues for network isolation.
- [x] **5.4 Secret Resolution**: Just-In-Time (JIT) secret fetching at the edge.

  **Constraints:**  
  - Must not introduce new execution states  
  - Must not block or authorize tool execution  
  - Must not inspect or reason over payload contents  
  - Must not begin until Phase 4 acceptance criteria are satisfied


---


## ✅ Phase 6: Verifiability & Admissibility
*Goal: Transform documented guarantees into mechanically verifiable artifacts that survive hostile audit.*
*Note: This phase MUST NOT alter execution semantics, state machines, or authority rules.*


- [x] **6.1 Commitment Artifact Implementation**: Implemented a concrete, inspectable commitment artifact emitted atomically with authority transitions.
- [x] **6.2 Artifact Storage & Log Independence**: Ensured artifacts are independent of operational logs and databases; deleting the DB must not invalidate artifacts.
- [x] **6.3 Offline Verification Tooling**: Enabled third-party verification without Gantral access via a standalone CLI/library (`gantral-verify`).
- [x] **6.4 Authority-Only Replay Enforcement**: Guarantee replay depends solely on authority artifacts, excluding agent memory and logs.
- [x] **6.5 Fail-Closed Guarantees**: Eliminate ambiguous execution paths; execution must terminate on missing/partial artifacts or hash mismatches.

**Stop Conditions (Non-Negotiable):**
- Artifact emission is non-atomic
- Replay requires Gantral services
- Logs or dashboards are treated as evidence
- Execution continues under ambiguity

---

## ✅ Phase 7: Production Hardening
*Goal: Production-grade artifact durability.*

- [x] **7.1 Immutable Artifact Storage**: Enforce write-once storage immutability for all authority commitment artifacts.
- [x] **7.2 Deterministic Upgrade Protocol**: Ensure that structural changes to the artifact schema or hashing logic do not invalidate the historical artifact chain during replay, guaranteeing backward compatibility.
- [x] **7.3 Multi-Tenant Isolation**: Bind `teamid` immutably to all execution instances and commitment artifacts, and enforce hard rejection of any cross-tenant artifact access or mixing.
- [x] **7.4 Replay Determinism Validation**: Prove that the Replay Verifier strictly matches the TRD output states (`VALID`, `INVALID`, `INCONCLUSIVE`)
- [x] **7.5 Fail-Closed Chaos Validation**: Prove that Gantral's core execution and artifact emission pipelines strictly fail-closed under infrastructure chaos or data corruption
- [x] **7.6 Example Architecture Hardening**: Prove that external users (and our own examples) interact with Gantral *only* through the official SDK, without ever leaking internal kernel dependencies (`core`, `internal`, `adapters`, `pkg`) into external integration code.

---

## Final Reminder

Gantral is an **execution authority layer**, not:
- A workflow engine
- An agent framework
- An autonomy platform


Determinism is guaranteed by the runtime.  
Authority is enforced by Gantral.  
**Human accountability is final.**
