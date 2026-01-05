# 5. Federated Identity & Team Derivation

*   **Status**: Accepted
*   **Date**: 2026-01-05
*   **Deciders**: Gantral Core Team

## Context
Enterprise organizations already have complex Identity Providers (IdP) like Okta, Azure AD, or Active Directory. They strictly oppose "Shadow IT" user directories where admins must manually add/remove users.

Maintaining a separate user database in Gantral creates **synchronization drift** (e.g., an employee leaves but retains access) and increases security liability.

## Decision
We will implement **Strict Identity Federation**.

*   **No Local Users:** Gantral will not store passwords or manage user profiles.
*   **Token Trust:** We will trust OIDC Tokens from configured upstream IdPs.
*   **Derived Context:** Team membership and Roles will be **derived at runtime** from Token Claims (Groups, Departments) or Identity Bindings.

## Consequences
*   **Positive:**
    *   **Instant Compliance:** User access is instantly revoked when disabled in the corporate AD.
    *   **Zero Admin Overhead:** No "Add User to Team" workflows to build or manage.
    *   **SSO Native:** Works out-of-the-box with existing enterprise security.
*   **Negative:**
    *   **Dependency:** Gantral cannot run without an external IdP (though a mock IdP can be provided for local dev).
