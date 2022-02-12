#######################
## Final Production Image
#######################
FROM alpine:3 as artifact
RUN adduser -D zitadel
COPY cmd/zitadel/*.yaml /app/
RUN echo zitadel-${BUILDOS}-${BUILDARCH}
COPY .download/zitadel/zitadel-${BUILDOS}-${BUILDARCH} /app/zitadel
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