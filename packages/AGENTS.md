# ZITADEL Packages Guide for AI Agents

## Context
`packages/` contains shared TypeScript libraries used by frontend applications and external consumers.

## Main Packages
- **`packages/zitadel-proto`** (`@zitadel/proto`): generated protobuf TypeScript artifacts.
- **`packages/zitadel-client`** (`@zitadel/client`): higher-level client library built on generated proto/connect types.
- **`packages/zitadel-js`** (`@zitadel/zitadel-js`): isomorphic core SDK — framework-agnostic primitives for OIDC, session management, JWT/JWE handling, webhook verification, and ConnectRPC transport creation. Generates its own protobuf types from `proto/` using local `protoc-gen-es`.
- **`packages/zitadel-react`** (`@zitadel/react`): React hooks and context for client-side state management. Depends on `@zitadel/zitadel-js`.
- **`packages/zitadel-nextjs`** (`@zitadel/nextjs`): Next.js App Router integration — OIDC lifecycle (via `oauth4webapi`), middleware, server actions, v2 API access, and Actions v2 webhook handling. Depends on `@zitadel/zitadel-js` and `@zitadel/react`.
  - **`auth/oidc`** — OIDC redirect-based login ("add login to your app"). Env: `ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`, `ZITADEL_COOKIE_SECRET`.
  - **`auth/session`** — Session API for custom login UIs *(planned, not yet implemented)*.
  - **`middleware`** — Route protection middleware. Reads `ZITADEL_COOKIE_SECRET`.
  - **`api`** — ZITADEL v2 API client factory. Reads `ZITADEL_API_URL`.
  - **`webhook`** — Actions v2 webhook handler. Reads `ZITADEL_WEBHOOK_SECRET`, `ZITADEL_WEBHOOK_JWKS_ENDPOINT`, `ZITADEL_WEBHOOK_JWE_PRIVATE_KEY`.
  - **`server-action`** — Protected server action wrapper.
- **`packages/zitadel-angular`** (`@zitadel/angular`): Angular SDK — injectable services, guards, interceptors. Depends on `@zitadel/zitadel-js`.

## Verified Nx Targets
- **Proto generation**: `pnpm nx run @zitadel/proto:generate`
- **Client build**: `pnpm nx run @zitadel/client:build`
- **Client lint**: `pnpm nx run @zitadel/client:lint`
- **Client tests**: `pnpm nx run @zitadel/client:test`
- **JS SDK generate**: `pnpm nx run @zitadel/zitadel-js:generate`
- **JS SDK build**: `pnpm nx run @zitadel/zitadel-js:build`
- **JS SDK tests**: `pnpm nx run @zitadel/zitadel-js:test`
- **Next.js SDK build**: `pnpm nx run @zitadel/nextjs:build`
- **Next.js SDK tests**: `pnpm nx run @zitadel/nextjs:test`

## Package Dependency Graph

```
@zitadel/proto          (standalone — generated from proto/)
    ↑
@zitadel/client         (depends on @zitadel/proto)

@zitadel/zitadel-js     (standalone — generates own protos, depends on @connectrpc + jose)
    ↑
@zitadel/react          (depends on @zitadel/zitadel-js)
    ↑
@zitadel/nextjs         (depends on @zitadel/react + @zitadel/zitadel-js + oauth4webapi)

@zitadel/angular        (depends on @zitadel/zitadel-js)
```

## Workflow Notes
- When changing `proto/`, regenerate `@zitadel/proto` **and** `@zitadel/zitadel-js` (both have `generate` targets that read from `proto/`).
- `@zitadel/zitadel-js` uses local `protoc-gen-es` (same as `@zitadel/proto`) — no BSR remote plugins.
- Keep package exports and public typings stable unless a breaking release is explicitly intended.
- `@zitadel/nextjs` requires `next >=15` and `react >=18` as peer dependencies.

## Environment Variables
A shared `.env.example` lives at `packages/.env.example` with three sections:
1. **OIDC** — `ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`, `ZITADEL_COOKIE_SECRET`
2. **API** — `ZITADEL_API_URL`, `ZITADEL_SERVICE_USER_TOKEN` (or private key JWT vars)
3. **Actions** — `ZITADEL_WEBHOOK_SECRET`, `ZITADEL_WEBHOOK_JWKS_ENDPOINT`, `ZITADEL_WEBHOOK_JWE_PRIVATE_KEY`

All SDK modules auto-read these env vars as fallbacks when options are not passed explicitly.

## Module Naming Convention
Auth modules live under `auth/` with submodules per integration pattern:
- `auth/oidc` — Redirect-based OIDC login
- `auth/session` — Session API for custom login UIs *(planned)*

Check types (password, passkey, TOTP, etc.) are **parameters** within `auth/session`, not separate modules.

