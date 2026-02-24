# ADR 0002: Use 3-Layer Backend Architecture (Handler, Service, Repository)

- **Status:** Accepted
- **Date:** 2026-02-24

## Context

The backend will be implemented by 3 developers working task-by-task (endpoint-based ownership), where each developer implements a feature end-to-end.

To keep code maintainable and consistent across contributors, the project needs a clear application structure with separation of concerns.

## Decision

We will use a **3-layer architecture** with the following layers:

1. **Handler layer**
   - HTTP request/response handling
   - request parsing and validation
   - calling services
   - returning API responses

2. **Service layer**
   - business logic
   - use-case orchestration
   - validation/business rules beyond transport concerns

3. **Repository layer**
   - data access
   - SQL/query execution
   - database persistence/retrieval

## Consequences

### Positive
- Clear separation of concerns
- Easier testing (service and repo can be tested independently)
- Consistent structure across APIs implemented by different developers
- Reduces coupling between HTTP and database logic

### Negative / Trade-offs
- More boilerplate compared to a simpler structure
- Requires discipline to avoid leaking logic across layers

## Rules / Guidelines

- Handlers should not contain business logic.
- Services should not depend on HTTP-specific types.
- Repositories should focus on persistence concerns only.
- Each API task should be implemented end-to-end across these layers.
