FROM golang:1.19.2

WORKDIR /go/src/github.com/zitadel/zitadel

# Buf plugins must be built for linux/amd64
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

# install tools
COPY tools tools
RUN ./tools/install.sh

COPY internal/protoc internal/protoc

# generate authoptions code from templates
RUN go-bindata \
  -pkg main \
  -prefix internal/protoc/protoc-gen-authoption \
  -o /go/src/github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/templates.gen.go \
  /go/src/github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/templates

RUN cat internal/protoc/protoc-gen-authoption/templates.gen.go

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go install /go/src/github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption/

RUN go install github.com/bufbuild/buf/cmd/buf@v1.8.0

RUN buf generate
