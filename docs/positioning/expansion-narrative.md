---
sidebar_position: 6
title: Expansion Narrative
---

# Expansion Narrative

Gantral is built as infrastructure. Infrastructure does not launch fully formed; it expands by owning the most critical control points first.

This document describes our architectural direction and intent.

---

## Horizon 1: The Authority Primitive (Current)

**Focus:**
- Human-in-the-Loop as a blocking state transition.
- Deterministic execution and audit.
- Developer Experience (SDKs, Local Runners).

**Why:**
Every material AI workflow eventually requires a human decision. Today, this is handled via Slack messages, emails, and ad-hoc dashboards. Gantral standardizes this boundary first because it is the highest-risk gap.

---

## Horizon 2: Federation & Scale (Next)

**Focus:**
- Federated Identity (OIDC) and Secret Resolution.
- Distributed Runners across team boundaries.
- Cross-team policy enforcement (e.g., "All Finance Agents need Risk Approval").

**Why:**
Once the primitive works for one team, it must scale to the enterprise. This requires solving the "Zero Trust" problem between the Control Plane and the Agents.

---

## Horizon 3: The Governance Layer

**Focus:**
- Advanced OPA Policy Libraries.
- Materiality-based routing.
- Compliance automation and historical replay.

**Why:**
At scale, enterprises need centralized governance without blocking teams. Policies must be enforced structurally (code), not socially (wiki pages).

---

## Horizon 4: Ecosystem Standards

**Focus:**
- Open protocols for "Agent Handoffs."
- Neutral Foundation-style governance.
- Broad integration with all major agent frameworks.

**Why:**
Execution control infrastructure must be neutral to be trusted long-term. We aim to be the standard "Sudo" command for the AI era.

---

## Note on Terminology

**Horizons** describe our strategic focus over time.
**Phases** (in the Build Plan) describe specific engineering milestones.

Gantral expands by deepening control, not by accumulating features.
