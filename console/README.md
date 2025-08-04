# Console Angular App

This is the ZITADEL Console Angular application.

## Development

### Prerequisites

- Node.js 18 or later
- pnpm (latest)

### Installation

```bash
pnpm install
```

### Proto Generation

The Console app uses **dual proto generation** with Turbo dependency management:

1. **`@zitadel/proto` generation**: Modern ES modules with `@bufbuild/protobuf` for v2 APIs
2. **Local `buf.gen.yaml` generation**: Traditional protobuf JavaScript classes for v1 APIs

The Console app's `turbo.json` ensures that `@zitadel/proto#generate` runs before the Console's own generation, providing both:

- Modern schemas from `@zitadel/proto` (e.g., `UserSchema`, `DetailsSchema`)
- Legacy classes from `src/app/proto/generated` (e.g., `User`, `Project`)

Generated files:

- **`@zitadel/proto`**: Modern ES modules in `login/packages/zitadel-proto/`
- **Local generation**: Traditional protobuf files in `src/app/proto/generated/`
  - TypeScript definition files (`.d.ts`)
  - JavaScript files (`.js`)
  - gRPC client files (`*ServiceClientPb.ts`)
  - OpenAPI/Swagger JSON files (`.swagger.json`)

To generate proto files:

```bash
pnpm run generate
```

This automatically runs both generations in the correct order via Turbo dependencies.

### Development Server

To start the development server:

```bash
pnpm start
```

This will:

1. Fetch the environment configuration from the server
2. Serve the app on the default port

### Building

To build for production:

```bash
pnpm run build
```

This will:

1. Generate proto files (via `prebuild` script)
2. Build the Angular app with production optimizations

### Linting

To run linting and formatting checks:

```bash
pnpm run lint
```

To auto-fix formatting issues:

```bash
pnpm run lint:fix
```

## Project Structure

- `src/app/proto/generated/` - Generated proto files (Angular-specific format)
- `buf.gen.yaml` - Local proto generation configuration
- `turbo.json` - Turbo dependency configuration for proto generation
- `prebuild.development.js` - Development environment configuration script

## Proto Generation Details

The Console app uses **dual proto generation** managed by Turbo dependencies:

### Dependency Chain

The Console app has the following build dependencies managed by Turbo:

1. `@zitadel/proto#generate` - Generates modern protobuf files
2. `@zitadel/client#build` - Builds the TypeScript gRPC client library
3. `console#generate` - Generates Console-specific protobuf files
4. `console#build` - Builds the Angular application

This ensures that the Console always has access to the latest client library and protobuf definitions.

### Legacy v1 API (Traditional Protobuf)

- Uses local `buf.gen.yaml` configuration
- Generates traditional Google protobuf JavaScript classes extending `jspb.Message`
- Uses plugins: `protocolbuffers/js`, `grpc/web`, `grpc-ecosystem/openapiv2`
- Output: `src/app/proto/generated/`
- Used for: Most existing Console functionality

### Modern v2 API (ES Modules)

- Uses `@zitadel/proto` package generation
- Generates modern ES modules with `@bufbuild/protobuf`
- Uses plugin: `@bufbuild/es` with ES modules and JSON types
- Output: `login/packages/zitadel-proto/`
- Used for: New user v2 API and services

### Dependency Management

The Console's `turbo.json` ensures proper execution order:

1. `@zitadel/proto#generate` runs first (modern ES modules)
2. Console's local generation runs second (traditional protobuf)
3. Build/lint/start tasks depend on both generations being complete

This approach allows the Console app to use both v1 and v2 APIs while maintaining proper build dependencies.

## Legacy Information

This project was originally generated with Angular CLI version 8.3.20 and has been updated over time.
