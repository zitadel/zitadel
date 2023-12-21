package oidc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/jonboulle/clockwork"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// keySetCache implements oidc.KeySet for Access Token verification.
// Public Keys are cached in a 2-dimensional map of Instance ID and Key ID.
// When a key is not present the queryKey function is called to obtain the key
// from the database.
type keySetCache struct {
	mtx          sync.RWMutex
	instanceKeys map[string]map[string]query.PublicKey
	queryKey     func(ctx context.Context, keyID string, current time.Time) (query.PublicKey, error)
	clock        clockwork.Clock
}

// newKeySet initializes a keySetCache and starts a purging Go routine,
// which runs once every purgeInterval.
// When the passed context is done, the purge routine will terminate.
func newKeySet(background context.Context, purgeInterval time.Duration, queryKey func(ctx context.Context, keyID string, current time.Time) (query.PublicKey, error)) *keySetCache {
	k := &keySetCache{
		instanceKeys: make(map[string]map[string]query.PublicKey),
		queryKey:     queryKey,
		clock:        clockwork.FromContext(background), // defaults to real clock
	}
	go k.purgeOnInterval(background, k.clock.NewTicker(purgeInterval))
	return k
}

func (k *keySetCache) purgeOnInterval(background context.Context, ticker clockwork.Ticker) {
	defer ticker.Stop()
	for {
		select {
		case <-background.Done():
			return
		case <-ticker.Chan():
		}

		// do the actual purging
		k.mtx.Lock()
		for instanceID, keys := range k.instanceKeys {
			for keyID, key := range keys {
				if key.Expiry().Before(k.clock.Now()) {
					delete(keys, keyID)
				}
			}
			if len(keys) == 0 {
				delete(k.instanceKeys, instanceID)
			}
		}
		k.mtx.Unlock()
	}
}

func (k *keySetCache) setKey(instanceID, keyID string, key query.PublicKey) {
	k.mtx.Lock()
	defer k.mtx.Unlock()

	if keys, ok := k.instanceKeys[instanceID]; ok {
		keys[keyID] = key
		return
	}

	k.instanceKeys[instanceID] = map[string]query.PublicKey{keyID: key}
}

func (k *keySetCache) getKey(ctx context.Context, keyID string) (_ *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()

	k.mtx.RLock()
	key, ok := k.instanceKeys[instanceID][keyID]
	k.mtx.RUnlock()

	if ok {
		if key.Expiry().After(k.clock.Now()) {
			return jsonWebkey(key), nil
		}
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-Zoh9E", "Errors.Key.ExpireBeforeNow")
	}

	key, err = k.queryKey(ctx, keyID, k.clock.Now())
	if err != nil {
		return nil, err
	}
	k.setKey(instanceID, keyID, key)
	return jsonWebkey(key), nil
}

// VerifySignature implements the oidc.KeySet interface.
func (k *keySetCache) VerifySignature(ctx context.Context, jws *jose.JSONWebSignature) (_ []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if len(jws.Signatures) != 1 {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-Gid9s", "Errors.Token.Invalid")
	}
	key, err := k.getKey(ctx, jws.Signatures[0].Header.KeyID)
	if err != nil {
		return nil, err
	}
	return jws.Verify(key)
}

func jsonWebkey(key query.PublicKey) *jose.JSONWebKey {
	return &jose.JSONWebKey{
		KeyID:     key.ID(),
		Algorithm: key.Algorithm(),
		Use:       key.Use().String(),
		Key:       key.Key(),
	}
}

// keySetMap is a mapping of key IDs to public key data.
type keySetMap map[string][]byte

// getKey finds the keyID and parses the public key data
// into a JSONWebKey.
func (k keySetMap) getKey(keyID string) (*jose.JSONWebKey, error) {
	pubKey, err := crypto.BytesToPublicKey(k[keyID])
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		Key:   pubKey,
		KeyID: keyID,
		Use:   domain.KeyUsageSigning.String(),
	}, nil
}

// VerifySignature implements the oidc.KeySet interface.
func (k keySetMap) VerifySignature(ctx context.Context, jws *jose.JSONWebSignature) ([]byte, error) {
	if len(jws.Signatures) != 1 {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-Eeth6", "Errors.Token.Invalid")
	}
	key, err := k.getKey(jws.Signatures[0].Header.KeyID)
	if err != nil {
		return nil, err
	}
	return jws.Verify(key)
}

const (
	locksTable = "projections.locks"
	signingKey = "signing_key"
	oidcUser   = "OIDC"

	retryBackoff   = 500 * time.Millisecond
	retryCount     = 3
	lockDuration   = retryCount * retryBackoff * 5
	gracefulPeriod = 10 * time.Minute
)

// SigningKey wraps the query.PrivateKey to implement the op.SigningKey interface
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

// PublicKey wraps the query.PublicKey to implement the op.Key interface
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

// KeySet implements the op.Storage interface
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

// SignatureAlgorithms implements the op.Storage interface
func (o *OPStorage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	key, err := o.SigningKey(ctx)
	if err != nil {
		logging.WithError(err).Warn("unable to fetch signing key")
		return nil, err
	}
	return []jose.SignatureAlgorithm{key.SignatureAlgorithm()}, nil
}

// SigningKey implements the op.Storage interface
func (o *OPStorage) SigningKey(ctx context.Context) (key op.SigningKey, err error) {
	err = retry(func() error {
		key, err = o.getSigningKey(ctx)
		if err != nil {
			return err
		}
		if key == nil {
			return zerrors.ThrowInternal(nil, "test", "test")
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
	var position float64
	if keys.State != nil {
		position = keys.State.Position
	}
	return nil, o.refreshSigningKey(ctx, o.signingKeyAlgorithm, position)
}

func (o *OPStorage) refreshSigningKey(ctx context.Context, algorithm string, position float64) error {
	ok, err := o.ensureIsLatestKey(ctx, position)
	if err != nil || !ok {
		return zerrors.ThrowInternal(err, "OIDC-ASfh3", "cannot ensure that projection is up to date")
	}
	err = o.lockAndGenerateSigningKeyPair(ctx, algorithm)
	if err != nil {
		return zerrors.ThrowInternal(err, "OIDC-ADh31", "could not create signing key")
	}
	return zerrors.ThrowInternal(nil, "OIDC-Df1bh", "")
}

func (o *OPStorage) ensureIsLatestKey(ctx context.Context, position float64) (bool, error) {
	maxSequence, err := o.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return position >= maxSequence, nil
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
		if zerrors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.OnError(err).Debug("initial lock failed")
		return err
	}

	return o.command.GenerateSigningKeyPair(setOIDCCtx(ctx), algorithm)
}

func (o *OPStorage) getMaxKeySequence(ctx context.Context) (float64, error) {
	return o.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner(authz.GetInstance(ctx).InstanceID()).
			AwaitOpenTransactions().
			AllowTimeTravel().
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			EventTypes(
				keypair.AddedEventType,
			).
			Or().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
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
		err = retryable()
		if err == nil {
			return nil
		}
		time.Sleep(retryBackoff)
	}
	return err
}
