
# Phase 4 Verification Runbook

This document is the **Master Test Plan** for verifying the Phase 4 Consumer Patterns:
1. **Persistent Pause** (`examples/persistent-agent`)
2. **Split Agent Pattern** (`examples/split-agent`)

## Prerequisites
- Docker & Docker Compose installed.
- Ports 8080, 8081, 7233, 5432 must be free.

---

## 1. Persistent Agent Demo (Pattern A)
**Goal**: Verify that an agent can hibernate (Exit Code 3) and resume from a saved checkpoint after human approval.

### Setup & run
```bash
cd examples/persistent-agent
docker compose up --build -d
```
*Wait for services to stabilize (check `docker compose logs -f runner`).*

### Test Case 1.1: Standard Approval Flow
1. **Trigger Execution**:
   ```bash
   ./scripts/trigger.sh
   # Copy the Execution ID from the output (or use 'latest' logic if implemented)
   ```
2. **Verify Pause**:
   - Check logs: `docker compose logs runner`
   - Expect: `Agent requested Hibernation`.
   - Check status: `./scripts/status.sh <ID>` -> Should be `SUSPENDED` / `WAITING_FOR_HUMAN` (or `PENDING` decision).
3. **Approve**:
   ```bash
   ./scripts/approve.sh <ID>
   ```
4. **Verify Resume**:
   - Check logs: `docker compose logs -f runner`
   - Expect: `Approval Granted! Resuming...` followed by `Task Complete`.
   - Check status: `./scripts/status.sh <ID>` -> Should be `COMPLETED`.

### Test Case 1.2: Red Team - Runner Restart
1. Trigger a new execution.
2. Wait for it to Hibernate (logs show `Exit code 3`).
3. **Kill the Runner**: `docker compose restart runner`
4. **Approve** the execution via script.
5. **Expectation**: Runner comes back up, polls Gantral, sees `APPROVED`, and resumes the agent correctly using the checkpoint on disk.

### Reset
```bash
docker compose down -v
```

---

## 2. Split Agent Demo (Pattern B)
**Goal**: Verify that an "agent" can be split into two separate processes (Pre/Post) with a file-based handoff, surviving a complete process termination.

### Setup & Run
```bash
cd examples/split-agent
docker compose up --build -d
```

### Test Case 2.1: Standard Split Flow
1. **Trigger Execution**:
   ```bash
   ./scripts/trigger.sh
   ```
2. **Verify Pre-Agent Success**:
   - Logs: `Agent-Pre finished successfully`.
   - Logs: `Marking PRE_DONE`.
   - Status: `./scripts/status.sh <ID>` -> `RUNNING` (waiting for decision).
3. **Approve**:
   ```bash
   ./scripts/approve.sh <ID>
   ```
4. **Verify Post-Agent Launch**:
   - Logs: `Decision APPROVED... Launching Post-Agent`.
   - Logs: `Agent-Post finished successfully`.
   - Status: `COMPLETED`.

### Test Case 2.2: Red Team - State Loss
1. Trigger execution -> Wait for `PRE_DONE`.
2. **Restart Runner**: `docker compose restart runner`
   - *Note: In-memory state `local_state` is lost.*
3. **Approve** execution.
4. **Expectation**:
   - Runner sees `WAITING_FOR_HUMAN` or `REQUEST_APPROVAL` status.
   - Runner logic checks status and skips re-running `agent_pre`.
   - Runner sees Approval (after step 3) and runs `agent_post`.
   - **Failure Mode**: If Runner immediately re-runs `agent_pre`, verify `agent_pre` is safe to re-run (idempotent).

### Reset
```bash
docker compose down -v
```

---

## 3. Common Failure Modes

### Gantral Core Unreachable
- **Symptom**: Runner logs `Error polling Gantral: Connection refused`.
- **Fix**: Ensure `gantral-core` container is healthy. Check `docker compose ps`.

### Permission Denied (Volume Maps)
- **Symptom**: Agent fails with `Permission denied` writing to `checkpoint/` or `handoff/`.
- **Fix**: Ensure host directories are writable by the Docker user (usually root in standard docker, but check UID mapping). `chmod -R 777 agent/checkpoint` if needed (for demo only).
