package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/caos/zitadel/internal/api/http"
)

func GetHeader(ctx context.Context, headername string) string {
	return metautils.ExtractIncoming(ctx).Get(headername)
}

func GetAuthorizationHeader(ctx context.Context) string {
	return GetHeader(ctx, http.Authorization)
}
