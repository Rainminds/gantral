# Technology Stack & Standards

## Core Stack
- **Language:** Go 1.25+ (Core).
- **Runtime:** **Temporal** (Workflow Orchestration & Durability).
- **Database:** PostgreSQL 16 (Event-sourced tables).
- **Logs:** PostgreSQL (Audit Trail).
- **Auth:** OIDC (Federated Identity).

## API Standards
- **Primary:** REST (JSON over HTTP).
- **Schema:** Handlers defined in `adapters/primary/http`.
- **Auth:** Bearer Token (JWT).

## Repository Structure (Monorepo)
- `docs/`: Guides & Documentation.
- `specs/`:
    - `adr/`: Architectural Decision Records.
- `core/`:
    - `engine/`: State machine logic.
    - `policy/`: Rule evaluation.
    - `activities/`: Temporal Activities.
    - `workflows/`: Temporal Workflows.
- `adapters/`:
    - `primary/http`: REST API handlers.
    - `secondary/postgres`: Data access.
- `cmd/`:
    - `server/`: Control Plane Binary.
    - `worker/`: Execution Plane Binary.
- `infra/`: Docker/K8s configs & Migrations.
- `tests/`: Integration & E2E Suites.

## Deployment
- **Dev:** `make up` (Docker Compose for Infra), `go run` (Services).
- **Prod:** Kubernetes (Helm charts for Gantral + Temporal).
- **Observability:** OpenTelemetry (Metrics, Traces).
