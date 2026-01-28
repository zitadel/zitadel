# @zitadel/docs

Documentation for Zitadel, built with [Fumadocs](https://fumadocs.dev) and [Next.js](https://nextjs.org).

## Getting Started

Ensure you have followed the [root quick start](../../CONTRIBUTING.md#quick-start) to set up dependencies.

### Local Development

Start the development server:

```bash
pnpm nx run @zitadel/docs:dev
```

The site will be available at http://localhost:3000.

### Scripts

Key scripts for documentation workflows:

| Script | Description |
| :--- | :--- |
| `dev` | Starts the development server. |
| `build` | Builds the production application. |
| `generate:fetch` | Fetches remote tags and referenced content. |
| `generate` | Runs all generation steps (`generate:fetch`, `generate:grpc`, `generate:api`, `generate:indices`). |
| `check:links` | Checks for broken links. |

## Contributing

### Build Process

The docs build process automatically handling the following steps via `generate`:
1.  Downloads required protoc plugins.
2.  Generates gRPC documentation from proto files.
3.  Generates API documentation from OpenAPI specs.
4.  Generates index files for directory structures.

### Style Guide

- **Variables**: Use environment variables in code snippets where possible.
- **Embedded Content**: Use `_filename.mdx` for content embedded in other pages (not indexed individually).
- **Code Embedding**: Use the `file` property in code blocks to embed code from the repo.
- **Voice**: Use active voice and sentence case for titles.

Refer to the [Google Developer Style Guide](https://developers.google.com/style) for general guidelines.

### Pull Requests

Use `docs(<scope>): <short summary>` for PR titles.
Pass quality checks before submitting:

```bash
pnpm nx run @zitadel/docs:build
pnpm nx run @zitadel/docs:check:links
```
