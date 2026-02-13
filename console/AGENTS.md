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

## Development Commands
- **Dev Server**: `pnpm nx run @zitadel/console:dev`
- **Build**: `pnpm nx run @zitadel/console:build`
- **Lint**: `pnpm nx run @zitadel/console:lint`
- **Test**: The `@zitadel/console` project currently has no `test` target configured in Nx. Tests are run via Karma; add a `test` target in the project configuration before using `nx run @zitadel/console:test`.
