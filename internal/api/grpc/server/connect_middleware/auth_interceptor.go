package connect_middleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func AuthorizationInterceptor(verifier authz.APITokenVerifier, systemUserPermissions authz.Config, authConfig authz.Config) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return authorize(ctx, req, handler, verifier, systemUserPermissions, authConfig)
		}
	}
}

func authorize(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc, verifier authz.APITokenVerifier, systemUserPermissions authz.Config, authConfig authz.Config) (_ connect.AnyResponse, err error) {
	authOpt, needsToken := verifier.CheckAuthMethod(req.Spec().Procedure)
	if !needsToken {
		return handler(ctx, req)
	}

	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authToken := req.Header().Get(http.Authorization)
	if authToken == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("auth header missing"))
	}

	orgID, orgDomain := orgIDAndDomainFromRequest(req)
	ctxSetter, err := authz.CheckUserAuthorization(authCtx, req, authToken, orgID, orgDomain, verifier, systemUserPermissions.RolePermissionMappings, authConfig.RolePermissionMappings, authOpt, req.Spec().Procedure)
	if err != nil {
		return nil, err
	}
	span.End()
	return handler(ctxSetter(ctx), req)
}

func orgIDAndDomainFromRequest(req connect.AnyRequest) (id, domain string) {
	orgID := req.Header().Get(http.ZitadelOrgID)
	oz, ok := req.Any().(OrganizationFromRequest)
	if ok {
		id = oz.OrganizationFromRequestConnect().ID
		domain = oz.OrganizationFromRequestConnect().Domain
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
	OrganizationFromRequestConnect() *Organization
}
