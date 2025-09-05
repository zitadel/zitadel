# ZITADEL-Docs

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

The documentation is part of the ZITADEL monorepo and uses **pnpm** and **Nx** for development and build processes.

## Quick Start

```bash
# From the repository root
pnpm install
pnpm add -g nx

# Start development server (with Nx)
nx run @zitadel/docs:dev

# Or serve a production build
nx run @zitadel/docs:start
```

The site will be available at http://localhost:3000

To regenerate and rebuild the docs, rerun the `nx run` command.

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

If you are introducing new APIs (gRPC), you need to add a new entry to `docusaurus.config.js` under the `plugins` section.

## Build Process

The documentation build process automatically:

1. **Downloads required protoc plugins** - Ensures `protoc-gen-connect-openapi` is available
2. **Generates gRPC documentation** - Creates API docs from proto files
3. **Generates API documentation** - Creates OpenAPI specification docs
4. **Copies configuration files** - Includes configuration examples
5. **Builds the Docusaurus site** - Generates the final static site

## Container Image

If you just want to start docusaurus locally without installing node you can fallback to our container image.
Execute the following commands from the repository root to build and start a local version of ZITADEL

```shell
docker build -f docs/Dockerfile . -t zitadel-docs
```

```shell
docker run -p 8080:8080 zitadel-docs
```
