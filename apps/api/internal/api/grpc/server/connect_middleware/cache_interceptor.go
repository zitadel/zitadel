package connect_middleware

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"

	_ "github.com/zitadel/zitadel/internal/statik"
)

func NoCacheInterceptor() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			headers := map[string]string{
				"cache-control": "no-store",
				"expires":       time.Now().UTC().Format(http.TimeFormat),
				"pragma":        "no-cache",
			}
			resp, err := handler(ctx, req)
			if err != nil {
				return nil, err
			}
			for key, value := range headers {
				resp.Header().Set(key, value)
			}
			return resp, err
		}
	}
}
