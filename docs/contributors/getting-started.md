---
sidebar_position: 1
title: Getting Started
---

# Getting Started with Gantral

Welcome to the Gantral project.

Gantral is an open-source **AI Execution Control Plane**. We focus on execution semantics, human oversight, and auditability for AI-enabled workflows.

This guide explains how to explore the project, run the demos, and contribute responsibly.

---

## 🎯 Who This Project Is For

Gantral is primarily intended for:
- **Platform Engineers** building internal AI platforms.
- **Infrastructure Architects** designing control planes.
- **Security & Governance Teams** auditing AI behavior.
- **Contributors** interested in distributed systems and execution authority.

> **Note:** If you are looking to build agents, optimize prompts, or experiment with LLM reasoning, Gantral is likely not the tool you need. We integrate *with* those tools; we do not replace them.

---

## 🚀 Quick Start (Run the Code)

The best way to understand Gantral is to see it enforce authority.

### 1. Run the "Persistent Agent" Demo
We have a fully dockerized example that demonstrates a "Zero CPU" hibernation pattern.

```bash
git clone https://github.com/Rainminds/gantral.git
cd gantral/examples/persistent-agent

# Start the stack (Core, Temporal, Runner, Agent)
docker compose up
```

Once running, use the included scripts to interact with the authority layer:

```bash
./scripts/trigger.sh   # Creates a workflow
./scripts/status.sh    # Observe "WAITING_FOR_HUMAN" state
./scripts/approve.sh   # Grant authority
```

---

## 📚 Essential Reading

Before contributing code, please ground yourself in the architecture:

1.  **[What is Gantral?](../positioning/what-is-gantral.md)** – High-level philosophy.
2.  **[Technical Reference (TRD)](../architecture/trd.md)** – The "Constitution" of the project.
3.  **[Consumer Guide](../guides/example-consumer-integration.md)** – How agents are expected to behave.

---

## 🛠️ Repository Structure

Gantral follows a **specs-first model**. Code that diverges from documented semantics is considered a bug.

| Folder | Purpose |
| :--- | :--- |
| `core/` | The Authority Service (API, State Machine). |
| `examples/` | Reference implementations and demos. |
| `specs/` | The Single Source of Truth for architecture. |
| `docs/` | Project documentation (this site). |

---

## 🤝 Before You Contribute

This project prioritizes **correctness and clarity over speed**.

1.  **Read the Invariants:** Understand why we don't store agent memory.
2.  **Check the Roadmap:** Ensure your feature isn't a "Non-Goal."
3.  **Start Small:** Fix a bug or improve docs before proposing architectural changes.

*See [How to Contribute](./how-to-contribute.md) for the full PR process.*
