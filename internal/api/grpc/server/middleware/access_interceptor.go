package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func AccessStorageInterceptor(svc *logstore.Service[*record.AccessLog]) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}
		reqMd, _ := metadata.FromIncomingContext(ctx)

		resp, handlerErr := handler(ctx, req)

		interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()

		var respStatus uint32
		grpcStatus, ok := status.FromError(handlerErr)
		if ok {
			respStatus = uint32(grpcStatus.Code())
		}

		resMd, _ := metadata.FromOutgoingContext(ctx)
		instance := authz.GetInstance(ctx)
		domainCtx := http_util.DomainContext(ctx)

		r := &record.AccessLog{
			LogDate:         time.Now(),
			Protocol:        record.GRPC,
			RequestURL:      info.FullMethod,
			ResponseStatus:  respStatus,
			RequestHeaders:  reqMd,
			ResponseHeaders: resMd,
			InstanceID:      instance.InstanceID(),
			ProjectID:       instance.ProjectID(),
			RequestedDomain: domainCtx.RequestedDomain(),
			RequestedHost:   domainCtx.RequestedHost(),
		}

		svc.Handle(interceptorCtx, r)
		return resp, handlerErr
	}
}
