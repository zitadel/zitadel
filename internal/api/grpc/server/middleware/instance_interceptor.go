package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName, externalDomain string, explicitInstanceIdServices ...string) grpc.UnaryServerInterceptor {
	translator, err := i18n.NewZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, headerName, externalDomain, translator, explicitInstanceIdServices...)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, headerName, externalDomain string, translator *i18n.Translator, idFromRequestsServices ...string) (_ interface{}, err error) {
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
	return addInstanceByRequestedHost(interceptorCtx, req, handler, verifier, translator, externalDomain, headerName)
}

func addInstanceByID(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, id string) (interface{}, error) {
	ctx = authz.WithInstanceID(ctx, id)
	instance, err := verifier.InstanceByID(ctx)
	if err != nil {
		notFoundErr := new(zerrors.ZitadelError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to set instance using id %s: %w", id, notFoundErr))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByDomain(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, domain string) (interface{}, error) {
	instance, err := verifier.InstanceByHost(ctx, domain)
	if err != nil {
		notFoundErr := new(zerrors.NotFoundError)
		if errors.As(err, &notFoundErr) {
			notFoundErr.Message = translator.LocalizeFromCtx(ctx, notFoundErr.GetMessage(), nil)
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("unable to set instance using domain %s: %w", domain, notFoundErr))
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func addInstanceByRequestedHost(ctx context.Context, req interface{}, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, translator *i18n.Translator, externalDomain, headerName string) (interface{}, error) {
	host, err := hostFromContext(ctx, headerName)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	instance, err := verifier.InstanceByHost(ctx, host)
	if err != nil {
		origin := zitadel_http.ComposedOrigin(ctx)
		logging.WithFields("origin", origin, "externalDomain", externalDomain).WithError(err).Error("unable to set instance")
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
