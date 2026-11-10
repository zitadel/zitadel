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
)

func AuthorizationInterceptor(verifier authz.APITokenVerifier, systemUserPermissions authz.Config, authConfig authz.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authorize(ctx, req, info, handler, verifier, systemUserPermissions, authConfig)
	}
}

func authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.APITokenVerifier, systemUserPermissions authz.Config, authConfig authz.Config) (_ interface{}, err error) {
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

	orgID, orgDomain := orgIDAndDomainFromRequest(authCtx, req)
	ctxSetter, err := authz.CheckUserAuthorization(authCtx, req, authToken, orgID, orgDomain, verifier, systemUserPermissions.RolePermissionMappings, authConfig.RolePermissionMappings, authOpt, info.FullMethod)
	if err != nil {
		return nil, err
	}
	span.End()
	return handler(ctxSetter(ctx), req)
}

func orgIDAndDomainFromRequest(ctx context.Context, req interface{}) (id, domain string) {
	orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID)
	oz, ok := req.(OrganizationFromRequest)
	if ok {
		id = oz.OrganizationFromRequest().ID
		domain = oz.OrganizationFromRequest().Domain
		if id != "" || domain != "" {
			return id, domain
		}
	}
	return orgID, domain
}

type Organization struct {
	ID     string
	Domain string
}

type OrganizationFromRequest interface {
	OrganizationFromRequest() *Organization
}
