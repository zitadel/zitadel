# ZITADEL Internal Backend Guide for AI Agents

## Context
`internal/` contains core backend domain logic for ZITADEL: commands, queries, repositories, eventstore integration, API service layers, and supporting infrastructure.

## Source of Truth
- **Go Toolchain**: Inspect root `go.mod` before Go work.
- **Architecture Pattern**: Relational data is the system of record; keep existing event writes that provide history/audit trails.
- **API Contract**: For API-facing schema decisions, follow `API_DESIGN.md` and `proto/AGENTS.md`.
- **Human Contribution Guide**: See `CONTRIBUTING.md` for setup, i18n workflow, and testing requirements.

## Architecture Overview

### V2 vs V3 Backend
ZITADEL has two backend architectures in use:

- **V2 (Current - `internal/`)**: Context-based API with eventstore patterns. Most of the current codebase.
- **V3 (New - `backend/v3/`)**: Hexagonal architecture with command pattern, repository pattern, and improved separation of concerns. See `backend/v3/doc.go` for details.

**When to use:**
- **V2**: Continue maintaining existing features in `internal/` using established patterns (commands, queries, eventstore).
- **V3**: New features should consider V3 architecture if starting fresh. Coordinate with maintainers on migration strategy.

### Key Packages
- **`internal/command/`**: Write operations (mutations, state changes)
- **`internal/query/`**: Read operations (queries, projections)
- **`internal/domain/`**: Domain models and business logic
- **`internal/repository/`**: Data access abstractions
- **`internal/eventstore/`**: Event persistence and streaming
- **`internal/api/`**: Transport adapters (gRPC/connectRPC handlers)
- **`internal/zerrors/`**: Structured error handling with gRPC status code mapping

## Command/Query Pattern (V2)

### Command Structure
Commands modify state and return events. The following is a **simplified pseudocode example** illustrating the pattern (see the actual implementation in `internal/command/user.go` for full details including domain policy and org-scoped username checks):

```go
func (c *Commands) ChangeUsername(ctx context.Context, orgID, userID, userName string) (*domain.ObjectDetails, error) {
    // 1. Validate inputs
    userName = strings.TrimSpace(userName)
    if orgID == "" || userID == "" || userName == "" {
        return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2N9fs", "Errors.IDMissing")
    }

    // 2. Load current state via write model
    existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
    if err != nil {
        return nil, err
    }

    // 3. Check preconditions
    if !isUserStateExists(existingUser.UserState) {
        return nil, zerrors.ThrowNotFound(nil, "COMMAND-5N9ds", "Errors.User.NotFound")
    }

    // 4. Apply business logic checks
    if existingUser.UserName == userName {
        return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6m9gs", "Errors.User.UsernameNotChanged")
    }

    // 5. Create aggregate and push events
    userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
    pushedEvents, err := c.eventstore.Push(ctx,
        user.NewUsernameChangedEvent(ctx, userAgg, existingUser.UserName, userName, ...))
    if err != nil {
        return nil, err
    }

    // 6. Update model and return
    err = AppendAndReduce(existingUser, pushedEvents...)
    if err != nil {
        return nil, err
    }
    return writeModelToObjectDetails(&existingUser.WriteModel), nil
}
```

**Key Principles:**
- **Structural input validation** belongs at the API/adapter layer (proto validate handles most of it; anything proto validate cannot express should be validated in the gRPC handler / converter, not in the command).
- Commands handle **business-rule enforcement** (e.g., "username not changed", "user not found") that must run against current state.
- Load write model to get current state
- Check business rules
- Push events to eventstore
- Return `*domain.ObjectDetails` (sequence, change date, resource owner)

### Query Structure
Queries read data without side effects. Located in `internal/query/`. Pattern:
- Use prepared SQL statements (see `.sql` files in `internal/query/`)
- Return view models, not domain models
- Support pagination, filtering, and sorting
- No state modification

## Error Handling

### Using `internal/zerrors/`
All errors MUST use the `zerrors` package for consistent error handling and gRPC status code mapping.

