grpc:
	go install github.com/go-bindata/go-bindata
	$$(go env GOPATH)/bin/go-bindata \
	-pkg main \
	-prefix internal/protoc/protoc-gen-authoption \
	-o internal/protoc/protoc-gen-authoption/templates.gen.go \
	internal/protoc/protoc-gen-authoption/templates
	go install github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption
	rm -rf .artifacts/grpc
	# This foreach is a workaround from a limitation of the authoption generator and only affects zitadel.
	# The authoption generator cannot work when passed *.proto but instead needs to have each file passed as {name}.proto
	for i in $$(find proto/zitadel -iname *.proto); do buf generate $${i}; done
	mv .artifacts/grpc/zitadel/auth.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/auth
	mv .artifacts/grpc/zitadel/admin.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/admin
	mv .artifacts/grpc/zitadel/management.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/management
	mv .artifacts/grpc/zitadel/system.pb.authoptions.go .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/system
	cp -rT .artifacts/grpc/github.com/zitadel/zitadel/pkg/grpc/ pkg/grpc/
	mkdir -p openapi/v2/zitadel
	cp -rT .artifacts/grpc/zitadel openapi/v2/zitadel

static:
	go install github.com/rakyll/statik@v0.1.7
	go generate internal/api/ui/login/statik/generate.go
	go generate internal/api/ui/login/static/generate.go
	go generate internal/notification/statik/generate.go
	go generate internal/statik/generate.go

assets:
	go generate openapi/statik/generate.go && \
    mkdir -p docs/apis/assets/ && \
    go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md

test: grpc static assets
	go test -race -v -coverprofile=profile.cov ./...

lint: grpc static assets
	golangci-lint run

generate: grpc static assets

tidy:
	go mod tidy

install:
	go mod download

build: grpc static assets
	go build