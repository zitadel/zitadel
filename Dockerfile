FROM node:20-slim

WORKDIR /app

# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --chown=nextjs:nodejs ./docker/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs ./docker/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs ./docker/apps/login/public ./apps/login/public

ENV HOSTNAME="0.0.0.0"
CMD node apps/login/server.js