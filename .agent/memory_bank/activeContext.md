# Active Context

## Current Phase: Phase 3 Integration & Phase 4 Prep
We have completed the **Core Runtime** aspects of Phase 3 (Temporal Integration, Adapters) and structurally implemented the **Federated Runner** model (Phase 5.3). We are now positioned to build the **SDKs** (Phase 3.3) and **Reference Implementations** (Phase 4).

## Recent Changes
- **Runtime:** Implemented `GantralExecutionWorkflow` (Phase 3.1) and `ExecutionActivities`.
- **Architecture:** Split `cmd/server` and `cmd/worker` to enforce Federated Execution (Phase 5.3 pulled forward).
- **Governance:** Implemented "Policy as Transition Guard" (Phase 2.2).
- **Verification:** Passed Unit, Integration, and E2E tests for the Core Control Plane.
- **Documentation:** Updated to reflect the v9.0 Build Plan authority.

## Immediate Goals
1.  [ ] **Complete Phase 3 (Gaps):**
    *   [ ] **3.3 SDKs:** Build Python/TS clients to wrap the API.
    *   [ ] **3.4 Observability:** Add OpenTelemetry/Metrics export.
2.  [ ] **Execute Phase 4 (Developer Experience):**
    *   [ ] **4.1 Reference Proxy:** Simple agent connecting to Gantral.
    *   [ ] **4.4 Framework Impls:** CrewAI/LangGraph examples.

### Current Focus
- **Stability:** The Core Control Plane is "Feature Complete" for MVP.
- **Adoption:** The next critical step is making it usable via SDKs (Phase 3.3).

## Open Decisions
- **SDK Generation:** `openapi-generator` vs Hand-rolled? (Leaning hand-rolled for DX).
- **Identity Mock:** Need a local development mock for Phase 5.1 (Identity Federation).
