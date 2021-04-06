package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"os"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/key/model"
	key_view "github.com/caos/zitadel/internal/key/repository/view"
)

type KeyRepository struct {
	Commands                 *command.Commands
	Eventstore               *eventstore.Eventstore
	View                     *view.View
	SigningKeyRotationCheck  time.Duration
	SigningKeyGracefulPeriod time.Duration
	KeyAlgorithm             crypto.EncryptionAlgorithm
	KeyChan                  <-chan *model.KeyView
	Locker                   spooler.Locker
	lockID                   string
	currentKeyID             string
	currentKeyExpiration     time.Time
}

const (
	signingKey = "signing_key"
)

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string) {
	renewTimer := time.After(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case key := <-k.KeyChan:
				refreshed, err := k.refreshSigningKey(ctx, key, keyCh, algorithm)
				logging.Log("KEY-asd5g").OnError(err).Error("could not refresh signing key on key channel push")
				renewTimer = time.After(k.getRenewTimer(refreshed))
			case <-renewTimer:
				key, err := k.latestSigningKey()
				logging.Log("KEY-DAfh4").OnError(err).Error("could not check for latest signing key")
				refreshed, err := k.refreshSigningKey(ctx, key, keyCh, algorithm)
				logging.Log("KEY-DAfh4").OnError(err).Error("could not refresh signing key when ensuring key")
				renewTimer = time.After(k.getRenewTimer(refreshed))
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

func (k *KeyRepository) getRenewTimer(refreshed bool) time.Duration {
	duration := k.SigningKeyRotationCheck
	if refreshed {
		duration = k.currentKeyExpiration.Sub(time.Now().Add(k.SigningKeyGracefulPeriod + k.SigningKeyRotationCheck*2))
	}
	logging.LogWithFields("EVENT-dK432", "in", duration).Info("next signing key check")
	return duration
}

func (k *KeyRepository) latestSigningKey() (shortRefresh *model.KeyView, err error) {
	key, errView := k.View.GetActivePrivateKeyForSigning(time.Now().UTC().Add(k.SigningKeyGracefulPeriod))
	if errView != nil && !errors.IsNotFound(errView) {
		logging.Log("EVENT-GEd4h").WithError(errView).Warn("could not get signing key")
		return nil, errView
	}
	return key, nil
}

func (k *KeyRepository) ensureIsLatestKey(ctx context.Context) (bool, error) {
	sequence, err := k.View.GetLatestKeySequence()
	if err != nil {
		return false, err
	}
	events, err := k.getKeyEvents(ctx, sequence.CurrentSequence)
	if err != nil {
		logging.Log("EVENT-der5g").WithError(err).Warn("error retrieving new events")
		return false, err
	}
	if len(events) > 0 {
		logging.Log("EVENT-GBD23").Warn("view not up to date, retrying later")
		return false, nil
	}
	return true, nil
}

func (k *KeyRepository) refreshSigningKey(ctx context.Context, key *model.KeyView, keyCh chan<- jose.SigningKey, algorithm string) (refreshed bool, err error) {
	if key == nil {
		if k.currentKeyExpiration.Before(time.Now().UTC()) {
			logging.Log("EVENT-ADg26").Info("unset current signing key")
			keyCh <- jose.SigningKey{}
		}
		if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
			return false, err
		}
		logging.Log("EVENT-sdz53").Info("lock and generate signing key pair")
		err = k.lockAndGenerateSigningKeyPair(ctx, algorithm)
		logging.Log("EVENT-B4d21").OnError(err).Warn("could not create signing key")
		return false, err
	}

	if k.currentKeyID == key.ID {
		logging.Log("EVENT-Abb3e").Info("no new signing key")
		return false, nil
	}
	if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
		logging.Log("EVENT-HJd92").Info("signing key in view is not latest key")
		return false, err
	}
	signingKey, err := model.SigningKeyFromKeyView(key, k.KeyAlgorithm)
	if err != nil {
		logging.Log("EVENT-HJd92").WithError(err).Error("signing key cannot be decrypted -> immediate refresh")
		return k.refreshSigningKey(ctx, nil, keyCh, algorithm)
	}
	k.currentKeyID = signingKey.ID
	k.currentKeyExpiration = key.Expiry
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(signingKey.Algorithm),
		Key: jose.JSONWebKey{
			KeyID: signingKey.ID,
			Key:   signingKey.Key,
		},
	}
	logging.LogWithFields("EVENT-dsg54", "keyID", signingKey.ID).Info("refreshed signing key")
	return true, nil
}

func (k *KeyRepository) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	err := k.Locker.Renew(k.lockerID(), signingKey, k.SigningKeyRotationCheck*2)
	if err != nil {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		return err
	}
	return k.Commands.GenerateSigningKeyPair(ctx, algorithm)
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

func (k *KeyRepository) getKeyEvents(ctx context.Context, sequence uint64) ([]eventstore.EventReader, error) {
	return k.Eventstore.FilterEvents(ctx, key_view.KeyPairQuery(sequence))
}
