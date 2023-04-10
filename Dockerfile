
ARG NODE_VERSION=18
ARG GO_VERSION=1.19
# TODO add os and platform args

FROM node:${NODE_VERSION} as console-base
WORKDIR /zitadel/console
COPY console/package.json console/package-lock.json console/buf.gen.yaml ./
COPY proto ../proto
RUN npm ci && npm run generate
COPY console .

FROM console-base as console-lint
WORKDIR /zitadel/console
RUN npm run lint

FROM console-base as console-build
WORKDIR /zitadel/console
RUN npm run build

FROM golang:${GO_VERSION} as build

COPY --from=console-build /zitadel/console/dist/console /zitadel/console/dist/console

RUN ls -latr /zitadel/console/dist/console