grpc:
	go get github.com/go-bindata/go-bindata/v3/go-bindata
	go install github.com/go-bindata/go-bindata/v3/go-bindata
	~/go/bin/go-bindata \
	-pkg main \
	-prefix internal/protoc/protoc-gen-authoption \
	-o internal/protoc/protoc-gen-authoption/templates.gen.go \
	internal/protoc/protoc-gen-authoption/templates
	go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
	rm -rf .artifacts/grpc
	buf generate

grpc_copy:
	cp -rT .artifacts/grpc/go-client/pkg/grpc pkg/grpc
	cp -rT .artifacts/grpc/go-client/openapi openapi
	cp -rT .artifacts/grpc/go-client/internal internal

test:
	go test -race -v -coverprofile=profile.cov ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

install:
	go mod download

build:
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/api/ui/login/static/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go
