# Integration tests

This directory contains **integration tests only** (end-to-end tests that exercise the HTTP API).

## What belongs here
- Tests that send HTTP requests to the API (router + middleware + handlers)
- Tests that verify request/response behavior (status codes, JSON, auth, etc.)
- Tests that may use real infrastructure in test mode (e.g., a test DB)

## What does NOT belong here
- Unit tests for services, validators, utilities
- Pure business-logic tests (those live next to the code under `internal/...`)

## Running
Run all tests:
```bash
go test ./...

