Template: When implementing a feature:
1.  **Consult Specs:** Read `specs/` to understand the domain model and invariants.
2.  **Define API:** Create or update OpenAPI spec in `api/openapi.yaml`.
3.  **Define Interfaces:** Create Go interfaces in `core/` for the new functionality.
4.  **Implement Logic:** Write the core logic, ensuring it adheres to the "Instance-First" and "HITL as State" invariants.
5.  **Implement Policy:** If applicable, add policy evaluation hooks.
6.  **Add Tests:** Create unit tests in the same package and integration tests in `tests/`.
