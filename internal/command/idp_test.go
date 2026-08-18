package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	idp_repo "github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_checkAuthorizationParameters(t *testing.T) {
	tests := []struct {
		name       string
		parameters map[string]string
		want       map[string]string
		wantErr    bool
	}{
		{
			name:       "empty",
			parameters: nil,
			want:       nil,
		},
		{
			name:       "valid",
			parameters: map[string]string{"acr_values": "AAL2_OR_AAL3_ANY", "prompt": "login"},
			want:       map[string]string{"acr_values": "AAL2_OR_AAL3_ANY", "prompt": "login"},
		},
		{
			name:       "normalized",
			parameters: map[string]string{" ACR_Values ": " loa3 "},
			want:       map[string]string{"acr_values": "loa3"},
		},
		{
			name:       "empty value to remove a default",
			parameters: map[string]string{"prompt": ""},
			want:       map[string]string{"prompt": ""},
		},
		{
			name:       "reserved parameter",
			parameters: map[string]string{"redirect_uri": "https://evil.com"},
			wantErr:    true,
		},
		{
			name:       "empty name",
			parameters: map[string]string{"  ": "value"},
			wantErr:    true,
		},
		{
			name:       "name needing escaping",
			parameters: map[string]string{"acr values": "loa3"},
			wantErr:    true,
		},
		{
			name:       "duplicate after normalization",
			parameters: map[string]string{"prompt": "login", "PROMPT": "consent"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkAuthorizationParameters(tt.parameters)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, zerrors.IsErrorInvalidArgument(err))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_checkForwardedParameters(t *testing.T) {
	tests := []struct {
		name       string
		parameters []string
		want       []string
		wantErr    bool
	}{
		{
			name:       "empty",
			parameters: nil,
			want:       nil,
		},
		{
			name:       "sorted and deduplicated",
			parameters: []string{"max_age", "acr_values", "ACR_VALUES"},
			want:       []string{"acr_values", "max_age"},
		},
		{
			name:       "reserved parameter",
			parameters: []string{"acr_values", "state"},
			wantErr:    true,
		},
		{
			name:       "empty name",
			parameters: []string{""},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkForwardedParameters(tt.parameters)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, zerrors.IsErrorInvalidArgument(err))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOIDCIDPWriteModel_NewChanges_authorizationParameters(t *testing.T) {
	tests := []struct {
		name       string
		current    map[string]string
		forwarded  []string
		updated    map[string]string
		newForward []string
		wantChange bool
	}{
		{
			name:       "unchanged",
			current:    map[string]string{"acr_values": "loa3"},
			updated:    map[string]string{"acr_values": "loa3"},
			wantChange: false,
		},
		{
			name:       "nil and empty are equal",
			current:    nil,
			updated:    map[string]string{},
			wantChange: false,
		},
		{
			name:       "value changed",
			current:    map[string]string{"acr_values": "loa3"},
			updated:    map[string]string{"acr_values": "loa2"},
			wantChange: true,
		},
		{
			name:       "cleared",
			current:    map[string]string{"acr_values": "loa3"},
			updated:    nil,
			wantChange: true,
		},
		{
			name:       "forwarded parameters changed",
			forwarded:  []string{"acr_values"},
			newForward: []string{"acr_values", "max_age"},
			wantChange: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &OIDCIDPWriteModel{
				Name:                    "name",
				Issuer:                  "issuer",
				ClientID:                "clientID",
				AuthorizationParameters: tt.current,
				ForwardedParameters:     tt.forwarded,
			}
			changes, err := wm.NewChanges("name", "issuer", "clientID", "", nil, nil, false, false, tt.updated, tt.newForward, wm.Options)
			require.NoError(t, err)
			assert.Equal(t, tt.wantChange, len(changes) > 0)

			// applying the changes must yield the updated configuration
			event := &idp_repo.OIDCIDPChangedEvent{}
			for _, change := range changes {
				change(event)
			}
			wm.reduceChangedEvent(event)
			assert.Len(t, wm.AuthorizationParameters, len(tt.updated))
			for name, value := range tt.updated {
				assert.Equal(t, value, wm.AuthorizationParameters[name])
			}
			assert.Equal(t, tt.newForward, wm.ForwardedParameters)
		})
	}
}
