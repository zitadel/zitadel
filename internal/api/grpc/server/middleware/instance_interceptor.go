package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_v3 "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
)

const (
	HTTP1Host = "x-zitadel-http1-host"
)

func InstanceInterceptor(verifier authz.InstanceVerifier, externalDomain string, translator *i18n.Translator, explicitInstanceIdServices ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, externalDomain, translator, explicitInstanceIdServices...)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, externalDomain string, translator *i18n.Translator, idFromRequestsServices ...string) (_ interface{}, err error) {
	interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	for _, service := range idFromRequestsServices {
		if !strings.HasPrefix(service, "/") {
			service = "/" + service
		}
		if strings.HasPrefix(info.FullMethod, service) {
			withInstanceIDProperty, ok := req.(interface {
				GetInstanceId() string
			})
			if !ok {
				return handler(ctx, req)
			}
			return addInstanceByID(interceptorCtx, req, handler, verifier, translator, withInstanceIDProperty.GetInstanceId())
		}
	}
	explicitInstanceRequest, ok := req.(interface {
		GetInstance() *object_v3.Instance
	})
	if ok {
		instance := explicitInstanceRequest.GetInstance()
		if id := instance.GetId(); id != "" {
			return addInstanceByID(interceptorCtx, req, handler, verifier, translator, id)
		}
		if domain := instance.GetDomain(); domain != "" {
			return addInstanceByDomain(interceptorCtx, req, handler, verifier, translator, domain)
		}
	}
	return addInstanceByRequestedHost(interceptorCtx, req, handler, verifier, translator, externalDomain)
}

func addInstanceByID(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, id string) (interface{}, error) {
	instance, err := verifier.InstanceByID(ctx, id)
	if err != nil {
		notFoundErr := new(zerrors.ZitadelError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, status.Error(codes.NotFound, fmt.Errorf("unable to set instance using id %s: %w", id, notFoundErr).Error())
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByDomain(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, domain string) (interface{}, error) {
	instance, err := verifier.InstanceByHost(ctx, domain, "")
	if err != nil {
		notFoundErr := new(zerrors.NotFoundError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, status.Error(codes.NotFound, fmt.Errorf("unable to set instance using domain %s: %w", domain, notFoundErr).Error())
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByRequestedHost(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, externalDomain string) (interface{}, error) {
	requestContext := zitadel_http.DomainContext(ctx)
	if requestContext.InstanceHost == "" {
		logging.WithFields("origin", requestContext.Origin(), "externalDomain", externalDomain).Error("unable to set instance")
		return nil, status.Error(codes.NotFound, "no instanceHost specified")
	}
	instance, err := verifier.InstanceByHost(ctx, requestContext.InstanceHost, requestContext.PublicHost)
	if err != nil {
		origin := zitadel_http.DomainContext(ctx)
		logging.WithFields("origin", requestContext.Origin(), "externalDomain", externalDomain).WithError(err).Error("unable to set instance")
		zErr := new(zerrors.ZitadelError)
		if errors.As(err, &zErr) {
			zErr.SetMessage(translator.LocalizeFromCtx(ctx, zErr.GetMessage(), nil))
			zErr.Parent = err
			return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s): %s", origin, externalDomain, zErr))
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s)", origin, externalDomain))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}
