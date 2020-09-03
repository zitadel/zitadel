# This Stage prepares the user in the container and copies the files
FROM alpine:latest as prepare
RUN adduser -D zitadel
COPY .artifacts/zitadel-linux-amd64 /app/zitadel
COPY cmd/zitadel/*.yaml /app
RUN chmod a+x /app/zitadel

# This Stage is intended as production image
FROM scratch as final
COPY --from=prepare /etc/passwd /etc/passwd
COPY --from=prepare /app /
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadel"]
