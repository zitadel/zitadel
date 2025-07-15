FROM login-client AS login-standalone-builder
# Copy package.json files first for better dependency caching
COPY apps/login/package.json ./apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./apps/login

# Copy source code
COPY apps/login ./apps/login

# Build the standalone application
RUN cd apps/login && \
    NEXT_PUBLIC_BASE_PATH=/ui/v2/login \
    NEXT_OUTPUT_MODE=standalone \
    pnpm build

FROM scratch AS login-standalone-out
COPY --from=login-standalone-builder /build/apps/login/.next/standalone /
COPY --from=login-standalone-builder /build/apps/login/.next/static /apps/login/.next/static
COPY --from=login-standalone-builder /build/apps/login/public /apps/login/public

FROM node:20-alpine AS login-standalone
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY ./scripts/entrypoint.sh ./
COPY ./scripts/healthcheck.js ./
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/apps/login/public ./apps/login/public
USER nextjs
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
CMD ["/bin/sh", "-c", "node ./healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["./entrypoint.sh"]
