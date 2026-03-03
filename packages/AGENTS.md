# ZITADEL Packages Guide for AI Agents

## Context
`packages/` contains shared TypeScript libraries used by frontend applications and external consumers.

## Main Packages
- **`packages/zitadel-proto`** (`@zitadel/proto`): generated protobuf TypeScript artifacts.
- **`packages/zitadel-client`** (`@zitadel/client`): higher-level client library built on generated proto/connect types.
- **`packages/zitadel-js`** (`@zitadel/zitadel-js`): isomorphic core SDK — framework-agnostic primitives for OIDC (wrapped via `oauth4webapi`), session management, JWT/JWE handling, Actions webhook verification, and ConnectRPC transport creation. Generates its own protobuf types from `proto/` using local `protoc-gen-es`.
- **`packages/zitadel-react`** (`@zitadel/react`): React hooks and context for client-side state management. Depends on `@zitadel/zitadel-js`.
- **`packages/zitadel-nextjs`** (`@zitadel/nextjs`): Next.js App Router integration — OIDC lifecycle, middleware, server actions, v2 API access, and Actions v2 webhook handling. Depends on `@zitadel/zitadel-js` and `@zitadel/react`.
  - **`auth/oidc`** — OIDC redirect-based login ("add login to your app"). Env: `ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`, `ZITADEL_COOKIE_SECRET`.
  - **`auth/session`** — Session API helper layer for custom login UIs.
  - **`auth/bearer-token`** — Canonical bearer-token helper lane.
  - **`middleware`** — Route protection middleware. Reads `ZITADEL_COOKIE_SECRET`.
  - **`api`** — ZITADEL v2 API client factory. Reads `ZITADEL_API_URL`.
  - **`actions/webhook`** — Canonical Actions v2 webhook handler. Reads `ZITADEL_WEBHOOK_SECRET`, `ZITADEL_WEBHOOK_JWKS_ENDPOINT`, `ZITADEL_WEBHOOK_JWE_PRIVATE_KEY`.
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

@zitadel/zitadel-js     (standalone — generates own protos, depends on @connectrpc + jose + oauth4webapi)
    ↑
@zitadel/react          (depends on @zitadel/zitadel-js)
    ↑
@zitadel/nextjs         (depends on @zitadel/react + @zitadel/zitadel-js)

