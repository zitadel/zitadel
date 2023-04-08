FROM golang:1.19-bullseye
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY internal/api/ui/login/static internal/api/ui/login/static
COPY internal/api/ui/login/statik internal/api/ui/login/statik
COPY internal/notification/static internal/notification/static
COPY internal/notification/statik internal/notification/statik
COPY internal/static internal/static
COPY internal/statik internal/statik

RUN go generate internal/api/ui/login/statik/generate.go \
    && go generate internal/api/ui/login/static/generate.go \
    && go generate internal/notification/statik/generate.go \
    && go generate internal/statik/generate.go

COPY build/zitadel/generate-grpc.sh build/zitadel/generate-grpc.sh
COPY internal/protoc internal/protoc
COPY openapi/statik openapi/statik
COPY internal/api/assets/generator internal/api/assets/generator
COPY internal/config internal/config
COPY internal/errors internal/errors

RUN build/zitadel/generate-grpc.sh && \
    go generate openapi/statik/generate.go && \
    mkdir -p docs/apis/assets/ && \
    go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md