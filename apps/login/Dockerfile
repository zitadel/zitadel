FROM node:22-alpine
WORKDIR /app
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
# If /.env-file/.env is mounted into the container, its variables are made available to the server before it starts up.
RUN mkdir -p /.env-file && touch /.env-file/.env && chown -R nextjs:nodejs /.env-file

COPY --chown=nextjs:nodejs .next/standalone ./
COPY --chown=nextjs:nodejs scripts/entrypoint.sh /app/entrypoint.sh
COPY --chown=nextjs:nodejs scripts/healthcheck.js /app/healthcheck.js
RUN chmod +x /app/entrypoint.sh

USER nextjs
ENV HOSTNAME="::" \
    PORT="3000" \
    NODE_ENV="production" \
    NODE_OPTIONS="--use-openssl-ca" \
    SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt" \
    SSL_CERT_DIR="/etc/ssl/certs" \
    OTEL_SERVICE_NAME="zitadel-login" \
    OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"

# TODO: Check healthy, not ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "node /app/healthcheck.js http://localhost:${PORT}/ui/v2/login/healthy"]
ENTRYPOINT ["/app/entrypoint.sh", "node", "apps/login/server.js"]
