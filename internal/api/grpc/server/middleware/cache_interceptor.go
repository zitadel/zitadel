package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/caos/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	_ "github.com/caos/zitadel/internal/statik"
)

func NoCacheInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		headers := map[string]string{
			"cache-control": "no-store",
			"expires":       time.Now().UTC().Format(http.TimeFormat),
			"pragma":        "no-cache",
		}
		header := metadata.New(headers)
		for key, value := range headers {
			header.Append(runtime.MetadataHeaderPrefix+key, value)
		}
		err := grpc.SendHeader(ctx, header)
		logging.Log("MIDDLE-efh41").OnError(err).WithField("req", info.FullMethod).Warn("cannot send cache-control on grpc response")
		resp, err := handler(ctx, req)
		return resp, err
	}
}
