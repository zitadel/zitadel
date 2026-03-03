# @zitadel/zitadel-js

ZITADEL JavaScript SDK — isomorphic core. Zero framework dependencies. Runs in Node.js ≥18, Edge runtimes, and browsers.

## Installation

```sh
npm install @zitadel/zitadel-js
```

## Documentation

For detailed documentation and API references, please visit the [ZITADEL documentation](https://zitadel.com/docs).

## Breaking migration: lane-first imports

`@zitadel/zitadel-js` now uses explicit lane entrypoints.

| Previous import | New import |
| --- | --- |
| `import { createOIDCAuthorizationUrl, exchangeOIDCAuthorizationCode, discoverOIDCAuthorizationServer, refreshOIDCTokens, createOIDCEndSessionUrl, generatePKCE, generateState } from "@zitadel/zitadel-js"` | `import { createOIDCAuthorizationUrl, exchangeOIDCAuthorizationCode, discoverOIDCAuthorizationServer, refreshOIDCTokens, createOIDCEndSessionUrl, generatePKCE, generateState } from "@zitadel/zitadel-js/auth/oidc"` |
| `import { isSessionValid, isSessionExpired } from "@zitadel/zitadel-js"` | `import { isSessionValid, isSessionExpired } from "@zitadel/zitadel-js/auth/session"` |
| `import { createAuthorizationBearerInterceptor } from "@zitadel/zitadel-js"` | `import { createBearerTokenInterceptor } from "@zitadel/zitadel-js/auth/bearer-token"` |

Root imports remain for shared transport/client primitives:
`createClientFor`, `createConnectTransport`, `createGrpcTransport`.

For API discoverability, use `@zitadel/zitadel-js/api/v1` (legacy API) and `@zitadel/zitadel-js/api/v2` (current API). `@zitadel/zitadel-js/v2` remains available as a compatibility alias.

Canonical taxonomy is lane-first: `auth/*`, `api/*`, and `actions/*`, plus root/core transport primitives.

For migration compatibility, `@zitadel/zitadel-js/api/bearer-token` remains an alias to `@zitadel/zitadel-js/auth/bearer-token`, and `@zitadel/zitadel-js/webhooks` remains an alias to `@zitadel/zitadel-js/actions/webhook`.

`@zitadel/client` and `@zitadel/proto` are not replaced in the current iteration; consolidation remains staged and will include explicit deprecation guidance before package lifecycle changes.
