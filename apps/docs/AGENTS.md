# ZITADEL Docs Guide for AI Agents

## Context
The **Docs App** (`apps/docs`) hosts the ZITADEL documentation. It has recently migrated to **Fumadocs**.

## Key Technology
- **Framework**: Fumadocs (built on Next.js).
- **Content**: MDX (Markdown + React Components).
- **Orchestration**: Nx.

## Content Conventions
- **Frontmatter**: Ensure standard Fumadocs frontmatter is used (title, description).
- **Links**: Use absolute paths or standard MDX linking.
- **Components**: Use only MDX components that are documented in the official Fumadocs MDX component reference and any project-specific MDX component documentation in this repo; do not introduce or rely on undocumented components.
- **API Docs**: Changes in `proto/` or API response schemas often require regeneration via docs generation targets.

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/docs:dev`
- **Build**: `pnpm nx run @zitadel/docs:build`
- **Generate**: `pnpm nx run @zitadel/docs:generate`
- **Lint**: `pnpm nx run @zitadel/docs:lint`
- **Check Links**: `pnpm nx run @zitadel/docs:check-links`
- **Check Types**: `pnpm nx run @zitadel/docs:check-types`
- **Test**: `pnpm nx run @zitadel/docs:test`

## Deployment

Docs are deployed to Vercel via `.github/workflows/deploy-docs.yml` using the `amondnet/vercel-action`. Vercel's own GitHub integration is **disabled** (`vercel.json` → `"github": { "enabled": false }`) because it requires every contributor who triggers a deployment to hold a paid Vercel team member seat — the GitHub Actions approach removes that restriction.

- **Production deploy**: triggered automatically on push to `main` (when Nx detects docs are affected), or manually via `workflow_dispatch` with `environment: production`.
- **Preview deploy**: triggered automatically on pull requests to `main`; a URL is posted as a PR comment.
- **Manual redeploy**: GitHub Actions → Deploy Docs → Run workflow.

The build runs `pnpm nx affected --target=build --projects=@zitadel/docs` on the GitHub runner (full Nx pipeline: fetch-remote-content → proto docs → API reference → check-links → build), then `vercel build --prebuilt` packages the `.next` output before upload.

Do **not** re-enable `github.enabled` in `vercel.json` — it would both cause double deployments and re-introduce the paid member seat requirement for contributors.

> **Secrets**: Three repository secrets are required: `VERCEL_TOKEN`, `VERCEL_ORG_ID`, and `VERCEL_PROJECT_ID_DOCS`. The token is a personal access token (currently belonging to **@fforootd (Florian)**). If deployments fail with auth errors, the token may have been revoked — any team member with Vercel access can generate a replacement at https://vercel.com/account/tokens and update `VERCEL_TOKEN` in GitHub → Settings → Secrets → Actions.
