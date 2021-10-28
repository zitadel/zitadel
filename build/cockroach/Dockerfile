FROM cockroachdb/cockroach:v21.1.8

RUN microdnf install curl wget tar gzip shadow-utils

RUN groupadd -g 1000 cockroach
RUN useradd -u 1000 -g cockroach cockroach

RUN wget https://storage.googleapis.com/oauth2l/latest/linux_amd64.tgz
RUN tar zxvf linux_amd64.tgz
RUN mv linux_amd64/oauth2l /usr/local/bin/oauth2l && rm -rf linux_amd64

USER 1000

ENTRYPOINT ["/cockroach/cockroach.sh"]
