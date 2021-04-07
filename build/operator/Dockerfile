#######################
## By default we build the prod enviroment
ARG ENV=prod

#######################
## Go base build
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:1.16.3 as go-base

WORKDIR src/github.com/caos/zitadel/
COPY go.mod go.sum ./
RUN go mod download


## Go test
FROM go-base as go-test
COPY . .
#Migrations for cockroach-secure
RUN go install github.com/rakyll/statik
RUN ./build/operator/prebuild.sh ./migrations

RUN go test -race -v -coverprofile=profile.cov ./operator/...

## Go test
FROM scratch as go-codecov
COPY --from=go-test /go/src/github.com/caos/zitadel/profile.cov profile.cov

## Go prod build
FROM go-test as prod-go-build


ARG ARCH=amd64
ARG OS=linux
ARG VERSION=none
RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -a -installsuffix cgo -ldflags "-extldflags '-static' -X main.Version=${VERSION}" -o zitadelctl cmd/zitadelctl/main.go

## Go dev build
FROM go-base as dev-go-build
RUN go get github.com/go-delve/delve/cmd/dlv

#######################
## Final Production Image
#######################
FROM alpine:latest as artifact
RUN adduser -D zitadel

ARG ARCH=amd64
ARG OS=linux
COPY --from=prod-go-build /go/src/github.com/caos/zitadel/zitadelctl /app/zitadelctl
RUN chmod a+x /app/zitadelctl

## Scratch Image
FROM scratch as final
COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /app /
USER zitadel
HEALTHCHECK NONE
ENTRYPOINT ["/zitadelctl"]
