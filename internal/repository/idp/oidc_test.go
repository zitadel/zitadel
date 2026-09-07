package idp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOIDCIDPChangedEvent_authorizationParameters makes sure resetting the parameters survives a
// marshal/unmarshal round trip. A nil map would be omitted from the payload and therefore be
// treated as "unchanged" when the event is replayed.
func TestOIDCIDPChangedEvent_authorizationParameters(t *testing.T) {
	tests := []struct {
		name          string
		changes       []OIDCIDPChanges
		wantParams    *map[string]string
		wantForwarded *[]string
	}{
		{
			name:    "no changes",
			changes: nil,
		},
		{
			name:          "set",
			changes:       []OIDCIDPChanges{ChangeOIDCAuthorizationParameters(map[string]string{"acr_values": "loa3"}), ChangeOIDCForwardedParameters([]string{"max_age"})},
			wantParams:    &map[string]string{"acr_values": "loa3"},
			wantForwarded: &[]string{"max_age"},
		},
		{
			name:          "reset with an empty value",
			changes:       []OIDCIDPChanges{ChangeOIDCAuthorizationParameters(map[string]string{}), ChangeOIDCForwardedParameters([]string{})},
			wantParams:    &map[string]string{},
			wantForwarded: &[]string{},
		},
		{
			name:          "reset with nil",
			changes:       []OIDCIDPChanges{ChangeOIDCAuthorizationParameters(nil), ChangeOIDCForwardedParameters(nil)},
			wantParams:    &map[string]string{},
			wantForwarded: &[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &OIDCIDPChangedEvent{ID: "idp-id"}
			for _, change := range tt.changes {
				change(event)
			}

			payload, err := json.Marshal(event)
			require.NoError(t, err)

			replayed := &OIDCIDPChangedEvent{}
			require.NoError(t, json.Unmarshal(payload, replayed))

			assert.Equal(t, tt.wantParams, replayed.AuthorizationParameters)
			assert.Equal(t, tt.wantForwarded, replayed.ForwardedParameters)
		})
	}
}

// TestOAuthIDPChangedEvent_authorizationParameters mirrors
// [TestOIDCIDPChangedEvent_authorizationParameters] for the generic OAuth 2.0 provider.
func TestOAuthIDPChangedEvent_authorizationParameters(t *testing.T) {
	event := &OAuthIDPChangedEvent{ID: "idp-id"}
	ChangeOAuthAuthorizationParameters(nil)(event)
	ChangeOAuthForwardedParameters(nil)(event)

	payload, err := json.Marshal(event)
	require.NoError(t, err)

	replayed := &OAuthIDPChangedEvent{}
	require.NoError(t, json.Unmarshal(payload, replayed))

	assert.Equal(t, &map[string]string{}, replayed.AuthorizationParameters)
	assert.Equal(t, &[]string{}, replayed.ForwardedParameters)
}
