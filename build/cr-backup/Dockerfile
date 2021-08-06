FROM cockroachdb/cockroach:v20.2.3

RUN microdnf install curl wget tar gzip shadow-utils

RUN wget https://storage.googleapis.com/oauth2l/latest/linux_amd64.tgz
RUN tar zxvf linux_amd64.tgz
RUN mv linux_amd64/oauth2l /usr/local/bin/oauth2l && rm -rf linux_amd64

COPY build/cr-backup/scripts/backup-cockroach.sh /scripts/backup.sh
RUN chmod +x /scripts/backup.sh

COPY build/cr-backup/scripts/restore-cockroach.sh /scripts/restore.sh
RUN chmod +x /scripts/restore.sh

COPY build/cr-backup/scripts/clean-db-cockroach.sh /scripts/clean-db.sh
RUN chmod +x /scripts/clean-db.sh
COPY build/cr-backup/scripts/clean-migration-cockroach.sh /scripts/clean-migration.sh
RUN chmod +x /scripts/clean-migration.sh
COPY build/cr-backup/scripts/clean-user-cockroach.sh /scripts/clean-user.sh
RUN chmod +x /scripts/clean-user.sh

RUN groupadd -g 1000 cockroach
RUN useradd -u 1000 -g cockroach cockroach

USER cockroach

ENTRYPOINT ["/cockroach/cockroach.sh"]
