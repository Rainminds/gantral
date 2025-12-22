# OSS Philosophy & Governance

## Open Source vs. Commercial

Gantral follows an **Open Core** model with a strict separation of concerns to maintain trust.

### Why Open Source?

For an **Execution Control Plane**, trust is paramount.
*   It touches approvals, audits, and budgets.
*   Closed control planes are "black boxes" that regulated industries cannot fully trust.
*   **The core must be open.**

### Gantral Core (Apache 2.0)
**Audience:** Platform teams, architects, contributors.
**Scope:**
*   Execution State Machine.
*   Instance Lifecycle APIs.
*   HITL Semantics.
*   Policy Engine.
*   Reference SDKs.

**License:** Apache 2.0 (Forever).
**Governance:** Maintainer-Led with public RFC process.

### Gantrio (Commercial Enterprise Platform)
**Audience:** Large enterprises requiring managed UI and support.
**Scope:**
*   Enterprise UI / Dashboards.
*   Org-level RBAC & Hierarchy.
*   SSO / SAML.
*   Compliance Reporting & Export.
*   Managed Hosting / SaaS.

## Governance Model

*   **Maintainer-Led:** Roadmap authority resides with core maintainers to ensure velocity in the early phase.
*   **RFC Process:** Major architectural changes follow a public "Request for Comments" (RFC) process.
*   **Contribution:** CLA required. We welcome community contributions to the Core.
