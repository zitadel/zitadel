# ZITADEL Login App Guide for AI Agents

## Context
The **Login App** (`apps/login`) provides the user interface for authentication flows (Login, Register, MFA, etc.). It is built with Next.js and React.

## Key Technology
- **Framework**: Next.js (React).
- **Styling**: TailwindCSS, configured via `apps/login/tailwind.config.mjs`.
- **Data Fetching**: Primarily server-side interaction with ZITADEL APIs via `@zitadel/client` or direct gRPC calls where applicable.
- **Language**: TypeScript.

## Architecture & Conventions
- **Routing**: Uses the Next.js App Router (routes are defined under `src/app/`).
- **Composability**: Components should be small and reusable.
- **State**: Critical authentication state is often managed via URL parameters (Auth Requests) and cookies/sessions.
- **Scope Rule**: For shared API typings and client behavior, also read `packages/AGENTS.md` and `proto/AGENTS.md`.

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/login:dev`
- **Build**: `pnpm nx run @zitadel/login:build`
- **Lint**: `pnpm nx run @zitadel/login:lint`
- **Test (all)**: `pnpm nx run @zitadel/login:test`
- **Test (unit)**: `pnpm nx run @zitadel/login:test-unit`
- **Test (integration)**: `pnpm nx run @zitadel/login:test-integration`

## Important: Always Use Nx, Never Run Package Scripts Directly

The unit tests depend on `@zitadel/proto` and `@zitadel/client`, which are generated/built packages whose outputs (`packages/zitadel-proto/{cjs,es,types}`, `packages/zitadel-client/dist`) are **not committed to git**.

`test-unit` declares `dependsOn: ["^build"]` in `project.json`, so Nx automatically runs `@zitadel/proto:generate` and `@zitadel/client:build` before the tests when invoked correctly.

**Always run tests from the repository root via Nx:**
```bash
pnpm nx run @zitadel/login:test-unit
```

**Never** invoke the test runner directly inside the package directory (e.g. `cd apps/login && pnpm test-unit`), as this bypasses Nx's dependency graph and the generated packages will be missing.
