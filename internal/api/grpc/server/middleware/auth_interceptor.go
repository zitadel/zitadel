package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	grpc_util "github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
)

func AuthorizationInterceptor(verifier *authz.TokenVerifier, authConfig authz.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authorize(ctx, req, info, handler, verifier, authConfig)
	}
}

func authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier *authz.TokenVerifier, authConfig authz.Config) (_ interface{}, err error) {
	authOpt, needsToken := verifier.CheckAuthMethod(info.FullMethod)
	if !needsToken {
		return handler(ctx, req)
	}

	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authToken := grpc_util.GetAuthorizationHeader(authCtx)
	if authToken == "" {
		return nil, status.Error(codes.Unauthenticated, "auth header missing")
	}

	var orgDomain string
	orgID := grpc_util.GetHeader(authCtx, http.ZitadelOrgID)
	if o, ok := req.(OrganisationFromRequest); ok {
		orgID = o.OrganisationFromRequest().GetOrgId()
		orgDomain = o.OrganisationFromRequest().GetOrgDomain()
	}

	ctxSetter, err := authz.CheckUserAuthorization(authCtx, req, authToken, orgID, orgDomain, verifier, authConfig, authOpt, info.FullMethod)
	if err != nil {
		return nil, err
	}
	span.End()
	return handler(ctxSetter(ctx), req)
}

type OrganisationFromRequest interface {
	OrganisationFromRequest() *object.Organisation
}
