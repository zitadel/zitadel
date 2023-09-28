package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", fmt.Errorf("cannot read metadata")
		}
		origins := md.Get(http.Origin)
		if len(origins) == 1 {
			ctx = http.WithOrigin(ctx, origins[0])
		}
		return handler(ctx, req)
	}
}
