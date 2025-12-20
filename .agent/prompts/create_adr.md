Template: When creating an Architectural Decision Record (ADR):
1.  **Identify Need:** Significant architectural choice with consequences.
2.  **Next Number:** Check `specs/adr/` for the next available number (e.g., `002`).
3.  **Copy Template:** `cp specs/adr/000-template.md specs/adr/XXX-short-title.md`.
4.  **Fill Content:**
    *   **Context:** What is the problem?
    *   **Decision:** What are we doing?
    *   **Consequences:** Pros, Cons, Risks.
5.  **Status:** Set to `Accepted` or `Proposed`.
6.  **Commit:** `docs(adr): add ADR-XXX title`.
