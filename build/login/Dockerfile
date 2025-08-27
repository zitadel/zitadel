FROM node:20-alpine
WORKDIR /runtime
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file
COPY --chown=nextjs:nodejs node_modules ./node_modules
COPY --chown=nextjs:nodejs apps/login/scripts ./
COPY --chown=nextjs:nodejs apps/login/public ./apps/login/public
COPY --chown=nextjs:nodejs apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs apps/login/.next/standalone ./apps/login

# Debug the final structure
USER nextjs
ENV HOSTNAME="0.0.0.0"
ENV PORT=3000
# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /runtime/healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["/runtime/entrypoint.sh"]
