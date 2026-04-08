package crypto

import (
	"context"
	"errors"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type mockKeyStorage struct {
	keys Keys
}

func (s *mockKeyStorage) ReadKeys() (Keys, error) {
	return s.keys, nil
}

func (s *mockKeyStorage) ReadKey(id string) (*Key, error) {
	return &Key{
		ID:    id,
		Value: s.keys[id],
	}, nil
}

func (*mockKeyStorage) CreateKeys(context.Context, ...*Key) error {
	return errors.New("mockKeyStorage.CreateKeys not implemented")
}

func newTestAESCrypto(t testing.TB) *AESCrypto {
	keyConfig := &KeyConfig{
		EncryptionKeyID:  "keyID",
		DecryptionKeyIDs: []string{"keyID"},
	}
	keys := Keys{"keyID": "ThisKeyNeedsToHave32Characters!!"}
	aesCrypto, err := NewAESCrypto(keyConfig, &mockKeyStorage{keys: keys})
	require.NoError(t, err)
	return aesCrypto
}

func TestAESCrypto_DecryptString(t *testing.T) {
	aesCrypto := newTestAESCrypto(t)
	const input = "SecretData"
	crypted, err := aesCrypto.Encrypt([]byte(input))
	require.NoError(t, err)

	type args struct {
		value []byte
		keyID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "unknown key id error",
			args: args{
				value: crypted,
				keyID: "foo",
			},
			wantErr: zerrors.ThrowNotFound(nil, "CRYPT-nkj1s", "unknown key id"),
		},
		{
			name: "ok",
			args: args{
				value: crypted,
				keyID: "keyID",
			},
			want: input,
		},
	}
	for _, tt := range tests {
		got, err := aesCrypto.DecryptString(tt.args.value, tt.args.keyID)
		require.ErrorIs(t, err, tt.wantErr)
		assert.Equal(t, tt.want, got)
	}
}

func FuzzAESCrypto_DecryptString(f *testing.F) {
	aesCrypto := newTestAESCrypto(f)
	tests := []string{
		" ",
		"SecretData",
		"FooBar",
		"HelloWorld",
	}
	for _, input := range tests {
		tc, err := aesCrypto.Encrypt([]byte(input))
		require.NoError(f, err)
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, value []byte) {
		got, err := aesCrypto.DecryptString(value, "keyID")
		if errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "CRYPT-23kH1", "cipher text block too short")) {
			return
		}
		if errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "CRYPT-hiCh0", "non-UTF-8 in decrypted string")) {
			return
		}
		require.NoError(t, err)
		assert.True(t, utf8.ValidString(got), "result is not valid UTF-8")
	})
}

func newTestAES256GCMCrypto(fallbackDecrypt func(value []byte, key string) ([]byte, error)) func(t testing.TB) EncryptionAlgorithm {
	keyConfig := &KeyConfig{
		EncryptionKeyID:  "keyID",
		DecryptionKeyIDs: []string{"keyID"},
	}
	keys := Keys{"keyID": "ThisKeyNeedsToHave32Characters!!"}
	return func(t testing.TB) EncryptionAlgorithm {
		aesCrypto, err := NewAES256GCMCrypto(
			keyConfig,
			&mockKeyStorage{keys: keys},
			WithAES256GCMCryptoFallbackDecrypt(fallbackDecrypt),
		)
		require.NoError(t, err)
		return aesCrypto
	}
}

func TestAES256GCMCrypto_EncryptString(t *testing.T) {
	tests := []struct {
		name       string
		config     *KeyConfig
		keyStorage KeyStorage
		value      string
		want       bool
		wantErr    error
	}{
		{
			name: "ok",
			config: &KeyConfig{
				EncryptionKeyID:  "keyID",
				DecryptionKeyIDs: []string{"keyID"},
			},
			keyStorage: &mockKeyStorage{keys: Keys{"keyID": "ThisKeyNeedsToHave32Characters!!"}},
			value:      "SecretData",
			want:       true,
			wantErr:    nil,
		},
		{
			name: "empty key error",
			config: &KeyConfig{
				EncryptionKeyID:  "keyID",
				DecryptionKeyIDs: []string{"keyID"},
			},
			keyStorage: &mockKeyStorage{keys: Keys{"keyID": ""}},
			value:      "SecretData",
			want:       false,
			wantErr:    zerrors.ThrowInternal(nil, "CRYPTO-Woox3", "Errors.Internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAES256GCMCrypto(tt.config, tt.keyStorage)
			require.NoError(t, err)
			got, gotErr := a.EncryptString([]byte(tt.value))
			require.ErrorIs(t, gotErr, tt.wantErr)
			if tt.want {
				assert.NotEmpty(t, got, "expected non-empty result")
			}
		})
	}
}

func TestAES256GCMCrypto_DecryptString(t *testing.T) {
	tests := []struct {
		name    string
		crypto  func(testing.TB) EncryptionAlgorithm
		payload func(*testing.T) []byte
		keyID   string
		want    string
		wantErr error
	}{
		{
			name:   "ok",
			crypto: newTestAES256GCMCrypto(nil),
			payload: func(t *testing.T) []byte {
				a := newTestAES256GCMCrypto(nil)(t)
				crypted, err := a.Encrypt([]byte("SecretData"))
				require.NoError(t, err)
				return crypted
			},
			keyID:   "keyID",
			want:    "SecretData",
			wantErr: nil,
		},
		{
			name:   "wrong key id",
			crypto: newTestAES256GCMCrypto(nil),
			payload: func(t *testing.T) []byte {
				a := newTestAES256GCMCrypto(nil)(t)
				crypted, err := a.Encrypt([]byte("SecretData"))
				require.NoError(t, err)
				return crypted
			},
			keyID:   "foo",
			wantErr: zerrors.ThrowNotFound(nil, "CRYPT-nkj1s", "unknown key id"),
		},
		{
			name:   "malformed encrypted value",
			crypto: newTestAES256GCMCrypto(nil),
			payload: func(t *testing.T) []byte {
				return []byte("malformed")
			},
			keyID:   "keyID",
			wantErr: zerrors.ThrowPreconditionFailed(nil, "CRYPT-ha6Oh", "malformed encypted value"),
		},
		{
			name:   "invalid value",
			crypto: newTestAES256GCMCrypto(nil),
			payload: func(t *testing.T) []byte {
				return []byte("eyJhbGciOiJBMjU2R0NNS1ciLCJlbmMiOiJBMjU2R0NNIiwiaXYiOiJfcUNRelpoVDF4bjJfUjJiIiwia2lkIjoia2V5SUQiLCJ0YWciOiJlcm1EV1oySE5mOGctMHRJa29zdHpRIn0.zGWGOxwPI8iq-ZSwmAqc1Ps8tltCp-g815nj7jY_m9Q.U_tpQDuz8SzrDEAD._YzsLNPi2BrmWA.G0I-gIfkWy6DNhqXPD_EXg")
			},
			keyID:   "keyID",
			wantErr: zerrors.ThrowUnauthenticated(nil, "CRYPT-OhN2u", "failed to decrypt value"),
		},
		{
			name:   "fallback decrypt",
			crypto: newTestAES256GCMCrypto(DecryptAES),
			payload: func(t *testing.T) []byte {
				a := newTestAESCrypto(t)
				crypted, err := a.Encrypt([]byte("SecretData"))
				require.NoError(t, err)
				return crypted
			},
			keyID: "keyID",
			want:  "SecretData",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.crypto(t)
			got, gotErr := a.DecryptString(tt.payload(t), tt.keyID)
			require.ErrorIs(t, gotErr, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
