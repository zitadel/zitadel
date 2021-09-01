#######################
## By default we build the prod enviroment
ARG ENV=prod

#######################
## Go base build
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:1.16 as go-base

WORKDIR src/github.com/caos/zitadel/
COPY go.mod go.sum ./
RUN go mod download

## Go prod build
FROM go-base as prod-go-build
COPY . .
RUN CGO_ENABLED=0 go build \
      -a \
      -installsuffix cgo \
      -o backupctl \
      ./cmd/backupctl/main.go

## Go dev build
FROM go-base as dev-go-build
RUN go get github.com/go-delve/delve/cmd/dlv

#######################
## Final Production Image
#######################
FROM alpine:latest as artifact

COPY --from=prod-go-build /go/src/github.com/caos/zitadel/backupctl /app/backupctl
RUN chmod a+x /app/backupctl

## CokroachDB Image for cockroach binary
FROM cockroachdb/cockroach:v20.2.3 as final

RUN microdnf install rsync

COPY --from=artifact /app /
HEALTHCHECK NONE
ENTRYPOINT ["/backupctl"]
