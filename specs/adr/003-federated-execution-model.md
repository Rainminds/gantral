# 3. Federated Execution Model (Runner Pattern)

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
Enterprises typically operate in **siloed environments** due to regulatory (GDPR/PCI), security (Air Gap), or data residency requirements.

A central SaaS control plane cannot simply "run the code" because it cannot (and should not) have direct network access to internal databases, private code repositories, or production infrastructure.

## Decision
We will adopt a **Federated Execution Model** (The "Runner" Pattern).

1.  **Central Control:** Gantral acts as the central **Control Plane**, managing authority, policy, and audit state.
2.  **Distributed Action:** Execution logic runs on distributed **Runners** deployed within the user's own infrastructure (VPC, On-Prem).
3.  **Pull-Based:** Runners pull tasks from central queues; Gantral never pushes commands or connects inbound to the user's environment.

## Consequences
*   **Positive:**
    *   **Data Residency:** Sensitive code and data never leave the customer's controlled environment.
    *   **Network Security:** No inbound ports required on customer firewalls (outbound-only polling).
    *   **Scalability:** Execution compute scales independently of the Control Plane.
*   **Negative:**
    *   **Deployment Friction:** Users must install and manage Runner binaries.
    *   **Version Skew:** We must support backward compatibility for older Runners connecting to a newer Control Plane.
