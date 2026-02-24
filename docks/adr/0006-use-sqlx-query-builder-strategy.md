# ADR 0006: Use SQLX for Database Access (with SQL/Query Builder Strategy)

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The project needs a database access approach that is lightweight, works well with PostgreSQL, and avoids the complexity of a full ORM for the current scope.

The team prefers explicit SQL and clear control over queries while still having ergonomic row scanning and mapping support in Go.

## Decision

We will use **SQLX** for database access in the repository layer.

SQLX will be the primary DB access library for:
- executing queries
- scanning rows into structs
- handling query helpers where needed

> Query-building style (raw SQL vs helper-based query builder) should remain consistent across the repository layer.

## Consequences

### Positive
- Lightweight compared to full ORM solutions
- Keeps SQL explicit and easier to reason about
- Good ergonomics for scanning and struct mapping
- Works well with PostgreSQL and repository pattern

### Negative / Trade-offs
- More manual query writing than ORM-based solutions
- Team must maintain query consistency and avoid duplicated SQL patterns

## Rules / Guidelines

- SQLX usage should be limited to repository layer.
- Business logic should not include SQL.
- Query style should be consistent across repositories.
- Shared query patterns should be documented if reused.
