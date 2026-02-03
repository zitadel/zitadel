# ZITADEL Monorepo Guide for AI Agents

## Mission & Context
ZITADEL is an open-source Identity Management System (IAM) written in Go and Angular/React. It provides secure login, multi-tenancy, and audit trails.

## Repository Structure Map
- **`apps/`**: Consumer-facing web applications.
  - **`login`**: Next.js (React) application for user authentication flows. See `apps/login/agents.md`.
  - **`docs`**: Documentation site built with Fumadocs (Next.js/MDX). See `apps/docs/agents.md`.
  - **`api`**: (Wait, check `internal/api` usually, but `apps/api` exists for build targets).
- **`console/`**: The Management Console application (Angular). See `console/agents.md`.
- **`internal/`**: Core backend logic (Go). Contains the EventStore, API implementations (gRPC/REST), and business logic.
- **`packages/`**: Shared TypeScript/JavaScript libraries used by the frontend apps.
- **`proto/`**: Protocol Buffer definitions.

## Technology Stack & Conventions
- **Orchestration**: Nx is used for build and task orchestration.
- **Package Manager**: pnpm.
- **Backend**: Go 1.24+.
  - **Communication**: gRPC is the primary communication protocol. REST is often capable via grpc-gateway.
  - **Pattern**: Event Sourcing is a core architectural pattern here.
- **Frontend**:
  - **Console**: Angular + RxJS.
  - **Login/Docs**: Next.js + React.

## Key Commands (AI Shortcuts)
Run these from the root using `pnpm nx`.

- **Build**: `pnpm nx run [PROJECT]:build` (e.g., `pnpm nx run @zitadel/console:build`)
- **Test**: `pnpm nx run [PROJECT]:test`
- **Lint**: `pnpm nx run [PROJECT]:lint`
- **Generate Code**: `pnpm nx run [PROJECT]:generate` (Important for proto generation)

## Documentation
- **Human Guide**: See `CONTRIBUTING.md` for setup and contribution details.
- **API Design**: See `API_DESIGN.md` for API specific guidelines.
