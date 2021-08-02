package authz

import (
	"context"
	"strings"
	"sync"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

const (
	BearerPrefix = "Bearer "
)

type TokenVerifier struct {
	authZRepo   authZRepo
	clients     sync.Map
	authMethods MethodMapping
}

type authZRepo interface {
	VerifyAccessToken(ctx context.Context, token, verifierClientID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error)
	VerifierClientID(ctx context.Context, name string) (clientID string, err error)
	SearchMyMemberships(ctx context.Context) ([]*Membership, error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	ExistsOrg(ctx context.Context, orgID string) error
	CheckOrgFeatures(ctx context.Context, orgID string, requiredFeatures ...string) error
}

func Start(authZRepo authZRepo) (v *TokenVerifier) {
	return &TokenVerifier{authZRepo: authZRepo}
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string, method string) (userID, clientID, agentID, prefLang, resourceOwner string, err error) {
	verifierClientID, err := v.clientIDFromMethod(ctx, method)
	if err != nil {
		return "", "", "", "", "", err
	}
	userID, agentID, clientID, prefLang, resourceOwner, err = v.authZRepo.VerifyAccessToken(ctx, token, verifierClientID)
	return userID, clientID, agentID, prefLang, resourceOwner, err
}

type client struct {
	id   string
	name string
}

func (v *TokenVerifier) RegisterServer(appName, methodPrefix string, mappings MethodMapping) {
	v.clients.Store(methodPrefix, &client{name: appName})
	if v.authMethods == nil {
		v.authMethods = make(map[string]Option)
	}
	for method, option := range mappings {
		v.authMethods[method] = option
	}
}

func prefixFromMethod(method string) (string, bool) {
	parts := strings.Split(method, "/")
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}

func (v *TokenVerifier) clientIDFromMethod(ctx context.Context, method string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	prefix, ok := prefixFromMethod(method)
	if !ok {
		return "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-GRD2Q", "Errors.Internal")
	}
	app, ok := v.clients.Load(prefix)
	if !ok {
		return "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-G2qrh", "Errors.Internal")
	}
	c := app.(*client)
	if c.id != "" {
		return c.id, nil
	}
	c.id, err = v.authZRepo.VerifierClientID(ctx, c.name)
	if err != nil {
		return "", caos_errs.ThrowPermissionDenied(err, "AUTHZ-ptTIF2", "Errors.Internal")
	}
	v.clients.Store(prefix, c)
	return c.id, nil
}
func (v *TokenVerifier) SearchMyMemberships(ctx context.Context) (_ []*Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.SearchMyMemberships(ctx)
}

func (v *TokenVerifier) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (_ string, _ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return v.authZRepo.ProjectIDAndOriginsByClientID(ctx, clientID)
}

func (v *TokenVerifier) ExistsOrg(ctx context.Context, orgID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.ExistsOrg(ctx, orgID)
}

func (v *TokenVerifier) CheckAuthMethod(method string) (Option, bool) {
	authOpt, ok := v.authMethods[method]
	return authOpt, ok
}

func verifyAccessToken(ctx context.Context, token string, t *TokenVerifier, method string) (userID, clientID, agentID, prefLan, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", "", "", "", "", caos_errs.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return t.VerifyAccessToken(ctx, parts[1], method)
}
