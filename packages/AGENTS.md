# ZITADEL Packages Guide for AI Agents

## Context
`packages/` contains shared TypeScript libraries used by frontend applications and external consumers.

## Main Packages
- **`packages/zitadel-proto`** (`@zitadel/proto`): generated protobuf TypeScript artifacts.
- **`packages/zitadel-client`** (`@zitadel/client`): higher-level client library built on generated proto/connect types.
- **`packages/zitadel-js`** (`@zitadel/zitadel-js`): isomorphic core SDK — framework-agnostic primitives for OIDC (wrapped via `oauth4webapi`), session management, JWT/JWE handling, webhook verification, and ConnectRPC transport creation. Generates its own protobuf types from `proto/` using local `protoc-gen-es`.
- **`packages/zitadel-react`** (`@zitadel/react`): React hooks and context for client-side state management. Depends on `@zitadel/zitadel-js`.
- **`packages/zitadel-nextjs`** (`@zitadel/nextjs`): Next.js App Router integration — OIDC lifecycle, middleware, server actions, v2 API access, and Actions v2 webhook handling. Depends on `@zitadel/zitadel-js` and `@zitadel/react`.
  - **`auth/oidc`** — OIDC redirect-based login ("add login to your app"). Env: `ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`, `ZITADEL_COOKIE_SECRET`.
  - **`auth/session`** — Session API helper layer for custom login UIs.
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

@zitadel/zitadel-js     (standalone — generates own protos, depends on @connectrpc + jose + oauth4webapi)
    ↑
@zitadel/react          (depends on @zitadel/zitadel-js)
    ↑
@zitadel/nextjs         (depends on @zitadel/react + @zitadel/zitadel-js)

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

## Package Ownership + Runtime Boundary Matrix

Use this matrix when deciding where a feature belongs and which runtime should execute it.

| Package | Primary role | Browser runtime | Server runtime |
| --- | --- | --- | --- |
| `@zitadel/zitadel-js` | Framework-agnostic core primitives (`oidc`, `session`, transport, shared helpers). | Yes (OIDC discovery/auth URL helpers, PKCE/state generation, token exchange without client secret). | Yes (`/node` token/JWT helpers, webhook verification, API transport factories). |
| `@zitadel/react` | React composition layer (context/hooks/components) over SDK primitives. | Yes (UI state and rendering helpers). | No direct server-only abstractions. |
| `@zitadel/nextjs` | Next.js App Router adapter (`auth/oidc`, `auth/session`, `api`, `webhook`, middleware). | Minimal browser surface (route entry points). | Yes (cookie/session handling, callback/token exchange, API clients, webhook routes). |
| `@zitadel/angular` | Angular adapter surface (provider/guard/interceptor patterns). | Target runtime is browser app integration. | Server helpers are limited; server-boundary handling should remain in dedicated BFF/server routes. |

### SPA server-boundary defaults

For SPA integrations (React/Angular/Vue/etc.), use these default ownership rules:

| Capability | Default owner |
| --- | --- |
| UI state, route guards, login button rendering, user-facing navigation | Browser app (`@zitadel/react` / framework app code). |
| Authorization callback endpoint, code exchange side effects, session cookie writes | Server/BFF route. |
| Service-user credentials, private key JWT generation, introspection credentials | Server/BFF only. |
| Webhook signature/JWT/JWE verification | Server-only (`@zitadel/zitadel-js/webhooks` or framework server adapters). |
| Access to management/system/event APIs with confidential credentials | Server-only (`api/*` lanes). |
| Calls to userinfo/resource endpoints using browser session state | Prefer server proxy in SPA+BFF architecture. |

Boundary rule for SPAs: if code needs confidential material (secret/private key/system token) or writes trusted session state, it belongs on the server/BFF side.

## Canonical Cross-SDK Capability Model
Use these capability lanes as the shared model for all SDKs (JS/Next.js today, Go and others as they are added).
Each public SDK surface should map to exactly one primary lane.

| Lane | Purpose | In scope | Out of scope | Current examples |
| --- | --- | --- | --- | --- |
| Embedded auth concepts | User-facing authentication building blocks for custom login UIs. | Session lifecycle, login step orchestration, callback completion, claim/session reads. | Protocol-specific redirect/discovery logic and IdP administration. | `@zitadel/nextjs/auth/session`, shared session/JWT helpers in `@zitadel/zitadel-js`. |
| Protocol integration | Standards-based protocol adapters between apps and ZITADEL. | OIDC/SAML auth URL handling, callback validation, code/token exchange, logout URL construction. | User-facing check UX orchestration and federation resource management. | `@zitadel/nextjs/auth/oidc`, OIDC wrappers in `@zitadel/zitadel-js`. |
| Federation management | Management of upstream identity provider federation resources. | CRUD for upstream IdPs and mapping/policy configuration through admin APIs. | Runtime end-user session/auth protocol execution. | Current usage via admin API clients; future SDK lane modules should stay isolated from `auth/*`. |
| Admin / events APIs | Administrative API access and event ingestion/verification surfaces. | Management/admin/system/event API clients, webhook validation/decryption, typed event payload handling. | Interactive end-user sign-in flows. | `@zitadel/nextjs/api`, `@zitadel/nextjs/webhook`, webhook verification utilities in `@zitadel/zitadel-js`. |

