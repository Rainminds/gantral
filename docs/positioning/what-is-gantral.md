# What is Gantral?

**The Kubernetes for AI Execution.**

Gantral is an open-source **AI Execution Control Plane** that enables enterprises to govern AI agents at scale. It sits above agent frameworks (LangChain, CrewAI, Vellum) and below enterprise processes (Jira, ServiceNow), providing a unified layer for authorization, audit, and Human-in-the-Loop (HITL) control.

## The Problem: "The Fox Guarding the Henhouse"

As teams adopt AI agents for coding, incidents, and workflows, a critical gap emerges: **The agent that performs an action is often the same agent that logs it.**

*   Who authorized the AI to merge that PR?
*   Why was the specific budget override approved?
*   Where is the immutable record of that decision?

Without a separate control plane, compliance depends on developer discipline and disparate logs across 50+ tools.

## The Solution: A Unified Control Layer

Gantral separates **Execution** from **Governance**.

1.  **Prevent Conflict of Interest:** Agents must request a "token" or permission state from Gantral to proceed.
2.  **Unified Ledger:** One standard audit trail for all AI actions, regardless of the underlying framework.
3.  **Decouple Policy from Code:** Update a spending limit policy *once* in Gantral, and it applies instantly to all agents.

## Strategic Positioning

> **"We don’t help you build agents.**
> **We help you run AI safely across your organization."**

Gantral enables you to:
*   Treat AI execution as **tier-0 infrastructure**.
*   Standardize **HITL (Human-in-the-Loop)** as a state transition, not a hacky Slack message.
*   Scale from 1 to 1,000 agents without losing control.
