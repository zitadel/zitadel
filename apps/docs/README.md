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

| Script                 | Description                                                                                                                  |
|:-----------------------|:-----------------------------------------------------------------------------------------------------------------------------|
| `dev`                  | Starts the development server.                                                                                               |
| `build`                | Builds the production application.                                                                                           |
| `fetch:remote-content` | Fetches remote tags and referenced content.                                                                                  |
| `generate`             | Runs all generation steps (`fetch:remote-content`, `generate:proto-docs`, `generate:api-reference`, `generate:index-pages`). |
| `check:links`          | Validates content integrity (broken links, missing frontmatter, schema errors).                                              |
| `check-types`          | Validates typescript types.                                                                                                  |
| `test`                 | Runs all validation steps (`check-types`, `check:links`).                                                                    |
| `lint`                 | checks for code style and syntax errors (ESLint).                                                                            |
| `clean`                | Cleans the build output and generated files.                                                                                 |

### Validation

*   **Code Quality**: Run `pnpm lint` to check for syntax and style issues in JS/TS/MDX files.
*   **Content Integrity**: Run `pnpm check:links` to validate content structure, including:
    *   Broken internal links
    *   Missing required front-matter (e.g., `title`)
    *   Image references

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

### Adding Content

All documentation content is located in the `content` directory. Note that the system strictly accepts **only `.mdx` files**.

To add a new page:
1.  Create a `.mdx` file in the appropriate subdirectory of [`content`](./content).
2.  Register the new page in the sidebar settings at [`lib/sidebar-data.ts`](./lib/sidebar-data.ts) to make it accessible in the navigation.


### Pull Requests

Use `docs(<scope>): <short summary>` for PR titles.
Pass quality checks before submitting:

```bash
pnpm nx run @zitadel/docs:build
pnpm nx run @zitadel/docs:check:links
```

You can also run specific steps individually:
*   `pnpm fetch:remote-content`
*   `pnpm generate:api-reference`
