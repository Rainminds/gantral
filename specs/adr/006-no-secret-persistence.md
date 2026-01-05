# 6. No Secret Persistence (Reference Only)

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
AI Agents require access to high-value credentials (API Keys, Database Passwords, SSH Keys). Storing these secrets in the Gantral Control Plane creates a **single point of catastrophic failure**. If the Gantral DB is compromised, every connected system is compromised.

## Decision
We will adopt a **"Zero-Knowledge" Secret Architecture**.

*   **Reference Only:** Gantral will store only **URIs/Pointers** to secrets (e.g., `vault://prod/db-password`), never the secret material itself.
*   **Edge Resolution:** Secrets are resolved **Just-In-Time** by the Runner, using the Runner's own identity to authenticate with the Secret Manager (Vault, AWS Secrets Manager).

## Consequences
*   **Positive:**
    *   **Reduced Blast Radius:** A compromise of the Control Plane reveals only where secrets *are*, not what they *are*.
    *   **Rotation Safety:** Secrets can be rotated in Vault without updating Gantral configurations.
*   **Negative:**
    *   **Setup Complexity:** Runners must be configured with identity/access to the Secret Manager.
