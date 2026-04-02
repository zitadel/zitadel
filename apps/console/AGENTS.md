# ZITADEL Console (Next.js) Guide for AI Agents

## Context
The **Console-Next App** (`apps/console`) is the new ZITADEL management console, replacing the Angular-based console. It provides the admin UI for managing users, projects, applications, organizations, and sessions.

## Key Technology
- **Framework**: Next.js 16 with App Router, React 19.
- **Styling**: TailwindCSS 4 + shadcn/ui component library.
- **Data Fetching**: Server-side via ConnectRPC to ZITADEL v2 APIs using `@zitadel/client` and `@zitadel/proto`.
- **Authentication**: Personal Access Token (PAT) via `ZITADEL_PAT` env var, injected as `Authorization: Bearer` header.
- **Language**: TypeScript.
- **Package Manager**: pnpm (monorepo with Nx orchestration).

## Architecture & Conventions

### Routing
Uses Next.js App Router. Routes under `app/`:
- `app/users/` — User list and detail
- `app/projects/` — Project list and detail
- `app/applications/` — Application list and detail
- `app/sessions/` — Session list
- `app/organizations/` — Organization list

### Server/Client Split Pattern
Pages follow a consistent pattern:
1. **Server component** (`page.tsx`) — fetches data via server actions, handles errors, passes JSON-safe data as props.
2. **Client component** (`*-client.tsx`) — renders UI, manages local state, handles pagination and re-fetching.

### API Layer (`lib/api/`)
All ZITADEL API calls are encapsulated as `"use server"` actions in `lib/api/`:
- `transport.ts` — ConnectRPC transport setup with PAT auth.
- `services.ts` — Service client factories.
- `users.ts`, `projects.ts`, `applications.ts`, `sessions.ts`, `organizations.ts` — Domain-specific CRUD.
- `user-actions.ts` — User lifecycle: lock, unlock, deactivate, reactivate, delete, password reset.
- `fetch-*.ts` — JSON-safe wrappers that call the above and return serializable data.

### Proto Usage
- Import schemas from `@zitadel/proto/zitadel/<domain>/v2/<service>_pb`.
- Use `create()` from `@zitadel/client` to build requests.
- Use `toJson()` to serialize responses for client components.

### Context Providers
- `lib/context/app-context.tsx` — Global state: current organization, available organizations.
- `lib/permissions/context.tsx` — Permission-based UI gating.
- `lib/deployment/context.tsx` — Self-hosted vs. cloud feature flags.

### UI Components
- shadcn/ui components in `components/ui/`.
- Domain components in `components/` (e.g., `components/users/`, `components/layout/`).
- Icons from `lucide-react`.

## Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| `ZITADEL_INSTANCE_URL` | Yes | ZITADEL gRPC/Connect endpoint URL |
| `ZITADEL_PAT` | Yes | Personal Access Token for API auth |
| `NEXT_PUBLIC_DEPLOYMENT_MODE` | No | `self-hosted` or `cloud` (default: `self-hosted`) |

## Verified Nx Targets
- **Dev Server**: `pnpm nx run console-next:dev` (alias: `pnpm nx dev console-next`)
- **Build**: `pnpm nx run console-next:build`
- **Lint**: `pnpm nx run console-next:lint`

## Scope Rules
- For shared API typings and client behavior, also read `packages/AGENTS.md` and `proto/AGENTS.md`.
- The root `AGENTS.md` contains domain model, multi-tenancy rules, and PR conventions.
