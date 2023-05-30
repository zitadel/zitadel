go_bin := "$$(go env GOPATH)/bin"
gen_authopt_path := "$(go_bin)/protoc-gen-authoption"
gen_zitadel_path := "$(go_bin)/protoc-gen-zitadel"

compile: core_build console_build
	cp -r console/dist/console internal/api/ui/console/static/
	go build -o zitadel-$$(go env GOOS)-$$(go env GOARCH) -ldflags="-s -w"

compile_pipeline:
	cp -r console/dist/console internal/api/ui/console/static/
	go build -o zitadel-$$(go env GOOS)-$$(go env GOARCH) -ldflags="-s -w"

core_dependencies:
	go mod download

core_static:
	go install github.com/rakyll/statik@v0.1.7
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/api/ui/login/static/resources/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go

core_assets:
	mkdir -p docs/apis/assets
	go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

core_api_generator:
ifeq (,$(wildcard $(gen_authopt_path)))
	go install internal/protoc/protoc-gen-authoption/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_authopt_path)
endif
ifeq (,$(wildcard $(gen_zitadel_path)))
	go install internal/protoc/protoc-gen-zitadel/main.go \
    && mv $$(go env GOPATH)/bin/main $(gen_zitadel_path)
endif

core_grpc_dependencies:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30 
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3 
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2 
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2 
	go install github.com/envoyproxy/protoc-gen-validate@v0.10.1

core_api: core_api_generator core_grpc_dependencies
	buf generate
	mkdir -p pkg/grpc
	cp -r .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/* pkg/grpc/
	mkdir -p openapi/v2/zitadel
	cp -r .artifacts/grpc/zitadel/ openapi/v2/zitadel

core_build: core_dependencies core_api core_static core_assets

console_dependencies:
	cd console && \
	yarn install --frozen-lockfile

console_client:
	cd console && \
	yarn generate

console_build: console_dependencies console_client
	cd console && \
	yarn build

clean:
	$(RM) .artifacts/grpc
	$(RM) $(gen_authopt_path)
	$(RM) $(gen_zitadel_path)

test:
	go test -race -v -coverprofile=profile.cov ./...
