FROM login-client AS login-standalone-builder
COPY apps/login ./apps/login
COPY packages/zitadel-tailwind-config packages/zitadel-tailwind-config
RUN pnpm exec turbo prune @zitadel/login --docker
WORKDIR /build/docker
RUN cp -r ../out/json/* .
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
RUN cp -r ../out/full/* .
RUN pnpm exec turbo run build:login:standalone

FROM scratch AS login-standalone-out
COPY --from=login-standalone-builder /build/docker/apps/login/.next/standalone /
COPY --from=login-standalone-builder /build/docker/apps/login/.next/static /apps/login/.next/static
COPY --from=login-standalone-builder /build/docker/apps/login/public /apps/login/public

FROM node:20-alpine AS login-standalone
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY ./scripts/entrypoint.sh ./
COPY ./scripts/healthcheck.js ./
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/docker/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/docker/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs --from=login-standalone-builder /build/docker/apps/login/public ./apps/login/public
USER nextjs
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
CMD ["/bin/sh", "-c", "node ./healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["./entrypoint.sh"]
