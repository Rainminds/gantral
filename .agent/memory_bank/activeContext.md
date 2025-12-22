# Active Context

## Current Phase: Initialization (Phase 1)
We are currently scaffolding the repository based on the TRD v1.1.


## Recent Changes
- Added `main.go` placeholder to satisfy `golangci-lint` "no go files" error.

## Immediate Goals
1.  [x] Initialize Repository & License (Done).
2.  [x] **Setup DevOps:** GitHub Actions & workflows.
3.  [ ] **Define API Spec:** Create `api/openapi.yaml` defining the `Workflow` and `Instance` resources.
4.  [ ] **Setup Infra:** Create `infra/docker-compose.yml` with Postgres and Redis.
5.  [ ] **Implement Core:** Start coding `core/engine`.

## Open Decisions
- Confirming Go module structure for the monorepo.
- Finalizing the exact JSON schema for the "Decision" object.


