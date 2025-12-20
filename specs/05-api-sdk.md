# APIs & SDKs

## Standards

### Primary API
- **Protocol:** REST (HTTP/1.1 or HTTP/2).
- **Format:** JSON.
- **Spec:** OpenAPI 3.1.
- **Versioning:** Semantic versioning in URL (e.g., `/api/v1/...`).

### Core API Groups
- `/workflows`: Manage templates.
- `/instances`: Manage executions (create, stop, resume).
- `/decisions`: Submit approvals/rejections.
- `/policies`: CRUD for governance rules.
- `/audit`: Read-only access to immutable logs.
- `/replay`: Deterministic replay triggers.

### Internal/Future API
- **Protocol:** gRPC.
- **Use Case:** High-volume internal communication between Control Plane components.

## SDKs
Thin wrappers around the API to facilitate integration.
- **Python:** For AI/Data Science teams.
- **Go:** For backend services and performance-critical agents.
- **TypeScript:** For frontend/Node.js integrations.
