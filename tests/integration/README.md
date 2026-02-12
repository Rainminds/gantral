# Integration Tests (Tier 2)

## Scope
Runtime & Integration tests requiring external dependencies:
- Temporal workflow runtime (test suite or local server)
- OPA policy engine (mock or embedded)
- PostgreSQL database (testcontainers or in-memory)
- OIDC identity provider (mock)
- Full Gantral API server

## Coverage (Test Inventory Sections)
- **Section C**: Policy engine integration (OPA ALLOW/DENY/REQUIRE_HUMAN, fail-closed)
- **Section G**: Identity & security (OIDC, JWT, role enforcement, team isolation)
- **Section H**: Temporal determinism (workflow replay, activities, signals, timers)
- **Section J**: End-to-end execution (trigger → HITL → approval → artifact → verify)
- **Section O**: Concurrency & atomicity (parallel instances, lock contention)
- **Section T**: Runner & federation (if implemented)
- **Section U**: Agent integration (if implemented)

## Runtime Target
2-5 minutes per merge to main (per Tier 2 strategy)

## Running Tests

### Prerequisites
- Docker (for testcontainers)
- Go 1.21+
- Make

### Quick Run
```bash
# Run all Tier 2 tests
make test-tier2

# Run specific test file
go test -tags=integration -v ./tests/integration/policy_integration_test.go

# Run with verbose output
go test -tags=integration -v ./tests/integration/...
```

### With Docker Compose (if needed)
```bash
docker-compose -f docker-compose.test.yml up -d
go test -tags=integration ./tests/integration/...
docker-compose -f docker-compose.test.yml down
```

## Test Structure
- `fixtures/`: Shared mocks, test workflows, fixture data
- `*_integration_test.go`: Integration test files (one per inventory section)
- `api/`: Web/API integration tests (if applicable)

## Writing New Integration Tests
1. Add `//go:build integration` build tag at top
2. Use `package integration_test`
3. Use testcontainers for external dependencies
4. Ensure deterministic despite external services
5. Keep total runtime <5 minutes
6. Follow naming convention: TestFeature_Scenario_ExpectedOutcome

## Dependencies
- Temporal test suite: `go.temporal.io/sdk/testsuite`
- Testcontainers: `github.com/testcontainers/testcontainers-go`
- Mock OPA: internal implementation in `fixtures/mock_opa.go`
- Mock IdP: internal implementation in `fixtures/mock_idp.go`
