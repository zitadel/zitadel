package connect_middleware

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func AccessStorageInterceptor(svc *logstore.Service[*record.AccessLog]) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {
			if !svc.Enabled() {
				return handler(ctx, req)
			}
			resp, handlerErr := handler(ctx, req)

			interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
			defer func() { span.EndWithError(err) }()

			var respStatus uint32
			if code := connect.CodeOf(handlerErr); code != connect.CodeUnknown {
				respStatus = uint32(code)
			}

			respHeader := http.Header{}
			if resp != nil {
				respHeader = resp.Header()
			}
			instance := authz.GetInstance(ctx)
			domainCtx := http_util.DomainContext(ctx)

			r := &record.AccessLog{
				LogDate:         time.Now(),
				Protocol:        record.GRPC,
				RequestURL:      req.Spec().Procedure,
				ResponseStatus:  respStatus,
				RequestHeaders:  req.Header(),
				ResponseHeaders: respHeader,
				InstanceID:      instance.InstanceID(),
				ProjectID:       instance.ProjectID(),
				RequestedDomain: domainCtx.RequestedDomain(),
				RequestedHost:   domainCtx.RequestedHost(),
			}

			svc.Handle(interceptorCtx, r)
			return resp, handlerErr
		}
	}
}
