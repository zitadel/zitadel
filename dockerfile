## This step downloads protoc and protoc-gen-grpc-web for later use
FROM alpine as protoc

RUN apk add tar curl
WORKDIR /.tmp

RUN wget -O protoc https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip \
    && unzip protoc \
    && wget -O bin/protoc-gen-grpc-web https://github.com/grpc/grpc-web/releases/download/1.2.0/protoc-gen-grpc-web-1.2.0-linux-x86_64

RUN chmod +x bin/protoc-gen-grpc-web

# This step downloads all protofiles
FROM alpine as proto

RUN apk add curl

WORKDIR /.tmp
RUN curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.4.0/validate/validate.proto --create-dirs -o validate/validate.proto  \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/annotations.proto --create-dirs -o protoc-gen-swagger/options/annotations.proto \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/openapiv2.proto --create-dirs -o protoc-gen-swagger/options/openapiv2.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto --create-dirs -o google/api/annotations.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/empty.proto --create-dirs -o google/protobuf/empty.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/timestamp.proto --create-dirs -o google/protobuf/timestamp.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/struct.proto --create-dirs -o google/protobuf/struct.proto

COPY pkg/grpc/admin/proto/admin.proto admin/proto/admin.proto
COPY pkg/grpc/auth/proto/auth.proto auth/proto/auth.proto
COPY pkg/grpc/management/proto/management.proto management/proto/management.proto
COPY pkg/grpc/message/proto/message.proto message/proto/message.proto
COPY internal/protoc/protoc-gen-authoption/authoption/options.proto authoption/options.proto

RUN ls -l

## With this step we prepare all node_modules, this helps caching the build
## Speed up this step by mounting your local node_modules directory
FROM node:12 as npminstaller

WORKDIR deps

COPY console/package.json console/package-lock.json ./

RUN npm install

## This step does build the angular code
FROM node:12 as npmbuilder

WORKDIR console

RUN mkdir .tmp

COPY console .
COPY --from=protoc /.tmp/bin /usr/local/bin/
COPY --from=proto /.tmp .tmp/protos/
COPY --from=npminstaller deps/node_modules node_modules/
COPY build/console build/console/

RUN build/console/generate-grpc.sh

## Override this when localy
ARG build=prodbuild
RUN npm run $build

RUN ls dist/console