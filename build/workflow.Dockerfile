# ##############################################################################
# core
# ##############################################################################

# #######################################
# download dependencies
# #######################################
FROM golang:buster AS core-deps

WORKDIR /go/src/github.com/zitadel/zitadel

COPY go.mod .
COPY go.sum .

RUN go mod download

# #######################################
# compile custom protoc plugins
# #######################################
FROM golang:buster AS core-api-generator

WORKDIR /go/src/github.com/zitadel/zitadel

COPY go.mod .
COPY go.sum .
COPY internal/protoc internal/protoc
COPY pkg/grpc/protoc/v2 pkg/grpc/protoc/v2

RUN go install internal/protoc/protoc-gen-authoption/main.go \
    && mv $(go env GOPATH)/bin/main $(go env GOPATH)/bin/protoc-gen-authoption \
	&& go install internal/protoc/protoc-gen-zitadel/main.go \
    && mv $(go env GOPATH)/bin/main $(go env GOPATH)/bin/protoc-gen-zitadel

# #######################################
# build backend stub
# #######################################
FROM golang:buster AS core-api

WORKDIR /go/src/github.com/zitadel/zitadel

COPY go.mod .
COPY go.sum .
COPY proto proto
COPY buf.*.yaml .
COPY Makefile Makefile
COPY --from=core-api-generator /go/bin /usr/local/bin

RUN make grpc

# #######################################
# generate code for login ui
# #######################################
FROM golang:buster AS core-login

WORKDIR /go/src/github.com/zitadel/zitadel

COPY Makefile Makefile
COPY internal/api/ui/login/static internal/api/ui/login/static
COPY internal/api/ui/login/statik internal/api/ui/login/statik
COPY internal/notification/static internal/notification/static
COPY internal/notification/statik internal/notification/statik
COPY internal/static internal/static
COPY internal/statik internal/statik

RUN make static

# #######################################
# generate code for assets
# #######################################
FROM golang:buster AS core-assets
WORKDIR /go/src/github.com/zitadel/zitadel

COPY go.mod .
COPY go.sum .
COPY Makefile Makefile
COPY internal/api/assets/generator internal/api/assets/generator
COPY internal/config internal/config
COPY internal/errors internal/errors
COPY --from=core-api /go/src/github.com/zitadel/zitadel/openapi/v2 openapi/v2

RUN make assets

# #######################################
# Gather all core files
# #######################################
FROM core-deps AS core-gathered

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY proto proto
COPY openapi openapi
COPY statik statik
COPY main.go main.go
COPY --from=core-api /go/src/github.com/zitadel/zitadel .
COPY --from=core-login /go/src/github.com/zitadel/zitadel .
COPY --from=core-assets /go/src/github.com/zitadel/zitadel/internal ./internal

# ##############################################################################
# build console
# ##############################################################################

# #######################################
# download console dependencies
# #######################################
FROM node:20-buster AS console-deps

WORKDIR /zitadel/console

COPY console/package.json .
COPY console/yarn.lock .

RUN yarn install --frozen-lockfile

# #######################################
# generate console client
# #######################################
FROM node:20-buster AS console-client

WORKDIR /zitadel/console

