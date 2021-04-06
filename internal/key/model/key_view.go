package model

import (
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

type KeyView struct {
	ID        string
	Private   bool
	Expiry    time.Time
	Algorithm string
	Usage     KeyUsage
	Key       *crypto.CryptoValue
	Sequence  uint64
}

type SigningKey struct {
	ID        string
	Algorithm string
	Key       interface{}
	Sequence  uint64
}

type PublicKey struct {
	ID        string
	Algorithm string
	Usage     KeyUsage
	Key       interface{}
}

type KeySearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn KeySearchKey
	Asc           bool
	Queries       []*KeySearchQuery
}

type KeySearchKey int32

const (
	KeySearchKeyUnspecified KeySearchKey = iota
	KeySearchKeyID
	KeySearchKeyPrivate
	KeySearchKeyExpiry
	KeySearchKeyUsage
)

type KeySearchQuery struct {
	Key    KeySearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type KeySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*KeyView
}

func (r *KeySearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return errors.ThrowInvalidArgument(nil, "SEARCH-8fn7f", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func SigningKeyFromKeyView(key *KeyView, alg crypto.EncryptionAlgorithm) (*SigningKey, error) {
	if key.Usage != KeyUsageSigning || !key.Private {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-5HBdh", "key must be private signing key")
	}
	keyData, err := crypto.Decrypt(key.Key, alg)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return nil, err
	}
	return &SigningKey{
		ID:        key.ID,
		Algorithm: key.Algorithm,
		Key:       privateKey,
		Sequence:  key.Sequence,
	}, nil
}

func PublicKeysFromKeyView(keys []*KeyView, alg crypto.EncryptionAlgorithm) ([]*PublicKey, error) {
	convertedKeys := make([]*PublicKey, 0, len(keys))
	for _, key := range keys {
		converted, err := PublicKeyFromKeyView(key, alg)
		if err != nil {
			logging.Log("MODEL-adB3f").WithError(err).Debug("cannot convert to public key") //TODO: change log level to warning when keys can be revoked
			continue
		}
		convertedKeys = append(convertedKeys, converted)
	}
	return convertedKeys, nil

}
func PublicKeyFromKeyView(key *KeyView, alg crypto.EncryptionAlgorithm) (*PublicKey, error) {
	if key.Private {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-dTZa2", "key must be public")
	}
	keyData, err := crypto.Decrypt(key.Key, alg)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.BytesToPublicKey(keyData)
	if err != nil {
		return nil, err
	}
	return &PublicKey{
		ID:        key.ID,
		Algorithm: key.Algorithm,
		Usage:     key.Usage,
		Key:       publicKey,
	}, nil
}
