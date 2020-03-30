package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/caos/zitadel/internal/api"
)

func GetHeader(ctx context.Context, headername string) string {
	return metautils.ExtractIncoming(ctx).Get(headername)
}

func GetAuthorizationHeader(ctx context.Context) string {
	return GetHeader(ctx, api.Authorization)
}
