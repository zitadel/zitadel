FROM login-base AS prune-for-docker
RUN pnpm install turbo --global
COPY . .
RUN turbo prune @zitadel/login --docker
FROM login-base AS installer
COPY --from=prune-for-docker /app/out/json/ .
RUN pnpm install --frozen-lockfile
COPY --from=prune-for-docker /app/out/full/ .
RUN NEXT_PUBLIC_BASE_PATH=/ui/v2/login NEXT_OUTPUT_MODE=standalone pnpm exec turbo run build

FROM login-platform AS login-standalone
WORKDIR /app
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY --chown=nextjs:nodejs --from=installer /app/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs --from=installer /app/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs --from=installer /app/apps/login/public ./apps/login/public
USER nextjs
ENV HOSTNAME="0.0.0.0"
CMD ["/bin/sh", "-c", " set -o allexport && . /.env-file/.env && set +o allexport && node apps/login/server.js"]
