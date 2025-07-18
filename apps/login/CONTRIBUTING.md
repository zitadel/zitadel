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

The repository zitadel/typescript is a read-only mirror of the git subtree at zitadel/zitadel/login.
To submit changes, please open a Pull Request [in the zitadel/zitadel repository](https://github.com/zitadel/zitadel/compare).

If you already made changes based on the zitadel/typescript repository, these changes are not lost.
Submitting them to the main repository is easy:

1. [Fork zitadel/zitadel](https://github.com/zitadel/zitadel/fork)
1. Clone your Zitadel fork git clone https://github.com/<your-owner>/zitadel.git
1. Change directory to your Zitadel forks root.
1. Pull your changes into the Zitadel fork by running make login_pull LOGIN_REMOTE_URL=https://github.com/<your-owner>/typescript.git LOGIN_REMOTE_BRANCH=<your-typescript-fork-branch>.
1. Push your changes and [open a pull request to zitadel/zitadel](https://github.com/zitadel/zitadel/compare)

Please consider the following guidelines when creating a pull request.

- The latest changes are always in `main`, so please make your pull request against that branch.
- pull requests should be raised for any change
- Pull requests need approval of a Zitadel core engineer @zitadel/engineers before merging
- We use ESLint/Prettier for linting/formatting, so please run `pnpm lint:fix` before committing to make resolving conflicts easier (VSCode users, check out [this ESLint extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [this Prettier extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) to fix lint and formatting issues in development)
- If you add new functionality, please provide the corresponding documentation as well and make it part of the pull request

### Setting up local environment

```sh
# Install dependencies. Developing requires Node.js v20
pnpm install

# Generate gRPC stubs
pnpm generate

# Start a local development server for the login and manually configure apps/login/.env.local
pnpm dev
```

The application is now available at `http://localhost:3000`

Configure apps/login/.env.local to target the Zitadel instance of your choice.
The login app live-reloads on changes, so you can start developing right away.

### <a name="latest"></a>Developing Against A Local Latest Zitadel Release

The following command uses Docker to run a local Zitadel instance and the login application in live-reloading dev mode.
Additionally, it runs a Traefik reverse proxy that exposes the login with a self-signed certificate at https://127.0.0.1.sslip.io
127.0.0.1.sslip.io is a special domain that resolves to your localhost, so it's safe to allow your browser to proceed with loading the page.

```sh
# Install dependencies. Developing requires Node.js v20
pnpm install

# Generate gRPC stubs
pnpm generate

# Start a local development server and have apps/login/.env.test.local configured for you to target the local Zitadel instance.
pnpm dev:local
```

Log in at https://127.0.0.1.sslip.io/ui/v2/login/loginname and use the following credentials:
**Loginname**: *zitadel-admin@zitadel.127.0.0.1.sslip.io*
**Password**: _Password1!_.

The login app live-reloads on changes, so you can start developing right away.

### <a name="local"></a>Developing Against A Locally Compiled Zitadel

To develop against a locally compiled version of Zitadel, you need to build the Zitadel docker image first.
Clone the [Zitadel repository](https://github.com/zitadel/zitadel.git) and run the following command from its root:

```sh
# This compiles a Zitadel binary if it does not exist at ./zitadel already and copies it into a Docker image.
# If you want to recompile the binary, run `make compile` first
make login_dev
```

Open another terminal session at zitadel/zitadel/login and run the following commands to start the dev server.

```bash
# Install dependencies. Developing requires Node.js v20
pnpm install

# Start a local development server and have apps/login/.env.test.local configured for you to target the local Zitadel instance.
NODE_ENV=test pnpm dev
```

Log in at https://127.0.0.1.sslip.io/ui/v2/login/loginname and use the following credentials:
**Loginname**: *zitadel-admin@zitadel.127.0.0.1.sslip.io*
**Password**: _Password1!_.

The login app live-reloads on changes, so you can start developing right away.

### Quality Assurance

Use `make` commands to test the quality of your code against a production build without installing any dependencies besides Docker.
Using `make` commands, you can reproduce and debug the CI pipelines locally.

```sh
# Reproduce the whole CI pipeline in docker
make login_quality
# Show other options with make
make help
```

Use `pnpm` commands to run the tests in dev mode with live reloading and debugging capabilities.

#### Linting and formatting

Check the formatting and linting of the code in docker

```sh
make login_lint
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
make login_test_unit
```

Run unit tests with live-reloading

```sh
pnpm test:unit
```

#### Running Integration Tests

Run the test in docker

```sh
make login_test_integration
```

Alternatively, run a live-reloading development server with an interactive Cypress test suite.
First, set up your local test environment.

```sh
# Install dependencies. Developing requires Node.js v20
pnpm install

# Generate gRPC stubs
pnpm generate

# Start a local development server and use apps/login/.env.test to use the locally mocked Zitadel API.
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

#### Running Acceptance Tests

To run the tests in docker against the latest release of Zitadel, use the following command:

:warning: The acceptance tests are not reliable at the moment :construction:

```sh
make login_test_acceptance
```

Alternatively, run can use a live-reloading development server with an interactive Playwright test suite.
Set up your local environment by running the commands either for [developing against a local latest Zitadel release](latest) or for [developing against a locally compiled Zitadel](compiled).

Now, in another terminal session, open the interactive Playwright acceptance test suite.

```sh
pnpm test:acceptance open
```

Show more options with Playwright

```sh
pnpm test:acceptance help
```
