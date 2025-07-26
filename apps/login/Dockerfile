FROM node:20-alpine AS base

FROM base AS build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack prepare pnpm@10.13.1 --activate && \
    apk update && \
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch --frozen-lockfile
COPY package.json ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm build:login:standalone

FROM scratch AS build-out
COPY --from=build /app/.next/standalone /
COPY --from=build /app/.next/static /.next/static
COPY public public

FROM base AS login-standalone
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY --chown=nextjs:nodejs ./scripts/ ./
COPY --chown=nextjs:nodejs --from=build-out / ./
USER nextjs
ENV HOSTNAME="0.0.0.0" \
    NEXT_PUBLIC_BASE_PATH="/ui/v2/login" \
    PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /runtime/healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["/runtime/entrypoint.sh"]
