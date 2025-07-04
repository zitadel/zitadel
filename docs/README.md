# ZITADEL-Docs

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

If you are introducing new APIs (gRPC), you need to add a new entry to `docusaurus.config.js` under the `plugins` section.

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
