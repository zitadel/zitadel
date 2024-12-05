# ZITADEL-Docs

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

## Knowledge Base Articles

Create a new .md file in the folder `/knowledge`

```md
---
description: "Add a description"
tags:
    - FAQ
    - subscription
    - downgrade
    - account-details
---

# Heading 1

## Overview

Provide an overview of the solution

## Solution

Normal Markdown text.
**Bold Text**
`code format``
[example@examplemail.com](mailto:example@examplemail.com)
```

Add images to `/src/img/knowledge`

Create a new category by creating a new folder. Add a `_category_.json` inside that folder and change the category title and other details.

## Installation

Install dependencies with

```
yarn install
```

then run 

```
yarn generate
```


## Local Development

Start a local development server with

```
yarn start
```

When working on the API docs, run a local development server with 

```
yarn start:api
```

## Container Image

If you just want to start docusaurus locally without installing node you can fallback to our container image.
Execute the following commands from the repository root to build and start a local version of ZITADEL 

```shell
docker build -f docs/Dockerfile . -t zitadel-docs
```

```shell
docker run -p 8080:8080 zitadel-docs
```
