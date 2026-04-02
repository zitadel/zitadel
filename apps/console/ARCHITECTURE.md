# ZITADEL Frontend Architecture Plan

## Target Monorepo Structure

```
apps/
├── console/        # Single-instance admin, self-hosted (AGPL-3.0)
├── login/          # Product login, ships in containers (MIT)
├── docs/           # Documentation (Apache 2.0)
├── cloud/          # zitadel.com — everything for cloud (Source Available)
│   └── app/
│       ├── (website)/     # Marketing pages (Sanity CMS)
│       ├── (docs)/        # Docs with authenticated API explorer (fumadocs)
│       ├── (auth)/        # Cloud login with customizations
│       ├── (console)/     # Multi-instance console (imports from apps/console)
│       ├── (admin)/       # Instance management, billing, usage
│       └── layout.tsx     # Shared: Mixpanel, auth, instance context
└── website/        # (placeholder — may merge into cloud)

packages/
├── zitadel-client/  # API client (MIT)
├── zitadel-proto/   # Proto definitions (MIT)
└── zitadel-ui/      # Shared design system (Apache 2.0)
```

> [!IMPORTANT]
> `apps/cloud` is **source-available but not open source**. `packages/zitadel-ui` is Apache 2.0 (industry standard for UI libraries). See `LICENSING.md`.

---

## Phase 1: Cloud App Skeleton ✅

- [x] Scaffold `apps/cloud` as Next.js app (port 3001)
- [x] `package.json`, `next.config.mjs`, `tsconfig.json`, `project.json`
- [x] `@console/*` path alias for cross-app imports
- [x] `(console)` route group with placeholder
- [x] Verified: `pnpm nx dev cloud` starts successfully

---

## Phase 2: Wire Console into Cloud

- [ ] Import console client components via `@console/*` path alias
- [ ] Create `InstanceContext` provider — resolves API target from URL params or cookie
- [ ] Modify server actions to accept instance URL + token as params (instead of env vars only)
- [ ] Cloud sidebar with instance switcher

---

## Phase 3: Cross-Cutting Cloud Features

- [ ] **Mixpanel** — single init in cloud root layout, track user journey across all sections
- [ ] **Auth state** — shared login session across console, docs, website
- [ ] **Authenticated API docs** — fumadocs routes with "Try it" that picks user's instance
- [ ] **Debug page** — `/debug` route (preview/dev only) for test instance configuration

---

## Phase 4: Cloud-Specific Pages

- [ ] Instance management (list, create, configure)
- [ ] Billing & subscriptions
- [ ] Usage metrics
- [ ] Cloud signup flow

---

## Phase 5: Design System Extraction

- [ ] Reconcile website theme tokens with console's shadcn tokens
- [ ] Extract shared components (`StatusBadge`, `TablePagination`, etc.) to `packages/zitadel-ui`
- [ ] All apps consume `@zitadel/ui`

---

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Cloud deployment | Single Next.js app | Shared Mixpanel, auth, instance context across all sections |
| Console import | `@console/*` path alias | No package publishing, direct cross-app imports |
| Login build problem | Separate routes in cloud app | No more build flags — cloud has own login routes |
| UI package license | Apache 2.0 | Industry standard (Supabase MIT, GitLab MIT, Grafana Apache) |
| Cloud/website license | Source available | Clear restriction in LICENSING.md |

---

## Open Questions

1. **Naming**: Should `apps/cloud` be `apps/web` or stay as `cloud`?
2. **Instance API**: Which API provides the instance list for the switcher?
3. **Auth flow**: Same OIDC as console, or separate for cloud?
4. **Tailwind version**: Website uses v4, console uses v4 — confirm alignment
