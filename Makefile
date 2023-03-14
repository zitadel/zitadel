grpc:
	rm -rf .artifacts/grpc
	docker build -f build/grpc/Dockerfile -t zitadel-base:local .
	docker build -f build/zitadel/Dockerfile . -t zitadel-go-base --target go-copy -o .artifacts/grpc/go-client

grpc_copy:
	make cp -rT .artifacts/grpc/go-client/pkg/grpc pkg/grpc
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