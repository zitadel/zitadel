FROM debian:latest as artifact
ARG TARGETOS TARGETARCH
ENV ZITADEL_ARGS=

WORKDIR /app

COPY build/entrypoint.sh .

COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH ./zitadel

RUN adduser -D zitadel && \
    chown zitadel zitadel && \
    chmod +x zitadel && \
    chown zitadel entrypoint.sh && \
    chmod +x entrypoint.sh

USER zitadel
# HEALTHCHECK NONE
ENTRYPOINT ["/app/entrypoint.sh"]