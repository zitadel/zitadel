# ZITADEL typescript with Turborepo and Changesets

This is an monorepo containing all typescript/javascript packages and applications for ZITADEL. Versioning and package publishing is handled by [Changesets](https://github.com/changesets/changesets) and fully automated with GitHub Actions.

**⚠️ This repo and packages are in alpha state and subject to change ⚠️**

The scope of functionality of this repo and packages is limited and under active development.
Once the package structure is set and all APIs are fully implemented we'll move this repo to beta state.
You can read the [contribution guide](/CONTRIBUTING.md) on how to contribute. Questions can be raised in our [Discord Channel](https://discord.gg/erh5Brh7jE) or as an issue in this repo.

## What's inside?

This Turborepo includes the following:

### Apps and Packages

- `login`: The new login UI powered by Next.js
- `@zitadel/server`: core components for establishing node client connection, grpc stub
- `@zitadel/client`: core components for establishing web client connection, grpc stub
- `@zitadel/react`: shared React utilities and components built with tailwindcss
- `@zitadel/next`: shared Next.js utilities
- `@zitadel/tsconfig`: shared `tsconfig.json`s used throughout the monorepo
- `eslint-config-zitadel`: ESLint preset

Each package and app is 100% [TypeScript](https://www.typescriptlang.org/).

### Utilities

This repo has some additional tools:

- [TypeScript](https://www.typescriptlang.org/) for static type checking
- [ESLint](https://eslint.org/) for code linting
- [Prettier](https://prettier.io) for code formatting

### Useful commands

- `pnpm generate` - Build proto stubs for server and client package
- `pnpm build` - Build all packages and the login app
- `pnpm dev` - Develop all packages and the login app
- `pnpm lint` - Lint all packages
- `pnpm changeset` - Generate a changeset
- `pnpm clean` - Clean up all `node_modules` and `dist` folders (runs each package's clean script)

## Versioning and Publishing packages

Package publishing has been configured using [Changesets](https://github.com/changesets/changesets). Here is their [documentation](https://github.com/changesets/changesets#documentation) for more information about the workflow.

The [GitHub Action](https://github.com/changesets/action) needs an `NPM_TOKEN` and `GITHUB_TOKEN` in the repository settings. The [Changesets bot](https://github.com/apps/changeset-bot) should also be installed on the GitHub repository.

Read the [changesets documentation](https://github.com/changesets/changesets/blob/main/docs/automating-changesets.md) for more information about this automation

### npm

If you want to publish package to the public npm registry and make them publicly available, this is already setup.

To publish packages to a private npm organization scope, **remove** the following from each of the `package.json`'s

```diff
- "publishConfig": {
-  "access": "public"
- },
```

### GitHub Package Registry

See [Working with the npm registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-npm-registry#publishing-a-package-using-publishconfig-in-the-packagejson-file)

### TODOs

- Buf setup to get grpc stub in the core package
- Decide whether a seperate client package is required to expose public client convenience methods only or generate a grpc-web output there
- Fix #/\* path in login application

### Run Login UI

To run the application make sure to install the dependencies with

```sh
pnpm install
```

then setup the environment for the login application which needs a `.env.local` in `/apps/login`.
Go to your instance and create a service user for the application having the IAM_OWNER manager role.
This user is required to have access to create users on your primary organization and reading policy data so it can be restricted to your personal use case but we'll stick with IAM_OWNER for convenience. Create a PAT and copy the value to paste it under the `ZITADEL_SERVICE_USER_TOKEN` key.
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
