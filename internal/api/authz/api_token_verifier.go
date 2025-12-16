package authz

import (
	"context"
	"sync"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// TODO: Define interfaces where they are accepted
type APITokenVerifier interface {
	AccessTokenVerifier
	SystemTokenVerifier
	RegisterServer(appName, methodPrefix string, mappings MethodMapping)
	CheckAuthMethod(method string) (Option, bool)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (_ string, _ []string, err error)
	ExistsOrg(ctx context.Context, id, domain string) (orgID string, err error)
	SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) (_ []*Membership, err error)
}

type ApiTokenVerifier struct {
	AccessTokenVerifier
	SystemTokenVerifier
	authZRepo   authZRepo
	clients     sync.Map
	authMethods MethodMapping
}

func StartAPITokenVerifier(authZRepo authZRepo, accessTokenVerifier AccessTokenVerifier, systemTokenVerifier SystemTokenVerifier) *ApiTokenVerifier {
	return &ApiTokenVerifier{
		authZRepo:           authZRepo,
		SystemTokenVerifier: systemTokenVerifier,
		AccessTokenVerifier: accessTokenVerifier,
	}
}

func (v *ApiTokenVerifier) RegisterServer(appName, methodPrefix string, mappings MethodMapping) {
	v.clients.Store(methodPrefix, &client{name: appName})
	if v.authMethods == nil {
		v.authMethods = make(map[string]Option)
	}
	for method, option := range mappings {
		v.authMethods[method] = option
	}
}

func (v *ApiTokenVerifier) SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) (_ []*Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.SearchMyMemberships(ctx, orgID, shouldTriggerBulk)
}

func (v *ApiTokenVerifier) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (_ string, _ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return v.authZRepo.ProjectIDAndOriginsByClientID(ctx, clientID)
}

func (v *ApiTokenVerifier) ExistsOrg(ctx context.Context, id, domain string) (orgID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.ExistsOrg(ctx, id, domain)
}

func (v *ApiTokenVerifier) CheckAuthMethod(method string) (Option, bool) {
	authOpt, ok := v.authMethods[method]
	return authOpt, ok
}
