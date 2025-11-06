FROM node:22-alpine

LABEL org.opencontainers.image.documentation="https://zitadel.com/docs" \
      org.opencontainers.image.vendor="ZITADEL" \
      org.opencontainers.image.description="ZITADEL Login - Identity infrastructure, simplified for you." \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="Zitadel Login" \
      org.opencontainers.image.url="https://github.com/zitadel/login" \
      org.opencontainers.image.vendor="ZITADEL"

WORKDIR /app
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file

COPY --chown=nextjs:nodejs .next/standalone ./

USER nextjs
ENV HOSTNAME="::" \
    PORT="3000" \
    NODE_ENV="production"

# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /app/healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["/app/entrypoint.sh", "node", "apps/login/server.js" ]
