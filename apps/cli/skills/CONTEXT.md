# ZITADEL CLI — Agent Context

> This file provides structured guidance for AI agents interacting with the ZITADEL CLI.
> It encodes invariants and best practices that cannot be inferred from `--help` alone.

## General Rules

1. **Always use `--dry-run` before mutating operations** — Review the JSON envelope to confirm method and payload before executing.
2. **Use `--output json`** for programmatic consumption — Table output is for humans.
3. **Use `--from-json` or `--request-json` for complex payloads** — Don't try to construct deeply nested objects via flags alone.
4. **Resource IDs are opaque strings** — Usually UUIDs. Never construct, parse, or modify them. Always use IDs returned by `list` or `get` commands.
5. **Use `describe` for schema introspection** — `zitadel describe <group> <verb>` returns request/response JSON schemas.

## Authentication

- Before any API call, run `zitadel-cli login` to set up a context.
- The CLI stores credentials in `~/.zitadel-cli/contexts/`.
- If you get `[unauthenticated]`, re-run `zitadel-cli login`.
- If you get `[permission_denied]`, the service account lacks grants.

## Common Workflows

### Create a Human User
```bash
# organizationId is REQUIRED for all mutating operations
zitadel-cli users create human \
  --organization-id "<org-id>" \
  --username "jane" \
  --given-name "Jane" \
  --family-name "Doe" \
  --email "jane@example.com"
```

### List Organizations
```bash
zitadel-cli orgs list --output json
```

### Get a Resource by ID
```bash
zitadel-cli users get <user-id>
```

### Complex Requests via JSON
```bash
zitadel-cli users create --request-json '{
  "username": "jane",
  "organizationId": "<org-id>",
  "human": {
    "profile": {"givenName": "Jane", "familyName": "Doe"},
    "email": {"email": "jane@example.com", "isVerified": true}
  }
}'
```

## Important Constraints

- **`organizationId` is required** — All mutating operations require an organization ID. Get it from `orgs list`.
- **IDs are immutable** — You cannot change a resource's ID after creation.
- **Organization isolation** — Resources belong to the organization of the authenticated context. Use `--context` to switch.
- **Flag aliases** — Nested fields support short aliases (e.g., `--given-name` for `--profile-given-name`). Use `--help` to see all flags.
- **Enum values** — Always use the proto enum name (e.g., `ACCESS_TOKEN_TYPE_BEARER`), not the numeric value. Use `describe` to see valid values.
- **Timestamps** — Returned as RFC 3339 strings in JSON mode (e.g., `2024-01-15T10:30:00Z`).
- **Pagination** — List commands accept `--limit` and `--offset` flags. Default limit varies by endpoint.

## Error Handling

Errors follow the pattern `[code] message`. Common codes:

| Code | Meaning | Action |
|------|---------|--------|
| `unauthenticated` | Token expired or missing | Re-run `zitadel-cli login` |
| `permission_denied` | Missing IAM/Org grants | Check service account roles |
| `not_found` | Resource doesn't exist | Verify the ID with a `list` command |
| `already_exists` | Duplicate resource | Use a unique identifier |
| `failed_precondition` | Business rule violation | Read the error message for details |
