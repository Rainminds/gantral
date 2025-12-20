# Technology Stack & Standards

## Core Stack
- **Language:** Go (Golang) for the core execution engine (performance/concurrency).
- **SDKs:** Python, Go, TypeScript (thin wrappers).
- **Database:** PostgreSQL 16 (Event-sourced tables).
- **Logs:** PostgreSQL (initially), ClickHouse (future) for immutable audit logs.
- **Cache:** Redis (optional, for state caching).

## API Standards
- **Primary:** REST (OpenAPI 3.1).
- **Future:** gRPC for high-volume internal communication.
- **Auth:** OAuth 2.0 / OIDC (Service accounts per agent).

## Repository Structure (Monorepo)
- `docs/`: Guides & TRDs.
- `specs/`: Technical specifications (SSOT).
- `core/`:
    - `engine/`: State machine logic.
    - `policy/`: Rule evaluation.
    - `hitl/`: Human interaction logic.
    - `audit/`: Logging.
- `api/`: OpenAPI specs.
- `sdk/`: Client libraries.
- `adapters/`: Integration code (Slack, Jira, GitHub).
- `infra/`: Docker/K8s configs.

## Deployment
- **Dev:** Docker Compose.
- **Prod:** Kubernetes (Helm charts).
- **Observability:** OpenTelemetry (Metrics, Traces).
