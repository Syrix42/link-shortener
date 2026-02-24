# Architecture Decision Records (ADR)

This directory contains architecture and technical decisions for the backend project.

## ADR List

- [ADR 0001 - Use PostgreSQL as the Primary Database](./0001-use-postgresql.md)
- [ADR 0002 - Use 3-Layer Backend Architecture (Handler, Service, Repository)](./0002-use-3-layer-architecture.md)
- [ADR 0003 - Use RBAC Authorization with Two Roles](./0003-use-rbac-two-roles.md)
- [ADR 0004 - Use JWT Access and Refresh Tokens for Authentication](./0004-use-jwt-access-refresh-tokens.md)
- [ADR 0005 - Use Goose for Database Migrations](./0005-use-goose-for-migrations.md)
- [ADR 0006 - Use SQLX for Database Access](./0006-use-sqlx-query-builder-strategy.md)

## Notes

- Update ADR status if a decision changes.
- Do not edit history silently; create a new ADR if replacing a major decision.
