# ZITADEL-Docs

This documentation page is built using [Docusaurus](https://docusaurus.io/).

## Quick Start

```bash
# From the repository root
pnpm install

# Start development server
nx run @zitadel/docs:start
```

The site will be available at http://localhost:3003

## Available Scripts

All scripts can be run from the repository root

```bash
# Build for production
nx run @zitadel/docs:build

# Generate API documentation and configuration docs
nx run @zitadel/docs:generate

# Lint and fix code
nx run @zitadel/docs:lint

# Serve production build locally
nx run @zitadel/docs:serve
```

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

If you are introducing new APIs (gRPC), you need to add a new entry to `docusaurus.config.js` under the `plugins` section.