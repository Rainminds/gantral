---
title: Enterprise Onboarding Playbook
---

# Enterprise Onboarding Playbook

This document describes how organizations typically introduce Gantral into production environments **safely, incrementally, and reversibly**.

It is not a deployment guide and does not prescribe tooling choices.  
It focuses on **organizational adoption of execution authority**, not feature enablement.

Gantral is an AI Execution Control Plane.  
Onboarding is intentionally deliberate.

---

## What You Are Introducing

Gantral introduces a **formal execution authority layer** into environments where AI-assisted workflows already exist.

Specifically, Gantral:

- determines whether execution may proceed
- enforces pauses for explicit human authority
- resumes or terminates execution based on recorded decisions
- produces a deterministic, replayable execution record

Gantral does **not**:

- build or host agents
- orchestrate workflows
- optimize models or tools
- store prompts, agent memory, or secrets
- make probabilistic or autonomous decisions

Execution authority is introduced without changing how work is performed.

---

## Where Gantral Fits

Gantral sits **above agents and automation** and **below organizational accountability**.

```

Agents / Tools
↓
Gantral (Execution Authority)
↓
Human Decision
↓
Audit & Replay

```

This placement is intentional and non-negotiable.

---

## Recommended Onboarding Sequence

Enterprise onboarding typically proceeds in four stages.  
Each stage builds trust before introducing enforcement.

---

### Stage 1 — Workflow Selection

Select **one or two workflows** with the following characteristics:

Good candidates:
- irreversible or production-impacting actions
- existing human approval rituals
- regulatory, financial, or security relevance

Poor candidates:
- exploratory or sandbox workflows
- low-risk automation
- performance-critical paths

Gantral is most effective where **human accountability already exists but is implicit**.

---

### Stage 2 — Shadow Deployment (No Enforcement)

Deploy Gantral in **observe-only mode**.

During this stage:
- execution proceeds without blocking
- Gantral records execution boundaries
- hypothetical approval points are identified
- no human interaction is required

This stage establishes baseline trust and visibility without disruption.

---

### Stage 3 — Controlled Enforcement

Enable explicit human-in-the-loop enforcement for:

- a small number of high-impact actions
- a single team or workflow

During this stage:
- execution pauses at defined authority boundaries
- approvals, rejections, and overrides are explicit
- execution resumes only after a recorded decision

Enforcement is narrow by design.

---

### Stage 4 — Operationalization

Expand usage when:

- multiple teams rely on AI-assisted execution
- approval patterns diverge across workflows
- policy changes outpace code changes
- audit and accountability questions increase

At this stage, Gantral becomes a **platform-level execution contract** rather than a workflow-specific tool.

---

## Operational Guarantees

Gantral guarantees:

- deterministic execution state transitions
- immutable, append-only execution history
- explicit human authority for governed actions
- replayable decision records independent of agent memory
- policy changes without redeploying agents

Gantral never:
- self-approves execution
- inspects agent reasoning or prompts
- stores raw credentials or secrets
- bypasses human authority
- mutates historical execution records

These guarantees are structural, not configurable.

---

## Security and Trust Posture (Summary)

From an onboarding perspective:

- Gantral is self-hosted
- Licensed under Apache 2.0
- Integrates with existing identity providers
- Stores references to secrets, not secrets themselves
- Produces tamper-evident execution records
- Supports air-gapped and regulated environments

Trust is established through **architecture and transparency**, not contractual assurances.

---

## Reversibility and Exit

Gantral onboarding is intentionally reversible.

If Gantral is removed:
- agents continue to function
- workflows continue to run
- execution authority reverts to prior conventions

No business logic is trapped inside Gantral.

This reversibility is essential to trust.

---

## Common Adoption Anti-Patterns

Avoid:

- enabling enforcement before observing real execution
- applying Gantral to low-risk exploratory workflows
- bypassing WAITING_FOR_HUMAN to preserve speed
- embedding approval logic inside agents
- treating Gantral as a workflow engine or policy system

These patterns erode governance rather than strengthening it.

---

## Final Perspective

Before Gantral:

> “We assume a human reviewed this.”

After Gantral:

> “We can demonstrate who authorized this execution,  
> under what conditions,  
> and replay that decision later.”

Enterprise onboarding succeeds when execution authority becomes **explicit, enforced, and defensible** — without changing how work gets done.
