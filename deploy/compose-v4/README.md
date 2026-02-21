# ZITADEL v4 Docker Compose Deployment Pack

A low-friction but production-aware single-node deployment for ZITADEL v4.

## 1. Architecture Summary

This pack follows the current ZITADEL v4 deployment model:

- `zitadel-api` (Go): HTTP APIs, gRPC, gRPC-gateway, OIDC/SAML, Console
- `zitadel-login` (Next.js): Login UI served at `/ui/v2/login`
- `postgres`: persistent state
- `proxy` (Traefik): path-based routing, HTTP/2, gRPC proxying

Optional services via profiles:

- `redis` (`cache` profile)
- `otel-collector` + `jaeger` (`observability` profile)

## 2. Why Traefik

Traefik is used because it provides:

- Native Docker label-based routing
- HTTP/2 and gRPC proxying with h2c upstream support
- Easy ACME/Let's Encrypt support
- Readable declarative config in one Compose stack

## 3. Quick Start (Mode 1: Local Dev, No TLS)

```bash
cd deploy/compose-v4
cp .env.example .env

docker compose --env-file .env -f docker-compose.yml up -d --wait
```

Expected endpoints:

- Login/UI entry: `http://localhost:8080/`
- API alias: `http://localhost:8080/api/...`
- Canonical paths are still available (for OIDC/issuer correctness)

## 4. Mode Matrix

