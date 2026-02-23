FROM node:22-alpine
WORKDIR /app
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file

COPY --chown=nextjs:nodejs .next/standalone ./

USER nextjs
ENV HOSTNAME="::" \
    PORT="3000" \
    NODE_ENV="production" \
    NODE_OPTIONS="--use-openssl-ca" \
    SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt" \
    SSL_CERT_DIR="/etc/ssl/certs" \
    ZITADEL_TLS_ENABLED="false"

# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/node", "/app/healthcheck.mjs"]
ENTRYPOINT ["/app/entrypoint.sh", "node", "apps/login/server.js" ]
