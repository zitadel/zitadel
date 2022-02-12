#######################
## Final Production Image
#######################
ARG TARGETOS
ARG TARGETARCH

FROM alpine:3 as artifact
RUN adduser -D zitadel
COPY cmd/zitadel/*.yaml /app/
RUN echo zitadel-${TARGETOS}-${TARGETARCH}
COPY .download/zitadel/zitadel-${TARGETOS}-${TARGETARCH} /app/zitadel
RUN chmod a+x /app/zitadel

#######################
## Scratch Image
#######################
FROM  scratch as final
COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app /
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadel"]