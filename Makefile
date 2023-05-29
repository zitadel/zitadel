# build-console:
# 	cd console
# 	yarn install --frozen-lockfile
#     yarn generate
#     yarn build

# build-core: install grpc static assets

# compile:
# 	rm internal/api/ui/console/static/*
# 	cp console/dist/console/* internal/api/ui/console/static/
# 	go build -o zitadel -ldflags="-s -w"

# install:
# 	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30
# 	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
# 	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
# 	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
# 	go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
# 	go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
# 	go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-zitadel
# 	go install github.com/rakyll/statik@v0.1.7

grpc:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	buf generate
	mkdir -p pkg/grpc
	mv .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc/
	mkdir -p openapi/v2/zitadel
	mv .artifacts/grpc/zitadel/ openapi/v2/zitadel
	rm -r .artifacts

static:
	go install github.com/rakyll/statik@v0.1.7
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/api/ui/login/static/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go

assets:
	go install github.com/rakyll/statik@v0.1.7
	go generate openapi/statik/generate.go && \
    mkdir -p docs/apis/assets/ && \
    go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

clean:
	rm -rf .artifacts/grpc

test:
	go test -race -v -coverprofile=profile.cov ./...
