# Execution State Machine

The execution of a Gantral Instance is governed by a strict state machine.

## Canonical States

- **CREATED:** Instance is initialized but not yet started.
- **RUNNING:** Workflow is actively executing steps.
- **WAITING_FOR_HUMAN:** Execution paused. A policy condition triggered a requirement for human intervention.
- **APPROVED:** Human authorized the action. Transitions back to RESUMED/RUNNING.
- **REJECTED:** Human denied the action. Transitions to TERMINATED or a remediation path.
- **OVERRIDDEN:** Human forced a decision, potentially bypassing policy logic.
- **RESUMED:** Transient state after approval/override before re-entering RUNNING.
- **COMPLETED:** Workflow finished successfully.
- **TERMINATED:** Workflow stopped due to error or rejection.
