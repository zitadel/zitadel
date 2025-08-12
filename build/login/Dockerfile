FROM node:20-alpine AS runtime

FROM runtime AS pnpm-base
RUN apk add --no-cache libc6-compat
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack prepare pnpm@10.13.1 --activate
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm add -g turbo@2.5.5

FROM pnpm-base AS pruner
WORKDIR /prune
COPY . .
RUN pnpm turbo prune @zitadel/login @zitadel/client @zitadel/proto --docker

FROM pnpm-base AS installer
WORKDIR /install
COPY --from=pruner /prune/out/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm fetch --frozen-lockfile
COPY --from=pruner /prune/out/json/ .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --ignore-scripts

FROM pnpm-base AS builder
WORKDIR /build
COPY --from=installer /install/ .
COPY --from=pruner /prune/out/full/ .
COPY proto ./proto
ENV CI=true
RUN --mount=type=cache,id=turbo,target=/build/.turbo/cache \
    --mount=type=cache,id=next,target=/build/apps/login/.next/cache \
    pnpm turbo build:login:standalone --cache-dir=/build/.turbo/cache

FROM scratch AS build-out
COPY /apps/login/public ./apps/login/public
COPY --from=builder /build/apps/login/.next/standalone ./
COPY --from=builder /build/apps/login/.next/static ./apps/login/.next/static

FROM runtime AS login-standalone
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY --chown=nextjs:nodejs apps/login/scripts ./
COPY --chown=nextjs:nodejs --from=build-out . .
# Debug the final structure
USER nextjs
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /runtime/healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["/runtime/entrypoint.sh"]
