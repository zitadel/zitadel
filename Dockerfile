
ARG NODE_VERSION=18
ARG GO_VERSION=1.19
# TODO add os and platform args

# -----

# FROM node:${NODE_VERSION} as console-base
# WORKDIR /zitadel/console
# COPY console/package.json console/package-lock.json console/buf.gen.yaml ./
# COPY proto ../proto
# RUN npm ci && npm run generate
# COPY console .

# FROM console-base as console-lint
# WORKDIR /zitadel/console
# RUN npm run lint

# FROM console-base as console-build
# WORKDIR /zitadel/console
# RUN npm run build

# -----

FROM golang:${GO_VERSION} as core-base
WORKDIR /zitadel
COPY go.mod go.sum buf.gen.yaml grpc.sh ui.sh ./
RUN go mod download
COPY internal/protoc/ internal/protoc/
COPY proto/ proto/
RUN bash grpc.sh
COPY internal/ internal/
RUN bash ui.sh
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY statik/ statik/
COPY openapi/ openapi/
COPY main.go LICENSE ./

FROM core-base as core-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2 \
    && golangci-lint run -n

FROM core-base as core-test
RUN ls -la pkg/grpc/action
RUN go test -race -v -coverprofile=profile.cov $(go list ./...)

FROM core-base as core-build
COPY --from=console-build /zitadel/console/dist/console internal/api/ui/console/static/
RUN go build -o zitadel

FROM scratch as core-export
COPY --from=core-build zitadel/zitadel zitadel