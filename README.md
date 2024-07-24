# ZITADEL TypeScript with Turborepo and Changesets

This repository contains all TypeScript and JavaScript packages and applications you need to create your own ZITADEL
Login UI.
The repo makes use of the [build system Turbo](https://turbo.build/repo) and
the [Changesets CLI for versioning the packages](https://github.com/changesets/changesets).

**⚠️ This repo and packages are in alpha state and subject to change ⚠️**

The scope of functionality of this repo and packages is limited and under active development.
Once the package structure is set and all APIs are fully implemented we'll move this repo to beta state.
You can read the [contribution guide](/CONTRIBUTING.md) on how to contribute.
Questions can be raised in our [Discord channel](https://discord.gg/erh5Brh7jE) or as
a [GitHub issue](https://github.com/zitadel/typescript/issues).

## Developing Your Own ZITADEL Login UI

We think the easiest path of getting up and running, is the following:

1. Fork and clone this repository
1. [Run the ZITADEL Cloud login UI locally](#run-login-ui)
1. Make changes to the code and see the effects live on your local machine
1. Study the rest of this README.md and get familiar and comfortable with how everything works.
1. Decide on a way of how you want to build and run your login UI.
   You can reuse ZITADEL Clouds way.
   But if you need more freedom, you can also import the packages you need into your self built application.

## Included Apps And Packages

- `login`: The login UI used by ZITADEL Cloud, powered by Next.js
- `@zitadel/node`: core components for establishing node client connection, grpc stub
- `@zitadel/client`: shared client utilities
- `@zitadel/proto`: shared protobuf types
- `@zitadel/react`: shared React utilities and components built with tailwindcss
- `@zitadel/next`: shared Next.js utilities
- `@zitadel/tsconfig`: shared `tsconfig.json`s used throughout the monorepo
- `eslint-config-zitadel`: ESLint preset

Each package and app is 100% [TypeScript](https://www.typescriptlang.org/).

### Login

The login is currently in a work in progress state.
The goal is to implement a login UI, using the session API of ZITADEL, which also implements the OIDC Standard and is
ready to use for everyone.

In the first phase we want to have a MVP login ready with the OIDC Standard and a basic feature set. In a second step
the features will be extended.

This list should show the current implementation state, and also what is missing.
You can already use the current state, and extend it with your needs.

- [x] Local User Registration (with Password)
- [x] User Registration and Login with external Provider
  - [x] Google
  - [ ] GitHub
  - [ ] GitHub Enterprise
  - [x] GitLab
  - [ ] GitLab Enterprise
  - [ ] Azure
  - [ ] Apple
  - [ ] Generic OIDC
  - [ ] Generic OAuth
  - [ ] Generic JWT
  - [ ] LDAP
  - [ ] SAML SP
- Multifactor Registration an Login
  - [x] Passkeys
  - [x] TOTP
  - [x] OTP: Email Code
  - [x] OTP: SMS Code
- [ ] Password Change/Reset
- [ ] Domain Discovery
- [x] Branding
- OIDC Standard
  - [x] Authorization Code Flow with PKCE
  - [x] AuthRequest `hintUserId`
  - [x] AuthRequest `loginHint`
  - [x] AuthRequest `prompt`
    - [x] Login
    - [x] Select Account
    - [ ] Consent
    - [x] Create
  - Scopes
    - [x] `openid email profile address``
    - [x] `offline access`
    - [ ] `urn:zitadel:iam:org:idp:id:{idp_id}`
    - [x] `urn:zitadel:iam:org:project:id:zitadel:aud`
    - [x] `urn:zitadel:iam:org:id:{orgid}`
    - [x] `urn:zitadel:iam:org:domain:primary:{domain}`
  - [ ] AuthRequest UI locales

## Tooling

- [TypeScript](https://www.typescriptlang.org/) for static type checking
- [ESLint](https://eslint.org/) for code linting
- [Prettier](https://prettier.io) for code formatting

## Useful Commands

- `pnpm generate` - Build proto stubs for server and client package
- `pnpm build` - Build all packages and the login app
- `pnpm test` - Test all packages and the login app
- `pnpm test:watch` - Rerun tests on file change
- `pnpm dev` - Develop all packages and the login app
- `pnpm lint` - Lint all packages
- `pnpm changeset` - Generate a changeset
- `pnpm clean` - Clean up all `node_modules` and `dist` folders (runs each package's clean script)

## Versioning And Publishing Packages

Package publishing has been configured using [Changesets](https://github.com/changesets/changesets).
Here is their [documentation](https://github.com/changesets/changesets#documentation) for more information about the
workflow.

The [GitHub Action](https://github.com/changesets/action) needs an `NPM_TOKEN` and `GITHUB_TOKEN` in the repository
settings. The [Changesets bot](https://github.com/apps/changeset-bot) should also be installed on the GitHub repository.

Read the [changesets documentation](https://github.com/changesets/changesets/blob/main/docs/automating-changesets.md)
for more information about this automation

### NPM

If you want to publish a package to the public npm registry and make them publicly available, this is already setup.

To publish packages to a private npm organization scope, **remove** the following from each of the `package.json`'s

```diff
- "publishConfig": {
-  "access": "public"
- },
```

### GitHub Package Registry

See [working with the npm registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-npm-registry#publishing-a-package-using-publishconfig-in-the-packagejson-file)

### Run Login UI

To run the application make sure to install the dependencies with

```sh
pnpm install
```

then setup the environment for the login application which needs a `.env.local` in `/apps/login`.
Go to your instance and create a service user for the application having the `IAM_OWNER` manager role.
This user is required to have access to create users on your primary organization and reading policy data so it can be
restricted to your personal use case but we'll stick with `IAM_OWNER` for convenience. Create a PAT and copy the value to
paste it under the `ZITADEL_SERVICE_USER_TOKEN` key.
The file should look as follows:

```
ZITADEL_API_URL=[yourinstanceurl]
ZITADEL_ORG_ID=[yourprimaryorg]
ZITADEL_SERVICE_USER_TOKEN=[yourserviceuserpersonalaccesstoken]
```

then generate the GRPC stubs with

```sh
pnpm generate
```

and then run it with

```sh
pnpm dev
```

Open the login application with your favorite browser at `localhost:3000`.

### Deploy to Vercel

To deploy your own version on Vercel, navigate to your instance and create a service user.
Copy its id from the overview and set it as ZITADEL_SERVICE_USER_ID.
Then create a personal access token (PAT), copy and set it as ZITADEL_SERVICE_USER_TOKEN, then navigate to your instance
settings and make sure it gets IAM_OWNER permissions.
Finally set your instance url as ZITADEL_API_URL. Make sure to set it without trailing slash.

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fzitadel%2Ftypescript&env=ZITADEL_API_URL,ZITADEL_SERVICE_USER_ID,ZITADEL_SERVICE_USER_TOKEN&root-directory=apps/login&envDescription=Setup%20a%20service%20account%20with%20IAM_OWNER%20membership%20on%20your%20instance%20and%20provide%20its%20id%20and%20personal%20access%20token.&project-name=zitadel-login&repository-name=zitadel-login)
