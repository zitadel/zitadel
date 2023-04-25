#ARG BUILDARCH=x86_64

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
FROM ${BUILDARCH}-base
ARG PROTOC_VERSION=3.18.0
ARG PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-${PROTOC_ARCH}.zip
ARG GRPC_WEB_VERSION=1.3.0
ARG GATEWAY_VERSION=2.15.1
ARG VALIDATOR_VERSION=0.6.2
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
    && curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v${VALIDATOR_VERSION}/validate/validate.proto --create-dirs -o include/validate/validate.proto  \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v${GATEWAY_VERSION}/protoc-gen-openapiv2/options/annotations.proto --create-dirs -o include/protoc-gen-openapiv2/options/annotations.proto \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v${GATEWAY_VERSION}/protoc-gen-openapiv2/options/openapiv2.proto --create-dirs -o include/protoc-gen-openapiv2/options/openapiv2.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto --create-dirs -o include/google/api/annotations.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto --create-dirs -o include/google/api/http.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto --create-dirs -o include/google/api/field_behavior.proto

#zitadel protos
COPY proto/ include/.
