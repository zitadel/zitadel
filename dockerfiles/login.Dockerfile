FROM node:20-alpine AS base

FROM base AS build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack prepare pnpm@9.1.2 --activate && \
    apk update && apk add --no-cache && \
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY pnpm-lock.yaml pnpm-workspace.yaml  ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm fetch --frozen-lockfile \
    --filter @zitadel/login \
    --filter @zitadel/client \
    --filter @zitadel/proto
COPY package.json ./
COPY apps/login/apps/login/package.json ./apps/login/apps/login/package.json
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/package.json
COPY packages/zitadel-client/package.json ./packages/zitadel-client/package.json
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile \
    --filter @zitadel/login \
    --filter @zitadel/client \
    --filter @zitadel/proto
COPY . .
RUN pnpm turbo build:login:standalone

FROM scratch AS build-out
COPY --from=build /app/apps/login/apps/login/.next/standalone /
COPY --from=build /app/apps/login/apps/login/.next/static /.next/static
COPY --from=build /app/apps/login/apps/login/public /public

FROM base AS login-standalone
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY apps/login/apps/login/scripts ./
COPY --chown=nextjs:nodejs --from=build-out . .
USER nextjs
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
CMD ["/bin/sh", "-c", "node ./healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["./entrypoint.sh"]
