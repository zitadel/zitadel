FROM debian:latest as artifact
ARG TARGETOS TARGETARCH
COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH /app/zitadel
# RUN adduser -D zitadel && \
#     chown zitadel /app/zitadel && \
#     chmod +x /app/zitadel

RUN chmod +x /app/zitadel

FROM scratch as final
COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app /app
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/app/zitadel"]