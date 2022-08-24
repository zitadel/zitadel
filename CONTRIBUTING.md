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
- If you found a mistake on our [docs page](https://docs.zitadel.com) or something is missing please read [the docs section](#contribute-docs)
- [Translate](#contribute-internationalization) and improve texts

Follow [@zitadel](https://twitter.com/zitadel) on twitter

## How to contribute

We strongly recommend to [talk to us](https://zitadel.com/contact) before you start contributing to streamline our and your work.

We accept contributions through pull requests. You need a github account for that. If you are unfamiliar with git have a look at Github's documentation on [creating forks](https://help.github.com/articles/fork-a-repo) and [creating pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork). Please draft the pull request as soon as possible. Go through the following checklist before you submit the final pull request:

### Submit a Pull Request (PR)

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

8.  On GitHub, [send a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/requesting-a-pull-request-review) to `zitadel:main`. Request review from one of the maintainers.

### Reviewing a Pull Request

The reviewers will provide you feedback and approve your changes as soon as they are satisfied. If we ask you for changes in the code, you can follow the [GitHub Guide](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/incorporating-feedback-in-your-pull-request) to incorporate feedback in your pull request.

<!-- TODO: how to do this via git -->
<!-- TODO: change commit message via git -->

### Commit Messages

Make sure you use [semantic release messages format](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type).

`<type>(<scope>): <short summary>`

#### Type

Must be one of the following:

- **feat**: New Feature
- **fix**: Bug Fix
- **docs**: Documentation

#### Scope

This is optional to indicate which component is affected. In doubt, leave blank (`<type>: <short summary>`)

#### Short Summary

Provide a brief description of the change.

## Contribute

The code consists of the following parts:

| name            | description                                                        | language                                                                    | where to find                                      |
| --------------- | ------------------------------------------------------------------ | --------------------------------------------------------------------------- | -------------------------------------------------- |
| backend         | Service that serves the grpc(-web) and RESTful API                 | [go](https://go.dev)                                                        | [API implementation](./internal/api/grpc)          |
| console         | Frontend the user interacts with after he is logged in             | [Angular](https://angular.io), [Typescript](https://www.typescriptlang.org) | [./console](./console)                             |
| login           | Server side rendered frontend the user interacts with during login | [go](https://go.dev), [go templates](https://pkg.go.dev/html/template)      | [./internal/api/ui/login](./internal/api/ui/login) |
| API definitions | Specifications of the API                                          | [Protobuf](https://developers.google.com/protocol-buffers)                  | [./proto/zitadel](./proto/zitadel)                 |
| docs            | Project documentation made with docusaurus                         | [Docusaurus](https://docusaurus.io/)                                        | [./docs](./docs)                                   |

Please validate and test the code before you contribute.

We add the label "good first issue" for problems we think are a good starting point to contribute to ZITADEL.

- [Issues for first time contributors](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [All issues](https://github.com/zitadel/zitadel/issues)


### Backend / Login

To keep the code clean and understandable we use [golangci-lint](https://golangci-lint.run). We recommend to format the code with this linter while working on ZITADEL to simplify the review process. The configuration is located [here](./.golangci.yaml).

To start the backend with a debugger run the [`main.go`-file](./main.go) located in the root of ZITADEL and provide the arguments and env-variables from below. Ensure that the database is running by running `docker compose -f ./build/local/docker-compose.yml up db`. For additional information please use the documentation of your IDE.

Make sure to use the following configurations:

<!-- TODO: document workflow -->

### Console

Change to the console directory

```bash
cd ./console
```

Run the database and the backend locally

```bash
docker compose --file ../e2e/docker-compose.yaml up --detach db zitadel
```

Console loads its environment from the file console/src/assets/environment.json.
Load it from your local target system.

```bash
curl -O ./src/assets/environment.json http://localhost:8080/ui/console/assets/environment.json
```

To generate source files from protos, run the following command
```
DOCKER_BUILDKIT=1 docker build -f ../build/console/Dockerfile . -t zitadel-npm-base --target npm-copy -o internal/api/ui/console/static
```

To run the console locally, run `npm install` and then `npm start` for a dev server. Navigate to http://localhost:4200/. The app will automatically reload if you change any of the source files.

You can now also run end-to-end tests interactively.
Open a new shell and change to the e2e directory

```bash
cd ./e2e
```

Start cypress and point it to your local dev server:

```bash
CYPRESS_BASE_URL=http://localhost:4200 npm start
```

### API Definitions

Ensure the provided code meets the [official style guide](https://developers.google.com/protocol-buffers/docs/style).

The following docker command builds the grpc stub into the correct folders:

```bash
docker build -f build/grpc/Dockerfile -t zitadel-base:local . \
    && docker build -f build/zitadel/Dockerfile . -t zitadel-go-base --target go-copy -o .
```

### Testing

<!-- TODO: how to run E2E tests -->

## Contribute Docs

Project documentation is made with docusaurus and is located under [./docs](./docs).
Please refer to the [README](./docs/README.md) for more information and local testing.

When making a pull request use `docs(<scope>): <short summary>` as title for the semantic release.
Scope can be left empty (omit the brackets) or refer to the top navigation sections.

## Contribute Internationalization

ZITADEL loads translations from four files:

- [Console texts](./console/src/assets/i18n)
- [Login interface](./internal/api/ui/login/static/i18n)
- [Email notification](./internal/notification/static/i18n)
- [Common texts](./internal/static/i18n) for success or error toasts

You may edit the texts in these files or create a new file for additional language support. Make sure you set the locale (ISO 639-1 code) as the name of the new language file.

## Want to start ZITADEL?

You can find an installation guide for all the different environments here:
[https://docs.zitadel.com/docs/guides/installation](https://docs.zitadel.com/docs/guides/installation)

## **Did you find a security flaw?**

- Please read [Security Policy](./SECURITY.md).
