---
slug: /
sidebar_position: 1
---

# Welcome to Gantral

**Gantral** is an open-source **AI Execution Control Plane**.

Gantral provides a neutral infrastructure layer that allows organizations to **govern, enforce, and audit execution-time decisions** in AI-assisted and agentic workflows — independent of agent frameworks, models, or orchestration tools.

Gantral does not build agents or workflows.  
It defines and enforces **execution authority**.

---

## What Problem Gantral Addresses

As AI systems are embedded into real operational workflows, organizations encounter execution-time governance gaps that existing tools do not address.

Human approval and accountability often exist only as convention:
- pull-request reviews
- chat messages
- informal checklists and runbooks

These mechanisms are not technically enforceable and do not produce a durable execution record.

At scale, these challenges typically appear as **operational fragmentation across teams** and a **loss of end-to-end decision traceability** between AI recommendations, human authority, and execution outcomes.

Gantral exists to make execution authority **explicit, enforceable, and replayable**.

---

## What Gantral Is — and Is Not

Gantral is:
- an execution authority layer
- a deterministic control plane for AI-assisted actions
- a system of record for execution-time decisions

Gantral is **not**:
- an agent framework
- a workflow engine
- a policy authoring system
- an AI governance UI
- an autonomous decision system

Gantral sits **above agents and automation** and **below organizational accountability**.

---

## Position Paper

Gantral is grounded in a vendor-neutral position paper that defines the **AI Execution Control Plane** as a missing infrastructure layer for execution-time governance in AI-assisted systems.

These documents explain *why* execution authority, human accountability, and replayable audit must be enforced at runtime — independent of agent frameworks, models, or workflows.

- **[Executive Summary](./positioning/ai-execution-control-plane-summary.md)**  
  A concise overview of the problem, core principles, and reference architecture.  
  Intended for platform leaders, architects, and decision-makers.

- **[AI Execution Control Plane — Position Paper](./positioning/ai-execution-control-plane.md)**  
  The full, non-normative position paper defining execution authority, determinism, and auditability.

---

## Getting Oriented

- **[What is Gantral?](./positioning/what-is-gantral.md)**  
  High-level overview and mental model.

- **[What Gantral Is Not](./positioning/what-gantral-is-not.md)**  
  Explicit non-goals and boundaries relative to agents and workflows.

- **[Product Specification (PRD)](./product/prd.md)**  
  Canonical product boundaries, responsibilities, and invariants.

---

## Adoption and Evaluation

Gantral is infrastructure.  
Adoption is intentionally **deliberate, incremental, and reversible**.

The adoption section describes how organizations:
- evaluate Gantral without disrupting existing workflows
- introduce execution authority safely
- establish trust before enforcement
- avoid common adoption anti-patterns

- **[Adoption Overview](./adoption/README.md)**  
  How organizations typically approach adoption.

- **[Design Partner Engagement](./adoption/design-partners.md)**  
  Structured, time-bound evaluation engagements.

- **[Enterprise Onboarding Playbook](./adoption/enterprise-onboarding.md)**  
  Practical guidance for introducing Gantral in production environments.

- **[Security & Trust (Adoption View)](./adoption/security-and-trust.md)**  
  Adoption-time security and trust considerations.

- **[Adoption Boundaries and Non-Goals](./adoption/adoption-boundaries.md)**  
  When Gantral is — and is not — the right solution.

---

## Integration Guides

- **[Consumer Guide](./guides/example-consumer-integration.md)**  
  Normative guide for agents and systems integrating with Gantral.

- **[Policy Integration](./guides/opa-integration.md)**  
  Integrating external policy engines as advisory inputs.

- **[Demo Walkthrough](./guides/demo.md)**  
  Running reference examples and end-to-end flows.

---

## Technical Architecture

- **[Technical Reference (TRD)](./architecture/trd.md)**  
  Authoritative technical reference, invariants, and execution semantics.

- **[Execution State Machine](./architecture/state-machine.md)**  
  Canonical execution lifecycle and guarantees.

- **[Implementation Guide](./architecture/implementation-guide.md)**  
  Reference implementation guidance aligned with the TRD.

---

## Roadmap and Project Status

- **[Phase-wise Build Plan](./product/phase-wise-build-plan.md)**  
  Current build status and execution milestones.

- **[Roadmap](./governance/roadmap.md)**  
  Planned evolution and areas of exploration.

---

## Governance and Contributing

Gantral is an open-source project developed in the open.

- **[Open Source Philosophy](./governance/oss-philosophy.md)**  
  Principles guiding scope, neutrality, and openness.

- **[How to Contribute](./contributors/how-to-contribute.md)**  
  Contribution guidelines for code and documentation.

---

## Executive Context (Optional Reading)

For readers evaluating Gantral from an enterprise, platform, risk, or regulatory perspective,
a small set of executive briefings provide higher-level context on:

- how execution authority changes after adoption
- how governance scales across teams
- how audit and accountability are enforced at runtime

These materials are **not required** to use or contribute to the open-source project.

- **[Executive Briefings](./executive/README.md)**

---

Gantral is an independent open-source project.  
It is not affiliated with the Cloud Native Computing Foundation (CNCF).

Design and governance choices are informed by CNCF principles, including
neutrality, composability, and explicit responsibility boundaries.
