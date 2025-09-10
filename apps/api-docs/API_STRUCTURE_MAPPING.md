# ZITADEL API Documentation Structure Mapping

This document explains how the ZITADEL API documentation app maps and displays the API structure, helping developers understand the organization and filtering logic.

## Overview

The ZITADEL API documentation app uses Scalar to render interactive API documentation from OpenAPI specifications generated from Protocol Buffer definitions. The app supports both v1 and v2 API architectures with intelligent filtering and organization.

## Directory Structure

### Source Files

```
.artifacts/openapi3/zitadel/
├── management.openapi.yaml           # v1 Management API
├── admin.openapi.yaml               # v1 Admin API
├── auth.openapi.yaml                # v1 Authentication API
├── system.openapi.yaml              # v1 System API
├── action/
│   └── v1/
│       └── action_service.openapi.yaml
├── user/
│   ├── v1/
│   │   └── user_service.openapi.yaml
│   └── v2/
│       ├── user_service.openapi.yaml      # v2 User API (contains endpoints)
│       ├── user.openapi.yaml              # Schema only (filtered out)
│       └── auth.openapi.yaml              # Schema only (filtered out)
├── org/
│   └── v2/
│       ├── org_service.openapi.yaml       # v2 Organization API
│       └── org.openapi.yaml               # Schema only (filtered out)
└── [other services]/
    └── v[version]/
        ├── [service]_service.openapi.yaml  # Service endpoints
        └── [resource].openapi.yaml         # Schema definitions
```

## File Filtering Logic

### Server-Side Filtering (`/api/openapi/route.ts`)

The API route recursively scans the `.artifacts/openapi3/zitadel/` directory and applies the following filters:

1. **Service Files**: Include files ending with `_service.openapi.yaml`

   - These contain actual API endpoints and operations
   - Found in nested directories like `user/v2/user_service.openapi.yaml`

2. **Top-Level v1 APIs**: Include root-level `.openapi.yaml` files

   - Legacy v1 APIs: `management.openapi.yaml`, `admin.openapi.yaml`, etc.
   - These are complete API specifications with endpoints

3. **Excluded Files**: Skip schema-only files
   - Files like `user.openapi.yaml`, `org.openapi.yaml`
   - Contain only data type definitions without endpoints

```typescript
// Server-side filtering logic
if (entry.endsWith("_service.openapi.yaml")) {
  // Include service files with endpoints
  files.push({ path: fullPath, relativePath: entryRelativePath });
} else if (entry.endsWith(".openapi.yaml") && relativePath === "") {
  // Include top-level v1 API files
  files.push({ path: fullPath, relativePath: entryRelativePath });
}
```

### Client-Side Filtering (`ApiReference.tsx`)

Additional filtering on the client ensures only services with actual endpoints are displayed:

```typescript
// Filter out specs with no endpoints (only schema definitions)
const specsWithEndpoints = data.specs.filter((spec) => {
  try {
    const parsed = yaml.load(spec.content) as any;
    return parsed.paths && Object.keys(parsed.paths).length > 0;
  } catch {
    return false;
  }
});
```

This removes any remaining files that passed server-side filtering but contain no API endpoints.

## Service Name Mapping

### v1 APIs (Legacy)

Direct mapping from filename to display name:

| File Name    | Display Name            | Description               |
| ------------ | ----------------------- | ------------------------- |
| `management` | Management API (v1)     | Core resource management  |
| `admin`      | Admin API (v1)          | Administrative operations |
| `auth`       | Authentication API (v1) | Authentication flows      |
| `system`     | System API (v1)         | System-level operations   |

### v2+ APIs (Versioned Services)

Dynamic mapping from nested path structure:

| File Path                  | Service Name               | Display Name  | Description                |
| -------------------------- | -------------------------- | ------------- | -------------------------- |
| `user/v2/user_service`     | `user/v2/user_service`     | User API V2   | User management operations |
| `org/v2/org_service`       | `org/v2/org_service`       | Org API V2    | Organization management    |
| `action/v1/action_service` | `action/v1/action_service` | Action API V1 | Custom actions             |

