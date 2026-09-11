# syntax=docker/dockerfile:1

FROM node:24-alpine AS builder

WORKDIR /workspace

RUN corepack enable

COPY . .

RUN pnpm install --frozen-lockfile

ARG NEXT_PUBLIC_BASE_PATH=/ui/v2/login
ENV NEXT_PUBLIC_BASE_PATH=${NEXT_PUBLIC_BASE_PATH}

RUN pnpm nx run @zitadel/login:build


FROM node:24-alpine AS runtime

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs \
    && adduser --system --uid 1001 nextjs \
    && mkdir -p /.env-file \
    && touch /.env-file/.env \
    && chown -R nextjs:nodejs /.env-file

COPY --from=builder --chown=nextjs:nodejs /workspace/apps/login/.next/standalone ./

USER nextjs

ENV HOSTNAME="::" \
    PORT="3000" \
    NODE_ENV="production" \
    NODE_OPTIONS="--use-openssl-ca --require /app/load-ssl-cert-dir.cjs" \
    SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt" \
    ZITADEL_TLS_ENABLED="false" \
    OTEL_SERVICE_NAME="zitadel-login" \
    OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /app/healthcheck.mjs \"${NEXT_PUBLIC_BASE_PATH:-/ui/v2/login}/ready\""]

ENTRYPOINT ["/app/entrypoint.sh", "node", "apps/login/server.js"]
