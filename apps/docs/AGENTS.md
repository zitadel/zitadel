# ZITADEL Docs Guide for AI Agents

## Context
The **Docs App** (`apps/docs`) hosts the ZITADEL documentation. It has recently migrated to **Fumadocs**.

## Key Technology
- **Framework**: Fumadocs (built on Next.js).
- **Content**: MDX (Markdown + React Components).
- **Orchestration**: Nx.

## Content Structure
```
apps/docs/
├── content/         # MDX documentation files
├── components/      # React components for docs
├── lib/             # Utility functions
│   └── sidebar-data.ts  # Navigation structure definition (sidebar ordering, categories)
├── scripts/         # Build and generation scripts (link checking, API/proto doc generation)
└── public/          # Static assets
```

## Content Conventions

### Frontmatter
Standard Fumadocs frontmatter is required for all MDX files:

```mdx
---
title: "Page Title"
description: "Brief description for SEO and navigation"
---
```

**Note:** Only `title` and `description` are used in frontmatter. Navigation structure (sidebar ordering, categories, links) is managed in `lib/sidebar-data.ts`, not through frontmatter fields like `sidebar_position` (Docusaurus) or `slug`.

### Links
- **Internal Links**: Use relative paths (`./other-page`) or absolute paths from docs root (`/guides/authentication`)
- **External Links**: Use standard Markdown syntax with `https://` URLs
- **Cross-references**: Link to API docs, code samples, and related guides

### MDX Components
Use only documented Fumadocs components and project-specific components:

**Standard Fumadocs Components:**
- `<Callout>` - Highlight important information (info, warning, danger, note)
- `<Tabs>` - Tabbed content for multi-language examples
- `<Steps>` - Numbered step-by-step instructions
- `<Accordion>` - Collapsible sections
- `<Cards>` - Grid of linked cards

**Project-Specific Components:**
Check `apps/docs/components/` for custom MDX components available in the project.

**Code Blocks:**
Use fenced code blocks with language identifiers:
````mdx
```typescript
// TypeScript code example
const example = "value";
```
````

### API Documentation
- **Generation**: API docs are generated from proto definitions
- **Workflow**: Changes in `proto/` require regeneration via `pnpm nx run @zitadel/docs:generate`
- **Location**: Generated API docs are placed in appropriate sections
- **Manual Edits**: Avoid editing generated API docs directly; modify proto comments instead

### Content Organization
- **Guides**: Step-by-step tutorials (`/guides/`)
- **Concepts**: Explanatory documentation (`/concepts/`)
- **API Reference**: Generated from proto (`/api/`)
- **Examples**: Code samples and use cases (`/examples/`)

### Link Checking
Run link checker before committing:
```bash
pnpm nx run @zitadel/docs:check-links
```

This validates:
- Internal links point to existing pages (including redirects)
- Anchor links target valid headings
- Image paths are correct (both Markdown and HTML img tags)
- All content files use .mdx extension (not .md)

**Note:** External links (http/https URLs) are NOT validated by this script. Manual verification or separate tools are needed for external link checking.

## Writing Guidelines

### Style
- **Be concise**: Prefer short sentences and paragraphs
- **Use active voice**: "Configure the instance" not "The instance should be configured"
- **Code examples**: Include working examples with context
- **Screenshots**: Use when helpful, but keep them up-to-date

### Terminology
Follow ZITADEL domain terminology from root `AGENTS.md`:
- **Instance**: The logical tenant/environment (NOT "example")
- **Organization**: Group within an Instance
- **Project**: Collection of applications
- See root `AGENTS.md` for full translation glossary

### Accessibility
- **Alt text**: Provide descriptive alt text for images
- **Heading hierarchy**: Use proper heading levels (h2, h3, h4)
- **Link text**: Use descriptive link text, avoid "click here"

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/docs:dev` (with hot reload)
- **Build**: `pnpm nx run @zitadel/docs:build` (production build)
- **Generate**: `pnpm nx run @zitadel/docs:generate` (regenerate API docs from proto)
- **Lint**: `pnpm nx run @zitadel/docs:lint`
- **Check Links**: `pnpm nx run @zitadel/docs:check-links` (validate all links)
- **Check Types**: `pnpm nx run @zitadel/docs:check-types` (TypeScript type checking)
- **Test**: `pnpm nx run @zitadel/docs:test`

## Workflow for Documentation Changes

1. **Content Changes**: Edit MDX files in appropriate directory
2. **Component Changes**: Modify React components if needed
3. **API Changes**: If proto changed, run `generate` target
4. **Validation**: Run `check-links` and `check-types`
5. **Preview**: Test with `dev` target
6. **Build**: Verify with `build` target before committing

## Cross-References
- **API Design**: See `API_DESIGN.md` for API documentation standards
- **Proto Definitions**: See `proto/AGENTS.md` for API changes that affect docs
- **Domain Terminology**: See root `AGENTS.md` for ZITADEL-specific terms and translations
- **Contributing**: See `CONTRIBUTING.md` for general contribution guidelines
