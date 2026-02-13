# ZITADEL Monorepo Guide for AI Agents

## Mission & Context
ZITADEL is an open-source Identity Management System (IAM) written in Go and Angular/React. It provides secure login, multi-tenancy, and audit trails.

## Read Order
1. Read this file first.
2. Read the nearest scoped `AGENTS.md` for the area you edit.
3. If multiple scopes apply, use the most specific path.

## Repository Structure Map
- **`apps/`**: Consumer-facing web applications.
  - **`login`**: Next.js authentication UI. See `apps/login/AGENTS.md`.
  - **`docs`**: Fumadocs documentation app. See `apps/docs/AGENTS.md`.
  - **`api`**: Backend Nx app target. See `apps/api/AGENTS.md`.
- **`console/`**: Angular Management Console. See `console/AGENTS.md`.
- **`internal/`**: Backend domain and service logic. See `internal/AGENTS.md`.
- **`proto/`**: API definitions. See `proto/AGENTS.md`.
- **`packages/`**: Shared TypeScript packages. See `packages/AGENTS.md`.
- **`tests/functional-ui/`**: Cypress functional UI tests. See `tests/functional-ui/AGENTS.md`.

## Technology Stack & Conventions
- **Orchestration**: Nx is used for build and task orchestration.
- **Package Manager**: pnpm.
- **Backend**:
  - **Go Version Source of Truth**: Inspect `go.mod` before Go work (`go` and optional `toolchain` directives).
  - **Communication**: For V2 APIs, connectRPC is the primary transport. gRPC and HTTP/JSON endpoints are also supported.
  - **Pattern**: The backend is transitioning to a relational design. Events are still persisted in a separate table for history/audit, but events are not the system of record.
- **Frontend**:
  - **Console**: Angular + RxJS.
  - **Login/Docs**: Next.js + React.

## Command Rules
Run commands from the repository root.

- Use verified Nx targets only.
- If target availability is unclear, run `pnpm nx show project <project>`.
- Do not assume all projects have `test`, `lint`, `build`, or `generate` targets.
- Known exception: `@zitadel/console` has no configured `test` target.

## Verified Common Targets
- `@zitadel/api`: `prod`, `build`, `generate`, `lint`, `test`, `test-unit`, `test-integration`
- `@zitadel/login`: `dev`, `build`, `lint`, `test`, `test-unit`, `test-integration`
- `@zitadel/docs`: `dev`, `build`, `generate`, `check-links`, `check-types`, `test`, `lint`
- `@zitadel/console`: `dev`, `build`, `generate`, `lint`

## Documentation
- **Human Guide**: See `CONTRIBUTING.md` for setup and contribution details.
- **API Design**: See `API_DESIGN.md` for API specific guidelines.
