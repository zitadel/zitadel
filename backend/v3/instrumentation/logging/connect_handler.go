package logging

import (
	"context"
	"log/slog"
	"slices"
	"strings"
	"time"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func NewConnectInterceptor(next connect.UnaryFunc, ignoredMethodSuffixes ...string) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
			return strings.HasSuffix(req.Spec().Procedure, s)
		}) {
			return next(ctx, req)
		}
		start := time.Now()
		ctx = NewCtx(ctx, StreamRequest)
		ctx = instrumentation.SetRequestID(ctx, start)

		resp, err := next(ctx, req)
		var code connect.Code
		if err != nil {
			code = connect.CodeOf(err)
		}
		spec := req.Spec()
		Info(ctx, "request served",
			slog.String("protocol", "connect"),
			slog.Any("domain", http_util.DomainContext(ctx)),
			slog.String("service", serviceFromRPCMethod(spec.Procedure)),
			slog.String("http_method", req.HTTPMethod()),
			slog.String("path", spec.Procedure),
			slog.Any("code", code),
			slog.Duration("duration", time.Since(start)),
		)
		return resp, err
	}
}
