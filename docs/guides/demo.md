---
sidebar_position: 1
title: End-to-End Demo
---

# Gantral End-to-End Demo (Phase 4 Verified)

This guide demonstrates the fully integrated **Gantral Control Plane** using the reference Docker environment.

It validates the core architectural promise: **"Zero CPU Hibernation"**. You will see an agent process start, hit a policy barrier, exit completely to save resources, and then resume only after you approve it.

## Prerequisites
*   Docker & Docker Compose
*   `git`
*   `curl` (optional, scripts provided)

---

## 1. Start the Environment

Clone the repository and navigate to the persistent agent example. This sets up Gantral Core, Temporal, and a specialized Runner.

```bash
git clone https://github.com/Rainminds/gantral.git
cd gantral/examples/persistent-agent

# Start the stack (detached mode)
docker compose up -d
```

*Wait ~10 seconds for services to become healthy.*

*   **Gantral Dashboard:** [http://localhost:8080](http://localhost:8080)
*   **Temporal UI:** [http://localhost:8081](http://localhost:8081)

---

## 2. Trigger the Agent

We will launch a workflow that simulates a "Production Deployment." The policy is configured to **REQUIRE_HUMAN** for this action.

Run the trigger script:

```bash
./scripts/trigger.sh
```

**What just happened?**
1.  A new execution instance was created (`CREATED`).
2.  The Runner picked up the task.
3.  The Agent process started, checked the policy, saved its state to disk, and **EXITED**.
4.  The Runner reported `SUSPENDED` to Gantral.

---

## 3. Observe "Zero CPU" State

Verify that the system is waiting for you, but consuming no compute resources for the agent.

```bash
./scripts/status.sh <INSTANCE_ID>
# Output: WAITING_FOR_HUMAN
```

**Verify in UI:**
*   Go to **[http://localhost:8080](http://localhost:8080)**. You will see the instance paused.
*   Go to **[http://localhost:8081](http://localhost:8081)** (Temporal). You will see the workflow is "Running" (sleeping), but there is no active Activity Worker processing it.

---

## 4. Approve Execution

Now, act as the "Manager" and grant authority.

```bash
./scripts/approve.sh <INSTANCE_ID>
```

**What happens next?**
1.  Gantral records your decision (auditable event).
2.  Gantral signals the Temporal workflow to wake up.
3.  The Runner receives a new task: "Resume Execution."
4.  The Runner launches a **NEW** agent process.
5.  The Agent loads its checkpoint and finishes the job.

---

## 5. Verify Completion

Check the final status:

```bash
./scripts/status.sh <INSTANCE_ID>
# Output: COMPLETED
```

You have just demonstrated a **federated, policy-controlled, auditable AI workflow** that survived a complete process restart.

---

## 6. Cleanup

To stop the environment:

```bash
docker compose down
```
