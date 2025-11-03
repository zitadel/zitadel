package connect_middleware

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	_ "github.com/zitadel/zitadel/internal/statik"
)

func ErrorHandler() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return toConnectError(ctx, req, handler)
		}
	}
}

func toConnectError(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc) (connect.AnyResponse, error) {
	resp, err := handler(ctx, req)
	return resp, gerrors.ZITADELToConnectError(err) // TODO !
}
