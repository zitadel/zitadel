FROM cockroachdb/cockroach:v20.2.3

RUN microdnf install curl wget tar gzip shadow-utils

RUN groupadd -g 1000 cockroach
RUN useradd -u 1000 -g cockroach cockroach

RUN wget https://storage.googleapis.com/oauth2l/latest/linux_amd64.tgz
RUN tar zxvf linux_amd64.tgz
RUN mv linux_amd64/oauth2l /usr/local/bin/oauth2l && rm -rf linux_amd64

COPY build/cockroach/scripts/backup-cockroach.sh /scripts/backup.sh
RUN chown 1000:1000 /scripts/backup.sh
RUN chmod u+xr /scripts/backup.sh

COPY build/cockroach/scripts/restore-cockroach.sh /scripts/restore.sh
RUN chown 1000:1000 /scripts/restore.sh
RUN chmod u+xr /scripts/restore.sh

COPY build/cockroach/scripts/clean-db-cockroach.sh /scripts/clean-db.sh
RUN chown 1000:1000 /scripts/clean-db.sh
RUN chmod u+xr /scripts/clean-db.sh
COPY build/cockroach/scripts/clean-migration-cockroach.sh /scripts/clean-migration.sh
RUN chown 1000:1000 /scripts/clean-migration.sh
RUN chmod u+xr /scripts/clean-migration.sh
COPY build/cockroach/scripts/clean-user-cockroach.sh /scripts/clean-user.sh
RUN chown 1000:1000 /scripts/clean-user.sh
RUN chmod u+xr /scripts/clean-user.sh

USER 1000

ENTRYPOINT ["/cockroach/cockroach.sh"]
