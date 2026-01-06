---
sidebar_position: 2
title: Long-Term Roadmap
---

# Project Roadmap & Vision

> **Note:** For the immediate engineering execution plan (Phases 1–5), please refer to the **[Phase-Wise Build Plan](../product/phase-wise-build-plan.md)**.

This document outlines the long-term strategic direction of Gantral. It defines our focus for future releases (post-v1.0), our commitment to stability, and the boundary between the Open Source project and commercial extensions.

---

## 🧭 Strategic Pillars

Our roadmap is guided by three core pillars:

1.  **Authority, Not Orchestration:** We will not build a workflow engine or an agent framework. We integrate with them.
2.  **Federation First:** Security boundaries must scale across teams. We prioritize features that enable decentralized execution.
3.  **Auditable by Default:** If an action cannot be proven, it should not be allowed.

---

## 🔮 Future Milestones (Post-v1.0)

These are high-priority initiatives planned for the **Post-Phase 5** timeline.

### 1. Expanded Runtime Support
*   **Goal:** Allow Gantral to run on substrates other than Temporal.
*   **Candidates:** Step Functions, Dauphin, or a lightweight in-memory option for local dev.
*   **Status:** *Research Phase*

### 2. Advanced Policy Context
*   **Goal:** Provide richer data to OPA policies during evaluation.
*   **Scope:** Inject historical context (e.g., "Has this agent failed significantly in the last 24h?") into the policy input payload.
*   **Status:** *Proposed*

### 3. "Break-Glass" Emergency Protocols
*   **Goal:** Standardized protocols for immediate suspension of all Agent Identity tokens across an organization.
*   **Use Case:** Security incident response where a compromised agent swarm must be halted globally.
*   **Status:** *RFC Pending*

---

## 🛑 Non-Goals (Out of Scope)

To maintain focus, the Gantral OSS project will **explicitly not build** the following:

*   **❌ Built-in Vector DB:** We are not a memory store.
*   **❌ Prompt Management UI:** Use existing tools (LangSmith, Portkey, etc.).
*   **❌ User Identity Provider:** We will always federate (OIDC). We will never store passwords.
*   **❌ Proprietary Model Gateways:** We do not broker LLM API keys.

---

## 🏢 Commercial Roadmap (Gantrio)

Some features are deemed "Enterprise Platform" capabilities and are developed under the commercial **Gantrio** roadmap. These will typically not be part of the OSS core.

*   **SAML / SCIM Integration:** Automated user provisioning.
*   **Multi-Region Governance Dashboard:** A unified "pane of glass" for global compliance.
*   **Long-Term Audit Archival:** Compliance storage beyond the immediate execution history.
*   **RBAC Policy UI:** Visual builders for non-technical risk officers.

---

## 📅 Release Cadence

We follow [Semantic Versioning](https://semver.org/).

*   **Major (v1.x -> v2.x):** Breaking API or architectural changes.
*   **Minor (v1.1 -> v1.2):** New features (e.g., new Adapter support) with backward compatibility.
*   **Patch (v1.1.1):** Bug fixes and security patches.

*For the latest release notes, see the [Releases](https://github.com/Rainminds/gantral/releases) page on GitHub.*
