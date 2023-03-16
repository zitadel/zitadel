grpc:
	go get github.com/go-bindata/go-bindata/v3/go-bindata
	go install github.com/go-bindata/go-bindata/v3/go-bindata
	~/go/bin/go-bindata \
	-pkg main \
	-prefix internal/protoc/protoc-gen-authoption \
	-o internal/protoc/protoc-gen-authoption/templates.gen.go \
	internal/protoc/protoc-gen-authoption/templates
	rm -rf .artifacts/grpc
	buf generate proto/zitadel/admin.proto
	buf generate proto/zitadel/system.proto
	buf generate proto/zitadel/management.proto
	buf generate proto/zitadel/auth.proto
	
grpc_copy:
	cp -rT .artifacts/grpc/go-client/pkg/grpc pkg/grpc
	cp -rT .artifacts/grpc/go-client/openapi openapi
	cp -rT .artifacts/grpc/go-client/internal internal

blub:
	protoc \
	-I=proto/ \
	-I=${HOME}/.cache/buf/v1/module/data/buf.build/envoyproxy/protoc-gen-validate/6607b10f00ed4a3d98f906807131c44a/ \
	-I=${HOME}/.cache/buf/v1/module/data/buf.build/grpc-ecosystem/grpc-gateway/a1ecdc58eccd49aa8bea2a7a9022dc27/ \
	-I=${HOME}/.cache/buf/v1/module/data/buf.build/googleapis/googleapis/75b4300737fb4efca0831636be94e517/ \
	-I=${HOME}/.cache/buf/v1/module/data/buf.build/googleapis/googleapis/62f35d8aed1149c291d606d958a7ce32/ \
	--authoption_out .artifacts/ proto/zitadel/system.proto proto/zitadel/admin.proto

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
