# Builder stage: Sets up the environment, installs dependencies, copies the Zitadel binary, and configures permissions for the application.
# This stage produces a runnable image that can be used for debugging.
FROM debian:latest AS builder
ARG TARGETPLATFORM

RUN apt-get update && apt-get install ca-certificates -y

COPY apps/api/entrypoint.sh /app/entrypoint.sh
COPY ./.artifacts/bin/${TARGETPLATFORM}/zitadel /app/zitadel

RUN useradd -s "" --home / zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

WORKDIR /app
ENV PATH="/app:${PATH}"

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

# Final stage: Creates a minimal container image with just the Zitadel binary and necessary files
FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/zitadel /app/zitadel

HEALTHCHECK NONE
EXPOSE 8080

USER zitadel
ENTRYPOINT ["/app/zitadel"]
