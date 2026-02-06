# GitHub Copilot Instructions for ZITADEL

You are working in the ZITADEL monorepo. This codebase uses specific conventions for its Go backend and Angular/Next.js frontends.

**CRITICAL**: Before suggesting complex changes, read the `agents.md` file in the root or the active application directory for architectural context.

## Key References
- `agents.md`: Root architecture map and global commands.
- `apps/login/agents.md`: Specifics for the Next.js Login UI.
- `apps/docs/agents.md`: Specifics for the Fumadocs documentation.
- `console/agents.md`: Specifics for the Angular Console.

## Behavior
- Use `pnpm nx` for all build and test commands.
- Respect the "Event Sourcing" pattern when working on the Go backend (`internal/`).
- Distinguish between `apps/login` (Next.js) and `console` (Angular) when suggesting frontend code.
