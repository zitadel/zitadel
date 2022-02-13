ARG GO_VERSION=1.17

#######################
## Go dependencies
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:${GO_VERSION} as go-dep
RUN mkdir -p src/github.com/caos/zitadel
WORKDIR /go/src/github.com/caos/zitadel

#download modules
COPY ./go.* .
RUN go mod download

# install tools
COPY tools ./tools
RUN ./tools/install.sh

ARG COCKROACH_BINARY=cockroach
RUN apt install openssl tzdata tar

# cockroach binary used to backup database
RUN mkdir /usr/local/lib/cockroach
RUN wget -qO- https://binaries.cockroachdb.com/cockroach-v21.2.5.linux-amd64.tgz \
    | tar  xvz && cp -i cockroach-v21.2.5.linux-amd64/cockroach /usr/local/bin/
RUN rm -r cockroach-v21.2.5.linux-amd64

#######################
## generates static files
#######################
FROM go-dep AS go-static
COPY internal/ui/login/static internal/ui/login/static
COPY internal/ui/login/statik internal/ui/login/statik
COPY internal/notification/static internal/notification/static
COPY internal/notification/statik internal/notification/statik
COPY internal/static internal/static
COPY internal/statik internal/statik

RUN go generate internal/ui/login/statik/generate.go \
    && go generate internal/ui/login/static/generate.go \
    && go generate internal/notification/statik/generate.go \
    && go generate internal/statik/generate.go

#######################
## generates grpc stub
#######################
FROM go-static AS go-stub
COPY --from=base /proto /proto
COPY --from=base /usr/local/bin /usr/local/bin/.

COPY build/zitadel/generate-grpc.sh build/zitadel/generate-grpc.sh
COPY internal/protoc internal/protoc
COPY openapi/statik openapi/statik
COPY internal/api/assets/generator internal/api/assets/generator
COPY internal/config internal/config
COPY internal/errors internal/errors

RUN build/zitadel/generate-grpc.sh \
    && go generate openapi/statik/generate.go \
    && mkdir -p docs/apis/assets/ \
    && go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

#######################
## Go base build
#######################
FROM go-stub as go-base
# copy remaining zitadel files
COPY . .

#######################
## copy for local dev
#######################
FROM scratch as go-copy
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/ui/login/statik/statik.go internal/ui/login/statik/statik.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/notification/statik/statik.go internal/notification/statik/statik.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/statik/statik.go internal/statik/statik.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/openapi/statik/statik.go openapi/statik/statik.go

COPY --from=go-stub /go/src/github.com/caos/zitadel/pkg/grpc pkg/grpc
COPY --from=go-stub /go/src/github.com/caos/zitadel/openapi/v2/zitadel openapi/v2/zitadel
COPY --from=go-stub /go/src/github.com/caos/zitadel/openapi/statik/statik.go openapi/statik/statik.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption/templates.gen.go internal/protoc/protoc-gen-authoption/templates.gen.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption/authoption/options.pb.go internal/protoc/protoc-gen-authoption/authoption/options.pb.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/docs/apis/proto docs/docs/apis/proto
COPY --from=go-stub /go/src/github.com/caos/zitadel/docs/apis/assets docs/docs/apis/assets

COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/api/assets/authz.go ./internal/api/assets/authz.go
COPY --from=go-stub /go/src/github.com/caos/zitadel/internal/api/assets/router.go ./internal/api/assets/router.go

#######################
## Go test
#######################
FROM go-base as go-test

# Migrations for cockroach-secure
RUN go install github.com/rakyll/statik \
    && go test -race -v -coverprofile=profile.cov $(go list ./... | grep -v /operator/)

#######################
## Go test results
#######################
FROM scratch as go-codecov
COPY --from=go-test /go/src/github.com/caos/zitadel/profile.cov profile.cov
