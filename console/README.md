# Zitadel Console

The Console is Zitadels management UI.

It is built using [Angular](https://angular.dev/) and part of the Zitadel monorepo.

To get started follow the [contributing quick start](../CONTRIBUTING.md#console).

### Building

To build for production:

```bash
pnpm nx run @zitadel/console:build
```

This will:

1. Generate proto files (via `prebuild` script)
2. Build the Angular app with production optimizations

### Linting

To run linting and formatting checks:

```bash
pnpm nx @zitadel/console:lint
```

To auto-fix formatting issues:

```bash
pnpm nx @zitadel/console:lint-fix
```

## Project Structure

- `src/app/proto/generated/` - Generated proto files (Angular-specific format)
- `buf.gen.yaml` - Local proto generation configuration
- `project.json` - Nx orchestration and caching for builds and tests
- `prebuild.development.js` - Development environment configuration script

### Dependency Chain

The Console app has the following build dependencies managed by Nx:

1. `@zitadel/proto:generate` - Generates the protobuf stubs
2. `@zitadel/client:build` - Builds the TypeScript gRPC client library
3. `@zitadel/console:generate` - Generates Console-specific protobuf stubs
4. `@zitadel/console:build` - Creates a production build from Console

This ensures that the Console always has access to the latest client library and protobuf definitions.


### Proto Generation Details

1. **`@zitadel/proto` generation**: Modern ES modules with `@bufbuild/protobuf` for v2 APIs
2. **Local `buf.gen.yaml` generation**: Traditional protobuf JavaScript classes for v1 APIs

The Console app calls Zitadel v1 and v2 APIs.
As long as the Console still calls v1 APIs, it needs to import client stubs from separate sources:
- [Source outputs from direct buf generation for v1 APIs](#v1-stubs)
- [@zitadel/client for v2 APIs](#v2-stubs)

### <a name="v1-stubs"></a>Legacy v1 API (Traditional Protobuf)

- Uses local `buf.gen.yaml` configuration
- Generates traditional Google protobuf JavaScript classes extending `jspb.Message`
- Uses plugins: `protocolbuffers/js`, `grpc/web`, `grpc-ecosystem/openapiv2`
- Output: `src/app/proto/generated/`
- Used for: Most existing Console functionality

### <a name="v2-stubs"></a>Modern v2 API (ES Modules)

- Uses `@zitadel/proto` package generation
- Generates modern ES modules with `@bufbuild/protobuf`
- Uses plugin: `@bufbuild/es` with ES modules and JSON types
- Output: `login/packages/zitadel-proto/`
- Used for: New user v2 API and services

### Dependency Management

The Console's `project.json` ensures proper execution order:

1. `@zitadel/proto:generate` runs first (modern ES modules)
2. Console's local generation runs second (traditional protobuf)
3. Build/lint/start tasks depend on both generations being complete

This approach allows the Console app to use both v1 and v2 APIs while maintaining proper build dependencies.

## Legacy Information

This project was originally generated with Angular CLI version 8.3.20 and has been updated over time.
