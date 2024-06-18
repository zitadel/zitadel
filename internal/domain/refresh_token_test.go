package domain

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type mockKeyStorage struct {
	keys crypto.Keys
}

func (s *mockKeyStorage) ReadKeys() (crypto.Keys, error) {
	return s.keys, nil
}

func (s *mockKeyStorage) ReadKey(id string) (*crypto.Key, error) {
	return &crypto.Key{
		ID:    id,
		Value: s.keys[id],
	}, nil
}

func (*mockKeyStorage) CreateKeys(_ context.Context, keys ...*crypto.Key) error {
	return errors.New("mockKeyStorage.CreateKeys not implemented")
}

func TestFromRefreshToken(t *testing.T) {
	const (
		userID  = "userID"
		tokenID = "tokenID"
	)

	keyConfig := &crypto.KeyConfig{
		EncryptionKeyID:  "keyID",
		DecryptionKeyIDs: []string{"keyID"},
	}
	keys := crypto.Keys{"keyID": "ThisKeyNeedsToHave32Characters!!"}
	algorithm, err := crypto.NewAESCrypto(keyConfig, &mockKeyStorage{keys: keys})
	require.NoError(t, err)

	refreshToken, err := NewRefreshToken(userID, tokenID, algorithm)
	require.NoError(t, err)

	invalidRefreshToken, err := algorithm.Encrypt([]byte(userID + ":" + tokenID))
	require.NoError(t, err)

	type args struct {
		refreshToken string
		algorithm    crypto.EncryptionAlgorithm
	}
	tests := []struct {
		name        string
		args        args
		wantUserID  string
		wantTokenID string
		wantToken   string
		wantErr     error
	}{
		{
			name:    "invalid base64",
			args:    args{"~~~", algorithm},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-BGDhn", "Errors.User.RefreshToken.Invalid"),
		},
		{
			name:    "short cipher text",
			args:    args{"DEADBEEF", algorithm},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "CRYPT-23kH1", "cipher text block too short"),
		},
		{
			name:    "incorrect amount of segments",
			args:    args{base64.RawURLEncoding.EncodeToString(invalidRefreshToken), algorithm},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-BGDhn", "Errors.User.RefreshToken.Invalid"),
		},
		{
			name:        "success",
			args:        args{refreshToken, algorithm},
			wantUserID:  userID,
			wantTokenID: tokenID,
			wantToken:   tokenID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, gotTokenID, gotToken, err := FromRefreshToken(tt.args.refreshToken, tt.args.algorithm)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantUserID, gotUserID)
			assert.Equal(t, tt.wantTokenID, gotTokenID)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}
