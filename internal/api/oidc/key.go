package oidc

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/jonboulle/clockwork"
	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var supportedWebKeyAlgs = []string{
	string(jose.EdDSA),
	string(jose.RS256),
	string(jose.RS384),
	string(jose.RS512),
	string(jose.ES256),
	string(jose.ES384),
	string(jose.ES512),
}

func supportedSigningAlgs(ctx context.Context) []string {
	if authz.GetFeatures(ctx).WebKey {
		return supportedWebKeyAlgs
	}
	return []string{string(jose.RS256)}
}

type cachedPublicKey struct {
	lastUse atomic.Int64 // unix micro time.
	expiry  *time.Time   // expiry may be nil if the key does not expire.
	webKey  *jose.JSONWebKey
}

func newCachedPublicKey(key *jose.JSONWebKey, expiry *time.Time, now time.Time) *cachedPublicKey {
	cachedKey := &cachedPublicKey{
		expiry: expiry,
		webKey: key,
	}
	cachedKey.setLastUse(now)
	return cachedKey
}

func (c *cachedPublicKey) setLastUse(now time.Time) {
	c.lastUse.Store(now.UnixMicro())
}

func (c *cachedPublicKey) getLastUse() time.Time {
	return time.UnixMicro(c.lastUse.Load())
}

func (c *cachedPublicKey) expired(now time.Time, validity time.Duration) bool {
	return c.getLastUse().Add(validity).Before(now)
}

// publicKeyCache caches public keys in a 2-dimensional map of Instance ID and Key ID.
// When a key is not present the queryKey function is called to obtain the key
// from the database.
type publicKeyCache struct {
	mtx          sync.RWMutex
	instanceKeys map[string]map[string]*cachedPublicKey

	// queryKey returns a public web key.
	// If the key does not have expiry, Time may be nil.
	queryKey func(ctx context.Context, keyID string) (*jose.JSONWebKey, *time.Time, error)
	clock    clockwork.Clock
}

// newPublicKeyCache initializes a keySetCache starts a purging Go routine.
// The purge routine deletes all public keys that are older than maxAge.
// When the passed context is done, the purge routine will terminate.
func newPublicKeyCache(background context.Context, maxAge time.Duration, queryKey func(ctx context.Context, keyID string) (*jose.JSONWebKey, *time.Time, error)) *publicKeyCache {
	k := &publicKeyCache{
		instanceKeys: make(map[string]map[string]*cachedPublicKey),
		queryKey:     queryKey,
		clock:        clockwork.FromContext(background), // defaults to real clock
	}
	go k.purgeOnInterval(background, k.clock.NewTicker(maxAge/5), maxAge)
	return k
}

