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
- **Components**: specific MDX components may be available for callouts, code blocks, etc.

## Development Commands
- **Dev Server**: `pnpm nx run @zitadel/docs:dev`
- **Build**: `pnpm nx run @zitadel/docs:build`
- **Lint**: `pnpm nx run @zitadel/docs:lint`
- **Check Links**: `pnpm nx run @zitadel/docs:check-links`
