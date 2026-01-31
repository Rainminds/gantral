# Security & Privacy

## Identity & Authentication
- **OAuth 2.0 / OIDC:** **Federation ONLY** (Okta, Azure AD, etc.). No local users.
- **Workload Identity:** Services use AWS IAM / K8s SA / equivalent. No static API keys.

## Audit Logs
- **Immutability:** Once written, audit logs CANNOT be modified.
- **Completeness:** Captures Input, Output, Decision, Policy Evaluation, and User Identity for every step.

## Secrets Management
- No secrets stored in plaintext.
- Integration with standard vaults (HashiCorp Vault, cloud providers).

## Adversary Model (Auditability)
We assume the **Database Operator is Untrusted**.
- The operational database (Postgres) is NOT the source of truth for audit.
- Only the **Cryptographic Artifact Log** (stored separately) is trusted.
- Verification must succeed even if the active database is wiped or tampered with.
