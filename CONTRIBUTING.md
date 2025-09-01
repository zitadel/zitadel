# Contributing to Zitadel

## Introduction

Thank you for your interest about how to contribute! As you might know there is more than code to contribute. You can find all information needed to start contributing here.

Please give us and our community the chance to get rid of security vulnerabilities by responsibly disclose this kind of issues by contacting [security@zitadel.com](mailto:security@zitadel.com).

The strongest part of a community is the possibility to share thoughts. That's why we try to react as soon as possible to your ideas, thoughts and feedback. We love to discuss as much as possible in an open space like in the [issues](https://github.com/zitadel/zitadel/issues) and [discussions](https://github.com/zitadel/zitadel/discussions) section here or in our [chat](https://zitadel.com/chat), but we understand your doubts and provide further contact options [here](https://zitadel.com/contact).

If you want to give an answer or be part of discussions please be kind. Treat others like you want to be treated. Read more about our code of conduct [here](CODE_OF_CONDUCT.md).

## What can I contribute?

For people who are new to Zitadel: We flag issues which are a good starting point to start contributing.
You find them [here](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
We add the label "good first issue" for problems we think are a good starting point to contribute to Zitadel.

- [Issues for first time contributors](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [All issues](https://github.com/zitadel/zitadel/issues)

Help shaping the future of Zitadel:

- Join our [chat](https://zitadel.com/chat) and discuss with us or others.
- Ask or answer questions in the [issues section](https://github.com/zitadel/zitadel/issues)
- Share your thoughts and ideas in the [discussions section](https://github.com/zitadel/zitadel/discussions)

Make Zitadel more popular and give it a ⭐

Follow [@zitadel](https://twitter.com/zitadel) on twitter

[Contribute](#how-to-contribute)

- [Contribute code](#contribute)
- If you found a mistake on our [docs page](https://zitadel.com/docs) or something is missing please read [the docs section](contribute-docs)
- [Translate](#contribute-internationalization) and improve texts

## How to contribute

We strongly recommend to [talk to us](https://zitadel.com/contact) before you start contributing to streamline our and your work.

We accept contributions through pull requests.
You need a github account for that.
If you are unfamiliar with git have a look at Github's documentation on [creating forks](https://help.github.com/articles/fork-a-repo) and [creating pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork).
Please draft the pull request as soon as possible.
Go through the following checklist before you submit the final pull request:

### Components

The code consists of the following parts:

| name            | description                                        | language                                                                                                  | where to find                                       | Development Guide                                  |
| --------------- | -------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------- | -------------------------------------------------- |
| backend         | Service that serves the grpc(-web) and RESTful API | [go](https://go.dev)                                                                                      | [API implementation](./internal/api/grpc)           | [Contribute to Backend](contribute-backend)        |
| API definitions | Specifications of the API                          | [Protobuf](https://developers.google.com/protocol-buffers)                                                | [./proto/zitadel](./proto/zitadel)                  | [Contribute to Backend](contribute-backend)        |
| console         | Frontend the user interacts with after log in      | [Angular](https://angular.io), [Typescript](https://www.typescriptlang.org)                               | [./console](./console)                              | [Contribute to Frontend](contribute-frontend)      |
| login           | Modern authentication UI built with Next.js        | [Next.js](https://nextjs.org), [React](https://reactjs.org), [TypeScript](https://www.typescriptlang.org) | [./login](./login)                                  | [Contribute to Frontend](contribute-frontend)      |
| docs            | Project documentation made with docusaurus         | [Docusaurus](https://docusaurus.io/)                                                                      | [./docs](./docs)                                    | [Contribute to Frontend](contribute-frontend)      |
| translations    | Internationalization files for default languages   | YAML                                                                                                      | [./console](./console) and [./internal](./internal) | [Contribute Translations](contribute-translations) |

Please follow the guides to validate and test the code before you contribute.

### Submit a pull request (PR)

1. [Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the [zitadel/zitadel](https://github.com/zitadel/zitadel) repository on GitHub
2. On your fork, commit your changes to a new branch

   `git checkout -b my-fix-branch main`

3. Make your changes following the [guidelines](#contribute) in this guide. Make sure that all tests pass.

4. Commit the changes on the new branch

   `git commit --all`

5. [Merge](https://git-scm.com/book/en/v2/Git-Branching-Basic-Branching-and-Merging) the latest commit of the `main`-branch

6. Push the changes to your branch on Github

   `git push origin my-fix-branch`

7. Use [Semantic Release commit messages](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type) to simplify creation of release notes. In the title of the pull request [correct tagging](#commit-messages) is required and will be requested by the reviewers.

8. On GitHub, [send a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/requesting-a-pull-request-review) to `zitadel:main`. Request review from one of the maintainers.

### Review a pull request

The reviewers will provide you feedback and approve your changes as soon as they are satisfied.
If we ask you for changes in the code, you can follow the [GitHub Guide](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/incorporating-feedback-in-your-pull-request) to incorporate feedback in your pull request.

<!-- TODO: how to do this via git -->
<!-- TODO: change commit message via git -->

### Commit messages

Make sure you use [semantic release messages format](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type).

`<type>(<scope>): <short summary>`

#### Type

Must be one of the following:

- **feat**: New Feature
- **fix**: Bug Fix
- **docs**: Documentation

#### Scope

This is optional to indicate which component is affected. In doubt, leave blank (`<type>: <short summary>`)

#### Short summary

Provide a brief description of the change.

### Quality assurance

Please make sure you cover your changes with tests before marking a Pull Request as ready for review:

- [ ] Integration tests against the gRPC server ensure that one or multiple API calls that belong together return the expected results.
- [ ] Integration tests against the gRPC server ensure that probable good and bad read and write permissions are tested.
- [ ] Integration tests against the gRPC server ensure that the API is easily usable despite eventual consistency.
- [ ] Integration tests against the gRPC server ensure that all probable login and registration flows are covered."
- [ ] Integration tests ensure that certain commands emit expected events that trigger notifications.
- [ ] Integration tests ensure that certain events trigger expected notifications.

### General Guidelines

#### Gender Neutrality and Inclusive Language

We are committed to creating a welcoming and inclusive community for all developers, regardless of their gender identity or expression. To achieve this, we are actively working to ensure that our contribution guidelines are gender-neutral and use inclusive language.

**Use gender-neutral pronouns**:
Don't use gender-specific pronouns unless the person you're referring to is actually that gender.
In particular, don't use he, him, his, she, or her as gender-neutral pronouns, and don't use he/she or (s)he or other such punctuational approaches. Instead, use the singular they.

**Choose gender-neutral alternatives**:
Opt for gender-neutral terms instead of gendered ones whenever possible.
Replace "policeman" with "police officer," "manpower" with "workforce," and "businessman" with "entrepreneur" or "businessperson."

**Avoid ableist language**:
Ableist language includes words or phrases such as crazy, insane, blind to or blind eye to, cripple, dumb, and others.
Choose alternative words depending on the context.

### API

Zitadel follows an API first approach. This means all features can not only be accessed via the UI but also via the API.
The API is designed to be used by different clients, such as web applications, mobile applications, and other services.
Therefore, the API is designed to be easy to use, consistent, and reliable.
Please check out the dedicated [API guidelines](./API_DESIGN.md) page when contributing to the API.

## <a name="backend"></a>Contribute Backend Code



### <a name="backend-requirements"></a> Backend Requirements

By executing the commands from this section, you run everything you need to develop the Zitadel backend locally.

> [!INFO]
> Some [dev containers are available](dev-containers) for remote development with docker and pipeline debugging in isolated environments.
> If you don't want to use one of the dev containers, you can develop the backend components directly on your local machine.
> To do so, proceed with installing the necessary dependencies.

Using [Docker Compose](https://docs.docker.com/compose/), you run a [PostgreSQL](https://www.postgresql.org/download/) on your local machine.
With [make](https://www.gnu.org/software/make/), you build a debuggable Zitadel binary and run it using [delve](https://github.com/go-delve/delve).
Then, you test your changes via the console your binary is serving at http://<span because="breaks the link"></span>localhost:8080 and by verifying the database.
Once you are happy with your changes, you run end-to-end tests and tear everything down.

Zitadel uses [golangci-lint v2](https://golangci-lint.run) for code quality checks. Please use [this configuration](.golangci.yaml) when running `golangci-lint`. We recommend to set golangci-lint as linter in your IDE.

The commands in this section are tested against the following software versions:

- [Docker version 20.10.17](https://docs.docker.com/engine/install/)
- [Go version 1.22](https://go.dev/doc/install)
- [Delve 1.9.1](https://github.com/go-delve/delve/tree/v1.9.1/Documentation/installation)

### <a name="build-and-run-zitadel"></a>Build and Run Zitadel

Make some changes to the source code, then run the database locally.

```bash
# You just need the db service to develop the backend against.
docker compose --file ./e2e/docker-compose.yaml up --detach db
```

Build the binary. This takes some minutes, but you can speed up rebuilds.

```bash
make compile
```

> Note: With this command, several steps are executed.
> For speeding up rebuilds, you can reexecute only specific steps you think are necessary based on your changes.  
> Generating gRPC stubs: `make core_api`  
> Running unit tests: `make core_unit_test`  
> Generating the console: `make console_build console_move`  
> Build the binary: `make compile`

You can now run and debug the binary in .artifacts/zitadel/zitadel using your favourite IDE, for example GoLand.
You can test if Zitadel does what you expect by using the UI at http://localhost:8080/ui/console.
Also, you can verify the data by running `psql "host=localhost dbname=zitadel sslmode=disable"` and running SQL queries.

### Run Local Unit Tests

To test the code without dependencies, run the unit tests:

```bash
make core_unit_test
```

### Run Local Integration Tests

Integration tests are run as gRPC clients against a running Zitadel server binary.
The server binary is typically [build with coverage enabled](https://go.dev/doc/build-cover).
It is also possible to run a Zitadel sever in a debugger and run the integrations tests like that. In order to run the server, a database is required.

In order to prepare the local system, the following will bring up the database, builds a coverage binary, initializes the database and starts the sever.

```bash
make core_integration_db_up core_integration_server_start
```

When this job is finished, you can run individual package integration test through your IDE or command-line. The actual integration test clients reside in the `integration_test` subdirectory of the package they aim to test. Integration test files use the `integration` build tag, in order to be excluded from regular unit tests.
Because of the server-client split, Go is usually unaware of changes in server code and tends to cache test results. Pas `-count 1` to disable test caching.

Example command to run a single package integration test:

```bash
go test -count 1 -tags integration ./internal/api/grpc/management/integration_test
```

To run all available integration tests:

```bash
make core_integration_test_packages
```

When you change any Zitadel server code, be sure to rebuild and restart the server before the next test run.

```bash
make core_integration_server_stop core_integration_server_start
```

To cleanup after testing (deletes the database!):

```bash
make core_integration_server_stop core_integration_db_down
```

The test binary has the race detector enabled. `core_core_integration_server_stop` checks for any race logs reported by Go and will print them along a `66` exit code when found. Note that the actual race condition may have happened anywhere during the server lifetime, including start, stop or serving gRPC requests during tests.

### Run Local End-to-End Tests

To test the whole system, including the console UI and the login UI, run the E2E tests.

```bash
# Build the production docker image
export Zitadel_IMAGE=zitadel:local GOOS=linux
make docker_image

# If you made changes in the e2e directory, make sure you reformat the files
pnpm turbo lint:fix --filter=e2e

# Run the tests
docker compose --file ./e2e/docker-compose.yaml run --service-ports e2e
```

When you are happy with your changes, you can cleanup your environment.

```bash
# Stop and remove the docker containers for zitadel and the database
docker compose --file ./e2e/docker-compose.yaml down
```

### Run Local End-to-End Tests Against Your Dev Server Console

If you also make [changes to the console](#console), you can run the test suite against your locally built backend code and frontend server.

```bash
# Install dependencies (from repository root)
pnpm install

# Run the tests interactively
pnpm run open:golangangular

# Run the tests non-interactively
pnpm run e2e:golangangular
```

When you are happy with your changes, you can cleanup your environment.

```bash
# Stop and remove the docker containers for zitadel and the database
docker compose --file ./e2e/docker-compose.yaml down
```

## Contribute Frontend Code

This repository uses **pnpm** as package manager and **Turbo** for build orchestration.
All frontend packages are managed as a monorepo with shared dependencies and optimized builds:

- [apps/login](contribute-login) (depends on packages/zitadel-client and packages/zitadel-proto)
- apps/login/integration
- apps/login/acceptance
- [console](contribute-console) (depends on packages/zitadel-client)
- packages/zitadel-client
- packages/zitadel-proto
- [docs](contribute-docs)

### <a name="frontend-dev-requirements"></a>Frontend Development Requirements

The frontend components are run in a [Node](https://nodejs.org/en/about/) environment and are managed using the pnpm package manager and the Turborepo orchestrator.

> [!INFO]
> Some [dev containers are available](dev-containers) for remote development with docker and pipeline debugging in isolated environments.
> If you don't want to use one of the dev containers, you can develop the frontend components directly on your local machine.
> To do so, proceed with installing the necessary dependencies.

We use **pnpm** as package manager and **Turbo** for build orchestration. Use angular-eslint/Prettier for linting/formatting.
VSCode users, check out [this ESLint extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [this Prettier extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) to fix lint and formatting issues during development.

The commands in this section are tested against the following software versions:

- [Docker version 20.10.17](https://docs.docker.com/engine/install/)
- [Node version v20.x](https://nodejs.org/en/download/)
- [pnpm version 9.x](https://pnpm.io/installation)

To run tests with Cypress, ensure you have installed the required [Cypress runtime dependencies](https://docs.cypress.io/guides/continuous-integration/introduction#Dependencies)

<details>
  <summary>Note for WSL2 on Windows 10</summary>
  Following the suggestions <a href="https://stackoverflow.com/questions/62641553/setup-cypress-on-wsl-ubuntu-for-windows-10">here </a> subsequently <a href="https://github.com/microsoft/WSL/issues/4106">here </a> may  need to XLaunch and configure your DISPLAY variable. Use at your own risk.

1. Install `VcXsrv Windows X Server`
2. Set the target of your shortcut to `"C:\Program Files\VcXsrv\xlaunch.exe" -ac`
3. In WSL2 run `export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0` to set your DISPLAY variable
4. When starting XLaunch, make sure to disable access control
</details>

### <a name="contribute-login"></a>Contribute to Login

The Login UI is a Next.js application that provides the user interface for authentication flows.
It's located in the `apps/login` directory and uses pnpm and Turbo for development.

To start developing the login, make sure your system has the [required system dependencies](frontend-dev-requirements) installed.

#### Development Setup

```bash
# Start from the root of the repository
# Start the database and Zitadel backend
docker compose --file ./apps/login/acceptance/docker-compose.yaml up --detach zitadel

# Install dependencies
pnpm install

# Option 1: Run login development server with Turbo (recommended)
pnpm turbo dev --filter=@zitadel/login

# Option 2: Build and serve login (production build)
pnpm turbo build --filter=@zitadel/login
cd ./login && pnpm start
```

The login UI is available at http://localhost:3000.

#### Login Architecture

The login application consists of multiple packages:

- `@zitadel/login` - Main Next.js application
- `@zitadel/client` - TypeScript client library for Zitadel APIs
- `@zitadel/proto` - Protocol buffer definitions and generated code

The build process uses Turbo to orchestrate dependencies:

1. Proto generation (`@zitadel/proto#generate`)
2. Client library build (`@zitadel/client#build`)
3. Login application build (`@zitadel/login#build`)

#### Pass Quality Checks

Reproduce the pipelines linting and testing for the login.

```bash
pnpm turbo quality --filter=./apps/login/* --filter=./packages/*
```

Fix the [quality checks](troubleshoot-frontend), add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

### <a name="contribute-console"></a>Contribute to Console

To start developing the console, make sure your system has the [required system dependencies](frontend-dev-requirements) installed.
Then, you need to decide which Zitadel instance you would like to target.
- The easiest starting point is to [configure your environment](console-dev-existing-zitadel) to use a [Zitadel cloud](https://zitadel.com) instance.
- Alternatively, you can [start a local Zitadel instance from scratch and develop against it](console-dev-local-zitadel).

#### <a name="console-dev-existing-zitadel"></a>Develop against an already running Zitadel instance

By default, `pnpm dev --filter=console` targets a Zitadel API running at http://localhost:8080.
To change this, export the link to your environment.json in your environment variables.

```bash
export ENVIRONMENT_JSON_URL=https://my-cloud-instance-abcdef.us1.zitadel.cloud/ui/console/assets/environment.json
```

Proceed [with configuring your console redirect URIs](console-redirect).

#### <a name="console-dev-local-zitadel"></a>Develop against a local Zitadel instance from scratch

By executing the commands from this section, you run everything you need to develop the console locally.
Using [Docker Compose](https://docs.docker.com/compose/), you run [PostgreSQL](https://www.postgresql.org/download/) and the [latest release of Zitadel](https://github.com/zitadel/zitadel/releases/latest) on your local machine.
You use the Zitadel container as backend for your console.

Run the database and the latest backend locally.

```bash
# Start from the root of the repository
# You just need the db and the zitadel services to develop the console against.
docker compose --file ./e2e/docker-compose.yaml up --detach zitadel
```

When Zitadel accepts traffic, navigate to http://localhost:8080/ui/console/projects?login_hint=zitadel-admin@zitadel.localhost and log in with  _Password1!_.

Proceed [with configuring your console redirect URIs](console-redirect).

#### <a name="console-redirect"></a> Configure Console redirect URI

To allow console access via http://localhost:4200, you have to configure the Zitadel backend.

1. Navigate to /ui/console/projects in your target Zitadel instance.
3. Select the _Zitadel_ project.
4. Select the _Console_ application.
5. Select _Redirect Settings_
6. Add _http://<span because="breaks the link"></span>localhost:4200/auth/callback_ to the _Redirect URIs_
7. Add _http://<span because="breaks the link"></span>localhost:4200/signedout_ to the _Post Logout URIs_
8. Select the _Save_ button

#### Develop

Run the local console development server.

```bash
# Install dependencies (from repository root)
pnpm install

# Option 1: Run console development server with live reloading and dependency rebuilds
pnpm turbo dev --filter=console

# Option 2: Build and serve console (production build)
pnpm turbo build --filter=console
pnpm turbo serve --filter=console
```

Navigate to http://localhost:4200/.
Make some changes to the source code and see how the browser is automatically updated.

#### Pass Quality Checks

Reproduce the pipelines linting and testing for the console.

```bash
pnpm turbo quality --filter=console --filter=e2e
```

Fix the [quality checks](troubleshoot-frontend), add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

### <a name="contribute-docs"></a>Contribute to Docs

Project documentation is made with Docusaurus and is located under [./docs](./docs). The documentation uses **pnpm** and **Turbo** for development and build processes.

#### Local Development

```bash
# Install dependencies (from repository root)
pnpm install

# Option 1: Run docs development server with Turbo (recommended)
pnpm turbo dev --filter=zitadel-docs

# Option 2: Build and serve docs (production build)
pnpm turbo build --filter=zitadel-docs
cd ./docs && pnpm serve
```

The docs build process automatically:

1. Downloads required protoc plugins
2. Generates gRPC documentation from proto files
3. Generates API documentation from OpenAPI specs
4. Copies configuration files
5. Builds the Docusaurus site

#### Local testing

The documentation server will be available at http://localhost:3000 with live reload for fast development feedback.

#### Style guide

- **Code with variables**: Make sure that code snippets can be used by setting environment variables, instead of manually replacing a placeholder.
- **Embedded files**: When embedding mdx files, make sure the template ist prefixed by "\_" (lowdash). The content will be rendered inside the parent page, but is not accessible individually (eg, by search).
- **Don't repeat yourself**: When using the same content in multiple places, save and manage the content as separate file and make use of embedded files to import it into other docs pages.
- **Embedded code**: You can embed code snippets from a repository. See the [plugin](https://github.com/saucelabs/docusaurus-theme-github-codeblock#usage) for usage.

Following the [Google style guide](https://developers.google.com/style) is highly recommended. Its clear and concise guidelines ensure consistency and effective communication within the wider developer community.

The style guide covers a lot of material, so their [highlights](https://developers.google.com/style/highlights) page provides an overview of its most important points. Some of the points stated in the highlights that we care about most are given below:

- Be conversational and friendly without being frivolous.
- Use sentence case for document titles and section headings.
- Use active voice: make clear who's performing the action.
- Use descriptive link text.

#### Docs pull request

When making a pull request use `docs(<scope>): <short summary>` as title for the semantic release.
Scope can be left empty (omit the brackets) or refer to the top navigation sections.

#### Pass Quality Checks

Reproduce the pipelines linting checks for the docs.

```bash
pnpm turbo quality --filter=docs
```

Fix the [quality checks](troubleshoot-frontend), add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

### <a name="troubleshoot-frontend"></a>Troubleshoot Frontend Quality Checks

To debug and fix failing tasks, execute them individually using the `--filter` flag.

We recommend to use [one of the dev containers](dev-containers) to reproduce pipeline issues.

```bash
# to reproduce linting error in the console:
pnpm lint --filter=console
# To fix them:
pnpm lint:fix --filter=console
```

More tasks that are runnable on-demand.
Some tasks have variants like `pnpm test:e2e:angulargolang`,
others support arguments and flags like `pnpm test:integration run --spec apps/login/integration/integration/login.cy.ts`.
For the turbo commands, check your options with `pnpm turbo --help`

| Command                   | Description                                              | Example                                                                                                                                                    |
| ------------------------- | -------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `pnpm turbo run generate` | Generate stubs from Proto files                          | Generate API docs: `pnpm turbo run generate --filter zitadel-docs`                                                                                         |
| `pnpm turbo build`        | Build runnable JavaScript code                           | Regenerate the proto stubs and build the @zitadel/client package: `pnpm turbo build --filter @zitadel/client`                                              |
| `pnpm turbo quality`      | Reproduce the pipeline quality checks                    | Run login-related quality checks `pnpm turbo quality --filter './apps/login/*' --filter './packages/*'`                                                    |
| `pnpm turbo lint`         | Check linting issues                                     | Check login-related linting issues for differences with main `pnpm turbo lint --filter=[main...HEAD] --filter .'/apps/login/**/*' --filter './packages/*'` |
| `pnpm turbo lint:fix`     | Fix linting issues                                       | Fix console-relevant linting issues `pnpm turbo lint:fix --filter console --filter './packages/*' --filter zitadel-e2e`                                    |
| `pnpm turbo test:unit`    | Run unit tests. Rerun on file changes                    | Run unit tests in all packages in and watch for file changes `pnpm turbo watch test:unit`                                                                  |
| `pnpm turbo test:e2e`     | Run the Cypress CLI for console e2e tests                | Test interactively against the console in a local dev server and Zitadel in a container: `pnpm turbo test:e2e:angular open`                                |
| `pnpm turbo down`         | Remove containers and volumes                            | Shut down containers from the integration test setup `pnpm turbo down`                                                                                     |
| `pnpm turbo clean`        | Remove downloaded dependencies and other generated files | Remove generated docs  `pnpm turbo clean --filter zitadel-docs`                                                                                            |

## <a name="dev-containers"></>Developing Zitadel with Dev Containers

You can use dev containers if you'd like to make sure you have the same development environment like the corresponding GitHub PR checks use.
The following dev containers are available:

- **.devcontainer/base/devcontainer.json**: Contains everything you need to run whatever you want.
- **.devcontainer/turbo-lint-unit/devcontainer.json**: Runs a dev container that executes frontent linting and unit tests and then exits. This is useful to reproduce the corresponding GitHub PR check. 
- **.devcontainer/turbo-lint-unit-debug/devcontainer.json**: Runs a dev container that executes frontent linting and unit tests in watch mode. You can fix the errors right away and have immediate feedback.
- **.devcontainer/login-integration/devcontainer.json**: Runs a dev container that executes login integration tests and then exits. This is useful to reproduce the corresponding GitHub PR check.
- **.devcontainer/login-integration-debug/devcontainer.json**: Runs a dev container that spins up the login in a hot-reloading dev server and executes login integration tests interactively. You can fix the errors right away and have immediate feedback.

You can also run the GitHub PR checks locally in dev containers without having to connect to a dev container.
 

The following pnpm commands use the [devcontainer CLI](https://github.com/devcontainers/cli/) and exit when the checks are done.
The minimal system requirements are having Docker and the devcontainers CLI installed.
If you don't have the node_modules installed already, you need to install the devcontainers CLI manually. Run `npm i -g @devcontainers/cli@0.80.0`. Alternatively, the [official Microsoft VS Code extension for Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) offers a command `Dev Containers: Install devcontainer CLI`


```bash
npm run devcontainer:lint-unit
npm run devcontainer:integration:login
```

If you don't have NPM installed, copy and execute the scripts from the package.json directly.

To connect to a dev container to have full IDE support, follow the instructions provided by your code editor/IDE to initiate the dev container.
This typically involves opening the "Command Palette" or similar functionality and searching for commands related to "Dev Containers" or "Remote Containers".
The quick start guide for VS Code can found [here](https://code.visualstudio.com/docs/devcontainers/containers#_quick-start-open-an-existing-folder-in-a-container)

For example, to build and run the Zitadel binary in a dev container, connect your IDE to the dev container described in .devcontainer/base/devcontainer.json.
Run the following commands inside the container to start Zitadel.

```bash
make compile && ./zitadel start-from-init --masterkey MasterkeyNeedsToHave32Characters --tlsMode disabled
```

Zitadel serves traffic as soon as you can see the following log line:

`INFO[0001] server is listening on [::]:8080`


## <a name="contribute-translations"></a>Contribute Translations

Zitadel loads translations from four files:

- [Console texts](./console/src/assets/i18n)
- [Login interface](./internal/api/ui/login/static/i18n)
- [Email notification](./internal/notification/static/i18n)
- [Common texts](./internal/static/i18n) for success or error toasts

You may edit the texts in these files or create a new file for additional language support. Make sure you set the locale (ISO 639-1 code) as the name of the new language file.
Please make sure that the languages within the files remain in their own language, e.g. German must always be `Deutsch.
If you have added support for a new language, please also ensure that it is added in the list of languages in all the other language files.

You also have to add some changes to the following files:

- [Register Local File](./console/src/app/app.module.ts)
- [Add Supported Language](./console/src/app/utils/language.ts)
- [Customized Text Docs](./docs/docs/guides/manage/customize/texts.md)
- [Add language option](./internal/api/ui/login/static/templates/external_not_found_option.html)

## **Did you find a security flaw?**

- Please read [Security Policy](./SECURITY.md).

## Product management

The Zitadel Team works with an agile product management methodology.
You can find all the issues prioritized and ordered in the [product board](https://github.com/orgs/zitadel/projects/2/views/1).

### Sprint

We want to deliver a new release every second week. So we plan everything in two-week sprints.
Each Tuesday we estimate new issues and on Wednesday the last sprint will be reviewed and the next one will be planned.
After a sprint ends a new version of Zitadel will be released, and publish to [Zitadel Cloud](https://zitadel.cloud) the following Monday.

If there are some critical or urgent issues we will have a look at it earlier, than the two weeks.
To show the community the needed information, each issue gets attributes and labels.

### About the attributes

You can find the attributes on the project "Product Management".

#### State

The state should reflect the progress of the issue and what is going on right now.

- **No status**: Issue just got added and has to be looked at.
- **🧐 Investigating**: We are currently investigating to find out what the problem is, which priority it should have and what has to be implemented. Or we need some more information from the author.
- **📨 Product backlog**: If an issue is in the backlog, it is not currently being worked on. These are recorded so that they can be worked on in the future. Issues with this state do not have to be completely defined yet.
- **📝 Prioritized product backlog**: An issue with the state "Prioritized Backlog" is ready for the refinement from the perspective of the product owner (PO) to implement. This means the developer can find all the relevant information and acceptance criteria in the issue.
- **🔖 Ready**: The issue is ready to take into a sprint. Difference to "prioritized..." is that the complexity is defined by the team.
- **📋 Sprint backlog**: The issue is scheduled for the current sprint.
- **🏗 In progress**: Someone is working on this issue right now. The issue will get an assignee as soon as it is in progress.
- **❌ Blocked**: The issue is blocked until another issue is resolved/done.
- **👀 In review**: The issue is in review. Please add someone to review your issue or let us know that it is ready to review with a comment on your pull request.
- **✅ Done**: The issue is implemented and merged to main.

#### Priority

Priority shows you the priority the Zitadel team has given this issue. In general the higher the demand from customers and community for the feature, the higher the priority.

- **🌋 Critical**: This is a security issue or something that has to be fixed urgently, because the software is not usable or highly vulnerable.
- **🏔 High**: These are the issues the Zitadel team is currently focusing on and will be implemented as soon as possible.
- **🏕 Medium**: After all the high issues are done these will be next.
- **🏝 Low**: This is low in priority and will probably not be implemented in the next time or just if someone has some time in between.

#### Complexity

This should give you an indication how complex the issue is. It's not about the hours or effort it takes.
Everything that is higher than 8 should be split in smaller parts.

**1**, **2**, **3**, **5**, **8**, **13**

### About the labels

There are a few general labels that don't belong to a specific category.

- **good first issue**: This label shows contributors, that it is an easy entry point to start developing on Zitadel.
- **help wanted**: The author is seeking help on this topic, this may be from an internal Zitadel team member or external contributors.

#### Category

The category shows which part of Zitadel is affected.

- **category: backend**: The backend includes the APIs, event store, command and query side. This is developed in golang.
- **category: ci**: ci is all about continues integration and pipelines.
- **category: design**: All about the ux/ui of Zitadel
- **category: docs**: Adjustments or new documentations, this can be found in the docs folder.
- **category: frontend**: The frontend concerns on the one hand the Zitadel management console (Angular) and on the other hand the login (gohtml)
- **category: infra**: Infrastructure does include many different parts. E.g Terraform-provider, docker, metrics, etc.
- **category: translation**: Everything concerning translations or new languages

#### Language

The language shows you in which programming language the affected part is written

- **lang: angular**
- **lang: go**
- **lang: javascript**
