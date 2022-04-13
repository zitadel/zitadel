package oidc

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/telemetry/tracing"
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
				if !renewTimer.Stop() {
					<-renewTimer.C
				}
				checkAfter := o.resetTimer(renewTimer, true)
				logging.Infof("requested next signing key check in %s", checkAfter)
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
		logging.Infof("next signing key check in %s", checkAfter)
		return
	}
	if len(keys.Keys) == 0 {
		var sequence uint64
		if keys.LatestSequence != nil {
			sequence = keys.LatestSequence.Sequence
		}
		o.refreshSigningKey(ctx, keyCh, o.signingKeyAlgorithm, sequence)
		checkAfter := o.resetTimer(renewTimer, true)
		logging.Infof("next signing key check in %s", checkAfter)
		return
	}
	err = o.exchangeSigningKey(selectSigningKey(keys.Keys), keyCh)
	logging.OnError(err).Error("could not exchange signing key")
	checkAfter := o.resetTimer(renewTimer, err != nil)
	logging.Infof("next signing key check in %s", checkAfter)
}

func (o *OPStorage) resetTimer(timer *time.Timer, shortRefresh bool) (nextCheck time.Duration) {
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

func (o *OPStorage) refreshSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string, sequence uint64) {
	if o.currentKey != nil && o.currentKey.Expiry().Before(time.Now().UTC()) {
		logging.Info("unset current signing key")
		keyCh <- jose.SigningKey{}
	}
	ok, err := o.ensureIsLatestKey(ctx, sequence)
	if err != nil {
		logging.New().WithError(err).Error("could not ensure latest key")
		return
	}
	if !ok {
		logging.Warn("view not up to date, retrying later")
		return
	}
	err = o.lockAndGenerateSigningKeyPair(ctx, algorithm)
	logging.OnError(err).Warn("could not create signing key")
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
		logging.Info("no new signing key")
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
	logging.WithFields("keyID", key.ID()).Info("exchanged signing key")
	return nil
}

func (o *OPStorage) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	logging.Info("lock and generate signing key pair")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := o.locker.Lock(ctx, o.signingKeyRotationCheck*2, authz.GetInstance(ctx).InstanceID())
	err, ok := <-errs
	if err != nil || !ok {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.OnError(err).Warn("initial lock failed")
		return err
	}

	return o.command.GenerateSigningKeyPair(ctx, algorithm)
}

func (o *OPStorage) getMaxKeySequence(ctx context.Context) (uint64, error) {
	return o.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner("system"). //TODO: change with multi issuer
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			Builder(),
	)
}

func selectSigningKey(keys []query.PrivateKey) query.PrivateKey {
	return keys[len(keys)-1]
}
