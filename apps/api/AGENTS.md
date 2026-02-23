# ZITADEL API App Guide for AI Agents

## Context
The **API App** (`apps/api`) is the Nx application target for building and running the Go backend. Most backend implementation lives in `internal/`, while this project orchestrates build, generate, lint, and test workflows.

## Source of Truth
- **Go Toolchain**: Inspect root `go.mod` before Go work.
- **API Design Contract**: Follow `API_DESIGN.md` for service and resource conventions.
- **Domain Logic Location**: For implementation details, also read `internal/AGENTS.md`.

## Verified Nx Targets
- **Run API (prod profile)**: `pnpm nx run @zitadel/api:prod`
- **Build**: `pnpm nx run @zitadel/api:build`
- **Generate (all)**: `pnpm nx run @zitadel/api:generate`
- **Install Proto Plugins**: `pnpm nx run @zitadel/api:generate-install` — installs all Go-based proto plugins (`buf`, `protoc-gen-go`, `protoc-gen-connect-go`, `protoc-gen-openapiv2`, `protoc-gen-validate`, `protoc-gen-authoption`, etc.) to `.artifacts/bin/$(GOOS)/$(GOARCH)/`. Output is Nx-cached; only reruns when plugin versions or local protoc sources change.
- **Lint**: `pnpm nx run @zitadel/api:lint`
- **Test (all)**: `pnpm nx run @zitadel/api:test`
- **Test (unit)**: `pnpm nx run @zitadel/api:test-unit`
- **Test (integration)**: `pnpm nx run @zitadel/api:test-integration`

## Generation Notes
- `@zitadel/api:generate` can update generated, tracked files (stubs/assets/statik). Run it intentionally.
- API changes in `proto/` often require regenerating API, package, and docs artifacts.
- Proto plugins are installed to `.artifacts/bin/` — do not commit these binaries; they are declared as Nx outputs and restored from cache.
