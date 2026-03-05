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

This opens your browser for authentication, exchanges an authorization code using PKCE, and stores the token locally.

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
