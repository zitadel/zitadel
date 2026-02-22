# GitHub Copilot Instructions for ZITADEL

You are working in the ZITADEL monorepo. This codebase uses specific conventions for its Go backend and Angular/Next.js frontends.

**CRITICAL**: Read `AGENTS.md` in the root first, then the nearest scoped `AGENTS.md` for changed files.

## Key References
- `AGENTS.md`: Root architecture map and global commands.
- `apps/api/AGENTS.md`: API app workflows and backend orchestration targets.
- `apps/login/AGENTS.md`: Specifics for the Next.js Login UI.
- `apps/docs/AGENTS.md`: Specifics for the Fumadocs documentation.
- `console/AGENTS.md`: Specifics for the Angular Console.
- `internal/AGENTS.md`: Backend domain and event-sourcing boundaries.
- `proto/AGENTS.md`: API schema and generation guidance.
- `packages/AGENTS.md`: Shared client/proto package workflows.
- `tests/functional-ui/AGENTS.md`: Cypress functional UI test workflows.
- `deploy/compose/agents.md`: Docker Compose deployment invariants, file conventions, and rejected alternatives.

## Behavior
- Before Go-related work, inspect `go.mod` for required Go version/toolchain.
- Use verified Nx targets only; if unsure, run `pnpm nx show project <project>`.
- For backend changes, note we are in transition: relational data is becoming the system of record, while event writes are still required for history/audit.
- Respect terminology defined in the `Technical Glossary` section of `AGENTS.md` when generating UI text or documentation.
- Distinguish between `apps/login` (Next.js) and `console` (Angular) when suggesting frontend code.
- For deployment-related work, consult the directory-specific agents files (e.g., `deploy/compose/agents.md`) for invariants, file conventions, and rejected alternatives.
