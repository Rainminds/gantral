# Infrastructure & Deployment

## Data Stores
- **PostgreSQL 16:** Primary operational database. Stores metadata and indices only.
- **Object Storage:** (S3/GCS/MinIO) **Immutable** storage for Commitment Artifacts.
- **Redis:** (Optional) Non-authoritative caching layer.
- **ClickHouse:** (Future) Purpose-built store for massive scale immutable audit logs.

## Deployment Models
- **Dev:** `docker-compose`. Single box, easy start.
- **Prod:** Kubernetes (`k8s`). Scalable, resilient deployment using Helm charts.
- **CI/CD:** GitHub Actions.
- **Observability:** OpenTelemetry (OTEL) for metrics and distributed tracing.
