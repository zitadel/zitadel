# ZITADEL CLI

A command-line interface for managing ZITADEL instances.

## Installation

```bash
go install github.com/zitadel/zitadel/apps/cli@main
```

Or build from source:

```bash
cd apps/cli
go build -o zitadel ./...
```

## Quick Start

### Using a Personal Access Token (PAT)

```bash
export ZITADEL_TOKEN=your-pat-here
zitadel users list --instance mycompany.zitadel.cloud
```

### Interactive Browser Login (OIDC PKCE)

Create a **Native** application in your ZITADEL project, then:

```bash
zitadel login --instance mycompany.zitadel.cloud --client-id <native-app-client-id>
```

This opens your browser for authentication, exchanges an authorization code using PKCE, and stores the token locally. Refresh tokens are saved automatically — the CLI will transparently refresh expired access tokens without re-prompting.

## Multiple Contexts

The CLI supports multiple configured instances (contexts):

```bash
# After logging in, contexts are created automatically
zitadel login --instance prod.zitadel.cloud --client-id <id> --context prod
zitadel login --instance staging.zitadel.cloud --client-id <id> --context staging

# Switch between contexts
zitadel context use prod

# List all contexts
zitadel context list

# Show current context
zitadel context current
```

## Self-Hosted Instances

```bash
zitadel login --instance https://auth.internal --client-id <id> --context myinstance
zitadel context use myinstance
zitadel users list
```

## Commands

| Command            | Description                                   |
| ------------------ | --------------------------------------------- |
| `login`            | Authenticate via browser-based OIDC PKCE flow |
| `logout`           | Clear stored token for the active context      |
| `context list`     | List all configured contexts                   |
| `context use`      | Switch the active context                      |
| `context current`  | Show the active context                        |
| `users list`       | List users in the current instance             |
| `orgs current`     | Show the current organization                  |

## Agent-friendly workflows

The CLI supports machine-driven discovery and request composition:

```bash
# Discover all command groups, commands, and global flags
zitadel describe

# List all commands in a group with full metadata
zitadel describe users

# Inspect one command (flags, types, enum values, and JSON template)
zitadel describe users create human
```

The top-level `describe` output includes a `global_flags` array with `--from-json`, `--request-json`, `--dry-run`, `--output`, and `--context` — these are available on every command.

### JSON template

Every command's describe output includes a `json_template` field showing the full request shape with zero/placeholder values. This reveals nested fields (like `password`, `email.sendCode`) that aren't available as individual CLI flags:

```bash
zitadel describe users create human | jq .json_template
```

For variant commands (e.g., `set-email send-code`), the template is filtered to show **only the chosen variant's fields** — no noise from sibling branches:

```bash
# Only shows sendCode, not returnCode or isVerified:
zitadel describe users set-email send-code | jq .json_template
```

### Nested oneof flags

When a variant contains a small inner oneof (like `password_type`), its fields are promoted to CLI flags so you don't need `--request-json` for common operations:

```bash
# Set a plaintext password:
zitadel users create human --given-name Alice --family-name Doe \
  --email alice@example.com --password s3cret!

# Or use a pre-hashed password:
zitadel users create human --given-name Alice --family-name Doe \
  --email alice@example.com --hashed-password-hash '$2a$12$...'
```

These flags are mutually exclusive — you can't combine `--password` with `--hashed-password-hash`.

### Providing JSON payloads

Two options for complex requests:

**Inline JSON** (recommended for agents — no stdin piping needed):

```bash
zitadel users create human --request-json '{"username":"alice","human":{"profile":{"givenName":"Alice","familyName":"Doe"}}}' --dry-run
```

**Stdin JSON**:

```bash
cat request.json | zitadel users create human --from-json --dry-run
```

When `--from-json` or `--request-json` is set, required request fields can be supplied from JSON instead of individual command flags.

### Pagination

List commands expose `--offset`, `--limit`, and `--asc` flags:

```bash
zitadel users list --limit 10 --offset 0 --asc
```

### Dry-run and structured output

`--dry-run` prints the normalized request envelope without calling the API, which is useful for validating generated payloads.

When stdout is piped (non-TTY), output automatically switches to JSON and errors are emitted as structured JSON on stderr with machine-readable error codes.

## Configuration

Config is stored at `~/.config/zitadel/config.toml` (respects `$XDG_CONFIG_HOME`).

### Environment Variables

| Variable           | Description                            |
| ------------------ | -------------------------------------- |
| `ZITADEL_TOKEN`    | PAT — overrides any configured token   |
| `ZITADEL_INSTANCE` | Override the instance URL              |

## Development

### Regenerate proto stubs

Proto stubs are generated via `buf` and are gitignored. To regenerate:

```bash
# Using Nx
pnpm nx run @zitadel/cli:generate

# Or manually
cd apps/cli
PATH="../../.artifacts/bin/$(go env GOOS)/$(go env GOARCH):$PATH" buf generate ../../proto
```

### Build

```bash
pnpm nx run @zitadel/cli:build
```

### Test

```bash
cd apps/cli && go test ./...
```
