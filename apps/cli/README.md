# ZITADEL Management CLI

A command-line interface for managing ZITADEL resources via the v2 APIs. Designed for both humans and AI agents — every command supports JSON input/output, machine-readable discovery via `describe`, and structured error codes.

## Installation

Build from the repository root (recommended):

```bash
pnpm nx run @zitadel/cli:build
# Binary is output to: .artifacts/bin/<GOOS>/<GOARCH>/zitadel-cli
```

Or build manually:

```bash
cd apps/cli && go build -o zitadel-cli .
```

## Quick Start

### Using a Personal Access Token (PAT)

```bash
export ZITADEL_TOKEN=your-pat-here
zitadel-cli context add --instance mycompany.zitadel.cloud --token "$ZITADEL_TOKEN"
zitadel-cli users list
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

### Session management

| Command        | Description                                    |
| -------------- | ---------------------------------------------- |
| `login`        | Authenticate via browser-based OIDC PKCE flow  |
| `logout`       | Clear stored token for the active context      |
| `context`      | Manage CLI contexts (instances + credentials)  |
| `describe`     | Describe commands as machine-readable JSON     |

### Resource management (v2 API)

All 15 ZITADEL v2 service groups are available. Methods marked `[DEPRECATED]` in the help output are deprecated in the API and will print a warning at runtime.

| Command group    | Description                                      |
| ---------------- | ------------------------------------------------ |
| `actions`        | Actions and execution targets                    |
| `apps`           | OIDC, SAML, and API applications                 |
| `authorizations` | User authorizations (grants)                     |
| `features`       | Instance and organization feature flags          |
| `groups`         | User groups                                      |
| `idps`           | Identity provider links                          |
| `instances`      | ZITADEL instance management                      |
| `oidc`           | OIDC introspection and token exchange            |
| `orgs`           | Organizations                                    |
| `projects`       | Projects and project grants                      |
| `saml`           | SAML service provider metadata                   |
| `sessions`       | User sessions                                    |
| `settings`       | Instance and organization settings               |
| `users`          | Users, passkeys, MFA, PATs, keys, metadata       |
| `webkeys`        | Web keys for OIDC/SAML signing                   |

Use `zitadel-cli <group> --help` to see all subcommands for a group, and `zitadel-cli <group> <command> --help` for flag details.

## Agent-friendly workflows

The CLI supports machine-driven discovery and request composition:

```bash
# Discover all command groups, commands, and global flags
zitadel-cli describe

# List all commands in a group with full metadata
zitadel-cli describe users

# Inspect one command (flags, types, enum values, and JSON template)
zitadel-cli describe users create human
```

The top-level `describe` output includes a `global_flags` array with `--from-json`, `--request-json`, `--dry-run`, `--output`, and `--context` — these are available on every command.

### JSON template

Every command's describe output includes a `json_template` field showing the full request shape with zero/placeholder values. This reveals nested fields (like `password`, `email.sendCode`) that aren't available as individual CLI flags:

```bash
zitadel-cli describe users create human | jq .json_template
```

For variant commands (e.g., `set-email send-code`), the template is filtered to show **only the chosen variant's fields** — no noise from sibling branches:

```bash
# Only shows sendCode, not returnCode or isVerified:
zitadel-cli describe users set-email send-code | jq .json_template
```

### Nested oneof flags

When a variant contains a small inner oneof (like `password_type`), its fields are promoted to CLI flags so you don't need `--request-json` for common operations:

```bash
# Set a plaintext password:
zitadel-cli users create human --given-name Alice --family-name Doe \
  --email alice@example.com --password s3cret!

# Or use a pre-hashed password:
zitadel-cli users create human --given-name Alice --family-name Doe \
  --email alice@example.com --hashed-password-hash '$2a$12$...'
```

These flags are mutually exclusive — you can't combine `--password` with `--hashed-password-hash`.

### Providing JSON payloads

Two options for complex requests:

**Inline JSON** (recommended for agents — no stdin piping needed):

```bash
zitadel-cli users create human --request-json '{"username":"alice","human":{"profile":{"givenName":"Alice","familyName":"Doe"}}}' --dry-run
```

**Stdin JSON**:

```bash
cat request.json | zitadel-cli users create human --from-json --dry-run
```

When `--from-json` or `--request-json` is set, required request fields can be supplied from JSON instead of individual command flags.

### Pagination

List commands expose `--offset`, `--limit`, and `--asc` flags:

```bash
zitadel-cli users list --limit 10 --offset 0 --asc
```

### Dry-run and structured output

`--dry-run` prints the normalized request envelope without calling the API, which is useful for validating generated payloads.

When stdout is piped (non-TTY), output automatically switches to JSON and errors are emitted as structured JSON on stderr with machine-readable error codes.

### Deprecated commands

Methods deprecated in the ZITADEL API are marked `[DEPRECATED]` in the help listing and print a warning to stderr when invoked. Prefer the suggested replacement from the API documentation.

## Configuration

Config is stored at `~/.config/zitadel/config.toml` (respects `$XDG_CONFIG_HOME`).

### Environment Variables

| Variable           | Description                            |
| ------------------ | -------------------------------------- |
| `ZITADEL_TOKEN`    | PAT — overrides any configured token   |
| `ZITADEL_INSTANCE` | Override the instance URL              |

## Transport

The CLI uses the **ConnectRPC** protocol (HTTP/1.1 or HTTP/2 with `application/connect+proto`). No gRPC dial is required — requests go through a plain `*http.Client` with a Bearer-token transport. This means the CLI works through standard HTTP proxies and does not require `h2c` or TLS-based gRPC.

## Development

### Regenerate proto stubs

CLI commands are generated from proto definitions. To regenerate:

```bash
pnpm nx run @zitadel/cli:generate
```

To add a new v2 service, add an entry to `v2ServiceFilter` in `internal/protoc/protoc-gen-zitadelcli/main.go` and the corresponding `--path` flag in `apps/cli/project.json`. See `apps/cli/AGENTS.md` for full guidance.

### Build

```bash
pnpm nx run @zitadel/cli:build
```

### Test

```bash
pnpm nx run @zitadel/cli:test
# or directly: cd apps/cli && go test ./...
```
