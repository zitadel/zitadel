ARG NODE_VERSION=22
ARG GO_VERSION=1.23

## Console
FROM node:${NODE_VERSION} AS console-base
WORKDIR /app
COPY console/package.json console/yarn.lock console/buf.gen.yaml  ./
COPY proto/ ../proto/
RUN yarn install && yarn generate
COPY console/ .
COPY docs/frameworks.json ../docs/frameworks.json

FROM console-base AS console-build
RUN yarn build

FROM console-base AS console-lint
RUN yarn lint

FROM nginx:stable-alpine AS console-image
RUN rm -rf /usr/share/nginx/html/*
COPY .build/console/nginx.conf /etc/nginx/nginx.conf
COPY --from=console-build /app/dist /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]

FROM scratch AS console-output
COPY --from=console-build /app/dist/console .

## Core
FROM golang:${GO_VERSION} AS core-base
ARG SASS_VERSION=1.64.1
ARG GOLANG_CI_VERSION=1.64.5
RUN apt-get update && apt-get install -y npm && npm install -g sass@${SASS_VERSION}
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANG_CI_VERSION}
WORKDIR /app
COPY go.mod go.sum Makefile buf.gen.yaml buf.work.yaml main.go .golangci.yaml ./
COPY .git/ .git/
COPY cmd/ cmd/
COPY internal/ internal/
COPY openapi/ openapi/
COPY pkg/ pkg/
COPY proto/ proto/
COPY statik/ statik/
COPY --from=console-build /app/dist/console internal/api/ui/console/static
RUN make core_build

FROM core-base AS core-build
RUN make compile

FROM core-base AS core-lint
RUN make core_lint

FROM scratch AS core-output
COPY --from=core-build /app/zitadel .

FROM core-base AS core-image

FROM core-base AS core-unit-test
RUN make core_unit_test
