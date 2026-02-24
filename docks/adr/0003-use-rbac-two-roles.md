# ADR 0003: Use RBAC Authorization with Two Roles

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The system requires authorization controls to restrict access to certain actions and resources based on user privileges.

At the current stage of the project, authorization requirements are simple and can be represented with a small set of roles.

## Decision

We will use **RBAC (Role-Based Access Control)** with **2 roles**.

> Exact role names and permissions will be defined in application-level authorization rules (e.g., `admin`, `user`).

## Consequences

### Positive
- Simple and easy to implement
- Easy to understand and test
- Suitable for current project scope
- Can be expanded later if needed

### Negative / Trade-offs
- Less flexible than fine-grained permission systems
- May require refactoring if role complexity grows significantly

## Notes

- Role information should be available in authenticated request context.
- Protected endpoints should enforce role checks in a consistent way.
- Role-related claims should be included in token claims where applicable.
