package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/key/model"
	key_event "github.com/caos/zitadel/internal/key/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	oidcUser = "OIDC"
	iamOrg   = "IAM"
)

type KeyRepository struct {
	KeyEvents          *key_event.KeyEventstore
	View               *view.View
	SigningKeyRotation time.Duration
	KeyAlgorithm       crypto.EncryptionAlgorithm
}

func (k *KeyRepository) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	_, err := k.KeyEvents.GenerateKeyPair(ctx, model.KeyUsageSigning, algorithm)
	return err
}

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, errCh chan<- error, renewTimer <-chan time.Time) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				send, err := k.refreshSigningKey(ctx, keyCh)
				if send {
					errCh <- err
				}
				d := k.SigningKeyRotation
				if err != nil {
					d = d / 2
				}
				renewTimer = time.After(d)
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

func (k *KeyRepository) refreshSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey) (sendErr bool, err error) {
	key, errView := k.View.GetSigningKey()
	if errView != nil && !errors.IsNotFound(errView) {
		logging.Log("EVENT-GEd4h").WithError(errView).Warn("could not get signing key")
		return false, errView
	}
	var sequence *repository.CurrentSequence
	if key == nil {
		key = new(model.SigningKey)
		sequence, err = k.View.GetLatestKeySequence()
		if err != nil {
			return true, err
		}
		key.Sequence = sequence.CurrentSequence
	}
	events, err := k.KeyEvents.LatestKeyEvents(ctx, key.Sequence)
	if err != nil {
		logging.Log("EVENT-der5g").Warn("error retrieving new events")
		return true, err
	}
	var newKey *view_model.KeyView
	for _, event := range events {
		privateKey, publicKey, err := view_model.KeysFromPairEvent(event)
		if err != nil {
			return true, err
		}
		if privateKey.Expiry.Before(time.Now()) && publicKey.Expiry.Before(time.Now()) {
			continue
		}
		newKey = privateKey
	}
	if errView != nil && newKey == nil {
		keyCh <- jose.SigningKey{}
		return true, errView
	}
	if newKey != nil {
		key, err = model.SigningKeyFromKeyView(view_model.KeyViewToModel(newKey), k.KeyAlgorithm)
		if err != nil {
			return true, err
		}
	}
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(key.Algorithm),
		Key: jose.JSONWebKey{
			KeyID: key.ID,
			Key:   key.Key,
		},
	}
	return true, nil
}

func setOIDCCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: iamOrg})
}
