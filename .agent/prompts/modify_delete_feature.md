Template: When modifying or deleting a feature:

### Modification / Refactoring
1.  **Impact Analysis:** Identify all consumers (API clients, internal calls) using `grep` or IDE references.
2.  **Backward Compatibility:** Ensure changes do not break existing API contracts. If breaking, follow deprecation policy (versioning).
3.  **Update Implementation:** Modify `core/` logic.
4.  **Update Tests:** Update unit and integration tests to reflect new behavior. *Do not just delete failing tests.*
5.  **Verify:** Run full test suite to ensure no regressions.

### Deletion
1.  **Deprecation:** Verify the feature was deprecated in a previous release if public-facing.
2.  **References:** Remove all code references.
3.  **Database:** Create a precise plan for data migration or cleanup (if applicable).
4.  **Tests:** Remove associated tests.
5.  **Docs:** Remove relevant sections from `specs/` and `api/openapi.yaml`.
