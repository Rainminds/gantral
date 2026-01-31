---
slug: /
---

# Welcome to Gantral

**Gantral** is an open-source **AI Execution Control Plane**.

Gantral provides a neutral infrastructure layer that allows organizations to
**enforce, record, and verify execution-time authority**
in AI-assisted and agentic workflows — independent of agent frameworks,
models, or orchestration tools.

Gantral does not build agents or workflows.  
It defines and enforces **who is allowed to execute what, and when**.

---

## What Problem Gantral Addresses

As AI systems are embedded into real operational workflows,
organizations encounter execution-time governance gaps that existing tools do not address.

Human approval and accountability often exist only as convention:
- pull-request reviews
- chat messages
- informal checklists and runbooks

These mechanisms are not technically enforceable
and do not produce durable, replayable execution evidence.

At scale, this results in:
- fragmented execution authority across teams
- ambiguous accountability
- post-incident reconstruction based on logs and testimony

Gantral exists to make execution authority:

**explicit · enforced · replayable**

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

Gantral sits **above agents and automation**
and **below organizational accountability**.

---

## Verifiability (Start Here for Audit & Proof)

Gantral is designed so that execution decisions can be
**independently verified under hostile, post-incident conditions**.

The Verifiability section defines:
- the threat model
- the commitment artifact
- the replay protocol
- explicit failure semantics
- what Gantral does *not* claim

These documents are written for auditors, regulators,
and reviewers who do **not** trust systems, operators, or narratives.

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

---

## Position Paper

Gantral is grounded in a vendor-neutral position paper that defines the
**AI Execution Control Plane** as a missing infrastructure layer for
execution-time governance.

These documents explain *why* authority, determinism,
and replayable audit must be enforced at runtime.

- **[Executive Summary](./positioning/ai-execution-control-plane-summary.md)**  
  High-level framing for platform leaders and decision-makers.

- **[AI Execution Control Plane — Position Paper](./positioning/ai-execution-control-plane.md)**  
  Full, non-normative technical and conceptual foundation.

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
Adoption is intentionally **deliberate, incremental, and reversible**.

The adoption section describes how organizations:
- evaluate Gantral without disrupting workflows
- introduce execution authority safely
- scope deployments conservatively
- avoid adoption anti-patterns

- **[Adoption Overview](./adoption/README.md)**  
  How organizations typically approach evaluation.

- **[Design Partner Engagement](./adoption/design-partners.md)**  
  Structured, time-bound evaluation engagements.

- **[Enterprise Onboarding Playbook](./adoption/enterprise-onboarding.md)**  
  Practical guidance for production environments.

- **[Adoption Boundaries and Non-Goals](./adoption/adoption-boundaries.md)**  
  When Gantral is — and is not — the right solution.

---

## Integration Guides

- **[Consumer Integration Guide](./guides/example-consumer-integration.md)**  
  Normative guide for systems integrating with Gantral.

- **[Policy Integration (OPA)](./guides/opa-integration.md)**  
  Using external policy engines as advisory inputs.

- **[Auditor Verification Guide](./guides/auditor-verification.md)**  
  How to verify an execution decision years later.

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

---

## Governance

- **[Policy vs Authority](./governance/policy-vs-authority.md)**  
  Why policy evaluation is advisory and authority is enforced as state.

- **[Roadmap](./governance/roadmap.md)**  
  Planned evolution and areas of exploration.

- **[Open Source Philosophy](./governance/oss-philosophy.md)**  
  Principles guiding neutrality and scope.

---

## Contributing

Gantral is developed in the open.

- **[How to Contribute](./contributors/how-to-contribute.md)**  
  Contribution guidelines for code and documentation.

---

## Executive Context (Optional Reading)

For readers evaluating Gantral from an enterprise,
platform, risk, or regulatory perspective,
a small set of executive briefings provide higher-level context.

These materials are **not required** to use or contribute to the project.

- **[Executive Briefings](./executive/README.md)**

---

Gantral is an independent open-source project.  
It is not affiliated with the Cloud Native Computing Foundation (CNCF).

Design and governance choices are informed by CNCF principles,
including neutrality, composability, and explicit responsibility boundaries.
