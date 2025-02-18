ARG NODE_VERSION=22
ARG GO_VERSION=1.23

## Console Base
FROM node:${NODE_VERSION} AS console-base
WORKDIR /app
COPY console/package.json console/yarn.lock console/buf.gen.yaml  ./
COPY proto/ ../proto/
RUN yarn install && yarn generate

## Console Build
FROM console-base AS console-build
COPY console/ .
COPY docs/frameworks.json ../docs/frameworks.json
RUN yarn build

## Console Image
FROM nginx:stable-alpine AS console
RUN rm -rf /usr/share/nginx/html/*
COPY .build/console/nginx.conf /etc/nginx/nginx.conf
COPY --from=console-build /app/dist /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]

## Core Base
FROM golang:${GO_VERSION} AS core-base
ARG SASS_VERSION=1.64.1
RUN apt-get update && apt-get install -y npm && npm install -g sass@${SASS_VERSION}
WORKDIR /app
COPY go.mod go.sum Makefile buf.gen.yaml buf.work.yaml main.go ./
COPY cmd/ cmd/
COPY internal/ internal/
COPY openapi/ openapi/
COPY pkg/ pkg/
COPY proto/ proto/
COPY statik/ statik/
COPY --from=console-build /app/dist/console internal/api/ui/console/static
RUN ls -la proto/zitadel
RUN make core_build

## Core Unit Test
FROM core-base AS core-unit-test
RUN make core_unit_test
