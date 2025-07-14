# Unified Domains Table Implementation

This implementation provides a unified domains table (`zitadel.domains`) that consolidates both organization and instance domains into a single table structure.

## Architecture

The implementation follows Zitadel's established patterns:

### Database Layer
- **Migration 61**: Creates the unified domains table with proper constraints
- **Table Structure**: Uses nullable `org_id` to distinguish between instance domains (NULL) and organization domains (NOT NULL)
- **Soft Deletes**: Implements `deleted_at` for data preservation
- **Unique Constraints**: Ensures domain uniqueness within instance/organization scope

### Domain Layer
- **Interfaces**: Clean separation between instance and organization domain operations
- **Models**: Unified domain model with optional organization ID

### Repository Layer
- **Implementation**: Single repository handling both instance and organization domains
- **Transactions**: Atomic operations for primary domain changes
- **Query Building**: Type-safe SQL generation using squirrel

### Projection Layer
- **Event Handling**: Processes both org and instance domain events
- **Data Synchronization**: Maintains consistency with event sourcing

## Event Mapping

### Organization Domain Events
- `org.domain.added` → Creates domain record with org_id
- `org.domain.verification.added` → Updates validation_type
- `org.domain.verified` → Sets is_verified = true
- `org.domain.primary.set` → Manages primary domain flags
- `org.domain.removed` → Soft deletes domain
- `org.removed` → Soft deletes all org domains

### Instance Domain Events
- `instance.domain.added` → Creates domain record with org_id = NULL, is_verified = true
- `instance.domain.primary.set` → Manages primary domain flags
- `instance.domain.removed` → Soft deletes domain
- `instance.removed` → Soft deletes all instance domains

## Usage Examples

### Instance Domain Operations
```go
// Add instance domain (always verified)
domain, err := repo.AddInstanceDomain(ctx, "instance-123", "api.example.com")

// Set primary instance domain
err := repo.SetInstanceDomainPrimary(ctx, "instance-123", "api.example.com")

// Remove instance domain
err := repo.RemoveInstanceDomain(ctx, "instance-123", "api.example.com")

// List instance domains
criteria := v2domain.DomainSearchCriteria{
    InstanceID: &instanceID,
}
pagination := v2domain.DomainPagination{
    Limit: 10,
    SortBy: v2domain.DomainSortFieldDomain,
    Order: database.SortOrderAsc,
}
list, err := repo.List(ctx, criteria, pagination)
```

### Organization Domain Operations
```go
// Add organization domain
domain, err := repo.AddOrganizationDomain(ctx, "instance-123", "org-456", "company.com", domain.OrgDomainValidationTypeHTTP)

// Verify organization domain
err := repo.SetOrganizationDomainVerified(ctx, "instance-123", "org-456", "company.com")

// Set primary organization domain
err := repo.SetOrganizationDomainPrimary(ctx, "instance-123", "org-456", "company.com")

// Find specific domain
criteria := v2domain.DomainSearchCriteria{
    Domain: &domainName,
    InstanceID: &instanceID,
}
domain, err := repo.Get(ctx, criteria)
```

## Testing

The implementation includes comprehensive tests:

### Repository Tests
- CRUD operations for both instance and organization domains
- Transaction behavior verification
- Error handling and edge cases
- SQL query parameter validation

### Projection Tests
- Event reduction logic for all supported events
- Column and condition verification
- Multi-statement handling for primary domain changes
- Soft delete behavior

## Migration Strategy

This table is designed to work alongside existing tables initially:

1. **Phase 1** (This PR): Create unified table and maintain via projections
2. **Phase 2** (Future): Migrate query logic to use unified table
3. **Phase 3** (Future): Deprecate separate org_domains2 and instance_domains tables

## Performance Considerations

- **Indexing**: The unique constraint provides efficient domain lookups
- **Queries**: Nullable org_id allows efficient filtering between domain types
- **Pagination**: Supports sorting by created_at, updated_at, and domain name
- **Soft Deletes**: WHERE deleted_at IS NULL conditions optimize active domain queries

## Validation

The table enforces:
- Domain length between 1-255 characters
- Non-negative validation_type values
- Foreign key integrity to instances and organizations
- Unique domain constraints per instance/organization scope
- Automatic updated_at timestamp management