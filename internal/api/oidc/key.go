package oidc

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/keypair"
)

const (
	locksTable = "projections.locks"
	signingKey = "signing_key"
)

func (o *OPStorage) GetKeySet(ctx context.Context) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	keys, err := o.query.ActivePublicKeys(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	webKeys := make([]jose.JSONWebKey, len(keys.Keys))
	for i, key := range keys.Keys {
		webKeys[i] = jose.JSONWebKey{
			KeyID:     key.ID(),
			Algorithm: key.Algorithm(),
			Use:       key.Use().String(),
			Key:       key.Key(),
		}
	}
	return &jose.JSONWebKeySet{Keys: webKeys}, nil
}

func (o *OPStorage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey) {
	renewTimer := time.NewTimer(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-o.keyChan:
				checkAfter := o.resetTimer(renewTimer, true)
				logging.Log("OIDC-dK432").Infof("requested next signing key check in %s", checkAfter)
			case <-renewTimer.C:
				o.getSigningKey(ctx, renewTimer, keyCh)
			}
		}
	}()
}

func (o *OPStorage) getSigningKey(ctx context.Context, renewTimer *time.Timer, keyCh chan<- jose.SigningKey) {
	keys, err := o.query.ActivePrivateSigningKey(ctx, time.Now().Add(o.signingKeyGracefulPeriod))
	if err != nil {
		checkAfter := o.resetTimer(renewTimer, true)
		logging.Log("OIDC-ASff").Infof("next signing key check in %s", checkAfter)
		return
	}
	if len(keys.Keys) == 0 {
		o.refreshSigningKey(ctx, keyCh, o.signingKeyAlgorithm, keys.LatestSequence)
		checkAfter := o.resetTimer(renewTimer, true)
		logging.Log("OIDC-ASDf3").Infof("next signing key check in %s", checkAfter)
		return
	}
	err = o.exchangeSigningKey(selectSigningKey(keys.Keys), keyCh)
	logging.Log("OIDC-aDfg3").OnError(err).Error("could not exchange signing key")
	checkAfter := o.resetTimer(renewTimer, err != nil)
	logging.Log("OIDC-dK432").Infof("next signing key check in %s", checkAfter)
}

func (o *OPStorage) resetTimer(timer *time.Timer, shortRefresh bool) (nextCheck time.Duration) {
	//if !timer.Stop() {
	//	<-timer.C
	//}
	nextCheck = o.signingKeyRotationCheck
	defer func() { timer.Reset(nextCheck) }()
	if shortRefresh || o.currentKey == nil {
		return nextCheck
	}
	maxLifetime := time.Until(o.currentKey.Expiry())
	if maxLifetime < o.signingKeyGracefulPeriod+2*o.signingKeyRotationCheck {
		return nextCheck
	}
	return maxLifetime - o.signingKeyGracefulPeriod - o.signingKeyRotationCheck
}

func (o *OPStorage) refreshSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string, sequence *query.LatestSequence) {
	if o.currentKey != nil && o.currentKey.Expiry().Before(time.Now().UTC()) {
		logging.Log("OIDC-ADg26").Info("unset current signing key")
		keyCh <- jose.SigningKey{}
	}
	ok, err := o.ensureIsLatestKey(ctx, sequence.Sequence)
	if err != nil {
		logging.Log("OIDC-sdz53").WithError(err).Error("could not ensure latest key")
		return
	}
	if !ok {
		logging.Log("EVENT-GBD23").Warn("view not up to date, retrying later")
		return
	}
	err = o.lockAndGenerateSigningKeyPair(ctx, algorithm)
	logging.Log("EVENT-B4d21").OnError(err).Warn("could not create signing key")
}

func (o *OPStorage) ensureIsLatestKey(ctx context.Context, sequence uint64) (bool, error) {
	maxSequence, err := o.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return sequence == maxSequence, nil
}

func (o *OPStorage) exchangeSigningKey(key query.PrivateKey, keyCh chan<- jose.SigningKey) (err error) {
	if o.currentKey != nil && o.currentKey.ID() == key.ID() {
		logging.Log("OIDC-Abb3e").Info("no new signing key")
		return nil
	}
	keyData, err := crypto.Decrypt(key.Key(), o.encAlg)
	if err != nil {
		return err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return err
	}
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(key.Algorithm()),
		Key: jose.JSONWebKey{
			KeyID: key.ID(),
			Key:   privateKey,
		},
	}
	o.currentKey = key
	logging.LogWithFields("OIDC-dsg54", "keyID", key.ID()).Info("exchanged signing key")
	return nil
}

func (o *OPStorage) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	logging.Log("OIDC-sdz53").Info("lock and generate signing key pair")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := o.locker.Lock(ctx, o.signingKeyRotationCheck*2)
	err, ok := <-errs
	if err != nil || !ok {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.Log("OIDC-Dfg32").OnError(err).Warn("initial lock failed")
		return err
	}

	return o.command.GenerateSigningKeyPair(ctx, algorithm)
}

func (o *OPStorage) getMaxKeySequence(ctx context.Context) (uint64, error) {
	return o.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner(domain.IAMID).
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			Builder(),
	)
}

func selectSigningKey(keys []query.PrivateKey) query.PrivateKey {
	return keys[len(keys)-1]
}
