# Gantral End-to-End Demo (Phase 3)

This guide demonstrates the fully integrated **Gantral Control Plane** running with Temporal as the execution backend. It validates:
- **Phase 1:** Canonical State Machine & Immutable History
- **Phase 2:** Governance Hardening (Policy Engine)
- **Phase 3:** Enterprise Integration (Temporal Workflow & Worker)

## Prerequisites
- Docker & Docker Compose
- Go 1.22+
- `curl` (for API requests)
- `jq` (optional, for pretty-printing JSON)

---

## 1. Start Infrastructure
Start the required backing services (PostgreSQL, Temporal, Temporal UI).

```bash
make up
```
*Wait for a few seconds for services to become healthy.*

- **PostgreSQL:** Port `5432`
- **Temporal Server:** Port `7233`
- **Temporal UI:** [http://localhost:8081](http://localhost:8081)

## 2. Start the Worker (Execution Plane)
The Worker is responsible for executing the workflow logic and activities (DB persistence).

Open a **new terminal**:
```bash
# Ensure your .env file is present or relying on defaults
go run cmd/worker/main.go
```
*You should see logs indicating the Worker started successfully.*

## 3. Start the API Server (Control Plane)
The API Server handles HTTP requests, policy evaluation triggers, and dashboard serving.

Open a **new terminal**:
```bash
go run cmd/server/main.go
```
*Server listening on [http://localhost:8080](http://localhost:8080).*

---

## 4. Trigger High-Materiality Workflow
We will trigger a workflow that mimics a high-risk operation (e.g., "Deploy to Production"). The policy is configured to **REQUIRE HUMAN APPROVAL**.

```bash
curl -X POST http://localhost:8080/instances \
  -H "Content-Type: application/json" \
  -d '{
    "workflow_id": "demo-deploy-v1",
    "trigger_context": {"requested_by": "alice", "environment": "production"},
    "policy": {
      "id": "policy-prod-001",
      "materiality": "HIGH",
      "requires_human_approval": true
    }
  }' | jq
```

**Expected Output:**
```json
{
  "id": "inst-...",
  "status": "PENDING"
}
```
*Note the `id` (e.g., `inst-123...`). You will use it below.*

---

## 5. Observe Execution (The "Pause")

### A. Temporal UI (The Runtime)
Navigate to [http://localhost:8081](http://localhost:8081).
- You will see the workflow `GantralExecutionWorkflow` in state **Running**.
- Click into it. You will see pending `Timer` (the storage-less wait) or Signal selector.
- *This confirms the infrastructure is holding the state reliably.*

### B. Gantral Dashboard (The Authority)
Navigate to [http://localhost:8080](http://localhost:8080).
- You will see the instance in state **WAITING_FOR_HUMAN**.
- The "Approve" and "Reject" buttons are visible.

### C. Audit Log (The Evidence)
Check what has been recorded so far:
```bash
# Replace <INSTANCE_ID> with your actual ID
curl http://localhost:8080/instances/<INSTANCE_ID>/audit | jq
```
*You should see events for creation and the policy decision to pause.*

---

## 6. Make a Human Decision
You can approve via the [Dashboard](http://localhost:8080) or via API.

**Via API:**
```bash
curl -X POST http://localhost:8080/instances/<INSTANCE_ID>/decisions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "APPROVE",
    "actor_id": "bob-admin",
    "justification": "Verified changes, clear to proceed"
  }'
```

---

## 7. Verify Completion

1.  **Dashboard:** The instance state updates to **APPROVED**.
2.  **Temporal UI:** The workflow status changes to **Completed**.
3.  **Audit Log:** A new event `DECISION_RECORDED` appears with `bob-admin`'s identity.

## 8. Persistence Test (Optional)
1. Stop the API Server (Ctrl+C).
2. Stop the Worker (Ctrl+C).
3. Restart both.
4. Refresh the [Dashboard](http://localhost:8080).
5. The `APPROVED` instance and its audit trail are still there.

---
**Success!** You have demonstrated a federated, policy-controlled, auditable AI workflow execution.
