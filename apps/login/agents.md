# ZITADEL Login App Guide for AI Agents

## Context
The **Login App** (`apps/login`) provides the user interface for authentication flows (Login, Register, MFA, etc.). It is built with Next.js and React.

## Key Technology
- **Framework**: Next.js (React).
- **Styling**: TailwindCSS (check project for specific config).
- **Data Fetching**: Primarily server-side interaction with ZITADEL APIs via `@zitadel/client` or direct gRPC calls where applicable.
- **Language**: TypeScript.

## Architecture & Conventions
- **Pages**: Uses Next.js App Router or Pages Router (Verify in `src/`).
- **Composability**: Components should be small and reusable.
- **State**: Critical authentication state is often managed via URL parameters (Auth Requests) and cookies/sessions.

## Development Commands
- **Dev Server**: `pnpm nx run @zitadel/login:dev`
- **Build**: `pnpm nx run @zitadel/login:build`
- **Lint**: `pnpm nx run @zitadel/login:lint`
- **Test**: `pnpm nx run @zitadel/login:test-unit`
