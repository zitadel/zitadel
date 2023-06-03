FROM debian:latest as artifact
ARG TARGETOS TARGETARCH
ENV ZITADEL_ARGS=

COPY build/entrypoint.sh /app/entrypoint.sh
COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH /app/zitadel

RUN adduser --disabled-password zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# USER zitadel
# HEALTHCHECK NONE
ENTRYPOINT ["/app/entrypoint.sh"]