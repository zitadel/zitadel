package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
)

func AccessLimitInterceptor(svc *logstore.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}
		instance := authz.GetInstance(ctx)
		remaining, err := svc.Limit(ctx, instance.InstanceID())
		if err != nil {
			logging.Warnf("failed to check whether requests should be limited: %s", err.Error())
			err = nil
		}

		resp, err := handler(ctx, req)
		if remaining != nil && *remaining == 0 {
			err = errors.ThrowResourceExhausted(nil, "QUOTA-vjAy8", "Quota.Access.Exhausted")
		}
		return resp, err
	}
}
func AccessStorageInterceptor(svc *logstore.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}

		reqMd, _ := metadata.FromIncomingContext(ctx)

		resp, handlerErr := handler(ctx, req)
		var respStatus uint32
		grpcErr, ok := status.FromError(handlerErr)
		if ok {
			respStatus = uint32(grpcErr.Code())
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

		if err := svc.Handle(ctx, record); err != nil {
			logging.Warnf("failed to handle access log: %s", err.Error())
			err = nil
		}

		return resp, handlerErr
	}
}
