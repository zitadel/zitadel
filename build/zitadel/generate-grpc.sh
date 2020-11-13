#! /bin/sh

set -eux

go generate pkg/grpc/auth/proto/generate.go
go generate pkg/grpc/admin/proto/generate.go
go generate pkg/grpc/management/proto/generate.go
go generate pkg/grpc/message/proto/generate.go
