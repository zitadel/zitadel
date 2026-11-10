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

Optional services via profiles: `redis` (`cache`), `otel-collector` (`observability`).

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
| `.env.test` | CI-only config | Used by NX targets: `test-run`, `test-e2e`, `test-full`, `stop` |
| `otel-collector-config.yaml` | OTEL Collector pipeline config | Logs traces to stdout; configure `OTEL_BACKEND_ENDPOINT` to forward to a backend |
| `traefik-local-tls.yml` | Traefik dynamic config for local certs | Referenced by local-tls overlay |
| `project.json` | NX project definition | Targets: `test-config`, `test-run`, `test-e2e`, `test`, `test-full`, `stop` |
| `AGENTS.md` | AI agent instructions for this directory | |

## Routing Rules

Traefik routes all traffic for `${ZITADEL_DOMAIN}` via Docker labels:

| Priority | Rule | Target | Middleware |
|----------|------|--------|------------|
| 400 | `Path(/)` | `zitadel-login` | `replacepath=/ui/v2/login/` |
| 250 | `PathPrefix(/ui/v2/login)` | `zitadel-login` | — |
| 200 | `PathPrefix(/api)` | `zitadel-api` | `stripprefix=/api` |
| 100 | Everything else (OIDC, SAML, gRPC, gRPC-web, API v2 REST, ...) | `zitadel-api` (h2c) | — |

No dedicated gRPC router is needed: Traefik's h2c backend scheme forwards gRPC and gRPC-web transparently. API v2 is served as REST/JSON via the gRPC-gateway at `/v2/...` paths.

Both `web` (HTTP) and `websecure` (HTTPS) entrypoints have identical router sets.

### Why this routing model

- `/api` alias exists for DX — tools can use `https://auth.example.com/api/...`
- Canonical paths (e.g., `/.well-known/openid-configuration`, `/oauth/v2/...`) must remain at root for OIDC/SAML protocol compliance
- gRPC, gRPC-web, and REST all share the catch-all router — no separate gRPC router is needed because the `h2c` backend scheme makes Traefik forward all protocols transparently over HTTP/2

## External Settings Invariant

`ZITADEL_EXTERNALDOMAIN`, `ZITADEL_EXTERNALPORT`, and `ZITADEL_EXTERNALSECURE` **must match the public URL** that users see. If they don't, ZITADEL returns "Instance not found" errors. This is the single most common deployment issue.

## Testing

Local NX targets for testing the compose stack:

| Target | What it does | Requires Docker? |
|--------|-------------|------------------|
| `test-config` | Validates all overlay combinations parse with `docker compose config` | No (just the CLI) |
| `test-run` | Builds local images (`@zitadel/api:pack` + `@zitadel/login:pack`), starts the stack with `docker compose up --wait` | Yes |
| `test-e2e` | Runs the full Playwright suite (`wiring.spec.ts` + `smoke.spec.ts`) against `localhost:8888` through Traefik: per-service wiring checks (login, console, OIDC, SAML, API v1 REST, gRPC h2c, gRPC-web, API v2 REST HTTP/1.1 + HTTP/2) and the browser login flow | Yes (stack must be running) |
| `test` | Lightweight — delegates to `test-config` only. Safe for `nx affected` | No |
| `test-full` | Full pipeline: `test-config` → `test-run` → Playwright wiring + browser tests → teardown | Yes |
| `stop` | Tears down the `zitadel-compose-test` stack and removes volumes | Yes |

## Rejected Alternatives

| Alternative | Why rejected |
|-------------|-------------|
| Strict `/api`-only rewrite | Breaks canonical protocol paths (OIDC, SAML) |
| `/grpc` path-prefix routing | gRPC clients/tools don't use path prefixes |
| Dedicated gRPC router (`HeaderRegexp(Content-Type, ^application/grpc.*)`) | Redundant: h2c backend handles gRPC, gRPC-web, and Connect-RPC natively; the catch-all covers all three |
| Single container (API + Login) | Not aligned with v4 architecture; Login is a separate Next.js process |
| Merged TLS configs | Each TLS mode must remain independently composable |
| `network_mode: service:` for Login | Fragile, port conflicts, doesn't work with Traefik routing |
