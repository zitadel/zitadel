ARG NODE_VERSION=16-alpine
ARG GO_VERSION=1.17-alpine

#######################
## With this step we prepare all node_modules, this helps caching the build
## Speed up this step by mounting your local node_modules directory
#######################
FROM --platform=$BUILDPLATFORM rvolosatovs/protoc:3 as angular-build
RUN apk add bash
WORKDIR /console

COPY console/package.json console/package-lock.json ./
RUN npm ci

COPY console .
COPY --from=base /proto /proto
COPY --from=base /usr/local/bin /usr/local/bin/.
COPY build/console build/console/
RUN build/console/generate-grpc.sh

#######################
## angular lint workspace and prod build
#######################
RUN npm install -g @angular/cli \
    && npm run lint \
    && npm run prodbuild

#######################
## Go dependencies
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} as go-build
ARG TARGETOS TARGETARCH
RUN mkdir -p src/github.com/caos/zitadel
WORKDIR /go/src/github.com/caos/zitadel

#download modules
COPY ./go.* .
RUN go mod download

# install tools
COPY tools ./tools
RUN ./tools/install.sh

#######################
## generates static files
#######################

## TODO Embed the console
COPY --from=angular-build /console/dist/console .

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
# copy remaining zitadel files
COPY . .

#######################
## Go test
#######################
# Migrations for cockroach-secure
RUN go install github.com/rakyll/statik \
    && go test -race -v -coverprofile=profile.cov $(go list ./... | grep -v /operator/)

#######################
## Go prod build
#######################
ARG VERSION=""

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION:-'dev'} -extldflags \"-static\"" -o zitadel main.go

#######################
## Go dev build
#######################
FROM go-build as dev-go-build
ENTRYPOINT [ "go", "run", "cmd/zitadel/main.go" ]

#######################
## Stage to export ZITADEL Binary
#######################
FROM --platform=$BUILDPLATFORM scratch as go-output
ARG TARGETOS TARGETARCH
COPY --from=go-build /go/src/github.com/caos/zitadel/zitadel .

#######################
## Stage to export the generated proto files
#######################
FROM scratch as go-client

## GO gRCP Client
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/ui/login/statik/statik.go internal/ui/login/statik/statik.go
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/notification/statik/statik.go internal/notification/statik/statik.go
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/statik/statik.go internal/statik/statik.go
COPY --from=go-build /go/src/github.com/caos/zitadel/openapi/statik/statik.go openapi/statik/statik.go

COPY --from=go-build /go/src/github.com/caos/zitadel/pkg/grpc pkg/grpc
COPY --from=go-build /go/src/github.com/caos/zitadel/openapi/v2/zitadel openapi/v2/zitadel
COPY --from=go-build /go/src/github.com/caos/zitadel/openapi/statik/statik.go openapi/statik/statik.go
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption/templates.gen.go internal/protoc/protoc-gen-authoption/templates.gen.go
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/protoc/protoc-gen-authoption/authoption/options.pb.go internal/protoc/protoc-gen-authoption/authoption/options.pb.go
COPY --from=go-build /go/src/github.com/caos/zitadel/docs/apis/proto docs/docs/apis/proto
COPY --from=go-build /go/src/github.com/caos/zitadel/docs/apis/assets docs/docs/apis/assets

COPY --from=go-build /go/src/github.com/caos/zitadel/internal/api/assets/authz.go ./internal/api/assets/authz.go
COPY --from=go-build /go/src/github.com/caos/zitadel/internal/api/assets/router.go ./internal/api/assets/router.go

#######################
## Stage to export the generated proto files
#######################
FROM scratch as js-client
COPY --from=angular-build /console/src/app/proto/generated ./console/src/app/proto/generated

#######################
## Go test results
#######################
FROM scratch as go-codecov
COPY --from=go-build /go/src/github.com/caos/zitadel/profile.cov profile.cov