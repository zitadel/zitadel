FROM node:20-alpine

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file

COPY --chown=nextjs:nodejs ./docker/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs ./docker/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs ./docker/apps/login/public ./apps/login/public

USER nextjs
ENV HOSTNAME="0.0.0.0"

CMD ["/bin/sh", "-c", " set -o allexport && . /.env-file/.env && set +o allexport && node apps/login/server.js"]
