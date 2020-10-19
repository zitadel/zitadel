package model

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
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
	Method model.SearchMethod
	Value  interface{}
}

type KeySearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*KeyView
}

func (r *KeySearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
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
	converted := make([]*PublicKey, len(keys))
	var err error
	for i, key := range keys {
		converted[i], err = PublicKeyFromKeyView(key, alg)
		if err != nil {
			return nil, err
		}
	}
	return converted, nil

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
