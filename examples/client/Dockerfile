#######################
## These steps set platform / arch type specific variables
#######################
FROM alpine AS arm64-base
ENV PROTOC_ARCH aarch_64

FROM alpine AS amd64-base
ENV PROTOC_ARCH x86_64

#######################
## This step downloads the protofiles, protoc and protoc-gen-grpc-web for later use
#######################
FROM ${BUILDARCH}-base as base
RUN apk add tar curl git
WORKDIR /.tmp
RUN protoc_url="https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-${PROTOC_ARCH}.zip" \
    && wget -O protoc ${protoc_url} \
    && unzip protoc
RUN curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.4.1/validate/validate.proto --create-dirs -o validate/validate.proto  \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/annotations.proto --create-dirs -o protoc-gen-swagger/options/annotations.proto \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v1.14.6/protoc-gen-swagger/options/openapiv2.proto --create-dirs -o protoc-gen-swagger/options/openapiv2.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto --create-dirs -o google/api/annotations.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto --create-dirs -o google/api/http.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/empty.proto --create-dirs -o google/protobuf/empty.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/timestamp.proto --create-dirs -o google/protobuf/timestamp.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/descriptor.proto --create-dirs -o google/protobuf/descriptor.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/duration.proto --create-dirs -o google/protobuf/duration.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/any.proto --create-dirs -o google/protobuf/any.proto \
    && curl https://raw.githubusercontent.com/protocolbuffers/protobuf/master/src/google/protobuf/struct.proto --create-dirs -o google/protobuf/struct.proto

RUN echo "a"
WORKDIR /zitadel
RUN git clone https://github.com/caos/zitadel \
    && cd zitadel \
    && git checkout origin/example-client

#######################
## Go base build
## Speed up this step by mounting your local go mod pkg directory
#######################
FROM golang:1.15 as zitadel-client
ARG GO_OUT
ARG GRPC_OUT

WORKDIR /go/src/github.com/caos/zitadel

# copy ZITADEL proto dependencies
COPY --from=base /.tmp /.tmp/protos
COPY --from=base /zitadel/zitadel .
COPY --from=base /.tmp/bin /usr/local/bin/

RUN go mod download

RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/envoyproxy/protoc-gen-validate

RUN go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN ./build/zitadel/generate-grpc-clients.sh /go/src /go/src github.com/caos/zitadel/examples/client/zitadel/message

FROM scratch as zitadel-copy
COPY --from=zitadel-client /go/src/github.com/caos/zitadel/pkg/grpc/admin/*.pb.go /zitadel/admin/
COPY --from=zitadel-client /go/src/github.com/caos/zitadel/pkg/grpc/auth/*.pb.go /zitadel/auth/
COPY --from=zitadel-client /go/src/github.com/caos/zitadel/pkg/grpc/management/*.pb.go /zitadel/management/
COPY --from=zitadel-client /go/src/github.com/caos/zitadel/pkg/grpc/message/*.pb.go /zitadel/message/
