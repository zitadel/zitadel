# ZITADEL Management CLI Guide for AI Agents

## Context
The **Management CLI** (`apps/cli`) is a command-line tool for managing ZITADEL resources via the v2 APIs. It is designed for both humans and AI agents, with `--from-json`/`--request-json`/`--dry-run`/`--output json` flags and a `describe` command for machine-readable introspection.

> **Note:** This is the *management API CLI*, distinct from the built-in ZITADEL server CLI (`zitadel init/setup/start`). The management CLI lives in `apps/cli/`; the server CLI entry point is `cmd/main.go`.

## Architecture

```
internal/protoc/protoc-gen-zitadelcli/
  main.go            ← protoc plugin (Go): reads proto descriptors, builds methodData, calls template
  cmd.go.tmpl        ← Go template: generates command functions, flag wiring, table output

apps/cli/
  gen/cmd_*.go       ← generated (gitignored): one file per v2 service group
  gen/registry.go    ← generated: Register() wires all service commands into the root
  gen/connect_error.go ← generated: FormatConnectError() parses connectRPC errors
  cmd/root.go        ← cobra root command, global flags, describe subcommand
  internal/auth/     ← OIDC PKCE login flow + PAT token source
  internal/client/   ← *http.Client factory with Bearer token transport
  internal/config/   ← XDG config file (contexts, tokens)
  internal/output/   ← JSON output helper
  main.go            ← entry point
```

The generator runs at `pnpm nx run @zitadel/cli:generate`:
1. Installs `protoc-gen-zitadelcli` binary via `go install`
2. Runs `buf generate --path <service>.proto` for each v2 service
3. Outputs `apps/cli/gen/cmd_*.go` (deleted and regenerated from scratch each run)

## Verified Nx Targets
- **Generate**: `pnpm nx run @zitadel/cli:generate` — regenerates `gen/cmd_*.go`
- **Build**: `pnpm nx run @zitadel/cli:build` — produces `.artifacts/bin/<GOOS>/<GOARCH>/zitadel-cli`
- **Lint**: `pnpm nx run @zitadel/cli:lint`
- **Test**: `pnpm nx run @zitadel/cli:test`

## Adding a New v2 Service

1. **Edit `v2ServiceFilter`** in `internal/protoc/protoc-gen-zitadelcli/main.go`:
   ```go
   "zitadel.myservice.v2": {
       resourceName: "myservices",
       resourceDesc: "description for --help",
   },
   ```

2. **Add `--path` flag** in `apps/cli/project.json` generate command:
   ```
   --path proto/zitadel/myservice/v2/myservice_service.proto
   ```

3. **Regenerate**: `pnpm nx run @zitadel/cli:generate`

The proto package name (e.g., `zitadel.myservice.v2`) is in the first line of the `.proto` file. The connect stubs must exist in `pkg/grpc/myservice/v2/myserviceconnect/`.

## Modifying Table Column Output

All column logic is in `extractMessageColumns()` in `main.go`. Column ordering follows this buckets rule:
1. **`ID`** — primary resource ID (renamed from `<RESOURCE> ID` to plain `ID`)
2. **`ORGANIZATION ID`** — owner org
3. **Other IDs** — foreign IDs (e.g., `USER ID`, `GRANTED ORGANIZATION ID`)
4. **Semantic fields** — strings, bools, enums, TYPE column (from oneofs)
5. **Timestamps** — `CREATION DATE`, `CHANGE DATE` last

To rename a column header, find the corresponding `columnDef` construction and change the `Header` field. To add a new column from a nested field, update the accessor path in `GoAccessor`.

## Oneof TYPE Column

When a message has a top-level oneof (e.g., `UserType` with `human`/`machine` variants), `extractMessageColumns()` generates a `TYPE` column. The value is the variant name in lowercase with hyphens (e.g., `human`, `machine`, `otp`, `u2f`, `otp-email`). This is handled by `IsOneofType bool` + `OneofVariants []oneofVariantColumn` on `columnDef`.

## Deprecated Methods

ZITADEL marks deprecated RPCs via the OpenAPI extension `openapiv2_operation.deprecated = true`. The generator detects this via `isMethodDeprecatedOpenAPI()` and:
- Prefixes the command's `Short` description with `[DEPRECATED]`
- Emits `fmt.Fprintln(os.Stderr, "Warning: ...")` at the top of the command's `RunE`

Standard proto `option deprecated = true` on methods causes them to be **skipped entirely** (not generated).

## Transport

**ConnectRPC only** — no `grpc.Dial` or raw gRPC. The client is constructed as:
```go
svcClient := userconnect.NewUserServiceClient(httpClient, baseURL)
```
where `httpClient` is a `*http.Client` with `authTransport` (injects `Authorization: Bearer <token>`). No `connect.WithGRPC()` or `connect.WithGRPCWeb()` options are passed — the default ConnectRPC protocol is used.

## Key Generator Conventions

- **Enum flags** include `(one of: VALUE_A, VALUE_B, ...)` in their help text, stripped of the common prefix (e.g., `USER_FIELD_NAME_` → `USER_NAME`, `EMAIL`, ...).
- **Well-known types** (`google.protobuf.Duration`, `Timestamp`, etc.) are not expanded into scalar flags in variant depth-1 expansion. Use `--request-json` for these fields.
- **String oneof variants** (e.g., `organization_id` in a `level` oneof) generate a subcommand with a positional argument rather than a flag.
- **Mutation responses** that only contain `details` show a `CHANGE DATE` column; responses with explicit fields show those fields.

## Schema Reference

The `describe` command outputs the full CLI schema as JSON. Agents should call:
```bash
zitadel-cli describe                    # all groups + global flags
zitadel-cli describe <group>            # all commands in a group
zitadel-cli describe <group> <command>  # flags, json_template, examples
```
