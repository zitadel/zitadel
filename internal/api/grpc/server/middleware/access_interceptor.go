package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zitadel/logging"

	"google.golang.org/grpc/codes"

	"github.com/zitadel/zitadel/internal/api/authz"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/access"
)

func AccessLimitInterceptor(svc *access.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		instance := authz.GetInstance(ctx)
		limit, err := svc.Limit(ctx, instance.InstanceID())
		if err != nil {
			logging.Warnf("failed to check whether requests should be limited: %s", err.Error())
			err = nil
		}

		resp, err := handler(ctx, req)
		if limit {
			err = status.Error(codes.ResourceExhausted, "quota for authenticated requests exceeded")
		}
		return resp, err
	}
}
func AccessStorageInterceptor(svc *access.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		reqMd, _ := metadata.FromIncomingContext(ctx)

		resp, handlerErr := handler(ctx, req)
		var respStatus uint32
		grpcErr, ok := status.FromError(handlerErr)
		if ok {
			respStatus = uint32(grpcErr.Code())
		}

		resMd, _ := metadata.FromOutgoingContext(ctx)
		instance := authz.GetInstance(ctx)

		// TODO: Why is the instance missing at some paths like /oauth, /.well-known and /ui? Should we fix that for the access logs?
		record := &logstore.AccessLogRecord{
			Timestamp:       time.Now(),
			Protocol:        logstore.GRPC,
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
