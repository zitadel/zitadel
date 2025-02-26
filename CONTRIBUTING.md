# Contributing to ZITADEL

## Introduction

Thank you for your interest about how to contribute! As you might know there is more than code to contribute. You can find all information needed to start contributing here.

Please give us and our community the chance to get rid of security vulnerabilities by responsibly disclose this kind of issues by contacting [security@zitadel.com](mailto:security@zitadel.com).

The strongest part of a community is the possibility to share thoughts. That's why we try to react as soon as possible to your ideas, thoughts and feedback. We love to discuss as much as possible in an open space like in the [issues](https://github.com/zitadel/zitadel/issues) and [discussions](https://github.com/zitadel/zitadel/discussions) section here or in our [chat](https://zitadel.com/chat), but we understand your doubts and provide further contact options [here](https://zitadel.com/contact).

If you want to give an answer or be part of discussions please be kind. Treat others like you want to be treated. Read more about our code of conduct [here](CODE_OF_CONDUCT.md).

## What can I contribute?

For people who are new to ZITADEL: We flag issues which are a good starting point to start contributing. You find them [here](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

Make ZITADEL more popular and give it a ⭐

Help shaping the future of ZITADEL:

- Join our [chat](https://zitadel.com/chat) and discuss with us or others.
- Ask or answer questions in the [issues section](https://github.com/zitadel/zitadel/issues)
- Share your thoughts and ideas in the [discussions section](https://github.com/zitadel/zitadel/discussions)

[Contribute](#how-to-contribute)

- [Contribute code](#contribute)
- If you found a mistake on our [docs page](https://zitadel.com/docs) or something is missing please read [the docs section](#contribute-docs)
- [Translate](#contribute-internationalization) and improve texts

Follow [@zitadel](https://twitter.com/zitadel) on twitter

## How to contribute

We strongly recommend to [talk to us](https://zitadel.com/contact) before you start contributing to streamline our and your work.

We accept contributions through pull requests.
You need a github account for that.
If you are unfamiliar with git have a look at Github's documentation on [creating forks](https://help.github.com/articles/fork-a-repo) and [creating pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork).
Please draft the pull request as soon as possible.
Go through the following checklist before you submit the final pull request:

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

## Contribute

The code consists of the following parts:

| name            | description                                                        | language                                                                    | where to find                                      |
| --------------- | ------------------------------------------------------------------ | --------------------------------------------------------------------------- | -------------------------------------------------- |
| backend         | Service that serves the grpc(-web) and RESTful API                 | [go](https://go.dev)                                                        | [API implementation](./internal/api/grpc)          |
| console         | Frontend the user interacts with after log in                      | [Angular](https://angular.io), [Typescript](https://www.typescriptlang.org) | [./console](./console)                             |
| login           | Server side rendered frontend the user interacts with during login | [go](https://go.dev), [go templates](https://pkg.go.dev/html/template)      | [./internal/api/ui/login](./internal/api/ui/login) |
| API definitions | Specifications of the API                                          | [Protobuf](https://developers.google.com/protocol-buffers)                  | [./proto/zitadel](./proto/zitadel)                 |
| docs            | Project documentation made with docusaurus                         | [Docusaurus](https://docusaurus.io/)                                        | [./docs](./docs)                                   |

Please validate and test the code before you contribute.

We add the label "good first issue" for problems we think are a good starting point to contribute to ZITADEL.

- [Issues for first time contributors](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [All issues](https://github.com/zitadel/zitadel/issues)

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

### Developing ZITADEL with Dev Containers

Follow the instructions provided by your code editor/IDE to initiate the development container. This typically involves opening the "Command Palette" or similar functionality and searching for commands related to "Dev Containers" or "Remote Containers". The quick start guide for VS Code can found [here](https://code.visualstudio.com/docs/devcontainers/containers#_quick-start-open-an-existing-folder-in-a-container)

When you are connected to the container run the following commands to start ZITADEL.

```bash
make compile && ./zitadel start-from-init --masterkey MasterkeyNeedsToHave32Characters --tlsMode disabled
```

ZITADEL serves traffic as soon as you can see the following log line:

`INFO[0001] server is listening on [::]:8080`

### Backend/login

By executing the commands from this section, you run everything you need to develop the ZITADEL backend locally.
Using [Docker Compose](https://docs.docker.com/compose/), you run a [CockroachDB](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker-mac.html) on your local machine.
With [make](https://www.gnu.org/software/make/), you build a debuggable ZITADEL binary and run it using [delve](https://github.com/go-delve/delve).
Then, you test your changes via the console your binary is serving at http://<span because="breaks the link"></span>localhost:8080 and by verifying the database.
Once you are happy with your changes, you run end-to-end tests and tear everything down.

ZITADEL uses [golangci-lint](https://golangci-lint.run) for code quality checks. Please use [this configuration](.golangci.yaml) when running `golangci-lint`. We recommend to set golangci-lint as linter in your IDE.

The commands in this section are tested against the following software versions:

- [Docker version 20.10.17](https://docs.docker.com/engine/install/)
- [Go version 1.22](https://go.dev/doc/install)
- [Delve 1.9.1](https://github.com/go-delve/delve/tree/v1.9.1/Documentation/installation)

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
You can test if ZITADEL does what you expect by using the UI at http://localhost:8080/ui/console.
Also, you can verify the data by running `cockroach sql --database zitadel --insecure` and running SQL queries.

#### Run Local Unit Tests

To test the code without dependencies, run the unit tests:

```bash
make core_unit_test
```

#### Run Local Integration Tests

Integration tests are run as gRPC clients against a running ZITADEL server binary.
The server binary is typically [build with coverage enabled](https://go.dev/doc/build-cover).
It is also possible to run a ZITADEL sever in a debugger and run the integrations tests like that. In order to run the server, a database is required.

The database flavor can **optionally** be set in the environment to `cockroach` or `postgres`. The default is `postgres`.

```bash
export INTEGRATION_DB_FLAVOR="cockroach"
```

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

When you change any ZITADEL server code, be sure to rebuild and restart the server before the next test run.

```bash
make core_integration_server_stop core_integration_server_start
```

To cleanup after testing (deletes the database!):

```bash
make core_integration_server_stop core_integration_db_down
```

The test binary has the race detector enabled. `core_core_integration_server_stop` checks for any race logs reported by Go and will print them along a `66` exit code when found. Note that the actual race condition may have happened anywhere during the server lifetime, including start, stop or serving gRPC requests during tests.

#### Run Local End-to-End Tests

To test the whole system, including the console UI and the login UI, run the E2E tests.

```bash
# Build the production docker image
export ZITADEL_IMAGE=zitadel:local GOOS=linux
make docker_image

# If you made changes in the e2e directory, make sure you reformat the files
make console_lint

# Run the tests
docker compose --file ./e2e/config/host.docker.internal/docker-compose.yaml run --service-ports e2e
```

When you are happy with your changes, you can cleanup your environment.

```bash
# Stop and remove the docker containers for zitadel and the database
docker compose --file ./e2e/config/host.docker.internal/docker-compose.yaml down
```

#### Run Local End-to-End Tests Against Your Dev Server Console

If you also make [changes to the console](#console), you can run the test suite against your locally built backend code and frontend server.
But you will have to install the relevant node dependencies.

```bash
# Install dependencies
(cd ./e2e && npm install)

# Run the tests interactively
(cd ./e2e && npm run open:golangangular)

# Run the tests non-interactively
(cd ./e2e && npm run e2e:golangangular)
```

When you are happy with your changes, you can cleanup your environment.

```bash
# Stop and remove the docker containers for zitadel and the database
docker compose --file ./e2e/config/host.docker.internal/docker-compose.yaml down
```

### Console

By executing the commands from this section, you run everything you need to develop the console locally.
Using [Docker Compose](https://docs.docker.com/compose/), you run [CockroachDB](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker-mac.html) and the [latest release of ZITADEL](https://github.com/zitadel/zitadel/releases/latest) on your local machine.
You use the ZITADEL container as backend for your console.
The console is run in your [Node](https://nodejs.org/en/about/) environment using [a local development server for Angular](https://angular.io/cli/serve#ng-serve), so you have fast feedback about your changes.

We use angular-eslint/Prettier for linting/formatting, so please run `yarn lint:fix` before committing. (VSCode users, check out [this ESLint extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [this Prettier extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) to fix lint and formatting issues in development)

Once you are happy with your changes, you run end-to-end tests and tear everything down.

The commands in this section are tested against the following software versions:

- [Docker version 20.10.17](https://docs.docker.com/engine/install/)
- [Node version v16.17.0](https://nodejs.org/en/download/)
- [npm version 8.18.0](https://docs.npmjs.com/try-the-latest-stable-version-of-npm)
- [Cypress runtime dependencies](https://docs.cypress.io/guides/continuous-integration/introduction#Dependencies)

<details>
  <summary>Note for WSL2 on Windows 10</summary>
  Following the suggestions <a href="https://stackoverflow.com/questions/62641553/setup-cypress-on-wsl-ubuntu-for-windows-10">here </a> subsequently <a href="https://github.com/microsoft/WSL/issues/4106">here </a> may  need to XLaunch and configure your DISPLAY variable. Use at your own risk.

1. Install `VcXsrv Windows X Server`
2. Set the target of your shortcut to `"C:\Program Files\VcXsrv\xlaunch.exe" -ac`
3. In WSL2 run `export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0` to set your DISPLAY variable
4. When starting XLaunch, make sure to disable access control
</details>

Run the database and the latest backend locally.

```bash
# Change to the console directory
cd ./console

# You just need the db and the zitadel services to develop the console against.
docker compose --file ../e2e/docker-compose.yaml up --detach zitadel
```

When the backend is ready, you have the latest zitadel exposed at http://localhost:8080.
You can now run a local development server with live code reloading at http://localhost:4200.
To allow console access via http://localhost:4200, you have to configure the ZITADEL backend.

1. Navigate to <http://localhost:8080/ui/console/projects>.
2. When prompted, login with _zitadel-admin@<span because="breaks the mailto"></span>zitadel.localhost_ and _Password1!_
3. Select the _ZITADEL_ project.
4. Select the _Console_ application.
5. Select _Redirect Settings_
6. Add _http://<span because="breaks the link"></span>localhost:4200/auth/callback_ to the _Redirect URIs_
7. Add _http://<span because="breaks the link"></span>localhost:4200/signedout_ to the _Post Logout URIs_
8. Select the _Save_ button

You can run the local console development server now.

```bash
# Install npm dependencies
yarn install

# Generate source files from Protos
yarn generate

# Start the server
yarn start

# If you don't want to develop against http://localhost:8080, you can use another environment
ENVIRONMENT_JSON_URL=https://my-cloud-instance-abcdef.zitadel.cloud/ui/console/assets/environment.json yarn start
```

Navigate to http://localhost:4200/.
Make some changes to the source code and see how the browser is automatically updated.
After making changes to the code, you should run the end-to-end-tests.
Open another shell.

```bash
# Reformat your console code
yarn lint:fix

# Change to the e2e directory
cd .. && cd e2e/

# If you made changes in the e2e directory, make sure you reformat the files here too
npm run lint:fix

# Install npm dependencies
npm install

# Run all e2e tests
npm run e2e:angular -- --headed
```

You can also open the test suite interactively for fast feedback on specific tests.

```bash
# Run tests interactively
npm run open:angular
```

If you also make [changes to the backend code](#backend--login), you can run the test against your locally built backend code and frontend server

```bash
npm run open:golangangular
npm run e2e:golangangular
```

When you are happy with your changes, you can format your code and cleanup your environment

```bash
# Stop and remove the docker containers for zitadel and the database
docker compose down
```

## Contribute docs

Project documentation is made with docusaurus and is located under [./docs](./docs).

### Local testing

Please refer to the [README](./docs/README.md) for more information and local testing.

### Style guide

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

### Docs pull request

When making a pull request use `docs(<scope>): <short summary>` as title for the semantic release.
Scope can be left empty (omit the brackets) or refer to the top navigation sections.

## Contribute internationalization

ZITADEL loads translations from four files:

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

## Want to start ZITADEL?

You can find an installation guide for all the different environments here:
[https://zitadel.com/docs/self-hosting/deploy/overview](https://zitadel.com/docs/self-hosting/deploy/overview)

## **Did you find a security flaw?**

- Please read [Security Policy](./SECURITY.md).

## Product management

The ZITADEL Team works with an agile product management methodology.
You can find all the issues prioritized and ordered in the [product board](https://github.com/orgs/zitadel/projects/2/views/1).

### Sprint

We want to deliver a new release every second week. So we plan everything in two-week sprints.
Each Tuesday we estimate new issues and on Wednesday the last sprint will be reviewed and the next one will be planned.
After a sprint ends a new version of ZITADEL will be released, and publish to [ZITADEL Cloud](https://zitadel.cloud) the following Monday.

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
- **👀 In review**: The issue is in review. Please add someone to review your issue or let us know that it is ready to review with a comment on your pull request.
- **✅ Done**: The issue is implemented and merged to main.

#### Priority

Priority shows you the priority the ZITADEL team has given this issue. In general the higher the demand from customers and community for the feature, the higher the priority.

- **🌋 Critical**: This is a security issue or something that has to be fixed urgently, because the software is not usable or highly vulnerable.
- **🏔 High**: These are the issues the ZITADEL team is currently focusing on and will be implemented as soon as possible.
- **🏕 Medium**: After all the high issues are done these will be next.
- **🏝 Low**: This is low in priority and will probably not be implemented in the next time or just if someone has some time in between.

#### Complexity

This should give you an indication how complex the issue is. It's not about the hours or effort it takes.
Everything that is higher than 8 should be split in smaller parts.

**1**, **2**, **3**, **5**, **8**, **13**

### About the labels

There are a few general labels that don't belong to a specific category.

- **good first issue**: This label shows contributors, that it is an easy entry point to start developing on ZITADEL.
- **help wanted**: The author is seeking help on this topic, this may be from an internal ZITADEL team member or external contributors.

#### Category

The category shows which part of ZITADEL is affected.

- **category: backend**: The backend includes the APIs, event store, command and query side. This is developed in golang.
- **category: ci**: ci is all about continues integration and pipelines.
- **category: design**: All about the ux/ui of ZITADEL
- **category: docs**: Adjustments or new documentations, this can be found in the docs folder.
- **category: frontend**: The frontend concerns on the one hand the ZITADEL management console (Angular) and on the other hand the login (gohtml)
- **category: infra**: Infrastructure does include many different parts. E.g Terraform-provider, docker, metrics, etc.
- **category: translation**: Everything concerning translations or new languages

#### Language

The language shows you in which programming language the affected part is written

- **lang: angular**
- **lang: go**
- **lang: javascript**
