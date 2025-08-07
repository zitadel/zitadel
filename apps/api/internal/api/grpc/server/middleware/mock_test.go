package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func emptyMockHandler(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func errorMockHandler(_ context.Context, req interface{}) (interface{}, error) {
	return nil, zerrors.ThrowInternal(nil, "test", "error")
}

type mockReq struct{}

func mockInfo(path string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: path,
	}
}
