ARG NODE_VERSION=14

#######################
## These steps set platform / arch type specific variables
#######################
FROM alpine:3 AS arm64-base
ENV PROTOC_ARCH aarch_64

FROM alpine:3 AS amd64-base
ENV PROTOC_ARCH x86_64

#######################
## This step sets up the folder structure,
## initalices go mods,
## downloads the protofiles,
## protoc and protoc-gen-grpc-web for later use
#######################
FROM ${BUILDARCH}-base AS base
ARG PROTOC_VERSION=3.13.0
ARG PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-${PROTOC_ARCH}.zip
ARG GRPC_WEB_VERSION=1.3.0
# no arm specific version available and x86 works fine at the moment:
ARG GRPC_WEB=protoc-gen-grpc-web-${GRPC_WEB_VERSION}-linux-x86_64

RUN apk add tar curl
WORKDIR /proto

#protoc
RUN apk add tar curl \
    && curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /proto include/* \
    && rm -f $PROTOC_ZIP \
    && curl -OL https://github.com/grpc/grpc-web/releases/download/${GRPC_WEB_VERSION}/${GRPC_WEB} \
    && mv ${GRPC_WEB} /usr/local/bin/protoc-gen-grpc-web \
    && chmod +x /usr/local/bin/protoc-gen-grpc-web \
    && curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.6.2/validate/validate.proto --create-dirs -o include/validate/validate.proto  \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v2.2.0/protoc-gen-openapiv2/options/annotations.proto --create-dirs -o include/protoc-gen-openapiv2/options/annotations.proto \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v2.2.0/protoc-gen-openapiv2/options/openapiv2.proto --create-dirs -o include/protoc-gen-openapiv2/options/openapiv2.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto --create-dirs -o include/google/api/annotations.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto --create-dirs -o include/google/api/http.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto --create-dirs -o include/google/api/field_behavior.proto

#zitadel protos
COPY proto/ include/.

#######################
## With this step we prepare all node_modules, this helps caching the build
## Speed up this step by mounting your local node_modules directory
#######################
FROM node:${NODE_VERSION} as npm-base
WORKDIR /console

COPY console/package.json console/package-lock.json ./
RUN npm ci --loglevel verbose

COPY console .
COPY --from=base /proto /proto
COPY --from=base /usr/local/bin /usr/local/bin/.
COPY build/console build/console/
RUN build/console/generate-grpc.sh

#######################
## copy for local dev
#######################
FROM scratch as npm-copy
COPY --from=npm-base /console/src/app/proto/generated ./console/src/app/proto/generated

######################
## Output generated files by running a container from this target mounting the console dir to /console and your ~/.cache dir to /cache
######################
FROM base as init-dev

RUN apk add rsync

COPY --from=npm-base /console/node_modules /builtconsole/node_modules
COPY --from=npm-base /console/src/app/proto/generated /builtconsole/src/app/proto/generated
COPY --from=npm-base /root/.cache /builtcache

ARG UID=1000
ARG GID=1000

# Create a group and user
RUN addgroup -S tmpgroup -g $GID && adduser -S tmpuser -G tmpgroup -u $UID

ENTRYPOINT [ "sh", "-c" ]
CMD [ "rsync -a --chown tmpuser:tmpgroup /builtconsole/ /console && rsync -a --chown tmpuser:tmpgroup /builtcache/ /cache" ]

#######################
## angular dev build
#######################
FROM npm-base as dev-angular-build
RUN npm install -g @angular/cli

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as prod-angular-build
RUN npm run lint
RUN npm run prodbuild

#######################
## Only Copy Assets
#######################
FROM scratch as prod-angular-export
COPY --from=prod-angular-build /console/dist/console .
