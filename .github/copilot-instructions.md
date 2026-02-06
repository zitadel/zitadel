# ZITADEL Development Guide for GitHub Copilot

## Repository Overview

ZITADEL is an open-source identity and access management platform with a modern tech stack:
- **Backend/API**: Go (1.24+)
- **Login UI**: Next.js/React
- **Management Console**: Angular
- **Documentation**: Docusaurus
- **Build System**: Nx monorepo with pnpm

## Architecture & Project Structure

This is an Nx monorepo with the following main projects:

### Applications (`apps/`)
- `@zitadel/api` - Go-based backend API (gRPC, REST, OpenID Connect, SAML)
- `@zitadel/login` - Next.js-based login interface
- `@zitadel/docs` - Docusaurus documentation site

### Console
- `@zitadel/console` - Angular-based management console UI (in `console/` directory)

### Packages
- `@zitadel/proto` - Protocol buffer definitions
- `@zitadel/client` - Client libraries

### Backend Structure
- `cmd/` - CLI commands and entry points
- `internal/` - Internal Go packages (main business logic)
- `backend/` - Backend services
- `pkg/` - Public Go packages
- `proto/` - Protocol buffer definitions

## Development Setup

### Prerequisites
- Node.js v22.x
- Go 1.24.x
- Docker (for database and services)
- pnpm (via Corepack: `corepack enable`)

### Initial Setup
```bash
pnpm install
pnpm nx run-many --target generate
```

## Common Development Commands

All commands follow the Nx pattern: `pnpm nx run PROJECT:TARGET`

### Development Servers
```bash
# Start API in development mode
pnpm nx run @zitadel/api:dev

# Start Login UI in development mode
pnpm nx run @zitadel/login:dev

# Start Console in development mode
pnpm nx run @zitadel/console:dev

# Start Documentation site
pnpm nx run @zitadel/docs:dev
```

### Building
```bash
# Build API (requires Console to be built first)
pnpm nx run @zitadel/api:build

# Build Login
pnpm nx run @zitadel/login:build

# Build Console
pnpm nx run @zitadel/console:build
```

### Code Generation
```bash
# Generate all code (proto, stubs, static files)
pnpm nx run @zitadel/api:generate

# Generate Go-specific files (Stringer, Enumer, Mockgen)
pnpm nx run @zitadel/api:generate-go
```

### Testing
```bash
# Run unit tests
pnpm nx run @zitadel/api:test-unit
pnpm nx run @zitadel/login:test-unit

# Run integration tests
pnpm nx run @zitadel/api:test-integration

# Stop integration test containers
pnpm nx run @zitadel/api:test-integration-stop
```

### Linting
```bash
# Lint Go code (uses golangci-lint)
pnpm nx run @zitadel/api:lint

# Lint Next.js Login
pnpm nx run @zitadel/login:lint

# Lint Angular Console
pnpm nx run @zitadel/console:lint
```

## Code Style & Conventions

### Go
- Follow standard Go conventions and idioms
- Use `golangci-lint` for linting (config in `.golangci.yaml`)
- Write unit tests alongside code (`_test.go` files)
- Integration tests go in `internal/integration/`
- Use event sourcing patterns (ZITADEL uses event sourcing)

### TypeScript/JavaScript
- Use Prettier for formatting
- Follow ESLint rules
- Use TypeScript strict mode
- Prefer functional components in React
- Use Angular best practices for Console

### Project-Specific Patterns
- **Event Sourcing**: ZITADEL uses event sourcing for data storage
- **Multi-tenancy**: Consider multi-tenant architecture in all features
- **API-First**: gRPC and REST APIs are the primary interfaces
- **Proto Files**: When modifying APIs, update `.proto` files in `proto/` directory
- **Console Embedding**: The Console is built and embedded into the API binary

## Code Generation Workflow

When modifying proto files or Go code that uses generation:

1. Update `.proto` files in `proto/` directory
2. Run `pnpm nx run @zitadel/api:generate-stubs` to generate gRPC/OpenAPI stubs
3. For Go code using `//go:generate` directives, run `pnpm nx run @zitadel/api:generate-go`
4. Generated files include:
   - gRPC stubs in `pkg/grpc/`
   - OpenAPI specs in `openapi/v2/zitadel/`
   - Statik embedded files
   - Asset routes

## Testing Guidelines

### Unit Tests
- Write alongside code in `*_test.go` files
- Use table-driven tests where appropriate
- Mock external dependencies
- Test coverage is tracked

### Integration Tests
- Located in `internal/integration/`
- Test against running API instance
- Use tag `integration` for build: `// +build integration`
- Run with `pnpm nx run @zitadel/api:test-integration`

### UI Tests
- Cypress tests for Console in `console/` 
- Vitest for Login unit tests
- Run functional UI tests with `pnpm nx run @zitadel/functional-ui:test`

## Database & Storage

- Primary database: PostgreSQL (v14+)
- Event sourcing pattern for data storage
- Database migrations handled automatically
- Integration tests use Docker containers for DB

## Security Considerations

- ZITADEL is an identity platform - security is paramount
- Follow OWASP best practices
- Validate all inputs
- Use secure defaults
- Report security issues to security@zitadel.com
- Never commit secrets or credentials

## API Development

### gRPC & REST
- Both gRPC and REST APIs are generated from proto definitions
- gRPC is the primary interface, REST via gRPC-gateway
- OpenAPI documentation is auto-generated

### Authentication Flows
- OpenID Connect certified
- SAML 2.0 support
- OAuth 2.x support  
- Passkeys/FIDO2 support
- Multi-factor authentication

## Console Development

The Management Console is an Angular application that:
- Is built separately but embedded in the API binary
- Located in `console/` directory
- Served by the API at runtime
- Built static files copied to `internal/api/ui/console/static/`

## Troubleshooting

### Nx Daemon Issues
If commands hang:
```bash
pnpm nx daemon --stop
```

### Clean Build
```bash
pnpm run clean:all  # Removes .nx, node_modules
```

### Generation Issues
Ensure all generators are installed:
```bash
pnpm nx run @zitadel/api:generate-install
```

## Documentation

- Main docs: https://zitadel.com/docs/
- API docs: https://zitadel.com/docs/apis/introduction
- Contributing guide: `CONTRIBUTING.md`
- Architecture: Nx monorepo with event sourcing backend

## Important Notes

- The API embeds the Console static files - always build Console before API for production
- Proto files are the source of truth for APIs
- Use Nx commands, not direct npm/go commands
- Integration tests require Docker
- The repo uses pnpm workspace with Nx
- Code generation is required after proto or certain Go file changes
