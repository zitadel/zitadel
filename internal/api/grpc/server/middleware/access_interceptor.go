package middleware

import (
	"context"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
)

func AccessStorageInterceptor(svc *logstore.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}

		interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()

		reqMd, _ := metadata.FromIncomingContext(ctx)

		resp, handlerErr := handler(ctx, req)
		var respStatus uint32
		grpcStatus, ok := status.FromError(handlerErr)
		if ok {
			respStatus = uint32(grpcStatus.Code())
		}

		resMd, _ := metadata.FromOutgoingContext(ctx)
		instance := authz.GetInstance(ctx)

		record := &access.Record{
			LogDate:         time.Now(),
			Protocol:        access.GRPC,
			RequestURL:      info.FullMethod,
			ResponseStatus:  respStatus,
			RequestHeaders:  http.Header(reqMd),
			ResponseHeaders: http.Header(resMd),
			InstanceID:      instance.InstanceID(),
			ProjectID:       instance.ProjectID(),
			RequestedDomain: instance.RequestedDomain(),
			RequestedHost:   instance.RequestedHost(),
		}

		svc.Handle(interceptorCtx, record)
		span.End() // TODO: Why?
		return resp, handlerErr
	}
}