**Error Kinds** (mapped to gRPC status codes — see `internal/zerrors/zerror.go` for the full list):
- `KindInvalidArgument` → `INVALID_ARGUMENT` — use for validation failures
- `KindNotFound` → `NOT_FOUND` — resource does not exist
- `KindAlreadyExists` → `ALREADY_EXISTS` — resource already exists (conflict)
- `KindPermissionDenied` → `PERMISSION_DENIED` — authorization failure
- `KindUnauthenticated` → `UNAUTHENTICATED` — authentication failure
- `KindPreconditionFailed` → `FAILED_PRECONDITION` — business rule violation
- `KindInternal` → `INTERNAL` — system / unexpected failures
- `KindUnavailable` → `UNAVAILABLE` — service temporarily unavailable

**Usage Pattern:**
```go
// Wrap existing error
if err != nil {
    return nil, zerrors.ThrowInternal(err, "COMMAND-mF9ds", "Errors.Internal")
}

// Create new error without wrapping
if userName == "" {
    return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2N9fs", "Errors.User.UsernameMissing")
}

// Error codes: "PACKAGE-ID123" format for traceability
// Message keys: "Errors.Domain.Reason" for i18n
```

## Testing Patterns

### Testify
Use `github.com/stretchr/testify` for assertions in all Go tests:
- **`assert`** (`github.com/stretchr/testify/assert`): non-fatal assertions — test continues on failure (`assert.Equal()`, `assert.NotZero()`, `assert.Nil()`, etc.).
- **`require`** (`github.com/stretchr/testify/require`): fatal gates — test stops immediately on failure (`require.NoError()`, `require.NotNil()`, etc.).

Prefer `require` when subsequent assertions are meaningless if the guarded condition fails.

### Unit Tests
Use table-driven tests throughout the Go codebase (command, query, API handlers, converters, etc.):

```go
func TestChangeUsername(t *testing.T) {
    tests := []struct {
        name    string
        args    args
        want    *domain.ObjectDetails
        wantErr bool
    }{
        {
            name: "success",
            args: args{orgID: "org1", userID: "user1", userName: "newname"},
            want: &domain.ObjectDetails{...},
            wantErr: false,
        },
        {
            name: "user not found",
            args: args{orgID: "org1", userID: "unknown", userName: "newname"},
            want: nil,
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests
- Located in `internal/integration/`
- Use test database fixtures
- Test full command → event → query flow
- Run via: `pnpm nx run @zitadel/api:test-integration`

## Repository Pattern
Repositories abstract data access:
- Define interfaces in domain packages
- Implement in `internal/repository/` or `internal/query/projection/`
- Use prepared statements for performance
- Support transactions where needed

## Security & Permissions
- **Authentication checks**: Use `internal/api/grpc/management/auth_checks.go` patterns
- **Permission checks**: Use `internal/command/permission_checks.go` patterns
- **Crypto operations**: Use `internal/crypto/` package for encryption/hashing
- **Never log sensitive data**: Passwords, tokens, secrets

## Observability
- **Logging**: Use `backend/v3/instrumentation/logging` package
- **Tracing**: Import `backend/v3/instrumentation/tracing` for distributed tracing
- **Context propagation**: Always pass `context.Context` through call chains

> **Note:** Existing V2 code (`internal/`) still uses `github.com/zitadel/logging` and the `internal/telemetry/tracing` shim (which delegates to the V3 tracing package). When working in V2 code, follow the existing import style; the V3 packages are the forward target.

## Boundary Rules
- Prefer implementing business behavior in command/query layers and repository packages, not in transport handlers.
- Avoid bypassing established event/repository flows with ad-hoc direct persistence patterns.
- Keep API/service adapters thin; place reusable domain behavior in internal domain packages.
- Database schema changes require coordination with maintainers (migration files, backwards compatibility).

## Validation Workflow
- Use API project targets to validate backend changes:
  - `pnpm nx run @zitadel/api:lint`
  - `pnpm nx run @zitadel/api:test-unit`
  - `pnpm nx run @zitadel/api:test-integration`

## Cross-References
- **API Design Principles**: See `API_DESIGN.md` for API-specific conventions
- **Development Setup**: See `CONTRIBUTING.md` for local environment, database setup
- **Proto Definitions**: See `proto/AGENTS.md` for API contract changes
- **V3 Architecture**: See `backend/v3/doc.go` for new hexagonal architecture details
- **Architecture Wiki**: Broader architecture documentation is maintained at https://github.com/zitadel/zitadel/wiki
