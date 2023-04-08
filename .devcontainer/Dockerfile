ARG NODE_VERSION=18
ARG GO_VERSION=1.19

FROM node:${NODE_VERSION}-bullseye AS console
WORKDIR /console
COPY console/package.json console/package-lock.json ./
RUN npm ci
RUN npm install -g @angular/cli
COPY proto/ /proto/
COPY console/buf.gen.yaml buf.gen.yaml
RUN npm run generate

# ----

FROM node:${NODE_VERSION}-bullseye AS docs
WORKDIR /docs
COPY docs/package.json docs/yarn.lock ./
RUN yarn --frozen-lockfile
COPY proto/ /proto/
COPY docs/buf.gen.yaml .
COPY docs/docusaurus.config.js .
COPY docs/sidebars.js .
RUN yarn run generate:grpc

# ----

FROM debian:bullseye AS image
RUN apt update && apt install -y curl
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt install -y nodejs
RUN npm install -g @angular/cli && npm install -g yarn
COPY --from=console console console
COPY --from=docs docs docs




