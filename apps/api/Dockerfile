FROM debian:latest AS artifact
ARG TARGETPLATFORM

ENV ZITADEL_ARGS=

RUN apt-get update && apt-get install ca-certificates -y

COPY apps/api/entrypoint.sh /app/entrypoint.sh
COPY ./.artifacts/${TARGETPLATFORM}/zitadel /app/zitadel

RUN useradd -s "" --home / zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

WORKDIR /app
ENV PATH="/app:${PATH}"

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

FROM scratch AS final

COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app/zitadel /app/zitadel

HEALTHCHECK NONE
EXPOSE 8080

USER zitadel
ENTRYPOINT ["/app/zitadel"]
