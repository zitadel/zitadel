## Docker Compose Deployment — AI Agent Instructions

### 1. Scope & Architecture

This directory provides a production-aware single-node Docker Compose deployment for ZITADEL.
The stack contains four core services and two optional profile-based services:

- **Core:** `zitadel-api` (Go API), `zitadel-login` (Next.js Login UI), `postgres`, `proxy` (Traefik)
- **Optional:** `redis` (cache profile), `otel-collector` + `jaeger` (observability profile)

### 2. File Conventions

| File | Purpose | Safe to modify? |
|------|---------|----------------|
| `docker-compose.yml` | Base stack — all modes start from this | Yes, with care |
| `docker-compose.mode-*.yml` | TLS mode overlays (one per mode) | Yes |
| `docker-compose.prodlike.yml` | Init/setup/start split overlay | Yes |
| `docker-compose.test.yml` | CI test overlay (local images, no direct ports) | Yes |
| `.env.example` | User-facing config template | Yes — update on version bumps |
| `.env.test` | CI-only test config | Internal only |
| `otel-collector-config.yaml` | OTEL Collector pipeline config | Yes |
| `traefik-local-tls.yml` | Traefik dynamic config for local TLS certs | Yes |
| `project.json` | NX project definition | Auto-managed by NX |
| `README.md` | Developer/contributor reference | Yes |

### 3. Key Invariants

- **Four TLS modes must remain independently composable.** Never merge mode-specific config into the base `docker-compose.yml`. Each `docker-compose.mode-*.yml` overlay must work when composed with the base file alone.
- **gRPC routing uses `Content-Type: application/grpc*` header matching**, not path prefixes. Do NOT introduce `/grpc` path routing.
- **`ZITADEL_EXTERNALDOMAIN`, `ZITADEL_EXTERNALPORT`, `ZITADEL_EXTERNALSECURE` must be consistent** with the actual public endpoint. Mismatches cause "Instance not found" errors — the single most common deployment issue.
- **Profiles (`cache`, `observability`) are opt-in** and must not affect the default stack behavior.
- **The `/api` path alias coexists with canonical ZITADEL paths.** Do NOT remove canonical root-level path routing (e.g., `/.well-known/`, `/oauth/v2/`). The `/api` prefix is a convenience alias.
- **Image versions are pinned in `.env.example`** via `ZITADEL_VERSION` and infrastructure image variables. When bumping versions, update `.env.example`.
- **CI tests must go through Traefik.** The test overlay must NOT expose direct container ports. All smoke test traffic flows through the proxy to validate routing end-to-end.

### 4. Rejected Alternatives

These designs were considered and explicitly rejected — do not re-propose them:

- **Single container** merging API + Login — not aligned with v4 architecture
- **`/grpc` path-prefix routing** — tool/client incompatibility risk
- **Strict `/api`-only rewrite model** — breaks canonical OIDC/SAML protocol paths
- **`network_mode: service:`** for Login — fragile, port conflicts, incompatible with Traefik routing
- **Merged TLS configurations** — each mode must be independently composable without side effects

### 5. Common Commands

| Task | Command |
|------|---------|
| Start (local dev) | `cp .env.example .env && docker compose up -d --wait` |
| Start (Let's Encrypt) | `docker compose --env-file .env -f docker-compose.yml -f docker-compose.mode-letsencrypt.yml up -d --wait` |
| Start (production-like) | `docker compose --env-file .env -f docker-compose.yml -f docker-compose.prodlike.yml up -d --wait` |
| Validate all configs | `docker compose --env-file .env.example -f docker-compose.yml -f <overlay> config > /dev/null` for each overlay |
| Run CI smoke test | `pnpm nx run @zitadel/compose:test` |
| Smoke check | `curl -sS http://localhost:8080/.well-known/openid-configuration` |

### 6. Terminology

Follow the root `agents.md` glossary. Key rules:

- **Instance** = a logical ZITADEL tenant/partition, NEVER "example"
- **System** = the entire ZITADEL installation/deployment
- See the [Technical Glossary](../../agents.md) for all user-facing text in any language
