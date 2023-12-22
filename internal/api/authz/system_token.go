package authz

import (
	"context"
	"crypto/rsa"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ SystemTokenVerifier = (*SystemTokenVerifierFromConfig)(nil)

type SystemTokenVerifierFromConfig struct {
	systemJWTProfile *op.JWTProfileVerifier
	systemUsers      map[string]Memberships
}

func StartSystemTokenVerifierFromConfig(issuer string, keys map[string]*SystemAPIUser) (*SystemTokenVerifierFromConfig, error) {
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
				return nil, errors.New("for system users, only the membership types System, IAM and Organization are supported")
			default:
				return nil, errors.New("unknown membership type")
			}
		}
	}
	return &SystemTokenVerifierFromConfig{
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

func (s *SystemTokenVerifierFromConfig) VerifySystemToken(ctx context.Context, token string, orgID string) (matchingMemberships Memberships, userID string, err error) {
	jwtReq, err := op.VerifyJWTAssertion(ctx, token, s.systemJWTProfile)
	if err != nil {
		return nil, "", err
	}
	systemUserMemberships, ok := s.systemUsers[jwtReq.Subject]
	if !ok {
		return nil, "", zerrors.ThrowPermissionDenied(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong")
	}
	matchingMemberships = make(Memberships, 0, len(systemUserMemberships))
	for _, membership := range systemUserMemberships {
		if membership.MemberType == MemberTypeSystem ||
			membership.MemberType == MemberTypeIAM && GetInstance(ctx).InstanceID() == membership.AggregateID ||
			membership.MemberType == MemberTypeOrganization && orgID == membership.AggregateID {
			matchingMemberships = append(matchingMemberships, membership)
		}
	}
	return matchingMemberships, jwtReq.Subject, nil
}

type systemJWTStorage struct {
	keys       map[string]*SystemAPIUser
	mutex      sync.Mutex
	cachedKeys map[string]*rsa.PublicKey
}

type SystemAPIUser struct {
	Path        string // if a path is specified, the key will be read from that path
	KeyData     []byte // else you can also specify the data directly in the KeyData
	Memberships Memberships
}

func (s *SystemAPIUser) readKey() (*rsa.PublicKey, error) {
	if s.Path != "" {
		var err error
		s.KeyData, err = os.ReadFile(s.Path)
		if err != nil {
			return nil, zerrors.ThrowInternal(err, "AUTHZ-JK31F", "Errors.NotFound")
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
		return nil, zerrors.ThrowNotFound(nil, "AUTHZ-asfd3", "Errors.User.NotFound")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	publicKey, err := key.readKey()
	if err != nil {
		return nil, err
	}
	s.cachedKeys[userID] = publicKey
	return &jose.JSONWebKey{KeyID: userID, Key: publicKey}, nil
}
