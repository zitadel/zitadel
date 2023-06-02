FROM alpine:3 as artifact
ENV ZITADEL_ARGS="start-from-init --masterkeyFromEnv"
ARG TARGETOS TARGETARCH
COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH /app/zitadel
RUN adduser -D zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/app/zitadel", ${ZITADEL_ARGS}]