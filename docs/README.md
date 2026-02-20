---
slug: /
---

# Welcome to Gantral

> **Execution Infrastructure for Deterministic Agentic AI**

Gantral is an open-source AI Execution Control Plane designed for organizations running agentic AI in real operational environments.

It eliminates policy–code duplication, enforces authority as canonical workflow state, and produces replayable, tamper-evident execution evidence — independent of agent frameworks, models, or orchestration tools.

Gantral does not build agents.  
Gantral does not replace workflow runtimes.  
Gantral ensures that execution decisions are deterministic, version-bound, and independently verifiable.

---

## Executive Overview

As agentic AI systems move into production, enterprises encounter three recurring risks:

- Slow policy changes that require workflow redeployment  
- Inconsistent governance across teams and environments  
- Inability to independently verify execution decisions  

Gantral addresses these risks by enforcing execution-time authority as infrastructure.

Authority becomes:

- Deterministic  
- Version-bound  
- Identity-bound  
- Cryptographically committed  
- Independently replayable  

---

## The Structural Challenges Gantral Solves

### 1. Operational Inefficiency

When approval thresholds and governance rules are embedded directly in workflow code:

- Policy updates require redeployment  
- Teams fork workflows for environment differences  
- Configuration drifts over time  
- Change velocity slows  

Gantral separates policy from workflow implementation.  
Policy updates no longer require redeploying workflows.

This increases operational agility while reducing change risk.

---

### 2. Governance Fragmentation

In many systems:

- Policy is evaluated in one service  
- Human approvals are recorded elsewhere  
- Execution continues independently  

Authority exists as documentation — not enforced state.

Gantral binds authority directly to canonical workflow state transitions.

Execution cannot proceed beyond governed boundaries without structural evidence.

---

### 3. Broken Chain of Custody

Post-incident reconstruction often depends on:

- Logs  
- Runtime access  
- Policy memory  
- Human testimony  

Gantral emits cryptographically chained commitment artifacts at each authority boundary.

Execution decisions become replayable and independently verifiable — even years later.

---

## What Gantral Is — and Is Not

Gantral is:

- An execution authority layer  
- Deterministic execution infrastructure  
- A version-bound system of record for execution-time decisions  

Gantral is not:

- An agent framework  
- A workflow engine  
- A policy authoring system  
- A governance UI  
- An autonomous decision system  

Gantral sits above agents and automation and below organizational accountability.

It defines and enforces who may execute what, and when — and records that determination in a replayable form.

---

## Verifiability (Start Here for Audit & Proof)

Gantral is designed so that execution decisions can be independently verified under adversarial, post-incident conditions.

These documents define:

- The minimum admissibility bar for execution-time authority  
- The mechanisms required to produce defensible execution evidence  

- **[Admissible Execution](./verifiability/admissible-execution.md)**  
  Normative definition of valid execution authority.

- **[Verifiability Overview](./verifiability/README.md)**  
  Entry point for independent verification.

- **[Threat & Adversary Model](./verifiability/threat-model.md)**  
  Explicit assumptions and out-of-scope risks.

- **[Commitment Artifact](./verifiability/commitment-artifact.md)**  
  The execution-time evidence object.

- **[Replay Protocol](./verifiability/replay-protocol.md)**  
  Deterministic third-party verification procedure.

- **[Failure Semantics](./verifiability/failure-semantics.md)**  
  Valid, invalid, and inconclusive outcomes.

- **[Explicit Non-Claims](./verifiability/non-claims.md)**  
  What Gantral deliberately does not guarantee.

Gantral provides verifiable authority evidence.  
It does not claim regulatory certification.

---

## Position Papers

Gantral is grounded in vendor-neutral position papers that define:

- The AI Execution Control Plane as a missing infrastructure layer  
- The admissibility standard for execution-time authority  

These documents explain why authority, determinism, and replayable audit must be enforced at runtime.

- **[Executive Summary](./positioning/ai-execution-control-plane-summary.md)**  
  High-level framing for platform leaders and decision-makers.

