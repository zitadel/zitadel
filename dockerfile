## Setup Base and generate proto stubs
FROM alpine as base

RUN apk add tar

RUN mkdir .tmp
WORKDIR /.tmp

RUN wget -O protoc https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip \
    && unzip protoc \
    && wget -O binprotoc-gen-grpc-web https://github.com/grpc/grpc-web/releases/download/1.2.0/protoc-gen-grpc-web-1.2.0-linux-x86_64

RUN ls -l /.tmp/bin

COPY pkg/grpc/*/proto/*.proto proto/

RUN ls -l proto

## Test Angular

## Build Angular

## Test Go

## Build Go