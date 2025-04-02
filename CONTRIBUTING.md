# Contributing

:attention: In this CONTRIBUTING.md you read about contributing to this very repository.
If you want to develop your own login UI, please refer [to the README.md](./README.md).

## Introduction

Thank you for your interest about how to contribute!

:attention: If you notice a possible **security vulnerability**, please don't hesitate to disclose any concern by contacting [security@zitadel.com](mailto:security@zitadel.com).
You don't have to be perfectly sure about the nature of the vulnerability.
We will give them a high priority and figure them out.

We also appreciate all your other ideas, thoughts and feedback and will take care of them as soon as possible.
We love to discuss in an open space using [GitHub issues](https://github.com/zitadel/typescript/issues),
[GitHub discussions in the core repo](https://github.com/zitadel/zitadel/discussions)
or in our [chat on Discord](https://zitadel.com/chat).
For private discussions,
you have [more contact options on our Website](https://zitadel.com/contact).

## Pull Requests

Please consider the following guidelines when creating a pull request.

- The latest changes are always in `main`, so please make your pull request against that branch.
- pull requests should be raised for any change
- Pull requests need approval of a ZITADEL core engineer @zitadel/engineers before merging
- We use ESLint/Prettier for linting/formatting, so please run `pnpm lint:fix` before committing to make resolving conflicts easier (VSCode users, check out [this ESLint extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [this Prettier extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) to fix lint and formatting issues in development)
- If you add new functionality, please provide the corresponding documentation as well and make it part of the pull request

## Setting Up The ZITADEL API

If you want to have a one-liner to get you up and running,
or if you want to develop against a ZITADEL API with the latest features,
or even add changes to ZITADEL itself at the same time,
you should develop against your local ZITADEL process.
However, it might be easier to develop against your ZITADEL Cloud instance
if you don't have docker installed
or have limited resources on your local machine.

### Developing Against Your Local ZITADEL Instance

```sh
# To have your service user key and environment file written with the correct ownership, export your current users ID.
export ZITADEL_DEV_UID="$(id -u)"

# Pull images
docker compose --file ./acceptance/docker-compose.yaml pull

# Run ZITADEL with local notification sink and configure ./apps/login/.env.local
pnpm run-sink
```

### Developing Against Your ZITADEL Cloud Instance

Configure your shell by exporting the following environment variables:

```sh
export ZITADEL_API_URL=<your cloud instance URL here>
export ZITADEL_ORG_ID=<your service accounts organization id here>
export ZITADEL_SERVICE_USER_TOKEN=<your service account personal access token here>
```

### Setting up local environment

```sh
# Install dependencies. Developing requires Node.js v20
pnpm install

# Generate gRPC stubs
pnpm generate

# Start a local development server
pnpm dev
```

The application is now available at `http://localhost:3000`

### Adding applications and IDPs

```sh
# OPTIONAL Run SAML SP
pnpm run-samlsp

# OPTIONAL Run OIDC RP
pnpm run-oidcrp

# OPTIONAL Run SAML IDP
pnpm run-samlidp

# OPTIONAL Run OIDC OP
pnpm run-oidcop
```

### Testing

You can execute the following commands `pnpm test` for a single test run or `pnpm test:watch` in the following directories:

- apps/login
- packages/zitadel-proto
- packages/zitadel-client
- packages/zitadel-node
- The projects root directory: all tests in the project are executed

In apps/login, these commands also spin up the application and a ZITADEL gRPC API mock server to run integration tests using [Cypress](https://www.cypress.io/) against them.
If you want to run the integration tests standalone against an environment of your choice, navigate to ./apps/login, [configure your shell as you like](# Developing Against Your ZITADEL Cloud Instance) and run `pnpm test:integration:run` or `pnpm test:integration:open`.
Then you need to lifecycle the mock process using the command `pnpm mock` or the more fine grained commands `pnpm mock:build`, `pnpm mock:build:nocache`, `pnpm mock:run` and `pnpm mock:destroy`.

That's it! 🎉
