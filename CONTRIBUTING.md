# Contributing to Zitadel

Zitadel is an open-source identity and access management platform built with a modern tech stack including Go (API), Next.js/React (Login), Angular (Console), and Docusaurus (Docs) - all orchestrated through an Nx monorepo with pnpm for efficient development workflows.

## Quick Start

1. Clone the repository: `git clone https://github.com/zitadel/zitadel` or [open it in a local Dev Container](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/zitadel/zitadel) or [create a GitHub Codespace](https://codespaces.new/zitadel/zitadel)
2. If you cloned the repository to your local machine, install the required development dependencies
   - [Node.js v22.x](https://nodejs.org/en/download/) - Required for UI development and to run development commands `pnpm nx ...`
   - [Go 1.24.x](https://go.dev/doc/install) - Required for API development
   - [Docker](https://docs.docker.com/engine/install/) - Required for supporting services like the development database and for tests.
   - [Cypress runtime dependencies](https://docs.cypress.io/guides/continuous-integration/introduction#Dependencies) - Required for Browser UI tests
   <details>
      <summary>WSL2 on Windows 10 users (click to expand)</summary>
      
      For Cypress tests on WSL2, you may need to configure X11 forwarding. Following suggestions [here](https://stackoverflow.com/questions/62641553/setup-cypress-on-wsl-ubuntu-for-windows-10) and [here](https://github.com/microsoft/WSL/issues/4106). Use at your own risk.

      1. Install `VcXsrv Windows X Server`
      2. Set shortcut target to `"C:\Program Files\VcXsrv\xlaunch.exe" -ac`
      3. In WSL2: `export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0`
      4. Disable access control when starting XLaunch
   </details>
3. Use [Corepack](https://pnpm.io/installation#using-corepack) to make sure you have [pnpm](https://pnpm.io/) installed in the correct version: `corepack enable`.
4. Install node module dependencies: `pnpm install`
5. Generate code `pnpm nx run-many --target generate`
6. Optionally, install the following VSCode plugins:
   - [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go) - For API development. Use golangci-lint v2 as linter.
   - [Angular Language Service](https://marketplace.visualstudio.com/items?itemName=Angular.ng-template) - For Console development
   - [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) - Code linting
   - [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) - Code formatting
   - [Nx Console](https://marketplace.visualstudio.com/items?itemName=nrwl.angular-console) - Nx task runner tooling

Jump to the dedicated sections for developing a specific project:

- [Contributing to the API](#contribute-to-api)
- [Contributing to the Login](#contribute-to-login)
- [Contributing to the Console](#contribute-to-console)
- [Contributing to the Docs](#contribute-to-docs)
- [Contributing translations](#contribute-translations)

## Development Commands Cheat Sheet

This repository contains multiple interconnected projects.
You can build and start any project with Nx commands.

| Task                          | Command                                     | Notes                        | Details                                                                                                                                                                                          |
| ----------------------------- | ------------------------------------------- | ---------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Production**                | `pnpm nx run PROJECT:prod`                  | Production server            |                                                                                                                                                                                                  |
| **Develop**                   | `pnpm nx run PROJECT:dev`                   | Development server           |                                                                                                                                                                                                  |
| **Generate**                  | `pnpm nx run PROJECT:generate`              | Generate .gitignored files   |                                                                                                                                                                                                  |
| **Generate Go Files**         | `pnpm nx run @zitadel/api:generate-go`      | Regenerate checked-in files  | This is needed to generate files using [Stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer), [Enumer](https://github.com/dmarkham/enumer) or [gomock](https://github.com/uber-go/mock) |
| **Test - Unit**               | `pnpm nx run PROJECT:test-unit`             | Run unit tests               |                                                                                                                                                                                                  |
| **Test - Integration**        | `pnpm nx run PROJECT:test-integration`      | Run integration tests        | Learn mnore about how to [debug API integration tests](#run-api-integration-tests)                                                                                                               |
| **Test - Integration Stop**   | `pnpm nx run PROJECT:test-integration-stop` | Stop integration containers  |                                                                                                                                                                                                  |
| **Test - Functional UI**      | `pnpm nx run @zitadel/functional-ui:test`   | Run functional UI tests      | Learn more about how to [develop the Console and opening the interactive Test Suite](#pass-console-quality-checks)                                                                               |
| **Test - Functional UI Stop** | `pnpm nx run @zitadel/functional-ui:stop`   | Run functional UI containers |                                                                                                                                                                                                  |
| **Test**                      | `pnpm nx run PROJECT:test`                  | Run all tests                |                                                                                                                                                                                                  |
| **Lint**                      | `pnpm nx run PROJECT:lint`                  | Check code style             |                                                                                                                                                                                                  |
| **Lint Fix**                  | `pnpm nx run PROJECT:lint-fix`              | Auto-fix style issues        |                                                                                                                                                                                                  |

Replace `PROJECT` with one of the following:

- `@zitadel/zitadel` (you can omit this root level project when using `pnpm nx run`, like `pnpm nx run db`)
- `@zitadel/api`
- `@zitadel/login`
- `@zitadel/console`
- `@zitadel/docs`
- `@zitadel/client`
- `@zitadel/proto`

Instead of the project names, you can also use their directory names for `PROJECT`, like `pnpm nx run login:dev`.
Alternatively, you can use the infix-notation, like `pnpm nx dev @zitadel/login` or `pnpm nx dev login`.
To stream all logs instead of opening the interactive terminal, disable the TUI with `pnpm nx --tui false ...`.
If a command is stuck because a process is already running, stop the Nx daemon and try again: `pnpm nx daemon --stop`.

## Introduction

Thank you for your interest in contributing! As you might know there is more than code to contribute. You can find all information needed to start contributing here.

Please give us and our community the chance to get rid of security vulnerabilities by responsibly disclosing these issues to [security@zitadel.com](mailto:security@zitadel.com).

The strongest part of a community is the possibility to share thoughts. That's why we try to react as soon as possible to your ideas, thoughts and feedback. We love to discuss as much as possible in an open space like in the [issues](https://github.com/zitadel/zitadel/issues) and [discussions](https://github.com/zitadel/zitadel/discussions) section here or in our [chat](https://zitadel.com/chat), but we understand your doubts and provide further contact options [here](https://zitadel.com/contact).

If you want to give an answer or be part of discussions please be kind. Treat others like you want to be treated. Read more about our code of conduct [here](CODE_OF_CONDUCT.md).

## What can I contribute?

For people who are new to Zitadel: We flag issues which are a good starting point to start contributing.
You can find them [here](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).
We add the label "good first issue" for problems we think are a good starting point to contribute to Zitadel.

- [Issues for first time contributors](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [All issues](https://github.com/zitadel/zitadel/issues)

Help shape the future of Zitadel:

- Join our [chat](https://zitadel.com/chat) and discuss with us or others.
- Ask or answer questions in the [issues section](https://github.com/zitadel/zitadel/issues)
- Share your thoughts and ideas in the [discussions section](https://github.com/zitadel/zitadel/discussions)

Make Zitadel more popular and give it a ⭐

Follow [@zitadel](https://twitter.com/zitadel) on twitter

[Contribute](#how-to-contribute)

- [Contribute API code](#contribute-to-api)
- [Contribute frontend code](#contribute-to-frontend)
- If you found a mistake on our [Docs page](https://zitadel.com/docs) or something is missing please read [the Docs section](#contribute-to-docs)
- [Translate](#contribute-translations) and improve texts

## How to contribute

We strongly recommend [talking to us](https://zitadel.com/contact) before you start contributing to streamline your work with ours.

We accept contributions through pull requests.
You need a github account for that.
If you are unfamiliar with git have a look at Github's documentation on [creating forks](https://help.github.com/articles/fork-a-repo) and [creating pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork).
Please draft the pull request as soon as possible.
Go through the following checklist before you submit the final pull request:

### Components

The code consists of the following parts:

| name               | description                                        | language                                                                                                  | where to find                                       | Development Guide                                   |
| ------------------ | -------------------------------------------------- | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------- | --------------------------------------------------- |
| API implementation | Service that serves the grpc(-web) and RESTful API | [go](https://go.dev)                                                                                      | [API implementation](./internal/api/grpc)           | [Contribute to API](#contribute-to-api)             |
| API definitions    | Specifications of the API                          | [Protobuf](https://developers.google.com/protocol-buffers)                                                | [./proto/zitadel](./proto/zitadel)                  | [Contribute to API](#contribute-to-api)             |
| Console            | Frontend the user interacts with after log in      | [Angular](https://angular.io), [Typescript](https://www.typescriptlang.org)                               | [./console](./console)                              | [Contribute to Frontend](#contribute-to-frontend)   |
| Login              | Modern authentication UI built with Next.js        | [Next.js](https://nextjs.org), [React](https://reactjs.org), [TypeScript](https://www.typescriptlang.org) | [./apps/login](./apps/login)                        | [Contribute to Frontend](#contribute-to-frontend)   |
| Docs               | Project documentation made with docusaurus         | [Docusaurus](https://docusaurus.io/)                                                                      | [./docs](./docs)                                    | [Contribute to Frontend](#contribute-to-frontend)   |
| translations       | Internationalization files for default languages   | YAML                                                                                                      | [./console](./console) and [./internal](./internal) | [Contribute Translations](#contribute-translations) |

Please follow the guides to validate and test the code before you contribute.

### Submitting a pull request (PR)

1. [Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the [zitadel/zitadel](https://github.com/zitadel/zitadel) repository on GitHub
2. On your fork, commit your changes to a new branch

   `git checkout -b my-fix-branch main`

3. Make your changes following the [guidelines](#how-to-contribute) in this guide. Make sure that all tests pass.

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

This is optional to indicate which component is affected. If in doubt, leave blank (`<type>: <short summary>`)

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

## Contribute to API

To start developing, make sure you followed the [quick start](#quick-start) steps.

### Develop the API

Optionally build the Console

```bash
pnpm nx run @zitadel/api:build-console
```

Optionally start the Login in another terminal

```bash
pnpm nx run @zitadel/login:prod
```

Run the local development database.

```bash
pnpm nx db
```

Start a debug session in your IDE.
For example, in VSCode, you can use a `launch.json` configuration like this.

```json
   {
      "name": "Debug Zitadel API",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "env": {
            "ZITADEL_DATABASE_POSTGRES_HOST": "${env:DEVCONTAINER_DB_HOST}"
      },
      "program": "main.go",
      "args": [
            "start-from-init",
            "--config",
            "${workspaceFolder}/apps/api/prod-default.yaml",
            "--steps",
            "${workspaceFolder}/apps/api/prod-default.yaml",
            "--masterkey",
            "MasterkeyNeedsToHave32Characters"
      ]
   }
```

If you have built the Console and started the Login, visit http://localhost:8080/ui/console?login_hint=zitadel-admin@zitadel.localhost and enter `Password1!` to log in.

Call the API using the generated [](./admin.pat) with [grpcurl](https://github.com/fullstorydev/grpcurl) or [grpcui](https://github.com/fullstorydev/grpcui), for example:

```bash
grpcurl -plaintext -H "Authorization: Bearer $(cat admin.pat)" localhost:8080 zitadel.user.v2.UserService.ListUsers
```

To connect to the database and explore Zitadel data, run `psql "host=${DEVCONTAINER_DB_HOST:-localhost} dbname=zitadel sslmode=disable"`.

### Run API Unit Tests

To test the code without dependencies, run the unit tests:

```bash
pnpm nx run @zitadel/api:test-unit
```

### Run API Integration Tests

API tests are run as gRPC clients against a running Zitadel server binary.
The server binary is [built with coverage enabled](https://go.dev/doc/build-cover).


```bash
pnpm nx run @zitadel/api:test-integration
```

To develop and run the test cases from within your IDE or by the command line, start only the database and the API.
The actual integration test clients reside in the `integration_test` subdirectory of the package they aim to test.
Integration test files use the `integration` build tag, in order to be excluded from regular unit tests.
Because of the server-client split, Go is usually unaware of changes in server code and tends to cache test results.
Pass `-count 1` to disable test caching.

Start the ephemeral database for integration tests.

```bash
pnpm nx run @zitadel/api:test-integration-run-db
```

In another terminal, start the API.

```bash
pnpm nx run @zitadel/api:test-integration-run-api
```

Example command to run a single package integration test:

```bash
go test -count 1 -tags integration ./internal/api/grpc/management/integration_test
```

To run all available integration tests:

```bash
go test -count 1 -tags integration -parallel 1 $(go list -tags integration ./... | grep -e \"integration_test\" -e \"events_testing\")
```

It is also possible to run the API in a debugger and run the integrations tests against it.

First, start the ephemeral database for integration tests.

```bash
pnpm nx run @zitadel/api:test-integration-run-db
```

When starting the debugger, make sure the Zitadel binary starts with `start-from-init --config=./apps/api/test-integration-api.yaml --steps=./apps/api/test-integration-api.yaml --masterkey=MasterkeyNeedsToHave32Characters"`

To cleanup after testing (deletes the ephemeral database!):

```bash
pnpm nx run @zitadel/devcontainer:compose down db-api-integration cache-api-integration
```

### Run Functional UI Tests

To test the whole system, including the Console UI and the Login UI, run the Functional UI tests.

```bash
# If you made changes in the tests/functional-ui directory, make sure you reformat the files
pnpm nx run @zitadel/functional-ui:lint-fix

# Run the tests
pnpm nx run @zitadel/functional-ui:test
```

## Contribute Frontend Code

This repository uses **pnpm** as package manager and **Nx** for build orchestration.

### Project Overview

Choose your contribution area:

- **[Login App](#contribute-to-login)** (Next.js/React) - Modern authentication flows
- **[Console](#contribute-to-console)** (Angular) - Admin dashboard and user management
- **[Docs](#contribute-to-docs)** (Docusaurus) - Project documentation
- **[Client Packages](#client-packages)** - Shared libraries for API communication

### Project Dependencies

```
apps/login → packages/zitadel-client → packages/zitadel-proto
console → packages/zitadel-client → packages/zitadel-proto
docs → (independent)
```

**Nx handles this automatically** - when you change `zitadel-proto`, Nx rebuilds dependent projects.

### Contribute to Login

The Login UI is a Next.js application that provides the user interface for authentication flows.
It is MIT-licensed, so you are free to change and deploy it as you like.
It's located in the `apps/login` directory and uses pnpm and Nx for development.
Get familiar with the [Login ui docs](https://zitadel.com/docs/guides/integrate/login-ui).

To start developing, make sure you followed the [quick start](#quick-start) steps.

#### Develop the Login against a local API

Run the local development database.

```bash
pnpm nx db
```

In another terminal, start the API

```bash
pnpm nx run @zitadel/api:prod
```

In another terminal, start the Login development server

```bash
pnpm nx run @zitadel/login:dev
```

Visit http://localhost:8080/ui/console?login_hint=zitadel-admin@zitadel.localhost and enter `Password1!` to log in.

Make some changes to the source code and see how the browser is automatically updated.

#### Develop against a Cloud instance

If you don't want to build and run a local API, you can just run the Login development server and point it to a cloud instance.

1. Create a personal access token and point your instance to your local Login, [as described in the Docs](https://zitadel.com/docs/self-hosting/manage/login-client).
2. Save the following file to `apps/login/.env.dev.local`

```env
ZITADEL_API_URL=https://[your-cloud-instance-domain]
ZITADEL_SERVICE_USER_TOKEN=[personal access token for an IAM Login Client]
```

3. Start the development server.

```bash
pnpm nx run @zitadel/login:dev
```

Visit http://localhost:8080/ui/console?login_hint=zitadel-admin@zitadel.localhost and enter `Password1!` to log in.

#### Login Architecture

The Login application consists of multiple packages:

- `@zitadel/login` - Main Next.js application
- `@zitadel/client` - TypeScript client library for Zitadel APIs
- `@zitadel/proto` - Protocol buffer definitions and generated code

The build process uses Nx and pnpm to orchestrate dependencies:

#### Pass Login Quality Checks

Reproduce the pipeline quality checks for the code you changed.

```bash
# Run Login-related linting builds and unit tests
pnpm nx run-many --projects @zitadel/login @zitadel/client @zitadel/proto --targets lint build test
```

Fix the quality checks, add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

#### <a name="login-deploy"></a>Deploy

- [![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fzitadel%2Fzitadel&env=ZITADEL_API_URL,ZITADEL_SERVICE_USER_ID,ZITADEL_SERVICE_USER_TOKEN&root-directory=apps/login&envDescription=Setup%20a%20service%20account%20with%20IAM_LOGIN_CLIENT%20membership%20on%20your%20instance%20and%20provide%20its%20id%20and%20personal%20access%20token.&project-name=zitadel-login&repository-name=zitadel-login)
- Build and deploy with Docker: `pnpm nx run @zitadel/login:build && docker build -t my-zitadel-login apps/login`
- Build and deploy with NodeJS: `pnpm nx run @zitadel/login:prod`

### Contribute to Console

To learn more about the Console, go to the Consoles [README.md](./console/README.md).

To start developing, make sure you followed the [quick start](#quick-start) steps.

#### Develop the Console against a local API

Run the local development database.

```bash
pnpm nx db
```

In another terminal, start the API

```bash
pnpm nx run @zitadel/api:prod
```

In another terminal, start the Login

```bash
pnpm nx run @zitadel/login:prod
```

Allow the API [to redirect to your dev server](#configure-console-dev-server-redirects).

In another terminal, start the Console development server

```bash
pnpm nx run @zitadel/console:dev
```

Visit http://localhost:4200/?login_hint=zitadel-admin@zitadel.localhost and enter `Password1!` to log in.

Make some changes to the source code and see how the browser is automatically updated.

#### Develop against a Cloud instance

If you don't want to build and run a local API, you can just run the console development server and point it to a cloud instance.

Save the following file to console/.env.local

```env
ENVIRONMENT_JSON_URL=https://[your-cloud-instance-domain]/ui/console/assets/environment.json
```

Start the development server.

```bash
pnpm nx run @zitadel/console:dev
```

Allow the API [to redirect to your dev server](#configure-console-dev-server-redirects).

Visit http://localhost:4200/?login_hint=zitadel-admin@zitadel.localhost and enter `Password1!` to log in.

#### Configure Console Dev Server Redirects

To allow Console access via http://localhost:4200, you have to configure the Zitadel API.

1. Navigate to http://localhost:8080/ui/console/projects.
2. Select the _ZITADEL_ project.
3. Select the _Console_ application.
4. Select _Redirect Settings_
5. Add _http://<span because="breaks the link"></span>localhost:4200/auth/callback_ to the _Redirect URIs_
6. Add _http://<span because="breaks the link"></span>localhost:4200/signedout_ to the _Post Logout URIs_
7. Select the _Save_ button

#### Pass Console Quality Checks

Run the quality checks for the code you changed.

```bash
# Run console-related linting builds and unit tests
pnpm nx run-many --projects @zitadel/console @zitadel/client @zitadel/proto @zitadel/functional-ui --targets lint build test
```

Run functional UI tests against a locally built API and a dev server Console.

Allow the API [to redirect to your dev server](#configure-console-dev-server-redirects).
Alternatively, create the file `tests/functional-ui/.env.open.local` with the following content:

```conf
CYPRESS_BASE_URL=http://localhost:8080/ui/console
```

```bash
# Run the API and the Console dev server
# Beware this doesn't work from within a dev container.
pnpm nx run @zitadel/functional-ui:open
```

Or run all tests to completion.

```bash
# Run the tests
pnpm nx run @zitadel/functional-ui:test
```

Fix the quality checks, add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

### Contribute to Client Packages

To start developing, make sure you followed the [quick start](#quick-start) steps.

**`@zitadel/proto`**: Protocol buffer definitions and generated TypeScript/JavaScript clients.

```bash
pnpm nx run @zitadel/proto:generate  # Regenerate after proto changes
```

**`@zitadel/client`**: High-level TypeScript client library with utilities for API interaction.

```bash
pnpm nx run @zitadel/client:build  # Build after changes
```

### Contribute to Docs

Project documentation is made with Docusaurus and is located under [./docs](./docs). The documentation uses **pnpm** and **Nx** for development and build processes.

To start developing, make sure you followed the [quick start](#quick-start) steps.

#### Local Development

```bash
# Start development server (recommended)
pnpm nx run @zitadel/docs:dev

# Or start production server
pnpm nx run @zitadel/docs:prod
```

The Docs build process automatically:

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

Verify the Docs build correctly.

```bash
pnpm nx run @zitadel/docs:build
```

Fix the quality checks, add new checks that cover your changes and mark your pull request as ready for review when the pipeline checks pass.

## Contribute Translations

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
- **category: ci**: ci is all about continuous integration and pipelines.
- **category: design**: All about the ux/ui of Zitadel
- **category: docs**: Adjustments or new documentations, this can be found in the docs folder.
- **category: frontend**: The frontend concerns on the one hand the Zitadel management Console (Angular) and on the other hand the Login (gohtml)
- **category: infra**: Infrastructure does include many different parts. E.g Terraform-provider, docker, metrics, etc.
- **category: translation**: Everything concerning translations or new languages

#### Language

The language shows you in which programming language the affected part is written

- **lang: angular**
- **lang: go**
- **lang: javascript**
