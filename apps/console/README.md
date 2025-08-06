# Console Angular App

This is the ZITADEL Console Angular application.

## Development

### Installation

```bash
pnpm install
```

This automatically runs both generations in the correct order via Turbo dependencies.

### Development Server

To start the development server:

```bash
nx run @zitadel/console:dev
```

This will:

1. Fetch the environment configuration from the server
2. Serve the app on the default port

### Building

To build for production:

```bash
nx run @zitadel/console:build
```

This will:

1. Generate proto files (via `prebuild` script)
2. Build the Angular app with production optimizations

### Linting

To run linting and formatting checks:

```bash
nx run @zitadel/console:lint
```

To auto-fix formatting issues:

```bash
nx run @zitadel/console:lint:fix