| Mode | Compose files | Proxy ports | TLS termination |
|---|---|---|---|
| Local Dev (default) | `docker-compose.yml` | `8080 -> 80` | none |
| Easy TLS (Let's Encrypt) | `docker-compose.yml` + `docker-compose.mode-letsencrypt.yml` | `80,443` | Traefik |
| External TLS | `docker-compose.yml` + `docker-compose.mode-external-tls.yml` | `80` | external LB/WAF/CDN |
| Local TLS (self-signed) | `docker-compose.yml` + `docker-compose.mode-local-tls.yml` | `443` | Traefik (local cert files) |

### Mode 2 command (Let's Encrypt)

```bash
docker compose --env-file .env \
  -f docker-compose.yml \
  -f docker-compose.mode-letsencrypt.yml \
  up -d --wait
```

### Mode 3 command (External TLS)

```bash
docker compose --env-file .env \
  -f docker-compose.yml \
  -f docker-compose.mode-external-tls.yml \
  up -d --wait
```

### Optional Local TLS command (self-signed cert)

Generate cert files first:

```bash
mkdir -p certs
openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout certs/local.key \
  -out certs/local.crt \
  -days 365 \
  -subj "/CN=localhost/O=ZITADEL Local"
```

Then run:

```bash
docker compose --env-file .env \
  -f docker-compose.yml \
  -f docker-compose.mode-local-tls.yml \
  up -d --wait
```

## 5. Routing and gRPC Handling

The proxy exposes one host (`ZITADEL_DOMAIN`) and routes:

- `/` -> rewritten to `/ui/v2/login/` -> `zitadel-login`
- `/ui/v2/login` -> `zitadel-login`
- `/api/*` -> `zitadel-api` with `/api` stripped
- `Content-Type: application/grpc...` -> `zitadel-api` over h2c
- all other non-login, non-root paths -> canonical `zitadel-api`

This keeps `/api` DX while preserving canonical ZITADEL paths required by issuer/OIDC behavior.

## 6. Forwarded Headers and Issuer Correctness

`zitadel-api` is configured with explicit external settings:

- `ZITADEL_EXTERNALDOMAIN`
- `ZITADEL_EXTERNALPORT`
- `ZITADEL_EXTERNALSECURE`

Traefik forwards host/proto/forwarded metadata. In external TLS mode, forwarded headers are trusted at Traefik entrypoint so upstream TLS termination can be represented correctly.

If these values are inconsistent with your public DNS/proxy chain, ZITADEL may return `Instance not found`.

## 7. Profiles

### Cache profile (Redis)

Run:

```bash
docker compose --env-file .env -f docker-compose.yml --profile cache up -d --wait
```

To actually use Redis for object caches, set these in `.env`:

- `ZITADEL_CACHES_CONNECTORS_REDIS_ENABLED=true`
- `ZITADEL_CACHES_INSTANCE_CONNECTOR=redis`
- `ZITADEL_CACHES_MILESTONES_CONNECTOR=redis`
- `ZITADEL_CACHES_ORGANIZATION_CONNECTOR=redis`

### Observability profile

Run:

```bash
docker compose --env-file .env -f docker-compose.yml --profile observability up -d --wait
```

Set API tracing exporter in `.env`:

- `ZITADEL_INSTRUMENTATION_TRACE_EXPORTER_TYPE=grpc`
- `ZITADEL_INSTRUMENTATION_TRACE_EXPORTER_ENDPOINT=otel-collector:4317`
- `ZITADEL_INSTRUMENTATION_TRACE_EXPORTER_INSECURE=true`

Jaeger UI: `http://localhost:16686`

Note: Login OTEL env vars are included as forward-compatible placeholders; current login images may ignore them.

## 8. Production-Like Single-Node vs Quickstart

Quickstart base uses `start-from-init` (minimal operator friction).

Production-like flow splits init/setup/start:

- `zitadel-init` (one-shot)
- `zitadel-setup` (one-shot)
- `zitadel-api` with `start`

Run production-like:

```bash
docker compose --env-file .env \
  -f docker-compose.yml \
  -f docker-compose.prodlike.yml \
  up -d --wait
```

## 9. Updating ZITADEL

To update ZITADEL to a new version:

1. Edit `.env` and bump the image versions:
   ```
   ZITADEL_API_IMAGE=ghcr.io/zitadel/zitadel:v4.11.0
   ZITADEL_LOGIN_IMAGE=ghcr.io/zitadel/zitadel-login:v4.11.0
   ```
2. Pull the new images:
   ```bash
   docker compose --env-file .env -f docker-compose.yml pull
   ```
3. Recreate the stack:
   ```bash
   docker compose --env-file .env -f docker-compose.yml up -d --wait
   ```

**How `start-from-init` handles upgrades**: The default quickstart command is idempotent.
The init phase skips resources (database, user, grants) that already exist.
The setup phase runs only the migration steps for the new version and leaves the rest untouched.
Log messages such as `role "zitadel" already exists` are expected and harmless.

**Production-like upgrades**: When using `docker-compose.prodlike.yml`, the `zitadel-setup` one-shot container
runs the new version's migrations before `zitadel-api` starts. This gives you a controlled, serialized upgrade:

```bash
docker compose --env-file .env \
  -f docker-compose.yml \
  -f docker-compose.prodlike.yml \
  up -d --wait
```

**FIRSTINSTANCE / DEFAULTINSTANCE settings**: Variables prefixed `ZITADEL_FIRSTINSTANCE_` and
`ZITADEL_DEFAULTINSTANCE_` are applied **once** during the very first setup of a new instance.
Changing them in `.env` after the first start has no effect on the running instance.
Use the Admin Console or Admin API to change instance settings on an existing installation.

## 10. Scaling and Externalization

### Scale API replicas

For production-like API-only scaling:

```bash
docker compose --env-file .env -f docker-compose.yml up -d --scale zitadel-api=2
```

Use a centralized cache backend (Redis/Postgres connector strategy) for multi-replica consistency.

### Externalize Postgres

- Point `ZITADEL_DATABASE_POSTGRES_HOST` to external DB host
- Keep admin/user credentials in secrets management
- Disable/remove local `postgres` service

### Externalize Redis

- Keep `redis` profile disabled
- Set `ZITADEL_CACHES_CONNECTORS_REDIS_ADDR` to external Redis
- Keep connector toggles explicit in `.env`

## 11. Kubernetes Migration Mapping

Compose concepts map directly to Kubernetes:

- `zitadel-init` / `zitadel-setup` -> Jobs (or Helm hooks)
- `zitadel-api` / `zitadel-login` -> separate Deployments
- Traefik container -> ingress controller + ingress routes
- `.env` settings -> ConfigMaps/Secrets
- external Postgres/Redis -> managed services

## 12. Tradeoffs and Rejected Alternatives

Chosen:

- `/api` alias + canonical paths (hybrid) for DX + protocol safety
- same-host gRPC routing via `Content-Type` matcher

Rejected:

- strict `/api`-only rewrite model (too brittle for canonical protocol paths)
- `/grpc` path-prefix routing (tool/client incompatibility risk)
- collapsing API and Login into one container (not aligned with v4 architecture)

## 13. Release Flow Placeholder

Version bump automation is intentionally deferred until the Nx release PR merges.

Current state:

- image versions are manually pinned in `.env.example`
- `make release-bump` is a placeholder target for future Nx integration

## 14. Validation Checklist

Render final Compose config for each variant:

```bash
docker compose --env-file .env -f docker-compose.yml config >/dev/null

docker compose --env-file .env -f docker-compose.yml -f docker-compose.mode-letsencrypt.yml config >/dev/null

docker compose --env-file .env -f docker-compose.yml -f docker-compose.mode-external-tls.yml config >/dev/null

docker compose --env-file .env -f docker-compose.yml -f docker-compose.mode-local-tls.yml config >/dev/null

docker compose --env-file .env -f docker-compose.yml -f docker-compose.prodlike.yml config >/dev/null
```

Smoke checks (default mode):

```bash
curl -i http://localhost:8080/
curl -sS http://localhost:8080/api/admin/v1/healthz
grpcurl -plaintext localhost:8080 zitadel.admin.v1.AdminService/Healthz
```
