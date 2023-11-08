package middleware

import (
	"context"
	errs "errors"
	"strings"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	HTTP1Host = "x-zitadel-http1-host"
)

func InstanceInterceptor(verifier authz.InstanceVerifier, explicitInstanceIdServices ...string) grpc.UnaryServerInterceptor {
	translator, err := newZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, translator, explicitInstanceIdServices...)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, idFromRequestsServices ...string) (_ interface{}, err error) {
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
	instance, err := verifier.InstanceByDomain(interceptorCtx, http_utils.RequestOriginFromCtx(ctx).Domain)
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
