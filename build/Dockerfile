FROM ubuntu:latest as artifact
ARG TARGETOS TARGETARCH
ENV ZITADEL_ARGS=

RUN adduser --disabled-password zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh
RUN apt-get install -y ca-certificates

COPY build/entrypoint.sh /app/entrypoint.sh
COPY zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH /app/zitadel

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

FROM scratch as final
COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app /app
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/app/zitadel"]