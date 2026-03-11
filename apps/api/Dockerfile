# Builder stage: Sets up the environment, installs dependencies, copies the Zitadel binary, and configures permissions for the application.
# This stage produces a runnable image that can be used for debugging.
FROM debian:latest AS builder
ARG TARGETPLATFORM
# BINARY allows callers to pass a differently-named local build artifact.
# Defaults to 'zitadel' for release builds; set to 'zitadel.local' for local dev builds.
ARG BINARY=zitadel

RUN apt-get update && apt-get install ca-certificates -y

COPY apps/api/entrypoint.sh /app/entrypoint.sh
COPY ./.artifacts/bin/${TARGETPLATFORM}/${BINARY} /app/zitadel

RUN useradd -s "" --home /app zitadel && \
    chown -R zitadel /app && \
    chmod +x /app/zitadel && \
    chmod +x /app/entrypoint.sh

WORKDIR /app
ENV PATH="/app:${PATH}"
ENV HOME="/app"

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

# Final stage: minimal debian image to support CGO-linked binary (DuckDB requires libc)
FROM debian:stable-slim AS final

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/zitadel /app/zitadel

# Give zitadel user write access to /app so DuckDB can create its .duckdb cache dir.
# Also create the default signals data directory so volume mounts inherit the ownership.
RUN chown -R zitadel /app && \
    mkdir -p /var/lib/zitadel/signals && \
    chown -R zitadel /var/lib/zitadel

HEALTHCHECK NONE
EXPOSE 8080

# DuckDB needs a writable home directory for its .duckdb extension cache.
ENV HOME="/app"
USER zitadel
ENTRYPOINT ["/app/zitadel"]
