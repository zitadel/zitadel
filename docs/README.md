# ZITADEL-Docs

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

The documentation is part of the ZITADEL monorepo and uses **pnpm** and **Nx** for development and build processes.

## Quick Start

```bash
# From the repository root
pnpm install

# Start development server (with Nx)
pnpm nx run docs:dev

# Or start directly from docs directory
cd docs && pnpm start
```

The site will be available at http://localhost:3000

## Available Scripts

All scripts can be run from the repository root using Nx:

```bash
# Development server with live reload
pnpm nx run docs:dev

# Build for production
pnpm nx run docs:build

# Generate API documentation and configuration docs
pnpm nx run docs:generate

# Lint and fix code
pnpm nx run docs:lint

# Serve production build locally
cd docs && pnpm serve
```

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

## Local Development

### Standard Development

```bash
# Install dependencies
pnpm install

# Start development server
pnpm start
```

### API Documentation Development

When working on the API docs, run a local development server with:

```bash
pnpm start:api
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
