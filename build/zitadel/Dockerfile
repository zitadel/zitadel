ARG GO_VERSION=1.19

#######################
## Go dependencies
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:${GO_VERSION} as go-dep
WORKDIR /go/src/github.com/zitadel/zitadel

#download modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# install tools
COPY tools ./tools
RUN ./tools/install.sh

#######################
## generates static files
#######################
FROM go-dep AS go-static
COPY internal/api/ui/login/static internal/api/ui/login/static
COPY internal/api/ui/login/statik internal/api/ui/login/statik
COPY internal/notification/static internal/notification/static
COPY internal/notification/statik internal/notification/statik
COPY internal/static internal/static
COPY internal/statik internal/statik

RUN go generate internal/api/ui/login/statik/generate.go \
    && go generate internal/api/ui/login/static/generate.go \
    && go generate internal/notification/statik/generate.go \
    && go generate internal/statik/generate.go

#######################
## generates grpc stub
#######################
FROM go-static AS go-stub
COPY --from=zitadel-base:local /proto /proto
COPY --from=zitadel-base:local /usr/local/bin /usr/local/bin/.

COPY build/zitadel/generate-grpc.sh build/zitadel/generate-grpc.sh
COPY internal/protoc internal/protoc
COPY openapi/statik openapi/statik
COPY internal/api/assets/generator internal/api/assets/generator
COPY internal/config internal/config
COPY internal/errors internal/errors

RUN build/zitadel/generate-grpc.sh && \
    go generate openapi/statik/generate.go && \
    mkdir -p docs/apis/assets/ && \
    go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

#######################
## Go base build
#######################
FROM go-stub as go-base
# copy remaining zitadel files
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY openapi openapi

#######################
## copy for local dev
#######################
FROM scratch as go-copy
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/api/ui/login/statik/statik.go internal/api/ui/login/statik/statik.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/notification/statik/statik.go internal/notification/statik/statik.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/statik/statik.go internal/statik/statik.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/openapi/statik/statik.go openapi/statik/statik.go

COPY --from=go-stub /go/src/github.com/zitadel/zitadel/pkg/grpc pkg/grpc
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/openapi/v2/zitadel openapi/v2/zitadel
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/openapi/statik/statik.go openapi/statik/statik.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/templates.gen.go internal/protoc/protoc-gen-authoption/templates.gen.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/authoption/options.pb.go internal/protoc/protoc-gen-authoption/authoption/options.pb.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/docs/apis/proto docs/docs/apis/proto
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/docs/apis/assets docs/docs/apis/assets

COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/api/assets/authz.go ./internal/api/assets/authz.go
COPY --from=go-stub /go/src/github.com/zitadel/zitadel/internal/api/assets/router.go ./internal/api/assets/router.go

#######################
## Go test
#######################
FROM go-base as go-test

ARG COCKROACH_BINARY=/usr/local/bin/cockroach
ARG COCKROACH_VERSION=v22.2.2
RUN apt install openssl tzdata tar

# cockroach binary used to backup database
RUN mkdir /usr/local/lib/cockroach
RUN wget -qO- https://binaries.cockroachdb.com/cockroach-${COCKROACH_VERSION}.linux-amd64.tgz \
    | tar  xvz && cp -i cockroach-${COCKROACH_VERSION}.linux-amd64/cockroach /usr/local/bin/
RUN rm -r cockroach-${COCKROACH_VERSION}.linux-amd64

# Migrations for cockroach-secure
RUN go install github.com/rakyll/statik \
    && go test -race -v -coverprofile=profile.cov $(go list ./... | grep -v /operator/)

#######################
## Go test results
#######################
FROM scratch as go-codecov
COPY --from=go-test /go/src/github.com/zitadel/zitadel/profile.cov profile.cov