- **[AI Execution Control Plane — Position Paper](./positioning/ai-execution-control-plane.md)**  
  Category definition and non-normative technical foundation.

- **[Admissible Execution — Invariants for Execution Authority](./verifiability/admissible-execution.md)**  
  Normative standard defining the minimum bar for defensible execution.

- **[Gantral — Implementation of an Admissible AI Execution Control Plane](./architecture/gantral-implementation-paper.md)**  
  Formal specification and open-source reference implementation.

---

## Getting Oriented

- **[What is Gantral?](./positioning/what-is-gantral.md)**  
  High-level overview and mental model.

- **[What Gantral Is Not](./positioning/what-gantral-is-not.md)**  
  Explicit non-goals and boundaries.

- **[Product Specification (PRD)](./product/prd.md)**  
  Canonical product responsibilities and invariants.

---

## Adoption and Evaluation

Gantral is infrastructure.  
Adoption should be deliberate, incremental, and reversible.

Organizations can:

- Introduce Gantral without disrupting existing workflows  
- Scope deployments to high-materiality boundaries first  
- Migrate authority enforcement progressively into canonical state  

- **[Adoption Overview](./adoption/README.md)**  
  How organizations evaluate and introduce Gantral safely.

- **[Design Partner Engagement](./adoption/design-partners.md)**  
  Structured, time-bound evaluation engagements.

- **[Enterprise Onboarding Playbook](./adoption/enterprise-onboarding.md)**  
  Practical guidance for production environments.

- **[Adoption Boundaries and Non-Goals](./adoption/adoption-boundaries.md)**  
  When Gantral is — and is not — the right solution.

---

## Integration Guides

- **[Consumer Integration Guide](./guides/example-consumer-integration.md)**  
  Normative guide for integrating systems with Gantral.

- **[Policy Integration (OPA)](./guides/opa-integration.md)**  
  Using external policy engines as advisory inputs.

- **[Auditor Verification Guide](./guides/auditor-verification.md)**  
  How to verify an execution decision independently.

- **[Demo Walkthrough](./guides/demo.md)**  
  Reference examples and end-to-end flows.

---

## Technical Architecture

- **[Technical Reference (TRD)](./architecture/trd.md)**  
  Authoritative execution semantics and invariants.

- **[Authority State Machine](./architecture/authority-state-machine.md)**  
  Canonical authority enforcement lifecycle.

- **[Implementation Guide](./architecture/implementation-guide.md)**  
  Reference implementation aligned with the TRD.

- **[Gantral — Implementation Paper](./architecture/gantral-implementation-paper.md)**  
  Formal specification and reference implementation of deterministic authority, policy separation, and verifiable chain-of-custody.

Gantral enforces canonical state transitions and deterministic authority semantics.

---

## Testing & Constitutional Enforcement

In Gantral, testing is not traditional QA.  
It is mechanical enforcement of the Admissible Execution standard.

The testing infrastructure validates:

- State transition correctness  
- Fail-closed behavior  
- Replay determinism  
- Artifact integrity  
- Version consistency  

- **[Testing Architecture Blueprint](testing/testing-architecture-blueprint.md)**  
  Structured test layering across unit, integration, and adversarial boundaries.

- **[Master Test Inventory](testing/master-test-inventory.md)**  
  Comprehensive coverage of state machine and replay invariants.

- **[Test Execution Strategy](testing/test-execution-strategy.md)**  
  Tiered execution model balancing rigor and velocity.

---

## Governance

- **[Policy vs Authority](./governance/policy-vs-authority.md)**  
  Why policy evaluation is advisory and authority is enforced as state.

- **[Roadmap](./governance/roadmap.md)**  
  Planned evolution and areas of exploration.

- **[Open Source Philosophy](./governance/oss-philosophy.md)**  
  Principles guiding neutrality and scope.

Gantral is a maintainer-led open-source project.  
It is not affiliated with the Cloud Native Computing Foundation (CNCF).

Design principles:

- Neutrality  
- Composability  
- Deterministic semantics  
- Explicit responsibility boundaries  

---

Gantral is an independent open-source project developed by Rainminds Solutions Pvt. Ltd.
