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