```typescript
const getServiceDisplayName = (serviceName: string): string => {
  if (serviceName.includes("/")) {
    const parts = serviceName.split("/");
    const service = parts[0]; // e.g., "user"
    const version = parts[1]; // e.g., "v2"

    return `${
      service.charAt(0).toUpperCase() + service.slice(1)
    } API ${version.toUpperCase()}`;
  }

  // Handle v1 services with direct mapping
  const nameMap = {
    management: "Management API (v1)",
    admin: "Admin API (v1)",
    // ...
  };

  return nameMap[serviceName] || serviceName;
};
```

## Service Sorting Logic

Services are sorted to provide a logical organization in the navigation:

1. **Primary Sort**: Alphabetical by service name (`action`, `admin`, `auth`, `management`, `org`, `user`)
2. **Secondary Sort**: Version preference (v2+ before v1)
   - `user/v2/user_service` appears before any v1 user APIs
   - `management` (v1) appears after `user/v2/user_service`

```typescript
const sortedSpecs = [...specs].sort((a, b) => {
  const aService = a.name.split("/")[0];
  const bService = b.name.split("/")[0];

  if (aService !== bService) {
    return aService.localeCompare(bService); // Alphabetical by service
  }

  // Same service: v2+ before v1
  const aIsV1 = !a.name.includes("/");
  const bIsV1 = !b.name.includes("/");

  if (aIsV1 && !bIsV1) return 1; // v1 after v2+
  if (!aIsV1 && bIsV1) return -1; // v2+ before v1

  return a.name.localeCompare(b.name);
});
```

## Default Service Selection

The app intelligently selects a default service to display:

1. **Preferred**: v2 User service (`user/v2/user_service`)
2. **Fallback**: v1 Management service (`management`)
3. **Last Resort**: First available service in sorted list

```typescript
const userV2Service = specsWithEndpoints.find((spec) =>
  spec.name.includes("user/v2/user_service")
);
const managementService = specsWithEndpoints.find(
  (spec) => spec.name === "management"
);

if (userV2Service) {
  setSelectedSpec(userV2Service.name);
} else if (managementService) {
  setSelectedSpec("management");
} else {
  setSelectedSpec(specsWithEndpoints[0].name);
}
```

## API Generation Process

### 1. Protocol Buffer Compilation

```bash
# Generate OpenAPI specs from .proto files
pnpm run generate
```

### 2. buf Configuration (`buf.gen.yaml`)

```yaml
version: v1
plugins:
  - plugin: buf.build/grpc-ecosystem/openapiv2
    out: .artifacts/openapi3
    opt:
      - openapi_naming_strategy=fqn
      - generate_unbound_methods=false
      - use_go_templates=true
```

### 3. Output Structure

Generated files follow the pattern:

- `{service_name}.openapi.yaml` for v1 APIs
- `{domain}/{version}/{service}_service.openapi.yaml` for v2+ APIs
- `{domain}/{version}/{resource}.openapi.yaml` for schema definitions

## Debugging and Development

### Console Logging

The app includes comprehensive logging for debugging:

```typescript
console.log("Selected spec:", selectedSpec);
console.log("Parsed spec:", parsedSpec);
console.log("Parsed spec paths:", parsedSpec.paths);
console.log("Number of paths:", Object.keys(parsedSpec.paths || {}).length);
```

### Common Issues

1. **Empty Navigation**: No services appear

   - **Cause**: Missing generated files
   - **Solution**: Run `pnpm run generate`

2. **Services Without Endpoints**: Services appear but show no content

   - **Cause**: Schema-only files passed filtering
   - **Solution**: Check client-side filtering logic

3. **Missing v2 Services**: Only v1 APIs appear
   - **Cause**: Server-side filtering not including `_service.openapi.yaml` files
   - **Solution**: Verify file naming convention and recursive directory scanning

## Future Considerations

1. **API Versioning**: As new versions are added, the sorting logic will automatically prefer newer versions
2. **Service Categories**: Consider grouping services by domain (user management, admin, etc.)
3. **Search Functionality**: Add service filtering/search in the navigation sidebar
4. **Customization**: Allow users to configure default service and display preferences

## File References

- **Main Component**: `src/components/ApiReference.tsx`
- **API Route**: `src/app/api/openapi/route.ts`
- **Generation Config**: `buf.gen.yaml`
- **Build Configuration**: `package.json` (scripts section)
- **Type Definitions**: TypeScript interfaces in component files

This mapping ensures that developers see a clean, organized view of the ZITADEL APIs while automatically filtering out non-functional schema definitions and maintaining logical service organization.
