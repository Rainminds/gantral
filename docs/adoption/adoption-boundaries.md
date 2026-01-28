---
sidebar_position: 5
title: Adoption Boundaries and Non-Goals
---

# Adoption Boundaries and Non-Goals

This document defines **when Gantral should be used**, **when it should not**, and **what it will explicitly never become**.

These boundaries are intentional.  
They protect execution correctness, organizational trust, and long-term viability.

Gantral is an AI Execution Control Plane.  
It exists to govern **execution authority**, not to optimize or accelerate execution.

---

## When Gantral Is a Good Fit

Gantral is appropriate when:

- AI-assisted actions have **irreversible or high-impact effects**
- human accountability already exists but is **implicit or informal**
- approvals occur through conventions (reviews, messages, checklists)
- audit or regulatory review requires **defensible execution records**
- policy changes occur more frequently than code changes
- multiple teams implement similar workflows with inconsistent governance

In these cases, Gantral formalizes authority without changing how work is performed.

---

## When Gantral Is Not a Good Fit

Gantral should **not** be used for:

- exploratory or sandbox experimentation
- low-risk or easily reversible automation
- purely advisory AI outputs
- performance-critical paths where pauses are unacceptable
- systems that cannot tolerate explicit human authority
- environments seeking full autonomy or “hands-off” execution

Gantral introduces intentional friction at points where **accountability matters**.

---

## Common Adoption Anti-Patterns

The following patterns undermine governance and should be avoided:

### Treating Gantral as a Workflow Engine

Gantral does not:
- orchestrate steps
- manage task execution
- optimize execution order

Using Gantral to sequence work conflates authority with orchestration.

---

### Embedding Approval Logic Inside Agents

Approval logic inside agents:
- couples reasoning and authority
- creates self-approval risks
- breaks replayability

Gantral exists to separate these concerns structurally.

---

### Enabling Enforcement Without Observation

Skipping shadow deployment:
- obscures real approval patterns
- increases resistance
- hides failure modes

Observation precedes enforcement by design.

---

### Bypassing Human-in-the-Loop for Speed

Optimizing away human authority:
- defeats the purpose of execution governance
- erodes audit credibility
- reintroduces informal approval paths

Speed is not a primary objective.

---

## What Gantral Will Explicitly Never Do

To preserve trust and determinism, Gantral will never:

- build, host, or manage AI agents
- orchestrate workflows or pipelines
- optimize models, tools, or execution paths
- make autonomous or probabilistic decisions
- approve or deny execution without a human decision
- store prompts, agent memory, or reasoning traces
- act as an identity provider
- function as a user experience or dashboard platform

These are structural non-goals.

---

## Boundaries Between Gantral and Surrounding Systems

Gantral’s responsibility ends at **execution authority**.

Surrounding systems remain responsible for:

- reasoning and planning (agents)
- task execution (runners, workflows)
- user experience and approvals (UIs)
- policy authoring and lifecycle management
- reporting, analytics, and optimization

Blurring these boundaries weakens governance rather than strengthening it.

---

## On Friction and Deliberateness

Gantral introduces **intentional friction**.

This friction is not a usability flaw.  
It is a governance guarantee.

If execution must pause for a human decision, the pause is a signal — not a failure.

---

## Adoption Principle

If your primary goal is:

- faster execution
- fewer approvals
- reduced human involvement

Gantral is likely the wrong tool.

If your goal is:

- explicit authority
- enforced accountability
- replayable decisions
- durable trust at scale

Gantral is designed for that purpose.

---

## Final Boundary

Gantral governs **whether execution may proceed**.

It does not decide **how work is done**.

Any adoption that compromises this distinction is incorrect.
