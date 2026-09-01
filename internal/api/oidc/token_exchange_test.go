package oidc

import (
	"slices"
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

	t.Run("exchange: OrgRoleIDScope further downscoping allowed", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		subjectScopes := []string{oidc.ScopeOpenID, orgA, orgB}
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgA},
			subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgA}, got)
	})

	t.Run("exchange: OrgRoleIDScope widening rejected", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		subjectScopes := []string{oidc.ScopeOpenID, orgA}
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgB},
			subjectScopes, subjectScopes, false)
		require.Error(t, err)
	})

	t.Run("exchange: omitting OrgRoleIDScope inherits subject filter", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile, orgA}
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID},
			subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgA}, got)
	})

	t.Run("actor path: subject filter is ceiling, actor filter ignored", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile, orgA}
		actorWithOtherFilter := append(slices.Clone(actorScopes), orgB)

		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgB},
			subjectScopes, actorWithOtherFilter, false)
		require.Error(t, err)

		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgA},
			subjectScopes, actorWithOtherFilter, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgA}, got)
	})

	t.Run("actor path: unfiltered subject ignores actor OrgRoleIDScope ceiling", func(t *testing.T) {
		// Subject-primary: only the subject's filter is a ceiling. An unfiltered
		// subject may still request any roles:id even if the actor is filtered.
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile}
		actorWithFilter := append(slices.Clone(actorScopes), orgA)
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgB},
			subjectScopes, actorWithFilter, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgB}, got)
	})

	t.Run("union path: empty-scope subject ignores actor OrgRoleIDScope ceiling", func(t *testing.T) {
		// Access/JWT subjects with an empty scope claim still use the union path.
		// An empty subject is unfiltered, so actor roles:id must not become a ceiling.
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		actorWithFilter := append(slices.Clone(actorScopes), orgB)
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgA},
			nil, actorWithFilter, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgA}, got)
	})

	t.Run("union path: empty-scope subject omit does not inherit actor filter", func(t *testing.T) {
		orgB := domain.OrgRoleIDScope + "orgB"
		actorWithFilter := append(slices.Clone(actorScopes), orgB)
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID},
			nil, actorWithFilter, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID}, got)
		assert.NotContains(t, got, orgB)
	})

	t.Run("exchange: bare OrgRoleIDScope request rejected", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile}
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, domain.OrgRoleIDScope},
			subjectScopes, subjectScopes, false)
		require.Error(t, err)
	})

	t.Run("impersonation: bare OrgRoleIDScope request rejected", func(t *testing.T) {
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, domain.OrgRoleIDScope},
			nil, actorScopes, true)
		require.Error(t, err)
	})

	t.Run("exchange: bare OrgRoleIDScope on subject remains a ceiling", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		subjectScopes := []string{oidc.ScopeOpenID, domain.OrgRoleIDScope}
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgA},
			subjectScopes, subjectScopes, false)
		require.Error(t, err)
	})

	t.Run("exchange: omit inherits bare OrgRoleIDScope ceiling", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile, domain.OrgRoleIDScope}
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID},
			subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, domain.OrgRoleIDScope}, got)
	})

	t.Run("exchange: empty request inherits bare OrgRoleIDScope from subject", func(t *testing.T) {
		subjectScopes := []string{oidc.ScopeOpenID, domain.OrgRoleIDScope}
		got, err := validateTokenExchangeScopes(client, nil, subjectScopes, subjectScopes, false)
		require.NoError(t, err)
		assert.Equal(t, subjectScopes, got)
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

	t.Run("impersonation: OrgRoleIDScope widening beyond actor rejected", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		orgB := domain.OrgRoleIDScope + "orgB"
		actorWithFilter := append(slices.Clone(actorScopes), orgA)
		_, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID, orgB},
			nil, actorWithFilter, true)
		require.Error(t, err)
	})

	t.Run("impersonation: omitting OrgRoleIDScope inherits actor filter", func(t *testing.T) {
		orgA := domain.OrgRoleIDScope + "orgA"
		actorWithFilter := append(slices.Clone(actorScopes), orgA)
		got, err := validateTokenExchangeScopes(client,
			[]string{oidc.ScopeOpenID},
			nil, actorWithFilter, true)
		require.NoError(t, err)
		assert.Equal(t, []string{oidc.ScopeOpenID, orgA}, got)
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
