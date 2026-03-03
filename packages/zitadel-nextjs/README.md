# @zitadel/nextjs

ZITADEL Next.js SDK — App Router integration with OIDC lifecycle.

## Installation

```sh
npm install @zitadel/nextjs
```

## Quick bootstrap for existing Next.js apps

```sh
npx @zitadel/nextjs add
```

By default, this command adds OIDC route handlers plus a dedicated test page at `/auth` (`route.ts`/`route.js` based on your project) and updates `.env.example` with required `ZITADEL_*` variables.

You can be explicit about scaffolds:

```sh
npx @zitadel/nextjs add --auth session --with-api --with-webhook
```

Use `--dry-run` to preview changes, `--cwd <path>` to target a specific app, `--source workspace` to install workspace-linked SDK dependencies, and `--skip-install` to avoid dependency installation.

For OIDC redirect testing, set these in `.env.local`:

```env
ZITADEL_ISSUER_URL=https://your-instance.zitadel.cloud
ZITADEL_CLIENT_ID=your-client-id
ZITADEL_CALLBACK_URL=http://localhost:3000/api/auth/callback
ZITADEL_COOKIE_SECRET=at-least-32-characters
ZITADEL_POST_LOGIN_URL=/auth
ZITADEL_POST_LOGOUT_URL=/auth
```

Then start your app and open `/auth`.

For fast local SDK iteration in this monorepo (recommended):

```sh
pnpm install
pnpm nx run @zitadel/nextjs-sdk-playground:dev-with-sdk
```

This starts watch builds for `@zitadel/zitadel-js`, `@zitadel/react`, and `@zitadel/nextjs` plus the Next.js playground at `tests/nextjs-sdk-playground`.

If you scaffold an app that already lives inside a pnpm workspace, you can link workspace dependencies directly:

```sh
npx @zitadel/nextjs add --source workspace
```

For local testing before publishing (external non-workspace app fallback):

```sh
pnpm nx run @zitadel/nextjs:build
node packages/zitadel-nextjs/dist/cli.cjs add --cwd /path/to/app --source npm --skip-install

# in the SDK repo, pack local tarballs
pnpm --filter @zitadel/zitadel-js pack --pack-destination /tmp
pnpm --filter @zitadel/react pack --pack-destination /tmp
pnpm --filter @zitadel/nextjs pack --pack-destination /tmp

# in your target app, install those tarballs
npm install /tmp/zitadel-zitadel-js-0.1.0.tgz /tmp/zitadel-react-0.1.0.tgz /tmp/zitadel-nextjs-0.1.0.tgz
```

## Modules

- `@zitadel/nextjs/auth/oidc`: OIDC redirect flow helpers (`signIn`, `handleCallback`, `signOut`)
- `@zitadel/nextjs/auth/session`: Session API helpers for custom login UIs (`createSession`, `setSession`, `getSession`, `deleteSession`, `createCallback`)
- `@zitadel/nextjs/api`: authenticated v2 API client factory
- `@zitadel/nextjs/auth/bearer-token`: canonical bearer token extraction/validation helpers for API route handlers
- `@zitadel/nextjs/actions/webhook`: canonical Actions v2 webhook route handler

## Migration notes: prefer lane entrypoints

Use lane-specific imports to keep UI and API concerns separated:

| Use case | Import |
| --- | --- |
| OIDC redirect login | `@zitadel/nextjs/auth/oidc` |
| Custom login UI via Session API | `@zitadel/nextjs/auth/session` |
| Authenticated management/system API client | `@zitadel/nextjs/api` |
| Bearer token validation in API routes | `@zitadel/nextjs/auth/bearer-token` |
| Actions v2 webhook handling | `@zitadel/nextjs/actions/webhook` |

## Documentation

For detailed documentation and API references, please visit the [ZITADEL documentation](https://zitadel.com/docs).
