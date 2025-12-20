Template: Before committing:
1.  **Build Check:** Run `make build` to ensure the project compiles.
2.  **Test Check:** Run `make test` to ensure no regressions.
3.  **Lint Check:** Run `golangci-lint run`.
4.  **Secret Scan:** manually verify no secrets or `.env` contents are being committed.
5.  **Clean:** Remove any temporary debug prints or commented-out code.
