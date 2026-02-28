# @zitadel/examples/nextjs

Next.js SDK demo lanes for OIDC, Session API, and ZITADEL v2 API calls.

## Demo lanes (UI routes)

- `/demo/oidc` — OIDC redirect lane (`/api/auth/signin`, `/api/auth/callback`, `/api/auth/signout`)
- `/demo/username-password` — Session API lane (`/api/demo/username-password/session`, `/api/demo/username-password/callback`)
- `/demo/signup` — Signup lane (`POST /api/demo/signup/password`)
- `/demo/org-registration` — Organization registration lane (`POST /api/demo/org-registration/register`)

## Required environment variables

Copy `.env.example` to `.env.local` and set:

### OIDC lane

- `ZITADEL_ISSUER_URL`
- `ZITADEL_CLIENT_ID`
- `ZITADEL_CALLBACK_URL` (for this app: `http://localhost:3000/api/auth/callback`)
- `ZITADEL_COOKIE_SECRET` (at least 32 chars)
- `ZITADEL_POST_LOGIN_URL` (recommended: `/demo/oidc`)
- `ZITADEL_POST_LOGOUT_URL` (recommended: `/demo/oidc`)

### API lanes (`signup`, `org-registration`, `username-password`)

- `ZITADEL_API_URL`
- and one auth option:
  - `ZITADEL_SERVICE_USER_TOKEN`, or
  - all of `ZITADEL_SERVICE_USER_KEY_ID`, `ZITADEL_SERVICE_USER_ID`, `ZITADEL_SERVICE_USER_PRIVATE_KEY`, or
  - an active OIDC session access token with API audience scope.

## Permissions and scope notes

- Signup lane (`userService.addHumanUser`) requires `user.write`.
- Org registration lane (`organizationService.addOrganization`) requires `org.create`.
- Session lane:
  - `createSession` requires `session.write`
  - `getSession` requires `session.read` (or own/current session token conditions)
  - `createCallback` requires `session.link`
- If API calls use the signed-in user token (instead of service-user credentials), request API audience scope `urn:zitadel:iam:org:project:id:zitadel:aud` in addition to OIDC scopes (`openid profile email`). Add `offline_access` if you also want refresh tokens.

## Local run/build

From the repository root:

```sh
pnpm install
cp examples/nextjs/.env.example examples/nextjs/.env.local
pnpm nx run @zitadel/examples/nextjs:dev
```

Open `http://localhost:3000` and select a demo lane.

Build command:

```sh
pnpm nx run @zitadel/examples/nextjs:build
```

For workspace-linked SDK watch mode:

```sh
pnpm nx run @zitadel/examples/nextjs:dev-with-sdk
```