# install buf
COPY --from=bufbuild/buf:latest /usr/local/bin/* /usr/local/bin/
ENV PATH="/usr/local/bin:${PATH}"

COPY console/package.json .
COPY console/buf.*.yaml .
COPY proto ../proto

RUN yarn generate

# #######################################
# Gather all console files
# #######################################
FROM console-deps as console-gathered

COPY --from=console-client /zitadel/console/src/app/proto/generated src/app/proto/generated

COPY console/src src
COPY console/angular.json .
COPY console/ngsw-config.json .
COPY console/tsconfig* .

# #######################################
# Build console
# #######################################
FROM console-gathered AS console
RUN yarn build

# ##############################################################################
# build the executable
# ##############################################################################

# #######################################
# build executable
# #######################################
FROM core-gathered AS compile

ARG GOOS 
ARG GOARCH

COPY --from=console /zitadel/console/dist/console internal/api/ui/console/static/

RUN go build -o zitadel -ldflags="-s -w -race" \
    && chmod +x zitadel

ENTRYPOINT [ "./zitadel" ]

# #######################################
# copy executable
# #######################################
FROM scratch AS copy-executable
ARG GOOS 
ARG GOARCH

COPY --from=compile /go/src/github.com/zitadel/zitadel/zitadel /.artifacts/zitadel

# ##############################################################################
#  tests
# ##############################################################################
FROM ubuntu/postgres:latest AS test-core-base

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        gcc \
        make \
        ca-certificates \
        gcc \
        && \
    update-ca-certificates; \
    rm -rf /var/lib/apt/lists/*

# install go
COPY --from=golang:latest /usr/local/go/ /usr/local/go/
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

WORKDIR /go/src/github.com/zitadel/zitadel

# default vars
ENV DB_FLAVOR=postgres
ENV POSTGRES_USER=zitadel
ENV POSTGRES_DB=zitadel
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_HOST_AUTH_METHOD=trust

ENV PGUSER=zitadel
ENV PGDATABASE=zitadel
ENV PGPASSWORD=postgres

ENV CGO_ENABLED=1

# copy zitadel files
COPY --from=core-deps /go/pkg/mod /root/go/pkg/mod
COPY --from=core-gathered /go/src/github.com/zitadel/zitadel .

# #######################################
# unit test core
# #######################################
FROM test-core-base AS test-core-unit
RUN go test -race -v -coverprofile=profile.cov ./...

# #######################################
# coverage output
# #######################################
FROM scratch AS coverage-core-unit
COPY --from=test-core-unit /go/src/github.com/zitadel/zitadel/profile.cov /coverage/

# #######################################
# integration test core
# #######################################
FROM test-core-base AS test-core-integration
ENV DB_FLAVOR=cockroach

# install cockroach
COPY --from=cockroachdb/cockroach:latest-v24.3 /cockroach/cockroach /usr/local/bin/
ENV COCKROACH_BINARY=/cockroach/cockroach

ENV ZITADEL_MASTERKEY=MasterkeyNeedsToHave32Characters

COPY build/core-integration-test.sh /usr/local/bin/run-tests.sh
RUN chmod +x /usr/local/bin/run-tests.sh

RUN run-tests.sh

# #######################################
# coverage output
# #######################################
FROM scratch AS coverage-core-integration
COPY --from=test-core-integration /go/src/github.com/zitadel/zitadel/profile.cov /coverage/

# ##############################################################################
#  linting
# ##############################################################################

# #######################################
# api
# #######################################
FROM bufbuild/buf:latest AS lint-api

COPY proto proto
COPY buf.*.yaml .

RUN buf lint

# #######################################
# console
# #######################################
FROM console-gathered AS lint-console

COPY console/.eslintrc.js .
COPY console/.prettier* .
RUN yarn lint

# #######################################
# core
# #######################################
FROM golangci/golangci-lint:latest AS lint-core
ARG LINT_EXIT_CODE=1

WORKDIR /go/src/github.com/zitadel/zitadel

COPY .golangci.yaml .
COPY .git/ .git/
COPY --from=core-deps /go/pkg/mod /go/pkg/mod
COPY --from=core-gathered /go/src/github.com/zitadel/zitadel .

RUN git fetch https://github.com/zitadel/zitadel main:main

RUN golangci-lint run \
    --timeout 10m \
    --config ./.golangci.yaml \
    --out-format=github-actions:report,colored-line-number \
    --issues-exit-code=${LINT_EXIT_CODE} \
    --concurrency=$(getconf _NPROCESSORS_ONLN)

# #######################################
# report output
# #######################################
FROM scratch AS lint-core-report
COPY --from=lint-core /go/src/github.com/zitadel/zitadel/report .