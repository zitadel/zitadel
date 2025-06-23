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

Configure apps/login/.env.local to target the Zitadel instance of your choice.
The login app live-reloads on changes, so you can start developing right away.

<!-- Console doesn't load

### Developing Against Your Local ZITADEL Instance

The following command uses Docker to run a local ZITADEL instance and the login application in live-reloading dev mode.
Additionally, it runs a Traefik reverse proxy that exposes the login at https://localhost with a self-signed certificate.

```sh
pnpm test:acceptance:setup
```
-->

### Quality Assurance

Use `make` commands to test the quality of your code without installing any dependencies besides Docker.
Using `make` commands, you can reproduce and debug the CI pipelines locally.
```sh
# Reproduce the whole CI pipeline in docker
make login-quality
# Show other options with make
make help
```

Use `pnpm` commands to run the tests in dev mode with live reloading and debugging capabilities.

#### Linting and formatting

Check the formatting and linting of the code in docker

```sh
make login-lint
```

Check the linting of the code using pnpm

```sh
pnpm lint
pnpm format
```

Fix the linting of your code

```sh
pnpm lint:fix
pnpm format:fix
```

#### Running Unit Tests

Run the tests in docker

```sh
make login-test-unit
```

Run unit tests with live-reloading

```sh
pnpm test:unit
```

#### Running Integration Tests

Run the test in docker

```sh
make login-test-integration
```

Open the Cypress test suite to run the integration tests in interactive mode.
First, set up your local test environment.
This runs a mock server in docker and the login application in dev mode with live-reloading enabled.

```sh
pnpm test:integration:setup
```

Now, in another terminal session, open the interactive Cypress integration test suite.

```sh
pnpm test:integration open
```

Show more options with Cypress

```sh
pnpm test:integration help
```
