package connect_middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_v3 "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
)

func InstanceInterceptor(verifier authz.InstanceVerifier, externalDomain string, translator *i18n.Translator, explicitInstanceIdServices ...string) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return setInstance(ctx, req, handler, verifier, externalDomain, translator, explicitInstanceIdServices...)
		}
	}
}

func setInstance(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, verifier authz.InstanceVerifier, externalDomain string, translator *i18n.Translator, idFromRequestsServices ...string) (_ connect.AnyResponse, err error) {
	interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	for _, service := range idFromRequestsServices {
		if !strings.HasPrefix(service, "/") {
			service = "/" + service
		}
		if strings.HasPrefix(req.Spec().Procedure, service) {
			withInstanceIDProperty, ok := req.Any().(interface {
				GetInstanceId() string
			})
			if !ok {
				return handler(ctx, req)
			}
			return addInstanceByID(interceptorCtx, req, handler, verifier, translator, withInstanceIDProperty.GetInstanceId())
		}
	}
	explicitInstanceRequest, ok := req.Any().(interface {
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

func addInstanceByID(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, verifier authz.InstanceVerifier, translator *i18n.Translator, id string) (connect.AnyResponse, error) {
	instance, err := verifier.InstanceByID(ctx, id)
	if err != nil {
		notFoundErr := new(zerrors.ZitadelError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unable to set instance using id %s: %w", id, notFoundErr))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByDomain(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, verifier authz.InstanceVerifier, translator *i18n.Translator, domain string) (connect.AnyResponse, error) {
	instance, err := verifier.InstanceByHost(ctx, domain, "")
	if err != nil {
		notFoundErr := new(zerrors.NotFoundError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unable to set instance using domain %s: %w", domain, notFoundErr))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByRequestedHost(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, verifier authz.InstanceVerifier, translator *i18n.Translator, externalDomain string) (connect.AnyResponse, error) {
	requestContext := zitadel_http.DomainContext(ctx)
	if requestContext.InstanceHost == "" {
		logging.WithFields("origin", requestContext.Origin(), "externalDomain", externalDomain).Error("unable to set instance")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no instanceHost specified"))
	}
	instance, err := verifier.InstanceByHost(ctx, requestContext.InstanceHost, requestContext.PublicHost)
	if err != nil {
		origin := zitadel_http.DomainContext(ctx)
		logging.WithFields("origin", requestContext.Origin(), "externalDomain", externalDomain).WithError(err).Error("unable to set instance")
		zErr := new(zerrors.ZitadelError)
		if errors.As(err, &zErr) {
			zErr.SetMessage(translator.LocalizeFromCtx(ctx, zErr.GetMessage(), nil))
			zErr.Parent = err
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unable to set instance using origin %s (ExternalDomain is %s): %s", origin, externalDomain, zErr.Error()))
		}
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unable to set instance using origin %s (ExternalDomain is %s)", origin, externalDomain))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}
