# Active Context

## Current Phase: Initialization (Phase 1)
We are currently scaffolding the repository based on the TRD v1.1.


## Recent Changes

- Added CNCF governance files (GOVERNANCE, CONTRIBUTING, SECURITY, etc).
- Initialized ADR system in `specs/adr/`.
- Initialized directory structure and modular specs (`specs/`).
- Initialized Go module (`github.com/Rainminds/gantral`).
- Added `.gitignore` and preserved empty directories with `.gitkeep`.
- Created Agent Prompt Templates (`implement_feature`, `run_tests`, `create_release`, etc) in `.agent/prompts/`.
- Initialized `docs/` structure for future Docusaurus migration.
- Configured GitHub Actions workflows (`ci.yml`, `dco.yml`, `security.yml`, `release.yml`) and `.goreleaser.yaml`.

## Immediate Goals
1.  [x] Initialize Repository & License (Done).
2.  [x] **Setup DevOps:** GitHub Actions & workflows.
3.  [ ] **Define API Spec:** Create `api/openapi.yaml` defining the `Workflow` and `Instance` resources.
4.  [ ] **Setup Infra:** Create `infra/docker-compose.yml` with Postgres and Redis.
5.  [ ] **Implement Core:** Start coding `core/engine`.

## Open Decisions
- Confirming Go module structure for the monorepo.
- Finalizing the exact JSON schema for the "Decision" object.


