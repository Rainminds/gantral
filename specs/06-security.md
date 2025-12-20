# Security & Privacy

## Identity & Authentication
- **OAuth 2.0 / OIDC:** Standard protocol for authentication.
- **Service Accounts:** Each Agent/Workflow Instance operates under a distinct Service Account identity.

## Audit Logs
- **Immutability:** Once written, audit logs CANNOT be modified.
- **Completeness:** Captures Input, Output, Decision, Policy Evaluation, and User Identity for every step.

## Secrets Management
- No secrets stored in plaintext.
- Integration with standard vaults (HashiCorp Vault, cloud providers).