Boundary rule: when a feature spans multiple lanes, compose modules across lanes instead of introducing a mixed abstraction.

### Maintainer lane triage (new features)
1. User-facing check/session orchestration for custom login UIs belongs to **Embedded auth concepts** (`auth/session`).
2. OIDC/SAML redirect/callback/token/logout protocol work belongs to **Protocol integration** (`auth/<protocol>`).
3. Upstream IdP CRUD/mapping/policy management belongs to **Federation management** (`idp/<protocol>`), even if currently reached via `api/*`.
4. Admin/system/event API access or inbound event verification belongs to **Admin / events APIs** (`api/*`, `webhook/*`).
5. If one feature needs multiple lanes, split it into lane-specific modules and compose.

## Module Naming Convention
Canonical cross-SDK module IDs are path-like and protocol-explicit:
- `auth/session` — session/check orchestration for embedded auth UIs
- `auth/oidc` — OIDC end-user auth protocol adapter
- `auth/saml` — SAML end-user auth protocol adapter
- `idp/oidc` — OIDC upstream IdP federation management
- `idp/saml` — SAML upstream IdP federation management
- `api/*` — typed API clients grouped by API family (for example `api/management`, `api/system`)
- `webhook/*` — inbound event/webhook verification grouped by source (for example `webhook/actions`)

Naming contract:
1. Each public SDK surface maps to exactly one canonical ID.
2. Keep the root segment in `{auth,idp,api,webhook}` and add a second segment for protocol/family when needed.
3. Check types (password, passkey, TOTP, etc.) stay as parameters within `auth/session`, not separate modules.
4. Start from the canonical ID first, then apply language-specific separators/casing without changing segment order or meaning.

### JS/Next.js ↔ Go Mapping Matrix
Use this as the maintainer reference for cross-SDK naming and rollout order.

| Canonical ID | JS/Next.js current | JS/Next.js planned | Go current | Go planned |
| --- | --- | --- | --- | --- |
| `auth/session` | `@zitadel/nextjs/auth/session`; shared session primitives in `@zitadel/zitadel-js`. | Keep `auth/session` as the primary auth UI lane and align aliases/docs to canonical naming. | No dedicated public Go SDK surface yet. | `auth/session` package as the embedded auth baseline. |
| `auth/oidc` | `@zitadel/nextjs/auth/oidc`; OIDC helpers in `@zitadel/zitadel-js`. | Keep `auth/oidc` as the protocol lane and preserve canonical naming in all wrappers. | No dedicated public Go SDK surface yet. | `auth/oidc` package for redirect-based OIDC flows. |
| `auth/saml` | Not exposed yet. | Add `auth/saml` as a peer lane to `auth/oidc` (same lane boundaries). | No dedicated public Go SDK surface yet. | `auth/saml` package for SAML end-user auth flows. |
| `idp/oidc` | No dedicated module; currently handled via management APIs (`@zitadel/nextjs/api`, `@zitadel/zitadel-js/v2`). | Add `idp/oidc` wrappers over management APIs without mixing with `auth/*`. | No dedicated public Go SDK surface yet. | `idp/oidc` package for upstream OIDC federation management. |
| `idp/saml` | No dedicated module; currently handled via management APIs (`@zitadel/nextjs/api`, `@zitadel/zitadel-js/v2`). | Add `idp/saml` wrappers over management APIs without mixing with `auth/*`. | No dedicated public Go SDK surface yet. | `idp/saml` package for upstream SAML federation management. |
| `api/*` | `@zitadel/nextjs/api` and `@zitadel/zitadel-js/v2`. | Expand family-specific surfaces (for example `api/management`, `api/system`) under the same root. | No dedicated public Go SDK surface yet. | `api/*` packages grouped by API family. |
| `webhook/*` | `@zitadel/nextjs/webhook`; verification utilities in `@zitadel/zitadel-js/webhooks`. | Keep webhook handlers under canonical `webhook/*` grouping (for example `webhook/actions`). | No dedicated public Go SDK surface yet. | `webhook/*` packages grouped by event source (starting with `webhook/actions`). |

Parity notes:
- Target parity is canonical-name parity first: each lane above should exist as one primary JS surface and one primary Go surface.
- Temporary gaps are allowed while Go SDK surfaces are introduced; JS may continue to use `api/*` as the interim path for `idp/*` and `auth/saml`.
- Temporary gaps/aliases must stay documented in this matrix with owning SDK + target milestone, and aliases must be deprecated until removed.

Do:
- Keep canonical roots and segment order aligned across SDKs.
- Extend existing roots (`api/*`, `webhook/*`) instead of inventing new top-level groups.
- If aliases are required for migration, keep them temporary and deprecate with a timeline.

Don't:
- Rename canonical surfaces per language when an equivalent canonical name exists (for example `login`, `federation`, `hooks`).
- Mix end-user auth execution (`auth/*`) and upstream federation management (`idp/*`) in one surface.
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
