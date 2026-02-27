package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	_ "github.com/zitadel/zitadel/internal/statik"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return toGRPCError(ctx, req, handler)
	}
}

func toGRPCError(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (_ interface{}, err error) {
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
		err = gerrors.ZITADELToGRPCError(ctx, err)
		cancel(err)
	}()
	return handler(ctx, req)
}
