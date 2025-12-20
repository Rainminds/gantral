Template: When running tests:
1.  **Unit Tests:** Run `make test` to execute all unit tests.
2.  **Integration Tests:** Ensure Docker environment is up (`make up`), then run `make test-integration` (if available) or specific test tags.
3.  **Linting:** Run `golangci-lint run` to check for style and common errors.
4.  **Fix Failures:** detailed analysis of any failures before proceeding.
