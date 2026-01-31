---
title: Open Source Philosophy
---

# Open Source Philosophy

Gantral is open source by design, not by distribution.

This project exists at the infrastructure layer of enterprise AI adoption. Infrastructure that controls execution, approvals, and audit must be inspectable to earn trust.

**Open source is a structural requirement, not a community growth tactic.**

---

## Why Gantral Is Open Source

Gantral sits at a sensitive control boundary:
- It enforces human authority over AI decisions.
- It records audit and execution evidence.
- It governs cost, policy, and accountability.

Enterprises—especially regulated ones—do not trust closed control planes for these responsibilities. Open source enables:
- **Verifiability** of execution semantics.
- **Independent security review** by risk teams.
- **Regulator and auditor confidence**.
- **Long-term ecosystem trust**.

Without openness, Gantral would not be viable in its intended environments.

---

## What “Open” Means in Practice

For Gantral, open source is the **operating model**, not a marketing artifact. It means:
- Core execution semantics are public.
- HITL logic is inspectable.
- Policy enforcement is transparent.
- Audit guarantees are verifiable.
- Specifications lead implementations.

---

## Decision-Making & Prioritization Principles

Gantral’s roadmap and features are guided by strict principles, not feature lists. We prioritize:

1.  **Execution Control Correctness:** Correctness > Speed.
2.  **Auditability & Replayability:** Evidence > Convenience.
3.  **Clear Authority Boundaries:** Humans must always be able to intervene.
4.  **Enterprise Trust:** Predictability > Novelty.

### How We Evaluate Features
A proposed feature is evaluated against these questions:
- Does this strengthen human authority?
- Does this preserve audit guarantees?
- Does this avoid coupling to specific vendors?
- Does this reduce, not increase, ambiguity?

**We explicitly avoid:**
- Feature accumulation without control rationale.
- "Agent Intelligence" improvements (we are a control plane, not a brain).
- Vendor-specific optimizations.

---

## Open Core Model

Gantral follows an **open core** model.
- **Open:** Trust-critical infrastructure, execution logic, and audit trails.
- **Commercial (Gantrio):** Enterprise SSO, RBAC UIs, multi-region dashboards, and compliance reporting.

This separation is intentional. It ensures the "Authority" layer remains transparent while the "Management" layer scales with enterprise needs.

---

## Long-Term Orientation

Gantral is designed to outlive individual vendors. The project favors:
- **Stability** over novelty.
- **Predictability** over rapid iteration.
- **Governance clarity** over speed.

These values are incompatible with closed, opaque systems. They are compatible with open infrastructure.
