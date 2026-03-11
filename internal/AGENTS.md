# ZITADEL Internal Backend Guide for AI Agents

## Context
`internal/` contains core backend domain logic for ZITADEL: commands, queries, repositories, eventstore integration, API service layers, and supporting infrastructure.

## Source of Truth
- **Go Toolchain**: Inspect root `go.mod` before Go work.
- **Architecture Pattern**: Relational data is the system of record; keep existing event writes that provide history/audit trails.
- **API Contract**: For API-facing schema decisions, follow `API_DESIGN.md` and `proto/AGENTS.md`.

## Behavioral Signal Package Map
- `internal/signals`: behavioral data plane; captures, stores, archives, and replays signals.
- `internal/detection`: decision plane; builds context from signal history, evaluates rules/classifiers, and produces findings plus allow/block/challenge decisions.
- Rules live in `internal/detection` because they describe how ZITADEL detects suspicious behavior from signals, not how signals are captured or stored.
- `internal/ratelimit`, `internal/captcha`, and `internal/llm` are supporting packages used by detection rules and decisions.

## Boundary Rules
- Prefer implementing business behavior in command/query layers and repository packages, not in transport handlers.
- Avoid bypassing established event/repository flows with ad-hoc direct persistence patterns.
- Keep API/service adapters thin; place reusable domain behavior in internal domain packages.

## Validation Workflow
- Use API project targets to validate backend changes:
  - `pnpm nx run @zitadel/api:lint`
  - `pnpm nx run @zitadel/api:test-unit`
  - `pnpm nx run @zitadel/api:test-integration`
