package eventstore

import (
	"context"
	"os"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/key/model"
	key_event "github.com/caos/zitadel/internal/key/repository/eventsourcing"
)

const (
	oidcUser = "OIDC"
	iamOrg   = "IAM"

	signingKey = "signing_key"
)

type KeyRepository struct {
	KeyEvents            *key_event.KeyEventstore
	View                 *view.View
	SigningKeyRotation   time.Duration
	KeyAlgorithm         crypto.EncryptionAlgorithm
	Locker               spooler.Locker
	lockID               string
	currentKeyID         string
	currentKeyExpiration time.Time
}

func (k *KeyRepository) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	_, err := k.KeyEvents.GenerateKeyPair(ctx, model.KeyUsageSigning, algorithm)
	return err
}

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string) {
	renewTimer := time.After(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				shortRefresh, _ := k.refreshSigningKey(ctx, keyCh, algorithm)
				d := k.SigningKeyRotation
				if shortRefresh {
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

func (k *KeyRepository) refreshSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string) (shortRefresh bool, err error) {
	key, expiration, errView := k.View.GetSigningKey(time.Now().UTC().Add(time.Duration(2) * k.SigningKeyRotation))
	if errView != nil && !errors.IsNotFound(errView) {
		logging.Log("EVENT-GEd4h").WithError(errView).Warn("could not get signing key")
		return true, errView
	}
	if key != nil {
		if k.currentKeyID == key.ID {
			return false, nil
		}
		k.currentKeyID = key.ID
		k.currentKeyExpiration = expiration
		keyCh <- jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(key.Algorithm),
			Key: jose.JSONWebKey{
				KeyID: key.ID,
				Key:   key.Key,
			},
		}
		return false, nil
	}
	if k.currentKeyExpiration.Before(time.Now().UTC()) {
		keyCh <- jose.SigningKey{}
	}
	sequence, err := k.View.GetLatestKeySequence()
	if err != nil {
		return true, err
	}
	events, err := k.KeyEvents.LatestKeyEvents(ctx, sequence.CurrentSequence)
	if err != nil {
		logging.Log("EVENT-der5g").WithError(err).Warn("error retrieving new events")
		return true, err
	}
	if len(events) > 0 {
		logging.Log("EVENT-GBD23").Warn("view not up to date, retrying later")
		return true, nil
	}

	// try create new, short rotation
	err = k.lockAndGenerateSigningKeyPair(ctx, algorithm)
	logging.Log("EVENT-B4d21").OnError(err).Warn("could not create signing key")
	return true, err
}

func (k *KeyRepository) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	err := k.Locker.Renew(k.lockerID(), signingKey, k.SigningKeyRotation/2)
	if err != nil {
		return err
	}
	return k.GenerateSigningKeyPair(ctx, algorithm)
}

func (k *KeyRepository) lockerID() string {
	if k.lockID == "" {
		var err error
		k.lockID, err = os.Hostname()
		if err != nil || k.lockID == "" {
			k.lockID, err = id.SonyFlakeGenerator.Next()
			logging.Log("EVENT-bsdf6").OnError(err).Panic("unable to generate lockID")
		}
	}
	return k.lockID
}

func setOIDCCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: iamOrg})
}
