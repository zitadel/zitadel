package eventstore

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAuthAlgorithm struct {
	decryptResult string
	decryptErr    error
}

func (m *mockAuthAlgorithm) EncryptToken(data string) (string, error) {
	return data, nil
}

func (m *mockAuthAlgorithm) DecryptToken(token string) (string, error) {
	return m.decryptResult, m.decryptErr
}

func (m *mockAuthAlgorithm) LegacyTokenEnabled() bool {
	return false
}

func TestTokenVerifierRepo_getTokenIDAndSubject(t *testing.T) {
	tests := []struct {
		name          string
		decryptResult string
		decryptErr    error
		wantTokenID   string
		wantSubject   string
		wantValid     bool
	}{
		{
			name:          "simple user ID without colons",
			decryptResult: "tokenID123:userID456",
			wantTokenID:   "tokenID123",
			wantSubject:   "userID456",
			wantValid:     true,
		},
		{
			name:          "user ID in URN format with colons",
			decryptResult: "sessionID-at_abc:urn:myorg:user:550e8400-e29b-41d4-a716-446655440000",
			wantTokenID:   "sessionID-at_abc",
			wantSubject:   "urn:myorg:user:550e8400-e29b-41d4-a716-446655440000",
			wantValid:     true,
		},
		{
			name:          "token without separator",
			decryptResult: "invalidtoken",
			wantValid:     false,
		},
		{
			name:       "decryption error",
			decryptErr: errors.New("decrypt failed"),
			wantValid:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TokenVerifierRepo{
				AuthAlgorithm: &mockAuthAlgorithm{
					decryptResult: tt.decryptResult,
					decryptErr:    tt.decryptErr,
				},
			}
			gotTokenID, gotSubject, gotValid := repo.getTokenIDAndSubject(context.Background(), "anytoken")
			assert.Equal(t, tt.wantTokenID, gotTokenID)
			assert.Equal(t, tt.wantSubject, gotSubject)
			assert.Equal(t, tt.wantValid, gotValid)
		})
	}
}
