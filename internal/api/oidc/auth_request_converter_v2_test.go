package oidc

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/command"
)

func TestAuthRequestV2_LogValue(t *testing.T) {
	tests := []struct {
		name     string
		req      *AuthRequestV2
		expected slog.Value
	}{
		{
			name:     "nil receiver",
			req:      nil,
			expected: slog.Value{},
		},
		{
			name:     "nil inner request",
			req:      &AuthRequestV2{},
			expected: slog.Value{},
		},
		{
			name: "populated auth request v2",
			req: &AuthRequestV2{
				CurrentAuthRequest: &command.CurrentAuthRequest{
					AuthRequest: &command.AuthRequest{
						ID:             "auth-v2-123",
						ClientID:       "client-789",
						Issuer:         "https://issuer.example.com",
						OrganizationID: "org-101",
						Scope:          []string{"openid", "profile"},
						RedirectURI:    "https://example.com/callback",
					},
				},
			},
			expected: slog.GroupValue(
				slog.String("id", "auth-v2-123"),
				slog.String("client_id", "client-789"),
				slog.String("issuer", "https://issuer.example.com"),
				slog.String("org_id", "org-101"),
				slog.Any("scopes", []string{"openid", "profile"}),
				slog.String("redirect_uri", "https://example.com/callback"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.LogValue()
			assert.Equal(t, tt.expected, got)
		})
	}
}
