FROM cockroachdb/cockroach:v20.2.3

RUN microdnf install curl wget tar gzip shadow-utils

RUN groupadd -g 1000 cockroach
RUN useradd -u 1000 -g cockroach cockroach

RUN wget https://storage.googleapis.com/oauth2l/latest/linux_amd64.tgz
RUN tar zxvf linux_amd64.tgz
RUN mv linux_amd64/oauth2l /usr/local/bin/oauth2l && rm -rf linux_amd64

ENV FLYWAY_VERSION 7.5.1

RUN mkdir /flyway \
  && wget https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/${FLYWAY_VERSION}/flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && tar -xzf flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && mv flyway-${FLYWAY_VERSION}/* /flyway \
  && rm flyway-commandline-${FLYWAY_VERSION}.tar.gz \
  && chown -R 1000:1000 /flyway

ENV PATH="/flyway:${PATH}"

COPY build/cr-backup/scripts/backup-cockroach.sh /scripts/backup.sh
RUN chown 1000:1000 /scripts/backup.sh
RUN chmod +x /scripts/backup.sh

COPY build/cr-backup/scripts/restore-cockroach.sh /scripts/restore.sh
RUN chown 1000:1000 /scripts/restore.sh
RUN chmod +x /scripts/restore.sh

COPY build/cr-backup/scripts/clean-db-cockroach.sh /scripts/clean-db.sh
RUN chown 1000:1000 /scripts/clean-db.sh
RUN chmod +x /scripts/clean-db.sh
COPY build/cr-backup/scripts/clean-migration-cockroach.sh /scripts/clean-migration.sh
RUN chown 1000:1000 /scripts/clean-migration.sh
RUN chmod +x /scripts/clean-migration.sh
COPY build/cr-backup/scripts/clean-user-cockroach.sh /scripts/clean-user.sh
RUN chown 1000:1000 /scripts/clean-user.sh
RUN chmod +x /scripts/clean-user.sh

USER cockroach

ENTRYPOINT ["/cockroach/cockroach.sh"]
