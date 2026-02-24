# ADR 0001: Use PostgreSQL as the Primary Database

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The backend requires a relational database for storing application data with reliable querying, indexing, and transactional support.

The project includes multiple API endpoints and will require structured schema changes over time. We need a database that works well with Go tooling and supports a migration workflow.

## Decision

We will use **PostgreSQL** as the primary database for the project.

## Consequences

### Positive
- Strong relational model and SQL support
- Mature ecosystem and tooling
- Good support in Go (drivers, sqlx, migration tools)
- Reliable transactions and indexing support
- Widely used in production systems

### Negative / Trade-offs
- Requires schema and migration management discipline
- Slightly more setup than simpler embedded databases

## Notes

- Database schema changes must be tracked through migration files.
- Local development configuration must be documented in `.env.example`.
