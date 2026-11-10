package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func emptyMockHandler(_ context.Context, req any) (any, error) {
	return req, nil
}

func errorMockHandler(_ context.Context, req any) (any, error) {
	return nil, zerrors.ThrowPreconditionFailed(nil, "test", "error")
}

func panicMockHandler(payload any) func(context.Context, any) (any, error) {
	return func(context.Context, any) (any, error) {
		panic(payload)
	}
}

type mockReq struct{}

func mockInfo(path string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: path,
	}
}
