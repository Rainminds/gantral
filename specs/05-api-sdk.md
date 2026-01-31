# APIs & SDKs

## Standards

### Primary API
- **Protocol:** REST (HTTP/1.1 or HTTP/2) - External/Public.
- **Format:** JSON.
- **Spec:** OpenAPI 3.1.
- **Versioning:** Semantic versioning in URL (e.g., `/api/v1/...`).

### Internal API
- **Protocol:** gRPC (Required for internal service-to-service).
- **Use Case:** High-volume internal communication between Control Plane components and Runners.

### Core API Groups
- `/workflows`: Manage templates.
- `/instances`: Manage executions (create, stop, resume).
- `/decisions`: Submit approvals/rejections.
- `/policies`: CRUD for governance rules.
- `/audit`: Read-only access to immutable logs.
- `/artifacts`: Retrieve cryptographic commitment artifacts.
- `/verify`: Online verification endpoint (use CLI for offline).
- `/replay`: Deterministic replay triggers.



## SDKs
Thin wrappers around the API to facilitate integration.
- **Python:** For AI/Data Science teams.
- **Go:** For backend services and performance-critical agents.
- **TypeScript:** For frontend/Node.js integrations.
