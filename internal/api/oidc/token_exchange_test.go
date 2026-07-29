package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
)

func testTokenExchangeClient() *Client {
	return &Client{}
}

func Test_validateTokenExchangeScopes(t *testing.T) {
	client := testTokenExchangeClient()
	actorScopes := []string{
		oidc.ScopeOpenID,
		oidc.ScopeProfile,
		ScopeProjectsRoles,
		ScopeProjectRolePrefix + "TOKEN_EXCHANGE_IMPERSONATOR",
		domain.ProjectIDScope + "288694447250664408" + domain.AudSuffix,
	}

	t.Run("exchange: subset of subject", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeOfflineAccess}
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, oidc.ScopeProfile},
			subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, oidc.ScopeProfile}, got)
	})

	t.Run("exchange: escalation rejected", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID}
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
			subjectScopes, subjectScopes, false)
		require.Error(t, err)
	})

	t.Run("impersonation: subject-data scope not on actor allowed", func(t *testing.T) {
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
			nil, actorScopes, true)
		require.NoError(t, err)
		assert.Contains(t, got, oidc.ScopeEmail)
	})

	t.Run("impersonation: authorization scope not on actor rejected", func(t *testing.T) {
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
			nil, actorScopes, true)
		require.Error(t, err)
	})

	t.Run("impersonation: authorization scope on actor allowed", func(t *testing.T) {
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, ScopeProjectsRoles, domain.ProjectIDScope + "288694447250664408" + domain.AudSuffix},
			nil, actorScopes, true)
		require.NoError(t, err)
		assert.Contains(t, got, ScopeProjectsRoles)
	})

	t.Run("impersonation: omitted scope inherits actor", func(t *testing.T) {
		got, err := validateTokenExchangeScopes(client, nil, nil, actorScopes, true)
		require.NoError(t, err)
		assert.Equal(t, []string{
			oidc.ScopeOpenID,
			oidc.ScopeProfile,
			ScopeProjectsRoles,
			domain.ProjectIDScope + "288694447250664408" + domain.AudSuffix,
		}, got)
	})

	t.Run("actor path with subject scopes uses exchange validator", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile}
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeEmail},
			subjectScopes, actorScopes, false)
		require.Error(t, err)
	})

	t.Run("actor path with JWT subject and empty scopes uses exchange validator", func(t *testing.T) {
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeEmail},
			nil, actorScopes, false)
		require.Error(t, err)
	})
}

func Test_isTokenExchangeAuthorizationScope(t *testing.T) {
	tests := []struct {
		scope string
		want  bool
	}{
		{oidc.ScopeOpenID, false},
		{oidc.ScopeEmail, false},
		{oidc.ScopeOfflineAccess, true},
		{ScopeProjectsRoles, true},
		{ScopeProjectRolePrefix + "admin", true},
		{domain.ProjectIDScope + "proj" + domain.AudSuffix, true},
		{ScopeUserMetaData, false},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			assert.Equal(t, tt.want, isTokenExchangeAuthorizationScope(tt.scope))
		})
	}
}