func (k *publicKeyCache) purgeOnInterval(background context.Context, ticker clockwork.Ticker, maxAge time.Duration) {
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
				if key.expired(k.clock.Now(), maxAge) {
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

func (k *publicKeyCache) setKey(instanceID, keyID string, cachedKey *cachedPublicKey) {
	k.mtx.Lock()
	defer k.mtx.Unlock()

	if keys, ok := k.instanceKeys[instanceID]; ok {
		keys[keyID] = cachedKey
		return
	}
	k.instanceKeys[instanceID] = map[string]*cachedPublicKey{keyID: cachedKey}
}

func (k *publicKeyCache) getKey(ctx context.Context, keyID string) (_ *cachedPublicKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()

	k.mtx.RLock()
	key, ok := k.instanceKeys[instanceID][keyID]
	k.mtx.RUnlock()

	if ok {
		key.setLastUse(k.clock.Now())
	} else {
		newKey, expiry, err := k.queryKey(ctx, keyID)
		if err != nil {
			return nil, err
		}
		key = newCachedPublicKey(newKey, expiry, k.clock.Now())
		k.setKey(instanceID, keyID, key)
	}

	return key, nil
}

func (k *publicKeyCache) verifySignature(ctx context.Context, jws *jose.JSONWebSignature, checkKeyExpiry bool) (_ []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		err = oidcError(err)
		span.EndWithError(err)
	}()

	if len(jws.Signatures) != 1 {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-Gid9s", "Errors.Token.Invalid")
	}
	key, err := k.getKey(ctx, jws.Signatures[0].Header.KeyID)
	if err != nil {
		return nil, err
	}
	if checkKeyExpiry && key.expiry != nil && key.expiry.Before(k.clock.Now()) {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-ciF4k", "Errors.Key.ExpireBeforeNow")
	}
	return jws.Verify(key.webKey)
}

type oidcKeySet struct {
	*publicKeyCache

	keyExpiryCheck bool
}

// newOidcKeySet returns an oidc.KeySet implementation around the passed cache.
// It is advised to reuse the same cache if different key set configurations are required.
func newOidcKeySet(cache *publicKeyCache, opts ...keySetOption) *oidcKeySet {
	k := &oidcKeySet{
		publicKeyCache: cache,
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

// VerifySignature implements the oidc.KeySet interface.
func (k *oidcKeySet) VerifySignature(ctx context.Context, jws *jose.JSONWebSignature) (_ []byte, err error) {
	return k.verifySignature(ctx, jws, k.keyExpiryCheck)
}

type keySetOption func(*oidcKeySet)

// withKeyExpiryCheck forces VerifySignature to check the expiry of the public key.
// Note that public key expiry is not part of the standard,
// but is currently established behavior of zitadel.
// We might want to remove this check in the future.
func withKeyExpiryCheck(check bool) keySetOption {
	return func(k *oidcKeySet) {
		k.keyExpiryCheck = check
	}
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
		Use:   crypto.KeyUsageSigning.String(),
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
			return zerrors.ThrowNotFound(nil, "OIDC-ve4Qu", "Errors.Internal")
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
		return PrivateKeyToSigningKey(SelectSigningKey(keys.Keys), o.encAlg)
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

func PrivateKeyToSigningKey(key query.PrivateKey, algorithm crypto.EncryptionAlgorithm) (_ op.SigningKey, err error) {
	keyData, err := crypto.Decrypt(key.Key(), algorithm)
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

func SelectSigningKey(keys []query.PrivateKey) query.PrivateKey {
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

func (s *Server) Keys(ctx context.Context, r *op.Request[struct{}]) (_ *op.Response, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if !authz.GetFeatures(ctx).WebKey {
		return s.LegacyServer.Keys(ctx, r)
	}

	keyset, err := s.query.GetWebKeySet(ctx)
	if err != nil {
		return nil, err
	}

	// Return legacy keys, so we do not invalidate all tokens
	// once the feature flag is enabled.
	legacyKeys, err := s.query.ActivePublicKeys(ctx, time.Now())
	logging.OnError(err).Error("oidc server: active public keys (legacy)")
	appendPublicKeysToWebKeySet(keyset, legacyKeys)

	resp := op.NewResponse(keyset)
	if s.jwksCacheControlMaxAge != 0 {
		resp.Header.Set(http_util.CacheControl,
			fmt.Sprintf("max-age=%d, must-revalidate", int(s.jwksCacheControlMaxAge/time.Second)),
		)
	}

	return resp, nil
}

func appendPublicKeysToWebKeySet(keyset *jose.JSONWebKeySet, pubkeys *query.PublicKeys) {
	if pubkeys == nil || len(pubkeys.Keys) == 0 {
		return
	}
	keyset.Keys = slices.Grow(keyset.Keys, len(pubkeys.Keys))

	for _, key := range pubkeys.Keys {
		keyset.Keys = append(keyset.Keys, jose.JSONWebKey{
			Key:       key.Key(),
			KeyID:     key.ID(),
			Algorithm: key.Algorithm(),
			Use:       key.Use().String(),
		})
	}
}

func queryKeyFunc(q *query.Queries) func(ctx context.Context, keyID string) (*jose.JSONWebKey, *time.Time, error) {
	return func(ctx context.Context, keyID string) (*jose.JSONWebKey, *time.Time, error) {
		if authz.GetFeatures(ctx).WebKey {
			webKey, err := q.GetPublicWebKeyByID(ctx, keyID)
			if err == nil {
				return webKey, nil, nil
			}
			if !zerrors.IsNotFound(err) {
				return nil, nil, err
			}
		}

		pubKey, err := q.GetPublicKeyByID(ctx, keyID)
		if err != nil {
			return nil, nil, err
		}
		return jsonWebkey(pubKey), gu.Ptr(pubKey.Expiry()), nil
	}
}
