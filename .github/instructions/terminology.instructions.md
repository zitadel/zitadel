---
applyTo: '{apps/docs/content/**/*.mdx,apps/docs/content/**/*.md,console/src/assets/i18n/*.json,apps/login/locales/*.json,internal/static/i18n/*.yaml,internal/static/i18n/*.yml,proto/**/*.proto}'
---

When reviewing user-facing naming changes, enforce the canonical terminology table in `TERMINOLOGY.md`.

### How to apply
- Cross-reference every changed user-facing string against the **Search for (discouraged)** column.
- If a discouraged term is found, flag it and suggest the canonical term from **Replace with / enforce**.
- For `proto/**/*.proto` files: check comments and `description:`/`summary:` text only; do not flag identifiers, field names, or enum values.
- Only apply rules where the **Scope** column matches the file type:
  - `UI` → `console/src/assets/i18n/`, `apps/login/locales/`
  - `Docs` → `apps/docs/content/`
  - `API` → `proto/**`, API description text
  - `Everywhere` → all of the above
- If a PR introduces a new policy, setting, or resource term not yet in the table, request that the contributor adds it to `TERMINOLOGY.md`.
