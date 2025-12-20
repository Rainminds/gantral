# Infrastructure & Deployment

## Data Stores
- **PostgreSQL 16:** Primary operational database. Stores state, policies, and active data. Uses event sourcing patterns.
- **Redis:** (Optional) Caching layer for high-speed state access.
- **ClickHouse:** (Future) Purpose-built store for massive scale immutable audit logs.

## Deployment Models
- **Dev:** `docker-compose`. Single box, easy start.
- **Prod:** Kubernetes (`k8s`). Scalable, resilient deployment using Helm charts.
- **Observability:** OpenTelemetry (OTEL) for metrics and distributed tracing.
