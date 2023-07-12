FROM --platform=$TARGETPLATFORM debian:latest as artifact
ENV ZITADEL_ARGS=
ARG TARGETPLATFORM

RUN apt-get update && apt-get install ca-certificates -y

COPY build/entrypoint.sh /app/entrypoint.sh
COPY zitadel /app/zitadel

RUN useradd -s "" --home / zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

WORKDIR /app
ENV PATH="/app:${PATH}"

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

FROM --platform=$TARGETPLATFORM scratch as final
ARG TARGETPLATFORM

COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app/zitadel /app/zitadel

HEALTHCHECK NONE

USER zitadel
ENTRYPOINT ["/app/zitadel"]