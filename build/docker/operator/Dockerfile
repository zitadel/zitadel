# This Stage prepares the user in the container and copies the files
FROM alpine:latest as prepare
RUN adduser -D zitadel
COPY .artifacts/zitadelctl-linux-amd64 /app/zitadelctl
RUN chmod a+x /app/zitadelctl
COPY migrations/cockroach/V1* /migrations/

# This Stage is intended as production image
FROM scratch as final
COPY --from=prepare /etc/passwd /etc/passwd
COPY --from=prepare /app /
COPY --from=prepare /migrations /migrations
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadelctl"]
