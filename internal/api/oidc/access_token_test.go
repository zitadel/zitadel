package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_splitDecryptedOpaqueToken(t *testing.T) {
	tests := []struct {
		name        string
		decrypted   string
		wantTokenID string
		wantSubject string
		wantErr     bool
	}{
		{
			name:        "user ID simple sans deux-points",
			decrypted:   "tokenID123:userID456",
			wantTokenID: "tokenID123",
			wantSubject: "userID456",
		},
		{
			name:        "user ID au format URN avec deux-points",
			decrypted:   "sessionID-at_abc:urn:myorg:user:550e8400-e29b-41d4-a716-446655440000",
			wantTokenID: "sessionID-at_abc",
			wantSubject: "urn:myorg:user:550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "token sans séparateur",
			decrypted: "invalidtoken",
			wantErr:   true,
		},
		{
			name:      "chaîne vide",
			decrypted: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenID, gotSubject, err := splitDecryptedOpaqueToken(tt.decrypted)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantTokenID, gotTokenID)
			assert.Equal(t, tt.wantSubject, gotSubject)
		})
	}
}
