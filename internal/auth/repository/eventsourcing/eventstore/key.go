package eventstore

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type KeyRepository struct {
	KeyEvents *key_event.KeyEventstore
}

func (k *KeyRepository) SaveKeyPair(ctx context.Context) error {
	key, err := a.createKeyPair()
	if err != nil {
		return err
	}
	_, err = k.repo.CreateKeyPair(ctx, key)
	return err
}

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, renewTimer <-chan time.Time) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				a.refreshSigningKey(ctx, keyCh, errCh)
				renewTimer = time.After(a.signingKeyRotation)
			}
		}
	}()
}

func (k *KeyRepository) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {

}
