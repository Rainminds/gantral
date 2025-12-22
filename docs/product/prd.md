# Product Reference Document

**Status:** Living Product Reference (v2)

**Purpose:** This document defines the core product vision for **Gantral**, the open-source AI Execution Control Plane.

## 1. The Core Problem

Large organizations are quietling adopting AI, but lacking **execution, control, and accountability**.

*   Teams invent their own approval/logging mechanisms.
*   No standard for Human-in-the-Loop (HITL).
*   No clear record of "Who authorized this AI action?".

**Enterprises do not lack AI tools. They lack an AI execution control plane.**

## 2. What Gantral Is

Gantral is an **AI Execution Control Plane**.

*   **Infrastructure** for running AI workflows with control.
*   **System of Record** for AI decisions.
*   **Kubernetes for AI Execution** (standardizes how it runs).

**Gantral is NOT:**
*   An agent builder (like LangChain/Vellum).
*   An observability tool (like LangFuse).
*   A "magic" autonomous dev platform.

## 3. Core Capabilities

### 3.1 Execution Plane
Standardizes how AI workflows run across teams.
*   **Processes** defined once (templates).
*   **Instances** provide isolated, auditable execution.
*   **Configuration** adapts processes per team.

### 3.2 Human-in-the-Loop (HITL)
HITL is a first-class state transition, not a UI feature.
*   `RUNNING` &rarr; `WAITING_FOR_HUMAN` &rarr; `APPROVED` / `REJECTED`.
*   Captures **Who**, **Why**, **When**, and **Context**.

### 3.3 Materiality Assessment
Governance framework to determine risk:
*   **High Materiality:** Mandatory HITL (e.g., financial decisions, prod code commits).
*   **Low Materiality:** Auto-approval OK (e.g., docs generation).

## 4. Initial Wedge: HITL for SDLC

Focus on **Code Review & Incident Management**.
1.  **AI Code Review:** High-risk changes require human sign-off.
2.  **Incident Response:** AI suggests remediation; human approves execution.
3.  **PR Risk Assessment:** AI scores risk; high score triggers blocking approval.

## 5. Personas

| Persona | Role | Needs |
|---------|------|-------|
| **Platform Engineer** | **Buyer/Owner** | Control, standardization, auditability. |
| **Engineering Manager** | **User** | Safety, visibility into team's AI risks. |
| **On-Call Engineer** | **User** | Faster incident resolution with safety rails. |
| **Compliance/GRC** | **Influencer** | "Who authorized this?" audit trails. |

## 6. Design Principles

1.  **Boring UX by Design:** Clarity > Delight.
2.  **Governance First:** Auditability is non-negotiable.
3.  **Process-Level:** Orchestrate workflows, not just agents.
4.  **No Magic:** Humans remain accountable.
