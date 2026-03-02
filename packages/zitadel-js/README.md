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
| `import { createAuthorizationBearerInterceptor } from "@zitadel/zitadel-js"` | `import { createBearerTokenInterceptor } from "@zitadel/zitadel-js/api/bearer-token"` |

Root imports remain for shared transport/client primitives:
`createClientFor`, `createConnectTransport`, `createGrpcTransport`.
