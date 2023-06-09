FROM ubuntu:latest as artifact
ARG TARGETOS TARGETARCH
ENV ZITADEL_ARGS=

WORKDIR /app
ENV PATH="/app:${PATH}"

RUN apt-get update && apt-get install ca-certificates -y

COPY build/entrypoint.sh /app/entrypoint.sh
ADD zitadel-$TARGETOS-$TARGETARCH/zitadel-$TARGETOS-$TARGETARCH.tar /app/
RUN mv /app/zitadel-$TARGETOS-$TARGETARCH/zitadel /app/zitadel && \
    rm -rf /app/zitadel-$TARGETOS-$TARGETARCH

RUN useradd -s "" --home / zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

USER zitadel
ENTRYPOINT ["entrypoint.sh"]

FROM scratch as final

WORKDIR /app
ENV PATH="/app:${PATH}"
HEALTHCHECK NONE

COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app/zitadel /app/zitadel

USER zitadel
ENTRYPOINT ["zitadel"]