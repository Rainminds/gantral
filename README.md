# Gantral

> **Deterministic Authority Infrastructure for Scaling AI**  
> *Structural enforcement that removes ambiguity and accelerates enterprise adoption.*

![Status](https://img.shields.io/badge/Status-Active-green)
![License](https://img.shields.io/badge/License-Apache_2.0-blue)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8)
![Runtime](https://img.shields.io/badge/Runtime-Temporal-black)
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.18682545.svg)](https://doi.org/10.5281/zenodo.18682545)

---

## 10-Second Definition

Gantral is an open-source **Execution Authority Kernel** for enterprise AI workflows.

Orchestration coordinates tasks.  
Gantral governs whether execution is admissible.

Gantral makes authority:

- Deterministic  
- Version-bound  
- Identity-bound  
- Cryptographically committed  
- Log-independent  
- Replayable  
- Fail-closed  

Gantral is infrastructure — not a dashboard.

---

# Why Gantral Exists

AI pilots succeed.

AI scaling stalls.

Not because models fail —  
but because authority fragments.

As AI systems move from experimentation into production workflows, they begin to:

- Move money  
- Modify infrastructure  
- Change access control  
- Influence regulatory posture  
- Trigger irreversible operational outcomes  

At that boundary:

AI can act.  
Humans remain accountable.

Most enterprise stacks today can:

- Orchestrate workflows  
- Evaluate policy  
- Apply guardrails  
- Log events  
- Produce dashboards  

Very few structurally enforce execution authority.

Authority is reconstructed from logs instead of enforced at runtime.

Gantral introduces a missing infrastructure layer:

**Execution Authority as canonical workflow state.**

---

# The Structural Barriers to Scaling AI

When AI moves beyond basic human-supervised automation, enterprises encounter structural friction:

### Policy-in-Code Drift
Authority thresholds embedded inside workflow logic.  
Every policy change requires redeployment.  
Version drift accumulates.

### Cross-Runtime Inconsistency
Different orchestrators.  
Different HITL semantics.  
Different escalation logic.  
No uniform authority model.

### Broken Chain of Custody
AI recommends.  
Human approves.  
Execution resumes elsewhere.  
Logs attempt reconstruction — but cannot prove.

### Non-Defendability
Cannot deterministically answer:

- Which workflow version ran?
- Which policy version applied?
- Who exercised authority?
- What context existed at decision time?

Logs reconstruct.  
High-impact AI requires proof.

### Black-Box Enforcement
Opaque runtime logic.  
Vendor-controlled semantics.  
Log-dependent replay.

Trust requires transparency and mechanical verification.

---

# Why Existing Systems Are Not Enough

Modern stacks include:

- Workflow orchestration engines  
- Agent frameworks  
- Policy engines (e.g., OPA)  
- Runtime guardrails  
- Observability and GRC systems  

Each solves an important problem.

But:

- Policy evaluation is not execution authority.
- Orchestration is not accountability.
- Logs are not admissible proof.
- Guardrails are not structural enforcement.

In most deployments:

Human approval is a pause in workflow code.

Authority is inferred from metadata.

Gantral makes authority structural.

Approval becomes a deterministic state transition — not a log entry.

---

# Where Gantral Sits in the Stack

Gantral integrates at execution authority boundaries.

Agents reason.  
Orchestration sequences tasks.  
Policy engines advise.  
Gantral enforces authority.

Gantral operates:

- At explicit authority checkpoints  
- At human-in-the-loop boundaries  
- At execution resumption points  

It does not replace:

- Orchestration  
- Agent frameworks  
- Policy engines  
- Observability systems  

It governs admissibility.

---

# What Gantral Enforces

Gantral enforces execution-time invariants:

- Canonical authority state machine  
- Explicit transition validity  
- Identity binding at approval boundaries  
- Policy version binding at decision time  
- Workflow version binding  
- Context snapshot binding  
- Atomic state transition + artifact emission  
- Tamper-evident artifact chains  
- Log-independent replay  
- Fail-closed semantics  

Authority is modeled as execution state — not metadata.

Illegal transitions terminate execution.

No ambient credentials are treated as authority.

---

# Deterministic Authority Model

Canonical authority state progression:

```

CREATED → RUNNING → WAITING_FOR_HUMAN
→ APPROVED / REJECTED / OVERRIDDEN
→ RESUMED → COMPLETED / TERMINATED

```

Transitions:

- Are explicitly enumerated  
- Must satisfy valid transition relations  
- Are atomic with artifact persistence  
- Fail closed if artifact write fails  
- Are identity-bound and version-bound  

Execution cannot proceed past an authority boundary without structural evidence.

---

# Policy Advisory Support (OPA Integration)

Gantral supports external policy engines in an advisory capacity.

Policy may be authored in:

- Open Policy Agent (OPA) / Rego  
- Custom risk engines  
- Enterprise policy services  

At authority checkpoints:

1. Policy is evaluated externally.
2. Policy returns advisory output.
3. The `policy_version_id` is bound to the execution instance.
4. Gantral enforces the authority transition structurally.

Policy remains advisory.  
Authority enforcement remains internal and deterministic.

Policy thresholds can evolve without modifying workflow code.

Gantral enforces authority as canonical state — not policy logic.

---

# Commitment Artifacts

At each authority transition, Gantral emits a commitment artifact that binds:

- Workflow version  
- Policy version  
- Validated identity  
- Context snapshot hash  
- Justification metadata  
- Recursive hash linkage  

Artifacts form an append-only, tamper-evident chain.

Modification invalidates downstream verification.

Authority history becomes cryptographically verifiable infrastructure.

---

# Log-Independent Replay

Gantral supports offline replay verification.

Replay validates:

- Hash-chain integrity  
- Authority state transition correctness  
- Version consistency  
- Policy binding  
- Identity binding  

Replay requires:

- No runtime access  
- No database access  
- No logs  
- No agent memory  

Output:

```

VALID / INVALID / INCONCLUSIVE

```

Authority becomes independently inspectable.

---

# Trust as Infrastructure

Gantral is:

- Open source (Apache 2.0)  
- Self-hosted  
- Infrastructure-grade  
- Deterministic by design  
- Transparent in enforcement semantics  
- Vendor-neutral  

Enterprises can:

- Inspect enforcement logic  
- Audit authority semantics  
- Verify replay independently  
- Avoid opaque runtime dependencies  

Trust is not claimed.

It is structurally enforced and mechanically verifiable.

---

# What Enterprises Gain

With Gantral, enterprises can:

- Enforce consistent authority semantics across workflows  
- Prevent policy-in-code drift  
- Improve audit defensibility  
- Produce replayable authority evidence  
- Reduce post-incident reconstruction effort  
- Separate intelligence from authority  
- Introduce explicit authority checkpoints in agentic systems  

Determinism replaces reconstruction.

Authority becomes infrastructure — not convention.

---

# When Gantral Becomes Necessary

Gantral becomes rational infrastructure when:

- Financial exposure is material  
- Regulatory scrutiny is plausible  
- Litigation risk exists  
- Decisions may be challenged years later  
- AI influences infrastructure or access control  
- Board-level explanation may be required  

If the answer to:

“Could we independently prove authority correctness without logs?”

is uncertain — structural authority is warranted.

---

# Design Principles

Gantral is built on five principles:

1. Authority must be structurally enforced, not reconstructed.
2. Authority must remain separate from intelligence.
3. Execution authority must be deterministic and replayable.
4. Replay must be log-independent.
5. Enforcement must fail closed.

---

# What Gantral Is NOT

Gantral is not:

- A workflow engine  
- An agent framework  
- A governance dashboard  
- A compliance certification tool  
- A GRC platform  
- A policy lifecycle manager  
- An autonomy orchestration system  

Gantral is the execution authority kernel.

Minimal by design.  
Infrastructure-grade by intent.

---

# Project Status

Gantral is actively developed with focus on:

- Deterministic authority semantics  
- Artifact immutability guarantees  
- Replay rigor  
- Policy separation discipline  
- Transparent enforcement boundaries  

---

# Design Partner Collaboration

Gantral collaborates with organizations exploring:

- Agentic AI in regulated environments  
- High-materiality execution workflows  
- Deterministic authority infrastructure  

If your platform team is evaluating AI execution at scale,  
we welcome structured engagement.

📩 abhishek@rainminds.com  

---

© 2025 Rainminds Solutions Pvt. Ltd.  
Licensed under the Apache License, Version 2.0.