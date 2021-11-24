package oidc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/keypair"
)

const (
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
				checkAfter := o.resetTimer(renewTimer, false)
				logging.Log("OIDC-dK432").Infof("requested next signing key check in %d", checkAfter)
			case <-renewTimer.C:
				keys, err := o.query.ActivePrivateSigningKey(ctx, time.Now())
				if err != nil {
					checkAfter := o.resetTimer(renewTimer, false)
					logging.Log("OIDC-dK432").Infof("requested next signing key check in %d", checkAfter)
					continue
				}
				if len(keys.Keys) == 0 {
					o.refreshSigningKey(ctx, keyCh, o.signingKeyAlgorithm)
					checkAfter := o.resetTimer(renewTimer, false)
					logging.Log("OIDC-ASDf3").Infof("tried refreshing signing key, next check in %d", checkAfter)
					continue
				}
				exchanged, err := o.exchangeSigningKey(selectSigningKey(keys.Keys), keyCh)
				logging.Log("OIDC-aDfg3").OnError(err).Error("could not exchange signing key")
				checkAfter := o.resetTimer(renewTimer, exchanged)
				logging.Log("OIDC-dK432").Infof("next signing key check in %d", checkAfter)
			}
		}
	}()
}

func (o *OPStorage) resetTimer(timer *time.Timer, exchanged bool) time.Duration {
	if timer.Stop() {
		<-timer.C
	}
	duration := o.signingKeyRotationCheck
	if exchanged && o.currentKey == nil {
		duration = o.currentKey.Expiry().Sub(time.Now().Add(o.signingKeyGracefulPeriod + o.signingKeyRotationCheck*2))
	}
	timer.Reset(duration)
	return duration
}

func (o *OPStorage) refreshSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string) {
	if o.currentKey != nil && o.currentKey.Expiry().Before(time.Now().UTC()) {
		logging.Log("OIDC-ADg26").Info("unset current signing key")
		keyCh <- jose.SigningKey{}
	}
	ok, err := o.ensureIsLatestKey(ctx)
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

func (o *OPStorage) ensureIsLatestKey(ctx context.Context) (bool, error) {
	var sequence uint64
	if o.currentKey != nil {
		sequence = o.currentKey.Sequence()
	}
	maxSequence, err := o.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return sequence == maxSequence, nil
}

func (o *OPStorage) exchangeSigningKey(key query.PrivateKey, keyCh chan<- jose.SigningKey) (refreshed bool, err error) {
	if o.currentKey != nil && o.currentKey.ID() == key.ID() {
		logging.Log("OIDC-Abb3e").Info("no new signing key")
		return false, nil
	}
	keyData, err := crypto.Decrypt(key.Key(), o.encAlg)
	if err != nil {
		return false, err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return false, err
	}
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(key.Algorithm()),
		Key: jose.JSONWebKey{
			KeyID: key.ID(),
			Key:   privateKey,
		},
	}
	o.currentKey = key
	logging.LogWithFields("OIDC-dsg54", "keyID", key.ID()).Info("refreshed signing key")
	return true, nil
}

func (o *OPStorage) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	logging.Log("OIDC-sdz53").Info("lock and generate signing key pair")
	err := o.locker.Renew(o.lockerID(), signingKey, o.signingKeyRotationCheck*2)
	if err != nil {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		return err
	}
	return o.command.GenerateSigningKeyPair(ctx, algorithm)
}

func (o *OPStorage) lockerID() string {
	if o.lockID != "" {
		return o.lockerID()
	}
	var err error
	o.lockID, err = os.Hostname()
	if err != nil || o.lockID == "" {
		o.lockID, err = id.SonyFlakeGenerator.Next()
		logging.Log("EVENT-bsdf6").OnError(err).Panic("unable to generate lockID")
	}
	return o.lockID
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
