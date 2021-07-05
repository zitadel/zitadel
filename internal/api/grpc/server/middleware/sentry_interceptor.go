package middleware

import (
	"context"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/caos/zitadel/internal/api/grpc/errors"
)

func SentryHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return sendErrToSentry(ctx, req, handler)
	}
}

func sendErrToSentry(ctx context.Context, req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	code, _, _, _ := errors.ExtractCaosError(err)
	if code == codes.Unknown || code == codes.Internal {
		sentry.CaptureException(err)
	}
	return resp, err
}
