# ADR 0005: Use Goose for Database Migrations

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The project requires a migration tool to manage database schema changes in a repeatable and team-friendly way.

Since multiple developers may introduce schema changes, migrations must be versioned and committed as part of the codebase.

## Decision

We will use **Goose** (`pressly/goose`) as the database migration tool.

## Consequences

### Positive
- Supports versioned SQL migrations
- Works well with PostgreSQL
- Team-friendly workflow for schema evolution
- Easy to integrate into local development / CI scripts

### Negative / Trade-offs
- Requires migration naming/version discipline
- Merge conflicts can happen when multiple migrations are created in parallel

## Team Rules for Migrations

- Migration files are a shared point and must be handled carefully.
- Discuss schema-impacting changes before merging.
- Do not edit already-merged migration files unless absolutely necessary.
- Prefer creating a new migration for follow-up schema changes.
