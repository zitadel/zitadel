package oidc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
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
	}
)

func queryKeyDB(_ context.Context, keyID string, current time.Time) (query.PublicKey, error) {
	if key, ok := keyDB[keyID]; ok {
		return key, nil
	}
	return nil, errors.New("not found")
}

func Test_keySetCache(t *testing.T) {
	background, cancel := context.WithCancel(
		clockwork.AddToContext(context.Background(), clock),
	)
	defer cancel()

	// create an empty keySet with a purge go routine, runs every Hour
	keySet := newKeySet(background, time.Hour, queryKeyDB)
	ctx := authz.NewMockContext("instanceID", "orgID", "userID")

	// query error
	_, err := keySet.getKey(ctx, "key9")
	require.Error(t, err)

	want := &jose.JSONWebKey{
		KeyID:     "key1",
		Algorithm: "alg",
		Use:       domain.KeyUsageSigning.String(),
	}

	// get key first time, populate the cache
	got, err := keySet.getKey(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, want, got)

	// move time forward
	clock.Advance(5 * time.Minute)
	time.Sleep(time.Millisecond)

	// key should still be in cache
	keySet.mtx.RLock()
	_, ok := keySet.instanceKeys["instanceID"]["key1"]
	require.True(t, ok)
	keySet.mtx.RUnlock()

	// the key is expired, should error
	_, err = keySet.getKey(ctx, "key1")
	require.Error(t, err)

	want = &jose.JSONWebKey{
		KeyID:     "key2",
		Algorithm: "alg",
		Use:       domain.KeyUsageSigning.String(),
	}

	// get the second key from DB
	got, err = keySet.getKey(ctx, "key2")
	require.NoError(t, err)
	assert.Equal(t, want, got)

	// move time forward
	clock.Advance(time.Hour)
	time.Sleep(time.Millisecond)

	// first key shoud be purged, second still present
	keySet.mtx.RLock()
	_, ok = keySet.instanceKeys["instanceID"]["key1"]
	require.False(t, ok)
	_, ok = keySet.instanceKeys["instanceID"]["key2"]
	require.True(t, ok)
	keySet.mtx.RUnlock()

	// get the second key from cache
	got, err = keySet.getKey(ctx, "key2")
	require.NoError(t, err)
	assert.Equal(t, want, got)

	// move time forward
	clock.Advance(10 * time.Hour)
	time.Sleep(time.Millisecond)

	// now the cache should be empty
	keySet.mtx.RLock()
	assert.Empty(t, keySet.instanceKeys)
	keySet.mtx.RUnlock()
}

func Test_keySetCache_VerifySignature(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	k := newKeySet(ctx, time.Second, queryKeyDB)

	tests := []struct {
		name string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
