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
