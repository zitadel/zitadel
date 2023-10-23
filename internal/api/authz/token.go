package authz

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/crypto"
	zitadel_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	BearerPrefix       = "Bearer "
	SessionTokenPrefix = "sess_"
	SessionTokenFormat = SessionTokenPrefix + "%s:%s"
)

type TokenVerifier struct {
	authZRepo        authZRepo
	clients          sync.Map
	authMethods      MethodMapping
	systemJWTProfile *op.JWTProfileVerifier
	systemUsers      map[string]Memberships
}

type MembershipsResolver interface {
	SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error)
}

type authZRepo interface {
	MembershipsResolver
	VerifyAccessToken(ctx context.Context, token, verifierClientID, projectID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error)
	VerifierClientID(ctx context.Context, name string) (clientID, projectID string, err error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	ExistsOrg(ctx context.Context, id, domain string) (string, error)
}

func Start(authZRepo authZRepo, issuer string, keys map[string]*SystemAPIUser) (*TokenVerifier, error) {
	systemUsers := make(map[string]Memberships, len(keys))
	for userID, key := range keys {
		if len(key.Memberships) == 0 {
			systemUsers[userID] = Memberships{{MemberType: MemberTypeSystem, Roles: []string{"SYSTEM_OWNER"}}}
			continue
		}
		for _, membership := range key.Memberships {
			switch membership.MemberType {
			case MemberTypeSystem, MemberTypeIAM, MemberTypeOrganization:
				systemUsers[userID] = key.Memberships
			case MemberTypeUnspecified, MemberTypeProject, MemberTypeProjectGrant:
				return nil, errors.New("for system users, only the membership types Organization, IAM and System are supported")
			default:
				return nil, errors.New("unknown membership type")
			}
		}
	}
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
		systemUsers: systemUsers,
	}, nil
}

func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (userID, clientID, agentID, prefLang, resourceOwner string, isSystemUser bool, err error) {
	userID, agentID, clientID, prefLang, resourceOwner, err = v.authZRepo.VerifyAccessToken(ctx, token, "", GetInstance(ctx).ProjectID())
	if err == nil || !zitadel_errors.IsUnauthenticated(err) {
		return userID, clientID, agentID, prefLang, resourceOwner, false, err
	}
	userID, sysTokenErr := v.verifySystemToken(ctx, token)
	if sysTokenErr == nil {
		isSystemUser = true
		err = nil
	}
	return userID, "", "", "", "", isSystemUser, err
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
	Path        string //if a path is specified, the key will be read from that path
	KeyData     []byte //else you can also specify the data directly in the KeyData
	Memberships Memberships
}

func (s *SystemAPIUser) readKey() (*rsa.PublicKey, error) {
	if s.Path != "" {
		var err error
		s.KeyData, err = os.ReadFile(s.Path)
		if err != nil {
			return nil, zitadel_errors.ThrowInternal(err, "AUTHZ-JK31F", "Errors.NotFound")
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
		return nil, zitadel_errors.ThrowNotFound(nil, "AUTHZ-asfd3", "Errors.User.NotFound")
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

func (v *TokenVerifier) SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) (_ []*Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.SearchMyMemberships(ctx, orgID, shouldTriggerBulk)
}

func (v *TokenVerifier) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (_ string, _ []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return v.authZRepo.ProjectIDAndOriginsByClientID(ctx, clientID)
}

func (v *TokenVerifier) ExistsOrg(ctx context.Context, id, domain string) (orgID string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return v.authZRepo.ExistsOrg(ctx, id, domain)
}

func (v *TokenVerifier) CheckAuthMethod(method string) (Option, bool) {
	authOpt, ok := v.authMethods[method]
	return authOpt, ok
}

func verifyAccessToken(ctx context.Context, token string, t *TokenVerifier) (userID, clientID, agentID, prefLan, resourceOwner string, isSystemUser bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", "", "", "", "", false, zitadel_errors.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return t.VerifyAccessToken(ctx, parts[1])
}

func SessionTokenVerifier(algorithm crypto.EncryptionAlgorithm) func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		decodedToken, err := base64.RawURLEncoding.DecodeString(sessionToken)
		if err != nil {
			return err
		}
		_, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		var token string
		token, err = algorithm.DecryptString(decodedToken, algorithm.EncryptionKeyID())
		spanPasswordComparison.EndWithError(err)
		if err != nil || token != fmt.Sprintf(SessionTokenFormat, sessionID, tokenID) {
			return zitadel_errors.ThrowPermissionDenied(err, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
		}
		return nil
	}
}
