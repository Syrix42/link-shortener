# ADR 0004: Use JWT Access and Refresh Tokens for Authentication

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The backend requires stateless authentication for APIs, while also supporting session continuity without forcing frequent re-login.

To balance security and usability, short-lived access tokens and refresh tokens are needed.

## Decision

We will use:

- **Access Token (JWT)** for API authentication
- **Refresh Token** for obtaining new access tokens

### Token Policy
- **Access token lifetime:** 15 minutes
- Refresh token lifetime: (to be defined by implementation / security requirements)

### Claims
JWTs must include correct and minimal claims required by the system, including (as applicable):
- subject/user identifier (`sub`)
- role / authorization-related claim(s)
- issued at (`iat`)
- expiration (`exp`)
- token type (if used, e.g., access/refresh distinction)

## Consequences

### Positive
- Short-lived access tokens reduce risk exposure
- Refresh flow improves user experience
- Stateless access-token verification works well for APIs
- Role claims support RBAC checks

### Negative / Trade-offs
- Token lifecycle handling adds implementation complexity
- Refresh token security/storage/revocation strategy must be handled carefully

## Notes / Implementation Guidance

- Access tokens should be validated on protected endpoints.
- Refresh endpoint must validate refresh token before issuing new access token.
- Claims must be consistent with RBAC design.
- Token signing method and secret/key management must be defined in implementation.
