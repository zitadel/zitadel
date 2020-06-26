package authz

import (
	"context"
	"strings"
	"sync"

	"github.com/caos/logging"

	//authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

type key int

const (
	permissionsKey key = 1
	dataKey        key = 2
)

type CtxData struct {
	UserID            string
	OrgID             string
	ProjectID         string
	AgentID           string
	PreferredLanguage string
}

func (ctxData CtxData) IsZero() bool {
	return ctxData.UserID == "" || ctxData.OrgID == ""
}

type Grants []*Grant

type Grant struct {
	OrgID string
	Roles []string
}

type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (string, string, string, error)
	ResolveGrant(ctx context.Context) (*Grant, error)
	GetProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}

type TokenVerificationSupplier interface {
	GetProjectIDByClientID(ctx context.Context, clientID string) (string, error)
}

type TokenVerifier2 struct {
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

func Start(authZRepo authZRepo) (v *TokenVerifier2) {
	return &TokenVerifier2{authZRepo: authZRepo}
}

func (v *TokenVerifier2) VerifyAccessToken(ctx context.Context, token string, method string) (string, string, string, error) {
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

func (v *TokenVerifier2) RegisterServer(appName, methodPrefix string, mappings MethodMapping) {
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

func (v *TokenVerifier2) clientIDFromMethod(ctx context.Context, method string) (string, error) {
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

func (v *TokenVerifier2) ResolveGrant(ctx context.Context) (*Grant, error) {
	return v.authZRepo.ResolveGrants(ctx)
}

func (v *TokenVerifier2) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return v.authZRepo.ProjectIDByClientID(ctx, clientID)
}

func (v *TokenVerifier2) CheckAuthMethod(method string) (Option, bool) {
	authOpt, ok := v.authMethods[method]
	return authOpt, ok
}

func VerifyTokenAndWriteCtxData(ctx context.Context, token, orgID string, t *TokenVerifier2, method string) (_ context.Context, err error) {
	userID, clientID, agentID, err := verifyAccessToken(ctx, token, t, method)
	if err != nil {
		return nil, err
	}
	projectID, err := t.GetProjectIDByClientID(ctx, clientID)
	logging.LogWithFields("AUTH-GfAoV", "clientID", clientID).OnError(err).Warn("could not read projectid by clientid")
	return context.WithValue(ctx, dataKey, CtxData{UserID: userID, OrgID: orgID, ProjectID: projectID, AgentID: agentID}), nil
}

func SetCtxData(ctx context.Context, ctxData CtxData) context.Context {
	return context.WithValue(ctx, dataKey, ctxData)
}

func GetCtxData(ctx context.Context) CtxData {
	ctxData, _ := ctx.Value(dataKey).(CtxData)
	return ctxData
}

func GetPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(permissionsKey).([]string)
	return ctxPermission
}
