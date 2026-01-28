---
sidebar_position: 1
title: Adoption Overview
---

# Gantral Adoption

This section describes how organizations typically **evaluate, adopt, and operationalize Gantral** as execution authority infrastructure.

These documents are **not product documentation**, **not sales material**, and **not guarantees of participation**.  
They exist to make Gantral’s adoption model **explicit, predictable, and inspectable** before any engagement begins.

Gantral is an AI Execution Control Plane.  
Its adoption is intentionally deliberate.

---

## Why These Documents Exist

As AI becomes embedded in real operational workflows, organizations consistently encounter two structural failures that are not solved by model capability, policy documents, or observability tooling.

### 1. Operational Fragmentation

As AI adoption grows:

- approval and risk rules drift into agent code, scripts, prompts, and team-specific runbooks
- similar workflows re-implement governance differently across teams
- policy changes require code changes in multiple places
- enforcement becomes inconsistent and difficult to reason about at scale

Governance exists, but it is **fragmented across tools, teams, and conventions**.

### 2. Broken Chain of Custody

Even when humans remain accountable:

- AI recommendations and human decisions are not captured as a single execution record
- the boundary where AI input ends and human authority begins is implicit
- approvals are enforced socially rather than technically
- audits reconstruct decisions after the fact rather than replaying them deterministically

Execution succeeds — but the **decision trail does not survive intact**.

Gantral exists to address both failures **at execution time**.

---

## What These Documents Are

The documents in this section describe:

- how organizations explore Gantral without disrupting existing workflows
- how execution authority is introduced safely and incrementally
- how trust is established before enforcement is enabled
- how adoption progresses from observation to operation

They reflect how Gantral is designed to be used in **real enterprise environments**, especially regulated or high-accountability contexts.

---

## What These Documents Are Not

These documents do **not**:

- describe features or roadmap
- define commercial terms or pricing
- promise timelines or outcomes
- replace architectural or technical documentation
- function as a sales or marketing funnel

Gantral is infrastructure.  
Adoption clarity matters more than conversion speed.

---

## Who Should Read This Section

This section is intended for:

- platform and infrastructure engineers
- staff and principal engineers responsible for execution correctness
- security, risk, and compliance reviewers
- technical leaders accountable for AI governance at organizational scale

If you are looking to **build, integrate, or extend Gantral**, refer to the architecture, guides, and product sections instead.

---

## How Adoption Typically Proceeds

Organizations generally adopt Gantral in stages:

1. **Exploration**  
   Understanding where AI already participates in execution and where authority is currently implicit.

2. **Shadow Deployment**  
   Observing execution flows without blocking actions or enforcing approvals.

3. **Controlled Enforcement**  
   Introducing explicit human-in-the-loop states for a small number of high-impact actions.

4. **Operationalization**  
   Standardizing execution authority across teams, workflows, and environments.

Each stage is intentional.  
Skipping steps increases risk and erodes trust.

---

## Participation and Scope

Some adoption paths, such as design partner engagements, are **limited and initiated by mutual agreement**.

Publishing this material publicly does not imply availability, prioritization, or commitment.

The goal is clarity — not scale.

---

## Related Documents

- **Design Partner Engagement**  
  Describes the structure and scope of collaborative evaluation engagements.

- **Enterprise Onboarding Playbook**  
  A concise, practical guide for introducing Gantral inside an organization.

- **Adoption Boundaries**  
  When Gantral is not the right solution and common anti-patterns to avoid.

---

## Final Note

Gantral exists to make execution authority **explicit, enforceable, and replayable**.

Adoption succeeds when organizations can confidently answer:

> “What was allowed to run,  
> who authorized it,  
> and why —  
> and can we prove it later?”

These documents exist to support that outcome.

If you are evaluating Gantral from a security or risk perspective,
start with the Security and Trust section before reviewing onboarding details.
