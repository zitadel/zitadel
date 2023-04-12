FROM alpine:3 as artifact
ARG TARGETOS TARGETARCH
COPY zitadel-core-$TARGETOS-$TARGETARCH/zitadel-core-$TARGETOS-$TARGETARCH /app/zitadel
RUN adduser -D zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel

FROM alpine:3 as final
COPY --from=artifact /app/zitadel zitadel
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadel"]