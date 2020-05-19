package eventstore

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/key/model"
	key_event "github.com/caos/zitadel/internal/key/repository/eventsourcing"
)

type KeyRepository struct {
	KeyEvents          *key_event.KeyEventstore
	View               *view.View
	signingKeyRotation time.Duration
}

func (k *KeyRepository) GenerateSigningKeyPair(ctx context.Context) error {
	_, err := k.KeyEvents.GenerateKeyPair(ctx, model.KeyUsageSigning)
	return err
}

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, renewTimer <-chan time.Time) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				k.refreshSigningKey(keyCh, errCh)
				renewTimer = time.After(k.signingKeyRotation)
			}
		}
	}()
}

func (k *KeyRepository) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {
	keys, err := k.View.GetActiveKeySet()
	if err != nil {
		return nil, err
	}
	webKeys := make([]jose.JSONWebKey, len(keys))
	for i, key := range keys {
		webKeys[i] = jose.JSONWebKey{KeyID: key.ID, Algorithm: key.Algorithm, Use: key.Usage.String(), Key: key.Key}
	}
	return &jose.JSONWebKeySet{Keys: webKeys}, nil
}

func (k *KeyRepository) refreshSigningKey(keyCh chan<- jose.SigningKey, errCh chan<- error) {
	key, err := k.View.GetSigningKey()
	if err != nil {
		errCh <- err
		return
	}
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(key.Algorithm),
		Key: jose.JSONWebKey{
			KeyID: key.ID,
			Key:   key.Key,
		},
	}
}
