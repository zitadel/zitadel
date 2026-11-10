# ZITADEL Proto Guide for AI Agents

## Context
`proto/` defines ZITADEL API contracts. Changes here affect generated clients, backend stubs, and docs API references.

## Source of Truth
- Follow `API_DESIGN.md` for naming, versioning, deprecations, and resource-oriented API design.
- Keep changes backward compatible within major versions unless a new major version is introduced.

## Verified Nx Targets
- **Generate TS Proto Package**: `pnpm nx run @zitadel/proto:generate`
- **Generate API Assets/Stubs**: `pnpm nx run @zitadel/api:generate`
- **Generate Docs Artifacts**: `pnpm nx run @zitadel/docs:generate`

## Workflow Notes
- After proto changes, validate dependent consumers (`@zitadel/client`, `@zitadel/api`, `@zitadel/docs`).
- If Go code is touched during generation or follow-up fixes, inspect root `go.mod` before running Go tooling.
