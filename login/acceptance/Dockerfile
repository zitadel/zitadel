FROM golang:1.24-alpine
RUN apk add curl jq
COPY setup.sh /setup.sh
RUN chmod +x /setup.sh
ENTRYPOINT [ "/setup.sh" ]
