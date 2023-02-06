[![Netlify Status](https://api.netlify.com/api/v1/badges/b82a23f7-d8c7-4025-af18-a46586e89ed0/deploy-status)](https://app.netlify.com/sites/zitadel-docs/deploys)

# Website

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

## Installation

```
yarn install
```

## Local Development

```
yarn start
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