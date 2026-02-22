# ZITADEL Packages Guide for AI Agents

## Context
`packages/` contains shared TypeScript libraries used by frontend applications and external consumers.

## Main Packages
- **`packages/zitadel-proto`** (`@zitadel/proto`): generated protobuf TypeScript artifacts.
- **`packages/zitadel-client`** (`@zitadel/client`): higher-level client library built on generated proto/connect types.

## Verified Nx Targets
- **Proto generation**: `pnpm nx run @zitadel/proto:generate`
- **Client build**: `pnpm nx run @zitadel/client:build`
- **Client lint**: `pnpm nx run @zitadel/client:lint`
- **Client tests**: `pnpm nx run @zitadel/client:test`

## Workflow Notes
- When changing `proto/`, regenerate `@zitadel/proto` first, then validate/build `@zitadel/client`.
- Keep package exports and public typings stable unless a breaking release is explicitly intended.
