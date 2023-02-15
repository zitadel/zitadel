package middleware

import (
	"context"
	errs "errors"
	"fmt"
	"strings"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	HTTP1Host = "x-zitadel-http1-host"
)

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName string, explicitInstanceIdServices ...string) grpc.UnaryServerInterceptor {
	translator, err := newZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, headerName, translator, explicitInstanceIdServices...)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, headerName string, translator *i18n.Translator, idFromRequestsServices ...string) (_ interface{}, err error) {
	interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()
	for _, service := range idFromRequestsServices {
		if !strings.HasPrefix(service, "/") {
			service = "/" + service
		}
		if strings.HasPrefix(info.FullMethod, service) {
			withInstanceIDProperty, ok := req.(interface{ GetInstanceId() string })
			if !ok {
				return handler(ctx, req)
			}
			ctx = authz.WithInstanceID(ctx, withInstanceIDProperty.GetInstanceId())
			instance, err := verifier.InstanceByID(ctx)
			if err != nil {
				notFoundErr := new(errors.NotFoundError)
				if errs.As(err, &notFoundErr) {
					notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
				}
				return nil, status.Error(codes.NotFound, err.Error())
			}
			return handler(authz.WithInstance(ctx, instance), req)
		}
	}

	host, err := hostFromContext(interceptorCtx, headerName)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	instance, err := verifier.InstanceByHost(interceptorCtx, host)
	if err != nil {
		notFoundErr := new(errors.NotFoundError)
		if errs.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, status.Error(codes.NotFound, err.Error())
	}
	span.End()
	return handler(authz.WithInstance(ctx, instance), req)
}

func hostFromContext(ctx context.Context, headerName string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("cannot read metadata")
	}
	host, ok := md[HTTP1Host]
	if ok && len(host) == 1 {
		if !isAllowedToSendHTTP1Header(md) {
			return "", fmt.Errorf("no valid host header")
		}
		return host[0], nil
	}
	host, ok = md[headerName]
	if !ok {
		return "", fmt.Errorf("cannot find header: %v", headerName)
	}
	if len(host) != 1 {
		return "", fmt.Errorf("invalid host header: %v", host)
	}
	return host[0], nil
}

// isAllowedToSendHTTP1Header check if the gRPC call was sent to `localhost`
// this is only possible when calling the server directly running on localhost
// or through the gRPC gateway
func isAllowedToSendHTTP1Header(md metadata.MD) bool {
	authority, ok := md[":authority"]
	return ok && len(authority) == 1 && strings.Split(authority[0], ":")[0] == "localhost"
}
