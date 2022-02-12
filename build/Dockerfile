#######################
## Final Production Image
#######################
FROM --platform=$BUILDPLATFORM alpine:3 as artifact
ARG TARGETOS TARGETARCH
RUN adduser -D zitadel
COPY cmd/defaults.yaml /app/
COPY .download/zitadel/$TARGETOS_$TARGETARCH/zitadel /app/zitadel
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