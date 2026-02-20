package authz

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
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
				cachedKeys: make(map[string]*SystemAPIPublicKey),
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
	mutex      sync.RWMutex
	cachedKeys map[string]*SystemAPIPublicKey
}

type SystemAPIUser struct {
	Path        string // if a path is specified, the key/cert will be read from that path
	KeyData     []byte // else you can also specify the data directly in the KeyData
	Memberships Memberships
	NotBefore   *time.Time
	NotAfter    *time.Time
}

type SystemAPIPublicKey struct {
	Data      *rsa.PublicKey
	NotBefore *time.Time
	NotAfter  *time.Time
}

func (s *SystemAPIUser) readKey() (*SystemAPIPublicKey, error) {
	if s.Path != "" {
		var err error
		s.KeyData, err = os.ReadFile(s.Path)
		if err != nil {
			return nil, zerrors.ThrowInternal(err, "AUTHZ-JK31F", "Errors.NotFound")
		}
	}

	// when an RSA key is provided, use the raw data
	key, err := crypto.BytesToPublicKey(s.KeyData)
	if err == nil {
		return &SystemAPIPublicKey{Data: key}, nil
	}

	// when x.509 cert is provided, parse it and extract RSA key
	block, _ := pem.Decode(s.KeyData)
	if block == nil {
		return nil, zerrors.ThrowInternal(err, "AUTHZ-FC8ohc", "Errors.SystemApiUser.CertDecodeFailed")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "AUTHZ-64nMHP", "Errors.SystemApiUser.CertParseFailed")
	}
	key = cert.PublicKey.(*rsa.PublicKey)
	return &SystemAPIPublicKey{
		Data:      key,
		NotBefore: &cert.NotBefore,
		NotAfter:  &cert.NotAfter,
	}, nil
}

func (s *systemJWTStorage) GetKeyByIDAndClientID(_ context.Context, _, userID string) (*jose.JSONWebKey, error) {
	s.mutex.RLock()
	key, ok := s.cachedKeys[userID]
	s.mutex.RUnlock()

	var err error
	if !ok {
		s.mutex.Lock()
		user, ok := s.keys[userID]
		if !ok {
			return nil, zerrors.ThrowNotFound(nil, "AUTHZ-asfd3", "Errors.User.NotFound")
		}
		key, err = user.readKey()
		if err != nil {
			return nil, err
		}
		s.cachedKeys[userID] = key
		s.mutex.Unlock()
	}

	now := time.Now().UTC()
	if key.NotBefore != nil && now.Before(*key.NotBefore) {
		return nil, zerrors.ThrowNotFound(nil, "AUTHZ-NiJstf", "Errors.User.NotBefore")
	}
	if key.NotAfter != nil && now.After(*key.NotAfter) {
		return nil, zerrors.ThrowNotFound(nil, "AUTHZ-CGmV4b", "Errors.User.NotBefore")
	}

	return &jose.JSONWebKey{KeyID: userID, Key: key}, nil
}
