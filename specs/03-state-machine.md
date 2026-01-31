# Execution State Machine

The execution of a Gantral Instance is governed by a strict state machine.

## Canonical Transitions (Executable Rules)

Allowed transitions ONLY:
- **CREATED** → **RUNNING**
- **RUNNING** → **WAITING_FOR_HUMAN**
- **WAITING_FOR_HUMAN** → **APPROVED** | **REJECTED** | **OVERRIDDEN**
- **APPROVED** | **OVERRIDDEN** → **RESUMED**
- **RESUMED** → **RUNNING**
- **RUNNING** → **COMPLETED** | **TERMINATED**

> **CRITICAL:** Any other transition MUST panic and terminate execution.
