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

## Server Actions & Server-Only Code (critical — read before adding or changing a server function)
- **`src/lib/actions/`**: only Server Actions. Every file starts with `"use server"`, and **every export is a public, unauthenticated POST endpoint**, whether or not a component calls it. Treat every argument as attacker-controlled: derive the user, session and policy (lifetimes, link hosts, IdP identities) on the server from the caller's signed cookie and ZITADEL, never from the arguments.
- **`src/lib/server/`**: server-only code. Every file imports `"server-only"`. Helpers that take identity, links or lifetimes as parameters belong here.
- **`src/lib/` (top level)**: code shared with client components only.
- `"use server"` is not allowed anywhere else, inline directives included. The `login/module-boundaries` ESLint rule (`checks/eslint-module-boundaries.mjs`) enforces the layout.
- The build fails when the published actions differ from `checks/server-actions.allowlist.json` (`checks/server-actions.mjs`). Adding an action means reviewing it and adding it there.

## Base Path & Redirects (critical — read before touching any redirect logic)
- **The login app is served under a non-root base path in production.** ZITADEL Cloud and the official container image are built with `NEXT_PUBLIC_BASE_PATH=/ui/v2/login` (see `apps/login/.env`; the value is inlined at build time by Next.js). Never assume the app lives at the domain root.
- **`constructUrl()` (`src/lib/server/service-url.ts`) prepends the base path and must ONLY ever receive root-relative paths** (`/loginname`, `/password?...`). Passing an absolute URL glues the base path onto it and produces a broken redirect like `https://<host>/ui/v2/loginhttps://login.microsoftonline.com/...`.
- **Some redirect targets ARE absolute URLs**: IdP authorize URLs returned by `startIdentityProviderFlow`, OIDC callback URLs from `createCallback`, and SAML endpoints from `createResponse`. When handling a redirect value that could be absolute, branch on `isExternalUrl()` first, validate absolute URLs with `isSafeRedirectUri()`, and pass them to `NextResponse.redirect` (server) / `window.location.href` (client) untouched. Canonical pattern: `resolveLoginHint` in `src/lib/server/flow-initiation.ts` and `handleServerActionResponse` in `src/lib/client-utils.ts`.
- **An empty base path masks this whole bug class**: with `NEXT_PUBLIC_BASE_PATH` unset, `new URL("" + absoluteUrl, origin)` still parses the absolute URL correctly, so local setups without the base path won't reproduce base-path bugs. Always test redirect changes with `NEXT_PUBLIC_BASE_PATH=/ui/v2/login` set, and verify the raw `Location` header (e.g. `curl -I`), not just browser behavior.

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/login:dev`
- **Build**: `pnpm nx run @zitadel/login:build`
- **Lint**: `pnpm nx run @zitadel/login:lint`
- **Test (all)**: `pnpm nx run @zitadel/login:test`
- **Test (unit)**: `pnpm nx run @zitadel/login:test-unit`
- **Test (integration)**: `pnpm nx run @zitadel/login:test-integration`
- **Pack (Docker)**: `pnpm nx run @zitadel/login:pack` — builds a local Docker image `zitadel/zitadel-login:local`. Requires Docker daemon.
