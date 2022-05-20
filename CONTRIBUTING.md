# How to contribute to ZITADEL

## Did you find a bug?

Please file an issue [here](https://github.com/zitadel/zitadel/issues/new?assignees=&labels=bug&template=bug_report.md&title=).

Bugs are evaluated every day as soon as possible.

## Enhancement

Do you miss a feature? Please file an issue [here](https://github.com/zitadel/zitadel/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=)

Enhancements are discussed and evaluated every Wednesday by the ZITADEL core team.

## Grab an Issuesa# Contributing to ZITADEL

## Introduction

Thank you for your interest about how to contribute! As you might know there is more than code to contribute. You can find all information needed to start contributing here.

Please give us and our community the chance to get rid of security vularbilities by responsibly disclose this kind of issues by contacting [security@zitadel.com](mailto:security@zitadel.com).

The strongest part of a community is the possibility to share thoughts. That's why we try to react as soon as possible to your ideas, thoughts and feedback. We love to discuss as much as possible in an open space like in the [issues](https://github.com/zitadel/zitadel/issues) and [discussions](https://github.com/zitadel/zitadel/discussions) section here or in our [chat](https://zitadel.com/chat), but we understand your doubts and provide further contact options [here](https://zitadel.com/contact).

If you want to give an answer or be part of discussions please be kind. Treat others like you want to be treated. Read more about our code of conduct [here](CODE_OF_CONDUCT.md).

## What can I contribute?

For people who are new to ZITADEL: We flag issues which are a good starting point to start contributing. You find them [here](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

Make ZITADEL more popular and give it a ⭐

Help shaping the future of ZITADEL by

- Join our [chat](https://zitadel.com/chat) and discuss with us or others.
- Ask or answer questions in the [issues section](https://github.com/zitadel/zitadel/issues)
- Share your thoughts and ideas in the [discussions section](https://github.com/zitadel/zitadel/discussions)

[Contribute](#how-to-contribute)

- [Code](#contribute-code)
- If you found a mistake on our [docs page](https://docs.zitadel.ch) or something is missing please read [the docs section](#contribute-docs)

Follow [@zitadel](https://twitter.com/zitadel) on twitter

## How to contribute

We strongly recomend to talk to us before you start contributing to streamline our and your work.

We accept contributions through pull requests. You need a github account for that. If you are unfamiliar with git have a look at Github's documentation on [creating forks](https://help.github.com/articles/fork-a-repo) and [creating pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork). Please draft the pull request as soon as possible. Go through the following checklist before you submit the final pull request:

1. Create a feature branch from the `main`-branch
1. Make your changes on the new branch
1. [Merge](https://git-scm.com/book/en/v2/Git-Branching-Basic-Branching-and-Merging) the lastet commit of the `main`-branch
1. Use [Semantic Release commit messages](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type) to simplify creation of release notes. In the title of the pull request correct tagging is required and will be requested by the reviewers.
1. Request a [review](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/requesting-a-pull-request-review) from one of the authors. The reviewers will provide you feedback and approve your changes as soon as they are satisfied.

## Contribute

The code consists of the following parts:

| name | description | language | where to find |
|---|---|---|---|
| backend | Service that serves the grpc(-web) and RESTful API  | [go](https://go.dev) | [API implementation](./internal/api/grpc) |
| console | Frontend the user inertacts with after he is logged in | [Angular](https://angular.io), [Typescript](https://www.typescriptlang.org) | [./console](./console) |
| login | Server side rendered frontend the user interacts with during login | [go](https://go.dev), [go templates](https://pkg.go.dev/html/template) | [./internal/api/ui/login](./internal/api/ui/login) |
| API definitions | Specifications of the API | [Protobuf](https://developers.google.com/protocol-buffers) | [./proto/zitadel](./proto/zitadel) |
| docs | Project documentation made with docusaurus | [Docusaurus](https://docusaurus.io/) | [./docs](./docs) |

Please validate and test the code before you contribute.

### Backend / Login

To keep the code clean and understandable we use [golangci-lint](https://golangci-lint.run). We recommend to format the code with this linter while working on ZITADEL to simplify the review process. The configuration is locaed [here](./.golangci.yaml).

To start the backend with a debugger run the [`main.go`-file](./main.go) located in the root of ZITADEL and provide the arguments and env-variables from below. Ensure that the database is running by running `docker compose -f ./build/local/docker-compose.yml up db`. For additional information please use the documentation of your IDE.

Make sure to use the following configurations:

TODO document workflow

### Console

<!-- TODO: ask maxpe for infos -->

### API Definitions

Ensure the provided code meets the [offical style guide](https://developers.google.com/protocol-buffers/docs/style).

The following docker command builds the grpc stub into the correct folders:

```bash
docker build -f build/grpc/Dockerfile -t zitadel-base:local . \
    && docker build -f build/zitadel/Dockerfile . -t zitadel-go-base --target go-copy -o .
```

## Contribute Docs

<!-- TODO: ask maxpe for infos -->

We add the label "good first issue" for problems we think are a good starting point to contribute to ZITADEL.

* [Issues for first time contributors](https://github.com/zitadel/zitadel/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* [All issues](https://github.com/zitadel/zitadel/issues)

### Make a PR

If you like to contribute fork the ZITADEL repository. After you implemented the new feature create a PullRequest in the ZITADEL reposiotry.

Make sure you use semantic release:

* feat: New Feature
* fix: Bug Fix
* docs: Documentation

## Want to start ZITADEL?

You can find an installation guide for all the different environments here:
[https://docs.zitadel.com/docs/guides/installation](https://docs.zitadel.com/docs/guides/installation)

## **Did you find a security flaw?**

* Please read [Security Policy](./SECURITY.md).
