#!/bin/bash

set -eux

go install github.com/rakyll/statik@v0.1.7

go generate internal/api/ui/login/statik/generate.go
# go generate internal/api/ui/login/static/generate.go # This looks like it should be go generate internal/api/ui/login/static/resources/generate.go but i used this from prod
go generate internal/notification/statik/generate.go
go generate internal/statik/generate.go

mkdir -p openapi/v2/zitadel
cp .artifacts/grpc/zitadel/*.swagger.json openapi/v2/zitadel
go generate openapi/statik/generate.go
mkdir -p docs/apis/assets/
go run internal/api/assets/generator/asset_generator.go -directory=internal/api/assets/generator/ -assets=docs/apis/assets/assets.md