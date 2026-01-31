---
title: End-to-End Demo
---

# Gantral End-to-End Demo (Execution Authority Demonstration)

This guide demonstrates Gantral’s **execution authority model** using the reference Docker environment.

The demo validates three core properties:

1. **Execution pauses at governed boundaries**
2. **Human authority is enforced as execution state**
3. **Execution can safely hibernate and resume without agent processes running**

This demo focuses on **execution semantics**, not full audit admissibility.

---

## What This Demo Proves (and What It Does Not)

### This Demo Proves

- Authority is enforced **before execution proceeds**
- Human-in-the-loop (HITL) is a **blocking execution state**
- Agent processes can **exit completely** during approval waits (zero CPU)
- Execution resumes via a **new process** after authority is granted
- Authority decisions are recorded as execution events

### This Demo Does NOT Prove

- Offline third-party verification
- Long-term artifact custody
- Cryptographic replay guarantees
- External auditor independence

Those properties are covered in the **Verifiability** documentation and will be demonstrated separately.

---

## Prerequisites

- Docker
- Docker Compose
- `git`
- `curl` (optional; helper scripts provided)

---

## 1. Start the Environment

Clone the repository and navigate to the demo directory:

```bash
git clone https://github.com/Rainminds/gantral.git
cd gantral/examples/persistent-agent
````

Start the reference stack:

```bash
docker compose up -d
```

Wait ~10 seconds for services to become healthy.

### Available Interfaces

* **Gantral API / UI:** [http://localhost:8080](http://localhost:8080)
* **Temporal UI:** [http://localhost:8081](http://localhost:8081)

---

## 2. Trigger an Execution Instance

This demo simulates a **material action** (e.g., production deployment).

The configured policy requires **explicit human authority**.

Trigger the workflow:

```bash
./scripts/trigger.sh
```

### What Happens Internally

1. A new execution instance is created (`CREATED`)
2. Execution enters `RUNNING`
3. Policy evaluation requires human authority
4. Execution transitions to `WAITING_FOR_HUMAN`
5. The agent checkpoints its state and **exits**
6. The runner reports `SUSPENDED`

At this point, **no agent process is running**.

---

## 3. Observe the Paused State (Zero CPU)

Check the execution state:

```bash
./scripts/status.sh <INSTANCE_ID>
# Expected: WAITING_FOR_HUMAN
```

### Verify via UI

* Gantral UI shows the instance waiting for authority
* Temporal UI shows the workflow **idle** (no active worker)

This demonstrates that:

* Execution state persists
* Compute resources are not consumed during long waits

---

## 4. Grant Human Authority

Act as the approving human and grant authority:

```bash
./scripts/approve.sh <INSTANCE_ID>
```

### What Happens Next

1. The authority decision is recorded as an execution event
2. The workflow runtime is signaled
3. The runner receives a resume task
4. A **new agent process** is launched
5. The agent restores its checkpoint
6. Execution continues

No previous process is reused.

---

## 5. Verify Completion

Check final state:

```bash
./scripts/status.sh <INSTANCE_ID>
# Expected: COMPLETED
```

You have now observed:

* Authority-enforced execution
* Process-level hibernation and resume
* Deterministic control over execution flow

---

## 6. Inspect Execution History (Optional)

You may inspect execution history via:

* Gantral API responses
* Temporal workflow history
* Runner logs

These show:

* State transitions
* Authority decisions
* Suspension and resume events

They are **execution records**, not long-term audit artifacts.

---

## 7. Cleanup

To stop the environment:

```bash
docker compose down
```

---

## Final Reminder

This demo is about **authority at execution time**.

Gantral’s core promise is not that execution is clever —
it is that execution is **allowed, paused, resumed, or stopped intentionally and provably**.