@zitadel/angular        (depends on @zitadel/zitadel-js)
```

## SDK Package Consolidation Plan (maintainer-only, planning)
- **Target public entrypoint**: `@zitadel/zitadel-js` is the primary long-term JS SDK entrypoint.
- **Current iteration scope**: planning and documentation only — **do not remove or replace** `@zitadel/client` or `@zitadel/proto` in this iteration.
- **Current repo state remains valid**: keep publishing/building `@zitadel/client` and `@zitadel/proto` while migration work is staged.

### Planned migration phases
1. **Phase 0 (now)**: document intent and guardrails; no package removals/replacements.
2. **Phase 1 (adoption)**: prefer `@zitadel/zitadel-js` in new SDK examples/surfaces while preserving `@zitadel/client` and `@zitadel/proto` behavior.
3. **Phase 2 (deprecation)**: add staged deprecation notices for `@zitadel/client` and `@zitadel/proto` with migration guidance to `@zitadel/zitadel-js`.
4. **Phase 3 (post-deprecation, future major)**: evaluate removal only after explicit release planning, migration coverage, and semver-major coordination.

### Guardrails for all phases
- Keep canonical taxonomy decisions (`auth/*`, `api/*`, `actions/*`) unchanged during consolidation.
- Do not introduce breaking changes to `@zitadel/client` or `@zitadel/proto` before an announced deprecation + major-release path.
- Keep Nx targets and generation paths for all three packages (`@zitadel/zitadel-js`, `@zitadel/client`, `@zitadel/proto`) operational until removal is explicitly approved.
- Every phase must include changelog + migration-note updates before changing package lifecycle state.
- Keep `@zitadel/zitadel-js` API discoverability explicit: `api/v1` and `api/v2` remain canonical API lanes.

## Workflow Notes
- When changing `proto/`, regenerate `@zitadel/proto` **and** `@zitadel/zitadel-js` (both have `generate` targets that read from `proto/`).
- `@zitadel/zitadel-js` uses local `protoc-gen-es` (same as `@zitadel/proto`) — no BSR remote plugins.
- Keep package exports and public typings stable unless a breaking release is explicitly intended.
- `@zitadel/nextjs` requires `next >=15` and `react >=18` as peer dependencies.
- In-repo runnable SDK example scope is currently Next.js-only (`examples/nextjs` and `tests/nextjs-sdk-playground`).

## Environment Variables
A shared `.env.example` lives at `packages/.env.example` with three sections:
1. **OIDC** — `ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`, `ZITADEL_COOKIE_SECRET`
2. **API** — `ZITADEL_API_URL`, `ZITADEL_SERVICE_USER_TOKEN` (or private key JWT vars)
3. **Actions** — `ZITADEL_WEBHOOK_SECRET`, `ZITADEL_WEBHOOK_JWKS_ENDPOINT`, `ZITADEL_WEBHOOK_JWE_PRIVATE_KEY`

All SDK modules auto-read these env vars as fallbacks when options are not passed explicitly.

## Package Ownership + Runtime Boundary Matrix

Use this matrix when deciding where a feature belongs and which runtime should execute it.

| Package | Primary role | Browser runtime | Server runtime |
| --- | --- | --- | --- |
| `@zitadel/zitadel-js` | Framework-agnostic core primitives (`auth/*`, `api/*`, `actions/*`, shared transport/helpers). | Yes (OIDC discovery/auth URL helpers, PKCE/state generation, token exchange without client secret). | Yes (`/node` token/JWT helpers, webhook verification, API transport factories). |
| `@zitadel/react` | React composition layer (context/hooks/components) over SDK primitives. | Yes (UI state and rendering helpers). | No direct server-only abstractions. |
| `@zitadel/nextjs` | Next.js App Router adapter (`auth/oidc`, `auth/session`, `auth/bearer-token`, `api`, `actions/webhook`, middleware). | Minimal browser surface (route entry points). | Yes (cookie/session handling, callback/token exchange, API clients, actions/webhook routes). |
| `@zitadel/angular` | Angular adapter surface (provider/guard/interceptor patterns). | Target runtime is browser app integration. | Server helpers are limited; server-boundary handling should remain in dedicated BFF/server routes. |

### SPA server-boundary defaults

For SPA integrations (React/Angular/Vue/etc.), use these default ownership rules:

| Capability | Default owner |
| --- | --- |
| UI state, route guards, login button rendering, user-facing navigation | Browser app (`@zitadel/react` / framework app code). |
| Authorization callback endpoint, code exchange side effects, session cookie writes | Server/BFF route. |
| Service-user credentials, private key JWT generation, introspection credentials | Server/BFF only. |
| Webhook signature/JWT/JWE verification | Server-only (`@zitadel/zitadel-js/actions/webhook` or framework server adapters). |
| Access to management/system/event APIs with confidential credentials | Server-only (`auth/bearer-token` lane). |
| Calls to userinfo/resource endpoints using browser session state | Prefer server proxy in SPA+BFF architecture. |

Boundary rule for SPAs: if code needs confidential material (secret/private key/system token) or writes trusted session state, it belongs on the server/BFF side.

## Canonical Cross-SDK Capability Model
Use these capability lanes as the shared model for all SDKs (JS/Next.js today, Go and others as they are added).
Each public SDK surface should map to exactly one primary lane.
Canonical JS SDK taxonomy groups are `auth/*`, `api/*`, and `actions/*`, plus root/core low-level primitives.
Canonical bearer-token lane ID is `auth/bearer-token`.

| Lane | Purpose | In scope | Out of scope | Current examples |
| --- | --- | --- | --- | --- |
| `auth/session` | User-facing authentication building blocks for custom login UIs. | Session lifecycle, login step orchestration, callback completion, claim/session reads. | Protocol-specific redirect/discovery logic and IdP administration. | `@zitadel/nextjs/auth/session`, shared session/JWT helpers in `@zitadel/zitadel-js`. |
| `auth/<protocol>` | Standards-based end-user auth protocol adapters between apps and ZITADEL. | OIDC/SAML auth URL handling, callback validation, code/token exchange, logout URL construction. | User-facing check UX orchestration and federation resource management. | `@zitadel/nextjs/auth/oidc`, OIDC wrappers in `@zitadel/zitadel-js`. |
| `api/idp/<protocol>` | Upstream identity provider federation management under the API lane. | CRUD for upstream IdPs and mapping/policy configuration through admin APIs. | Runtime end-user session/auth protocol execution. | Current usage via admin API clients; future modules should stay under `api/*` and isolated from `auth/*`. |
| `auth/bearer-token` | Canonical server-side bearer-token lane for confidential API credentials and validation helpers. | Service-user/private-key credential handling, bearer token acquisition/rotation, token validation helpers for route handlers. | Browser execution and interactive end-user sign-in flows. | `@zitadel/nextjs/auth/bearer-token`, `@zitadel/zitadel-js/auth/bearer-token`. |
| `api/*` | API-family client lane for management/system/event resources. | API family clients (for example management/system/event) that rely on confidential bearer-token credentials. | Interactive end-user sign-in flows and webhook verification handlers. | `@zitadel/nextjs/api`, `@zitadel/zitadel-js/api/v2`. |
| `actions/*` | Inbound event/webhook verification lane. | Actions webhook validation/decryption and typed event payload handling. | Bearer-token API client ownership and interactive end-user sign-in flows. | `@zitadel/nextjs/actions/webhook`, `@zitadel/zitadel-js/actions/webhook`. |

Boundary rule: when a feature spans multiple lanes, compose modules across lanes instead of introducing a mixed abstraction.

### Maintainer lane triage (new features)
1. User-facing check/session orchestration for custom login UIs belongs to **`auth/session`**.
2. OIDC/SAML redirect/callback/token/logout protocol work belongs to **`auth/<protocol>`**.
3. Upstream IdP CRUD/mapping/policy management belongs to **`api/idp/<protocol>`**, even if currently reached via broader `api/*` surfaces.
4. Confidential admin/system/event API credentials and bearer-token helpers belong to **`auth/bearer-token`**.
5. Inbound event/webhook verification belongs to **`actions/*`** (`actions/webhook` canonical).
6. If one feature needs multiple lanes, split it into lane-specific modules and compose.

## Module Naming Convention
Canonical cross-SDK module IDs are path-like and protocol-explicit:
- `auth/session` — session/check orchestration for embedded auth UIs
- `auth/oidc` — OIDC end-user auth protocol adapter
- `auth/saml` — SAML end-user auth protocol adapter
- `auth/bearer-token` — canonical server-side bearer-token lane for confidential API credentials/helpers
- `api/idp/oidc` — OIDC upstream IdP federation management under the API lane
- `api/idp/saml` — SAML upstream IdP federation management under the API lane
- `api/*` — API-family module names (for example `api/management`, `api/system`) that use `auth/bearer-token` credentials
- `actions/webhook` — canonical inbound Actions webhook verification module
- `actions/*` — inbound event/webhook verification grouped by source

Naming contract:
1. Each public SDK surface maps to exactly one canonical ID.
2. Keep the root segment in `{auth,api,actions}` (plus root/core low-level primitives) and add segments for protocol/family when needed.
3. Check types (password, passkey, TOTP, etc.) stay as parameters within `auth/session`, not separate modules.
4. Start from the canonical ID first, then apply language-specific separators/casing without changing segment order or meaning.

### JS/Next.js ↔ Go Mapping Matrix
Use this as the maintainer reference for cross-SDK naming and rollout order.

| Canonical ID | JS/Next.js current | JS/Next.js planned | Go current | Go planned |
| --- | --- | --- | --- | --- |
| `auth/session` | `@zitadel/nextjs/auth/session`; shared session primitives in `@zitadel/zitadel-js`. | Keep `auth/session` as the primary auth UI lane and align aliases/docs to canonical naming. | No dedicated public Go SDK surface yet. | `auth/session` package as the embedded auth baseline. |
| `auth/oidc` | `@zitadel/nextjs/auth/oidc`; OIDC helpers in `@zitadel/zitadel-js`. | Keep `auth/oidc` as the protocol lane and preserve canonical naming in all wrappers. | No dedicated public Go SDK surface yet. | `auth/oidc` package for redirect-based OIDC flows. |
| `auth/saml` | Not exposed yet. | Add `auth/saml` as a peer lane to `auth/oidc` (same lane boundaries). | No dedicated public Go SDK surface yet. | `auth/saml` package for SAML end-user auth flows. |
| `api/idp/oidc` | No dedicated module; currently handled via management APIs (`@zitadel/nextjs/api`, `@zitadel/zitadel-js/api/v2`). | Add `api/idp/oidc` wrappers over management APIs without mixing with `auth/*`. | No dedicated public Go SDK surface yet. | `api/idp/oidc` package for upstream OIDC federation management. |
| `api/idp/saml` | No dedicated module; currently handled via management APIs (`@zitadel/nextjs/api`, `@zitadel/zitadel-js/api/v2`). | Add `api/idp/saml` wrappers over management APIs without mixing with `auth/*`. | No dedicated public Go SDK surface yet. | `api/idp/saml` package for upstream SAML federation management. |
| `auth/bearer-token` | `@zitadel/nextjs/auth/bearer-token`, `@zitadel/zitadel-js/auth/bearer-token`. | Keep explicit `auth/bearer-token` ownership across SDK surfaces. | No dedicated public Go SDK surface yet. | `auth/bearer-token` packages with family-level `api/*` modules (for example `api/management`, `api/system`). |
| `actions/webhook` | `@zitadel/nextjs/actions/webhook`, `@zitadel/zitadel-js/actions/webhook`. | Keep webhook handlers under canonical `actions/*` grouping. | No dedicated public Go SDK surface yet. | `actions/webhook` packages grouped by event source. |

Parity notes:
- Target parity is canonical-name parity first: each lane above should exist as one primary JS surface and one primary Go surface.
- Temporary gaps are allowed while Go SDK surfaces are introduced; JS may continue to use `api/*` surfaces as interim paths for `api/idp/*` and `auth/saml`, but confidential bearer-token ownership is `auth/bearer-token`.
- Temporary gaps/aliases must stay documented in this matrix with owning SDK + target milestone, and aliases must be deprecated until removed.

Do:
- Keep canonical roots and segment order aligned across SDKs.
- Extend existing roots (`auth/bearer-token`, family-level `api/*` modules, and `actions/*`) instead of inventing new top-level groups.
- If aliases are required for migration, keep them temporary and deprecate with a timeline.

Don't:
- Rename canonical surfaces per language when an equivalent canonical name exists (for example `login`, `federation`, `hooks`).
- Mix end-user auth execution (`auth/*`) and upstream federation management (`api/idp/*`) in one surface.
- Encode transport/framework details into canonical module names.

Language compatibility notes:
- If a language cannot represent `/` directly, use a deterministic transform that preserves segment order and meaning.
- Allowed adaptations: `/` → `.`, `::`, `_`, or nested package/module directories; segment casing may follow language conventions.
- Always document the canonical ID alongside the adapted symbol/path (for example `auth/session` → `auth.session`, `Auth::Session`, `auth_session`).

## Maintainer SDK Governance Checklist (pre-ship)
- [ ] **Lane placement confirmed**: module/feature maps to exactly one lane from the canonical capability model.
- [ ] **Canonical naming confirmed**: canonical ID is used first; language-specific adaptation keeps segment order and meaning.
- [ ] **JS/Go mapping matrix updated**: update `JS/Next.js ↔ Go Mapping Matrix` current/planned status for the affected canonical ID.
- [ ] **Public docs updated**: update user-facing docs/examples/changelog for new or changed SDK surface before release.
- [ ] **Temporary alias/gap tracked**: document temporary aliases or parity gaps with owner SDK and target milestone (plus deprecation/removal plan for aliases).
- [ ] **Validation completed**: run lane-relevant tests, build/generate targets, and a minimal smoke flow for the new/changed surface.
