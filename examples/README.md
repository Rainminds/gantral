
# Gantral Integration Examples

This directory contains reference implementations for integrating agents with Gantral.

> **Developer Note:** These examples are designed to be run as self-contained Docker Compose environments.

## 1. Persistent Agent (Pattern A)
**Path:** [`persistent-agent/`](persistent-agent/)  
**Best For:** Advanced frameworks (CrewAI, LangGraph) that support native checkpointing.

This demo proves the **"Zero CPU"** architecture.
- An agent hits a sensitive step and checks Gantral.
- If approval is needed, the agent **saves state to disk** and **exits process (Exit Code 3)**.
- The `runner` container detects this, signals Gantral to SUSPEND, and frees up resources.
- Once approved, the `runner` launches a **fresh process** which loads the checkpoint and completes the task.

**[Quick Start]**
```bash
cd persistent-agent
docker compose up --build
# In new terminal:
./scripts/trigger.sh
./scripts/approve.sh <execution_id>
```

## 2. Split Agent (Pattern B)
**Path:** [`split-agent/`](split-agent/)  
**Best For:** Simple scripts or legacy bots that cannot magically "resume" mid-function.

This demo shows how to handle "dumb" agents.
- **Agent Pre**: Runs, prepares context, requests decision, and **terminates**.
- **Runner**: Waits (polling) for Human Approval.
- **Agent Post**: Launched *only* after approval, reading the passed context file.

**[Quick Start]**
```bash
cd split-agent
docker compose up --build
# In new terminal:
./scripts/trigger.sh
./scripts/approve.sh <execution_id>
```

## 3. Policy Library
**Path:** [`policies/`](policies/)

A standard library of drop-in Rego policies for common governance patterns:
- `basic_approval.rego`: Materiality-based checks.
- `multi_step.rego`: Dual-approval for critical ops.
- `timeout.rego`: Time-based auto-rejection.
- `auto_approve.rego`: Read-only whitelist.

See [`policies/README.md`](policies/README.md) for testing instructions.

## Verification
To run the full verification suite across these examples, refer to the [Master Verification Plan](VERIFICATION.md).
