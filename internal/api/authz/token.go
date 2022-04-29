package authz

import (
	"context"
	"crypto/rsa"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/zitadel/oidc/v2/pkg/op"
	"gopkg.in/square/go-jose.v2"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	BearerPrefix = "Bearer "
)

type TokenVerifier struct {
	authZRepo        authZRepo
	clients          sync.Map
	authMethods      MethodMapping
	systemJWTProfile op.JWTProfileVerifier
}

type authZRepo interface {
	VerifyAccessToken(ctx context.Context, token, verifierClientID, projectID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error)
	VerifierClientID(ctx context.Context, name string) (clientID, projectID string, err error)
	SearchMyMemberships(ctx context.Context) ([]*Membership, error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	ExistsOrg(ctx context.Context, orgID string) error
}

func Start(authZRepo authZRepo, systemAPI string, keys map[string]string) (v *TokenVerifier) {
	return &TokenVerifier{
		authZRepo:        authZRepo,
		systemJWTProfile: op.NewJWTProfileVerifier(&systemJWTStorage{keys: keys}, systemAPI, 1*time.Hour, time.Second),
	}
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string, method string) (userID, clientID, agentID, prefLang, resourceOwner string, err error) {
	if strings.HasPrefix(method, "/zitadel.system.v1.SystemService") {
		userID, err := v.verifySystemToken(ctx, token)
		if err != nil {
			return "", "", "", "", "", err
		}
		return userID, "", "", "", "", nil
	}
	verifierClientID, projectID, err := v.clientIDAndProjectIDFromMethod(ctx, method)
	if err != nil {
		return "", "", "", "", "", err
	}
	userID, agentID, clientID, prefLang, resourceOwner, err = v.authZRepo.VerifyAccessToken(ctx, token, verifierClientID, projectID)
	return userID, clientID, agentID, prefLang, resourceOwner, err
}

func (v *TokenVerifier) verifySystemToken(ctx context.Context, token string) (string, error) {
	jwtReq, err := op.VerifyJWTAssertion(ctx, token, v.systemJWTProfile)
	if err != nil {
		return "", err
	}
	return jwtReq.Subject, nil
}

type systemJWTStorage struct {
	keys       map[string]string
	mutex      sync.Mutex
	cachedKeys map[string]*rsa.PublicKey
}

type SystemAPIUser struct {
	Key   *rsa.PublicKey
	Roles string
}

func (s *systemJWTStorage) GetKeyByIDAndUserID(ctx context.Context, keyID, _ string) (*jose.JSONWebKey, error) {
	key, ok := s.cachedKeys[keyID]
	if ok {
		return &jose.JSONWebKey{KeyID: keyID, Key: key}, nil
	}
	path, ok := s.keys[keyID]
	if !ok {

	}
	defer s.mutex.Unlock()
	s.mutex.Lock()
	keyData, err := os.ReadFile(path)
	if err != nil {

	}
	publicKey, err := crypto.BytesToPublicKey(keyData)
	if err != nil {

	}
	s.cachedKeys[keyID] = publicKey
	return &jose.JSONWebKey{KeyID: keyID, Key: publicKey}, nil
}

type client struct {
	id        string
	projectID string
	name      string
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

func (v *TokenVerifier) clientIDAndProjectIDFromMethod(ctx context.Context, method string) (clientID, projectID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	prefix, ok := prefixFromMethod(method)
	if !ok {
		return "", "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-GRD2Q", "Errors.Internal")
	}
	app, ok := v.clients.Load(prefix)
	if !ok {
		return "", "", caos_errs.ThrowPermissionDenied(nil, "AUTHZ-G2qrh", "Errors.Internal")
	}
	c := app.(*client)
	if c.name == "System-API" {
		return c.id, c.projectID, nil
	}
	c.id, c.projectID, err = v.authZRepo.VerifierClientID(ctx, c.name)
	if err != nil {
		return "", "", caos_errs.ThrowPermissionDenied(err, "AUTHZ-ptTIF2", "Errors.Internal")
	}
	return c.id, c.projectID, nil
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
