# ZITADEL Docker Compose — Developer Reference

> **User-facing documentation:** [zitadel.com/docs/self-hosting/deploy/compose](https://zitadel.com/docs/self-hosting/deploy/compose)
>
> This README is the contributor/developer reference — architecture decisions, file conventions, and routing logic.

## Architecture

```
                 ┌─────────────────────────┐
  Browser ──────►│  Traefik (proxy)        │
                 │  Port 80 / 443          │
                 └───┬──────────┬──────────┘
                     │          │
          ┌──────────▼──┐  ┌───▼──────────┐
          │ zitadel-api  │  │ zitadel-login │
          │ Go :8080     │  │ Next.js :3000 │
          └──────┬───────┘  └──────────────┘
                 │
          ┌──────▼───────┐
          │  PostgreSQL   │
          └──────────────┘
```

Optional services via profiles: `redis` (`cache`), `otel-collector` + `jaeger` (`observability`).

## File Conventions

| File | Role | Notes |
|------|------|-------|
| `docker-compose.yml` | **Base stack** — all modes start from this | Must work standalone with `.env.example` |
| `docker-compose.mode-letsencrypt.yml` | TLS overlay: ACME HTTP challenge | Declares its own `letsencrypt` volume |
| `docker-compose.mode-external-tls.yml` | TLS overlay: upstream LB terminates TLS | Enables forwarded headers |
| `docker-compose.mode-local-tls.yml` | TLS overlay: self-signed certs | Mounts `./certs/` and `traefik-local-tls.yml` |
| `docker-compose.prodlike.yml` | Init/setup/start split | Uses YAML anchors for shared DB env |
| `docker-compose.test.yml` | CI smoke test overlay | Overrides images to `:local` tags |
| `.env.example` | User-facing config template | Copy to `.env` before first run |
| `.env.test` | CI-only config | Used by NX `@zitadel/compose:test` |
| `otel-collector-config.yaml` | OTEL Collector pipeline config | Traces only (OTLP → Jaeger) |
| `traefik-local-tls.yml` | Traefik dynamic config for local certs | Referenced by local-tls overlay |
| `project.json` | NX project definition | Targets: `test-config`, `test-run`, `test`, `stop` |
| `agents.md` | AI agent instructions for this directory | |

## Routing Rules

Traefik routes all traffic for `${ZITADEL_DOMAIN}` via Docker labels:

| Priority | Rule | Target | Middleware |
|----------|------|--------|------------|
| 400 | `Path(/)` | `zitadel-login` | `replacepath=/ui/v2/login/` |
| 300 | `HeadersRegexp(Content-Type, ^application/grpc.*)` | `zitadel-api` (h2c) | — |
| 250 | `PathPrefix(/ui/v2/login)` | `zitadel-login` | — |
| 200 | `PathPrefix(/api)` | `zitadel-api` | `stripprefix=/api` |
| 100 | Everything else (canonical ZITADEL paths) | `zitadel-api` | — |

Both `web` (HTTP) and `websecure` (HTTPS) entrypoints have identical router sets.

### Why this routing model

- `/api` alias exists for DX — tools can use `https://auth.example.com/api/...`
- Canonical paths (e.g., `/.well-known/openid-configuration`, `/oauth/v2/...`) must remain at root for OIDC/SAML protocol compliance
- gRPC uses `Content-Type` header matching because gRPC clients send to the root path (not a `/grpc` prefix)

## External Settings Invariant

`ZITADEL_EXTERNALDOMAIN`, `ZITADEL_EXTERNALPORT`, and `ZITADEL_EXTERNALSECURE` **must match the public URL** that users see. If they don't, ZITADEL returns "Instance not found" errors. This is the single most common deployment issue.

## CI Testing

The NX target `@zitadel/compose:test` runs:

1. **`test-config`** — validates all overlay combinations parse with `docker compose config`
2. **`test-run`** — builds local images (`@zitadel/api:pack` + `@zitadel/login:pack`), starts the stack with `docker compose up --wait`
3. **`test`** — curls 4 endpoints **through Traefik** (OIDC discovery, login health, API health, root redirect), then tears down

## Rejected Alternatives

| Alternative | Why rejected |
|-------------|-------------|
| Strict `/api`-only rewrite | Breaks canonical protocol paths (OIDC, SAML) |
| `/grpc` path-prefix routing | gRPC clients/tools don't use path prefixes |
| Single container (API + Login) | Not aligned with v4 architecture; Login is a separate Next.js process |
| Merged TLS configs | Each TLS mode must remain independently composable |
| `network_mode: service:` for Login | Fragile, port conflicts, doesn't work with Traefik routing |
