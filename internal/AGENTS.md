# ZITADEL Internal Backend Guide for AI Agents

## Context
`internal/` contains core backend domain logic for ZITADEL: commands, queries, repositories, eventstore integration, API service layers, and supporting infrastructure.

## Source of Truth
- **Go Toolchain**: Inspect root `go.mod` before Go work.
- **Architecture Pattern**: Relational data is the system of record; keep existing event writes that provide history/audit trails.
- **API Contract**: For API-facing schema decisions, follow `API_DESIGN.md` and `proto/AGENTS.md`.

## Boundary Rules
- Prefer implementing business behavior in command/query layers and repository packages, not in transport handlers.
- Avoid bypassing established event/repository flows with ad-hoc direct persistence patterns.
- Keep API/service adapters thin; place reusable domain behavior in internal domain packages.

## Validation Workflow
- Use API project targets to validate backend changes:
  - `pnpm nx run @zitadel/api:lint`
  - `pnpm nx run @zitadel/api:test-unit`
  - `pnpm nx run @zitadel/api:test-integration`

## Identity Signals Subsystem (Preview)

The `signals/` package provides identity-aware observability. See `signals/DESIGN.md` for the full architecture.

### Key Files
- `config.go` — Configuration structs (`IdentitySignalsConfig`)
- `emitter.go` — Fire-and-forget signal emission with channel buffering
- `ducklake_store.go` — DuckLake storage backend (requires CGO)
- `signal_interceptor.go` — HTTP/connectRPC middleware for request signals
- `event_hook.go` — Eventstore hook for event signals

### Boundary Rules
- Signal emission must be non-blocking (fire-and-forget through buffered channel)
- DuckDB code must be behind `//go:build cgo` tags with nocgo stubs
- All queries must be scoped by `instance_id` (tenant isolation)
- The signal API self-excludes to prevent recording loops
- ID extraction handles multiple JSON field name variants across aggregate types
