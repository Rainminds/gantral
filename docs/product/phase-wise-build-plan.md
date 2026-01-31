---
title: Roadmap & Build Status
sidebar_label: Roadmap
---

# Gantral Roadmap & Build Status

**Current Status:** Phase 4 Complete (Developer Experience Verified)
**Next Milestone:** Phase 5 (Federated Execution)

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

## 🚧 Phase 5: Federated Execution (In Progress)
*Goal: Enable secure, multi-team execution with zero trust.*

- [x] **5.1 Identity Federation**: Implemented OIDC/JWT Middleware with Dev Mode (HS256) and Production (RS256) support. Verified via Docker.
- [x] **5.2 Service Identity**: Implemented Multi-Verifier (Chain) and RBAC Middleware. Runners restricted to Polling; Users restricted to Decisions.
- [x] **5.3 Runner Protocol**: Pull-based task queues for network isolation.
- [x] **5.4 Secret Resolution**: Just-In-Time (JIT) secret fetching at the edge.
- [ ] **5.5 Evidence Capture & Tool Mediation (Non-Authoritative)**  
  Optional runner-side capability to capture execution evidence:
  - Tool inputs/outputs captured at the runner boundary  
  - Evidence stored externally and immutably  
  - Gantral stores **references only**, never raw payloads  
  - Evidence may be required by policy, but **never interpreted by Gantral**  

  **Constraints:**  
  - Must not introduce new execution states  
  - Must not block or authorize tool execution  
  - Must not inspect or reason over payload contents  
  - Must not begin until Phase 4 acceptance criteria are satisfied

---

## 🔮 Future: Gantrio (Commercial Layer)
*Note: These features are explicit non-goals for Gantral OSS.*

- [ ] Enterprise SSO (SAML)
- [ ] Role-Based Access Control (RBAC) UI
- [ ] Approval inboxes and escalation UX
- [ ] Cost attribution dashboards
- [ ] Compliance reporting and exports
- [ ] Managed hosting and support

---

## Final Reminder

Gantral is an **execution authority layer**, not:
- A workflow engine
- An agent framework
- An autonomy platform

Determinism is guaranteed by the runtime.  
Authority is enforced by Gantral.  
**Human accountability is final.**
