package oidc

import (
	"context"
	"fmt"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/v2/pkg/op"
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
	oidcUser   = "OIDC"

	retryBackoff   = 500 * time.Millisecond
	retryCount     = 3
	lockDuration   = retryCount * retryBackoff * 5
	gracefulPeriod = 10 * time.Minute
)

//SigningKey wraps the query.PrivateKey to implement the op.SigningKey interface
type SigningKey struct {
	algorithm jose.SignatureAlgorithm
	id        string
	key       interface{}
}

func (s *SigningKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *SigningKey) Key() interface{} {
	return s.key
}

func (s *SigningKey) ID() string {
	return s.id
}

//PublicKey wraps the query.PublicKey to implement the op.Key interface
type PublicKey struct {
	key query.PublicKey
}

func (s *PublicKey) Algorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Algorithm())
}

func (s *PublicKey) Use() string {
	return s.key.Use().String()
}

func (s *PublicKey) Key() interface{} {
	return s.key.Key()
}

func (s *PublicKey) ID() string {
	return s.key.ID()
}

//KeySet implements the op.Storage interface
func (o *OPStorage) KeySet(ctx context.Context) (keys []op.Key, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	err = retry(func() error {
		publicKeys, err := o.query.ActivePublicKeys(ctx, time.Now())
		if err != nil {
			return err
		}
		keys = make([]op.Key, len(publicKeys.Keys))
		for i, key := range publicKeys.Keys {
			keys[i] = &PublicKey{key}
		}
		return nil
	})
	return keys, err
}

//SignatureAlgorithms implements the op.Storage interface
func (o *OPStorage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	key, err := o.SigningKey(ctx)
	if err != nil {
		return nil, err
	}
	return []jose.SignatureAlgorithm{key.SignatureAlgorithm()}, nil
}

//SigningKey implements the op.Storage interface
func (o *OPStorage) SigningKey(ctx context.Context) (key op.SigningKey, err error) {
	err = retry(func() error {
		key, err = o.getSigningKey(ctx)
		if err != nil {
			return err
		}
		if key == nil {
			return errors.ThrowInternal(nil, "test", "test")
		}
		return nil
	})
	return key, err
}

func (o *OPStorage) getSigningKey(ctx context.Context) (op.SigningKey, error) {
	keys, err := o.query.ActivePrivateSigningKey(ctx, time.Now().Add(gracefulPeriod))
	if err != nil {
		return nil, err
	}
	if len(keys.Keys) > 0 {
		return o.privateKeyToSigningKey(selectSigningKey(keys.Keys))
	}
	var sequence uint64
	if keys.LatestSequence != nil {
		sequence = keys.LatestSequence.Sequence
	}
	return nil, o.refreshSigningKey(ctx, o.signingKeyAlgorithm, sequence)
}

func (o *OPStorage) refreshSigningKey(ctx context.Context, algorithm string, sequence uint64) error {
	ok, err := o.ensureIsLatestKey(ctx, sequence)
	if err != nil || !ok {
		return errors.ThrowInternal(err, "OIDC-ASfh3", "cannot ensure that projection is up to date")
	}
	err = o.lockAndGenerateSigningKeyPair(ctx, algorithm)
	if err != nil {
		return errors.ThrowInternal(err, "OIDC-ADh31", "could not create signing key")
	}
	return errors.ThrowInternal(nil, "OIDC-Df1bh", "")
}

func (o *OPStorage) ensureIsLatestKey(ctx context.Context, sequence uint64) (bool, error) {
	maxSequence, err := o.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return sequence == maxSequence, nil
}

func (o *OPStorage) privateKeyToSigningKey(key query.PrivateKey) (_ op.SigningKey, err error) {
	keyData, err := crypto.Decrypt(key.Key(), o.encAlg)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return nil, err
	}
	return &SigningKey{
		algorithm: jose.SignatureAlgorithm(key.Algorithm()),
		key:       privateKey,
		id:        key.ID(),
	}, nil
}

func (o *OPStorage) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	logging.Info("lock and generate signing key pair")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := o.locker.Lock(ctx, lockDuration, authz.GetInstance(ctx).InstanceID())
	err, ok := <-errs
	if err != nil || !ok {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.OnError(err).Warn("initial lock failed")
		return err
	}

	return o.command.GenerateSigningKeyPair(setOIDCCtx(ctx), algorithm)
}

func (o *OPStorage) getMaxKeySequence(ctx context.Context) (uint64, error) {
	return o.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner(authz.GetInstance(ctx).InstanceID()).
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			Builder(),
	)
}

func selectSigningKey(keys []query.PrivateKey) query.PrivateKey {
	return keys[len(keys)-1]
}

func setOIDCCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: authz.GetInstance(ctx).InstanceID()})
}

func retry(retryable func() error) (err error) {
	for i := 0; i < retryCount; i++ {
		time.Sleep(retryBackoff)
		err = retryable()
		if err == nil {
			return nil
		}
	}
	return err
}
