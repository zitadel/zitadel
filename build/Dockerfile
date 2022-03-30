#######################
## Final Production Image
#######################
FROM alpine:3 as artifact
COPY zitadel /app/zitadel
RUN adduser -D zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel

#######################
## Scratch Image
#######################
FROM scratch as final
COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app /
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadel"]

