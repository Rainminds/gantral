---
sidebar_position: 2
id: roadmap
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

- [ ] **5.1 Identity Federation**: Replace local users with OIDC (JWT) claims.
- [ ] **5.2 Service Identity**: Support AWS IAM / K8s SA for machine auth.
- [ ] **5.3 Runner Protocol**: Pull-based task queues for network isolation.
- [ ] **5.4 Secret Resolution**: Just-In-Time (JIT) secret fetching at the edge.

---

## 🔮 Future: Gantrio (Commercial Layer)
*Note: These features are non-goals for Gantral OSS.*

- [ ] Enterprise SSO (SAML)
- [ ] Role-Based Access Control (RBAC) UI
- [ ] Multi-Region Replication
- [ ] Compliance Reporting Dashboard
