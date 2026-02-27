package connect_middleware

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func ErrorHandler() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return toConnectError(ctx, req, handler)
		}
	}
}

func toConnectError(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc) (_ connect.AnyResponse, err error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer func() {
		if rec := recover(); rec != nil {
			recErr, ok := rec.(error)
			if !ok {
				recErr = fmt.Errorf("%v", rec)
			}
			if recErr != nil {
				err = zerrors.ThrowInternal(recErr, zerrors.IDRecover, "Errors.Internal")
			}
		}
		err = gerrors.ZITADELToConnectError(ctx, err)
		cancel(err)
	}()
	return handler(ctx, req)
}
