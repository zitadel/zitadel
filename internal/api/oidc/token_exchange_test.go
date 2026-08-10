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

	t.Run("exchange: OrgRoleIDScope downscoping allowed", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile}
		orgRoleScope := domain.OrgRoleIDScope + "388047065096336384"
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgRoleScope},
			subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgRoleScope}, got)
	})

	t.Run("impersonation: subject-data scope not on actor allowed", func(t *testing.T) {
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
			nil, actorScopes, true)
		require.NoError(t, err)
		assert.Contains(t, got, oidc.ScopeEmail)
	})

	t.Run("impersonation: OrgRoleIDScope not on actor allowed", func(t *testing.T) {
		orgRoleScope := domain.OrgRoleIDScope + "388047065096336384"
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgRoleScope},
			nil, actorScopes, true)
		require.NoError(t, err)
		assert.Contains(t, got, orgRoleScope)
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
		{domain.OrgRoleIDScope + "388047065096336384", false},
		{ScopeUserMetaData, false},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			assert.Equal(t, tt.want, isTokenExchangeAuthorizationScope(tt.scope))
		})
	}
}

func Test_isTokenExchangeRestrictionScope(t *testing.T) {
	tests := []struct {
		scope string
		want  bool
	}{
		{domain.OrgRoleIDScope + "388047065096336384", true},
		{domain.OrgRoleIDScope, true},
		{oidc.ScopeOpenID, false},
		{oidc.ScopeOfflineAccess, false},
		{ScopeProjectsRoles, false},
		{domain.OrgIDScope + "388047065096336384", false},
		{ScopeProjectRolePrefix + "admin", false},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			assert.Equal(t, tt.want, isTokenExchangeRestrictionScope(tt.scope))
		})
	}
}
