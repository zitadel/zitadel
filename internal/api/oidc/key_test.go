package oidc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type publicKey struct {
	id     string
	alg    string
	use    domain.KeyUsage
	seq    uint64
	expiry time.Time
	key    any
}

func (k *publicKey) ID() string {
	return k.id
}

func (k *publicKey) Algorithm() string {
	return k.alg
}

func (k *publicKey) Use() domain.KeyUsage {
	return k.use
}

func (k *publicKey) Sequence() uint64 {
	return k.seq
}

func (k *publicKey) Expiry() time.Time {
	return k.expiry
}

func (k *publicKey) Key() any {
	return k.key
}

var (
	clock = clockwork.NewFakeClock()
	keyDB = map[string]*publicKey{
		"key1": {
			id:     "key1",
			alg:    "alg",
			use:    domain.KeyUsageSigning,
			seq:    1,
			expiry: clock.Now().Add(time.Minute),
		},
		"key2": {
			id:     "key2",
			alg:    "alg",
			use:    domain.KeyUsageSigning,
			seq:    3,
			expiry: clock.Now().Add(10 * time.Hour),
		},
		"exp1": {
			id:     "key2",
			alg:    "alg",
			use:    domain.KeyUsageSigning,
			seq:    4,
			expiry: clock.Now().Add(-time.Hour),
		},
	}
)

func queryKeyDB(_ context.Context, keyID string) (query.PublicKey, error) {
	if key, ok := keyDB[keyID]; ok {
		return key, nil
	}
	return nil, errors.New("not found")
}

func Test_publicKeyCache(t *testing.T) {
	background, cancel := context.WithCancel(
		clockwork.AddToContext(context.Background(), clock),
	)
	defer cancel()

	// create an empty cache with a purge go routine, runs every minute.
	// keys are cached for at least 1 Hour after last use.
	cache := newPublicKeyCache(background, time.Hour, queryKeyDB)
	ctx := authz.NewMockContext("instanceID", "orgID", "userID")

	// query error
	_, err := cache.getKey(ctx, "key9")
	require.Error(t, err)

	// get key first time, populate the cache
	got, err := cache.getKey(ctx, "key1")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, keyDB["key1"], got.PublicKey)

	// move time forward
	clock.Advance(15 * time.Minute)
	time.Sleep(time.Millisecond)

	// key should still be in cache
	cache.mtx.RLock()
	_, ok := cache.instanceKeys["instanceID"]["key1"]
	require.True(t, ok)
	cache.mtx.RUnlock()

	// move time forward
	clock.Advance(50 * time.Minute)
	time.Sleep(time.Millisecond)

	// get the second key from DB
	got, err = cache.getKey(ctx, "key2")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, keyDB["key2"], got.PublicKey)

	// move time forward
	clock.Advance(15 * time.Minute)
	time.Sleep(time.Millisecond)

	// first key should be purged, second still present
	cache.mtx.RLock()
	_, ok = cache.instanceKeys["instanceID"]["key1"]
	require.False(t, ok)
	_, ok = cache.instanceKeys["instanceID"]["key2"]
	require.True(t, ok)
	cache.mtx.RUnlock()

	// get the second key from cache
	got, err = cache.getKey(ctx, "key2")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, keyDB["key2"], got.PublicKey)

	// move time forward
	clock.Advance(2 * time.Hour)
	time.Sleep(time.Millisecond)

	// now the cache should be empty
	cache.mtx.RLock()
	assert.Empty(t, cache.instanceKeys)
	cache.mtx.RUnlock()
}

func Test_oidcKeySet_VerifySignature(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cache := newPublicKeyCache(ctx, time.Second, queryKeyDB)

	tests := []struct {
		name string
		opts []keySetOption
		jws  *jose.JSONWebSignature
	}{
		{
			name: "invalid token",
			jws:  &jose.JSONWebSignature{},
		},
		{
			name: "key not found",
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "xxx",
					},
				}},
			},
		},
		{
			name: "verify error",
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "key1",
					},
				}},
			},
		},
		{
			name: "expired, no check",
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "exp1",
					},
				}},
			},
		},
		{
			name: "expired, with check",
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "exp1",
					},
				}},
			},
			opts: []keySetOption{
				withKeyExpiryCheck(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newOidcKeySet(cache, tt.opts...)
			_, err := k.VerifySignature(ctx, tt.jws)
			require.Error(t, err)
		})
	}
}

func Test_keySetMap_VerifySignature(t *testing.T) {
	tests := []struct {
		name string
		k    keySetMap
		jws  *jose.JSONWebSignature
	}{
		{
			name: "invalid signature",
			k: keySetMap{
				"key1": []byte("foo"),
			},
			jws: &jose.JSONWebSignature{},
		},
		{
			name: "parse error",
			k: keySetMap{
				"key1": []byte("foo"),
			},
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "key1",
					},
				}},
			},
		},
		{
			name: "verify error",
			k: keySetMap{
				"key1": []byte("-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsvX9P58JFxEs5C+L+H7W\nduFSWL5EPzber7C2m94klrSV6q0bAcrYQnGwFOlveThsY200hRbadKaKjHD7qIKH\nDEe0IY2PSRht33Jye52AwhkRw+M3xuQH/7R8LydnsNFk2KHpr5X2SBv42e37LjkE\nslKSaMRgJW+v0KZ30piY8QsdFRKKaVg5/Ajt1YToM1YVsdHXJ3vmXFMtypLdxwUD\ndIaLEX6pFUkU75KSuEQ/E2luT61Q3ta9kOWm9+0zvi7OMcbdekJT7mzcVnh93R1c\n13ZhQCLbh9A7si8jKFtaMWevjayrvqQABEcTN9N4Hoxcyg6l4neZtRDk75OMYcqm\nDQIDAQAB\n-----END RSA PUBLIC KEY-----\n"),
			},
			jws: &jose.JSONWebSignature{
				Signatures: []jose.Signature{{
					Header: jose.Header{
						KeyID: "key1",
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.k.VerifySignature(context.Background(), tt.jws)
			require.Error(t, err)
		})
	}
}
