FROM debian:latest as artifact
ARG TARGETOS TARGETARCH
ENV ZITADEL_ARGS=

COPY build/entrypoint.sh /app/entrypoint.sh
COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH /app/zitadel

# RUN adduser -D zitadel && \
#     chown zitadel zitadel && \
#     chmod +x zitadel && \
#     chown zitadel entrypoint.sh && \
#     chmod +x entrypoint.sh

RUN chmod +x /app/zitadel && \
    chmod +x /app/entrypoint.sh

USER zitadel
# HEALTHCHECK NONE
ENTRYPOINT ["/app/entrypoint.sh"]