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

## Deployment

Docs are deployed to [Vercel](https://vercel.com) via the [Deploy Docs](https://github.com/zitadel/zitadel/actions/workflows/deploy-docs.yml) GitHub Actions workflow using the **Vercel CLI** directly. Vercel's native GitHub integration is intentionally disabled — it requires every contributor who triggers a deployment to hold a paid Vercel team member seat. The GitHub Actions approach removes that restriction entirely.

| Trigger | Target |
|:--------|:-------|
| Push to `main` (when Nx detects docs are affected) | Production |
| Pull request to `main` | Preview URL (posted as a PR comment + GitHub deployment status) |
| Manual dispatch (`workflow_dispatch`) | Production or Preview (selectable, always builds — bypasses the Nx affected check) |

### Redeploying

To manually redeploy without a code change:

1. Open the [Deploy Docs](https://github.com/zitadel/zitadel/actions/workflows/deploy-docs.yml) workflow in GitHub Actions.
2. Click **Run workflow**.
3. Select the target branch (`main`) and environment (`production` or `preview`).
4. Click **Run workflow**.

### Required Secrets

The workflow requires three repository secrets:

| Secret | Description |
|:-------|:------------|
| `VERCEL_TOKEN` | A personal Vercel access token from a team member with Vercel project access. Generate a replacement at https://vercel.com/account/tokens. |
| `VERCEL_ORG_ID` | The Vercel team/org ID. |
| `VERCEL_PROJECT_ID_DOCS` | The Vercel project ID for the docs app. |

> **Note**: If the deployment stops working, check whether the token has been revoked or the Vercel account owner has changed. A new token from any team member with Vercel access can replace it in GitHub → Settings → Secrets → Actions → `VERCEL_TOKEN`.
