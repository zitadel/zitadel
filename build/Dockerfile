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
    chmod +x /app/entrypoint.sh && \
    mkdir /emptytmp

WORKDIR /app
ENV PATH="/app:${PATH}"

USER zitadel

HEALTHCHECK --interval=5s --timeout=120s --retries=6 \
    CMD ["/app/zitadel", "ready"]

ENTRYPOINT ["/app/entrypoint.sh"]

FROM --platform=$TARGETPLATFORM scratch as final
ARG TARGETPLATFORM

COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact --chown=zitadel:zitadel /app/zitadel /app/zitadel
COPY --from=artifact --chown=zitadel:zitadel /emptytmp /tmp

USER zitadel

# /app/zitadel ready is a healthcheck endpoint that immediately
HEALTHCHECK --interval=5s --timeout=600s --retries=3 \
    CMD ["/app/zitadel", "ready"]

ENTRYPOINT ["/app/zitadel"]