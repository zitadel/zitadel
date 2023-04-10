
ARG NODE_VERSION=18
ARG GO_VERSION=1.19
ARG CONSOLE_DIR=/zitadel/console
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

COPY --from=console-build ${CONSOLE_DIR}/dist/console ${CONSOLE_DIR}/dist/console

RUN ls -latr ${CONSOLE_DIR}/dist/console