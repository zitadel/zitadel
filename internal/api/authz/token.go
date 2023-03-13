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

func Start(authZRepo authZRepo, issuer string, keys map[string]*SystemAPIUser) (v *TokenVerifier) {
	return &TokenVerifier{
		authZRepo: authZRepo,
		systemJWTProfile: op.NewJWTProfileVerifier(
			&systemJWTStorage{
				keys:       keys,
				cachedKeys: make(map[string]*rsa.PublicKey),
			},
			issuer,
			1*time.Hour,
			time.Second,
		),
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
	userID, agentID, clientID, prefLang, resourceOwner, err = v.authZRepo.VerifyAccessToken(ctx, token, "", GetInstance(ctx).ProjectID())
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
	keys       map[string]*SystemAPIUser
	mutex      sync.Mutex
	cachedKeys map[string]*rsa.PublicKey
}

type SystemAPIUser struct {
	Path    string //if a path is specified, the key will be read from that path
	KeyData []byte //else you can also specify the data directly in the KeyData
}

func (s *SystemAPIUser) readKey() (*rsa.PublicKey, error) {
	if s.Path != "" {
		var err error
		s.KeyData, err = os.ReadFile(s.Path)
		if err != nil {
			return nil, caos_errs.ThrowInternal(err, "AUTHZ-JK31F", "Errors.NotFound")
		}
	}
	return crypto.BytesToPublicKey(s.KeyData)
}

func (s *systemJWTStorage) GetKeyByIDAndClientID(_ context.Context, _, userID string) (*jose.JSONWebKey, error) {
	cachedKey, ok := s.cachedKeys[userID]
	if ok {
		return &jose.JSONWebKey{KeyID: userID, Key: cachedKey}, nil
	}
	key, ok := s.keys[userID]
	if !ok {
		return nil, caos_errs.ThrowNotFound(nil, "AUTHZ-asfd3", "Errors.User.NotFound")
	}
	defer s.mutex.Unlock()
	s.mutex.Lock()
	publicKey, err := key.readKey()
	if err != nil {
		return nil, err
	}
	s.cachedKeys[userID] = publicKey
	return &jose.JSONWebKey{KeyID: userID, Key: publicKey}, nil
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
