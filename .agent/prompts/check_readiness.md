# Gantral Readiness Checklist (CI-Parity)
Target: Ensure 100% compliance with GitHub Actions and GitLab CI before committing.

## Setup
- Ensure local Postgres is running: `docker-compose up -d postgres`
- Ensure Temporal is running (for integration tests).

## Instructions for Agent
Execute the following steps sequentially and report any failures. Do not proceed to the next step if one fails.

1.  **Build Check (Core & Runner)**:
    - Run `make build` to verify the Go server compiles.
    - Run `cd examples/persistent-agent/runner && pip install -r requirements.txt && python -m compileall .` to verify the Python runner.

2.  **Lint & Architecture Check**:
    - Run `make lint` to verify Go code quality.
    - Run `./scripts/verify-architecture.sh` to ensure boundary rules are respected.

3.  **Tiered Testing**:
    - Run `make test-tier1` (Unit, Logic, Artifacts).
    - Run `make test-tier2` (Integration, Engine).
    - **Note**: Integration tests must use typed `engine.State` and terminal states for replay.

4.  **Coverage Validation**:
    - Run: `go test -tags=integration -coverpkg=./... -coverprofile=coverage.raw.out ./...`
    - Filter and check: `grep -v -E "infra/|cmd/|pkg/logger/|pkg/config/|sdk/|github.com/Rainminds/gantral/main.go" coverage.raw.out > coverage.out`
    - Verify: `TOTAL=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//'); if (( $(echo "$TOTAL < 70.0" | bc -l) )); then echo "Coverage $TOTAL% is below 70%"; exit 1; fi`

5.  **Documentation Integrity**:
    - Run `make docs-build` to catch MDX or Docusaurus compilation errors.

6.  **Compliance & Hygiene**:
    - Verify all staged commits are signed (DCO): `git commit -s` is required.
    - Secret Scan: Verify no `.env`, tokens, or private keys are being committed.
    - Clean: Remove `coverage.out`, `coverage.raw.out`, and orphaned `.diff` files.

7.  **Final Confirmation**:
    - Only if all above pass, mark as "Ready for Commit".
