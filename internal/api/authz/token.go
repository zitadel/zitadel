package authz

import (
	"context"
	"strings"
	"sync"

	caos_errs "github.com/caos/zitadel/internal/errors"
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
	VerifyAccessToken(ctx context.Context, token, clientID string) (string, string, error)
	VerifierClientID(ctx context.Context, name string) (string, error)
	ResolveGrants(ctx context.Context) (*Grant, error)
	ProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}

func Start(authZRepo authZRepo) (v *TokenVerifier) {
	return &TokenVerifier{authZRepo: authZRepo}
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string, method string) (string, string, string, error) {
	clientID, err := v.clientIDFromMethod(ctx, method)
	if err != nil {
		return "", "", "", err
	}
	userID, agentID, err := v.authZRepo.VerifyAccessToken(ctx, token, clientID)
	return userID, clientID, agentID, err
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

func (v *TokenVerifier) clientIDFromMethod(ctx context.Context, method string) (string, error) {
	prefix, ok := prefixFromMethod(method)
	if !ok {
		return "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-GRD2Q", "Errors.Internal")
	}
	app, ok := v.clients.Load(prefix)
	if !ok {
		return "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-G2qrh", "Errors.Internal")
	}
	var err error
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

func (v *TokenVerifier) ResolveGrant(ctx context.Context) (*Grant, error) {
	return v.authZRepo.ResolveGrants(ctx)
}

func (v *TokenVerifier) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return v.authZRepo.ProjectIDByClientID(ctx, clientID)
}

func (v *TokenVerifier) CheckAuthMethod(method string) (Option, bool) {
	authOpt, ok := v.authMethods[method]
	return authOpt, ok
}

func verifyAccessToken(ctx context.Context, token string, t *TokenVerifier, method string) (string, string, string, error) {
	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", "", "", caos_errs.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return t.VerifyAccessToken(ctx, parts[1], method)
}
