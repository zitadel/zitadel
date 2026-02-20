# ZITADEL Console Guide for AI Agents

## Context
The **Management Console** (`console/`) is the administrative interface for ZITADEL. It allows developers and administrators to configure organizations, projects, and users.

## Key Technology
- **Framework**: Angular.
- **Language**: TypeScript.
- **State Management**: Reactive patterns with RxJS.
- **UI Component Library**: Angular Material (see `console/package.json`, `@angular/material` ^20.2.14).

## Architecture & Conventions
- **Services**: Business logic should reside in injectable services, not components.
- **Modules**: Angular Modules (NgModule) are used for grouping features.
- **gRPC**: Heavy usage of gRPC-web or REST mappings to talk to the ZITADEL API.

## Verified Nx Targets
- **Dev Server**: `pnpm nx run @zitadel/console:dev`
- **Build**: `pnpm nx run @zitadel/console:build`
- **Lint**: `pnpm nx run @zitadel/console:lint`
- **Generate**: `pnpm nx run @zitadel/console:generate` — runs `buf generate` to produce TypeScript/JS proto stubs in `src/app/proto/generated/`. Automatically depends on `install-proto-plugins`.
- **Install Proto Plugins**: `pnpm nx run @zitadel/console:install-proto-plugins` — downloads `protoc-gen-grpc-web` v1.5.0, `protoc-gen-js` v3.21.4, and `protoc-gen-openapiv2` v2.22.0 pre-built binaries to `.artifacts/bin/`. No Go toolchain required. Output is Nx-cached.
- **Test**: The `@zitadel/console` project currently has no `test` target configured in Nx.
- **Functional UI Tests**: Use `pnpm nx run @zitadel/functional-ui:test` (see `tests/functional-ui/AGENTS.md`).
