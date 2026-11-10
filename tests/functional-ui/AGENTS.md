# ZITADEL Functional UI Test Guide for AI Agents

## Context
`tests/functional-ui` contains Cypress-based end-to-end tests focused on Management Console flows against a running ZITADEL API.

## Verified Nx Targets
- **Open interactive Cypress runner**: `pnpm nx run @zitadel/functional-ui:open`
- **Run test suite**: `pnpm nx run @zitadel/functional-ui:test`
- **Start test DB only**: `pnpm nx run @zitadel/functional-ui:run-db`
- **Start test API only**: `pnpm nx run @zitadel/functional-ui:run-api`
- **Stop test infra**: `pnpm nx run @zitadel/functional-ui:stop`

## Workflow Notes
- Functional UI tests are the primary test path for Console user journeys, since `@zitadel/console` has no Nx `test` target.
- These tests depend on API build/run orchestration; avoid changing API startup assumptions without updating this suite.
