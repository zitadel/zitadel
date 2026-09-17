//go:build integration

package oidc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
)

// TestServer_IDToken_Groups asserts the groups claim is present on the ID token
// when scope=groups is requested. userInfoToOIDC does not gate ScopeUserGroups
// on userInfoAssertion, so groups must surface in both the ID token and userinfo.
func TestServer_IDToken_Groups(t *testing.T) {
	clientID, _ := createClient(t, Instance)

	groupName1, groupName2 := integration.GroupName(), integration.GroupName()
	group1 := Instance.CreateGroup(CTXIAM, t, Instance.DefaultOrg.GetId(), groupName1)
	group2 := Instance.CreateGroup(CTXIAM, t, Instance.DefaultOrg.GetId(), groupName2)
	Instance.AddUsersToGroup(CTXIAM, t, group1.GetId(), []string{User.GetUserId()})
	Instance.AddUsersToGroup(CTXIAM, t, group2.GetId(), []string{User.GetUserId()})
	t.Cleanup(func() {
		// the User is shared in this package; other tests assert its exact memberships
		Instance.RemoveUsersFromGroup(CTXIAM, t, group1.GetId(), []string{User.GetUserId()})
		Instance.RemoveUsersFromGroup(CTXIAM, t, group2.GetId(), []string{User.GetUserId()})
	})

	scope := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess, oidc_api.ScopeUserGroups}
	tokens := getTokens(t, clientID, scope)
	require.NotNil(t, tokens.IDTokenClaims)
	require.NotNil(t, tokens.IDTokenClaims.Claims)

	// groups claim encoding per RFC 9068 / SCIM (RFC 7643 section 4.1.2)
	expected := []map[string]interface{}{
		{"value": group1.GetId(), "display": groupName1},
		{"value": group2.GetId(), "display": groupName2},
	}
	assert.Subset(t, tokens.IDTokenClaims.Claims[oidc_api.ClaimUserGroups], expected)
}

// TestServer_IDToken_Groups_NoScope asserts that the groups claim is absent
// from the ID token when scope=groups is not requested, even if the user has
// group memberships.
func TestServer_IDToken_Groups_NoScope(t *testing.T) {
	clientID, _ := createClient(t, Instance)

	group := Instance.CreateGroup(CTXIAM, t, Instance.DefaultOrg.GetId(), integration.GroupName())
	Instance.AddUsersToGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	t.Cleanup(func() {
		Instance.RemoveUsersFromGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	})

	scope := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess}
	tokens := getTokens(t, clientID, scope)
	require.NotNil(t, tokens.IDTokenClaims)

	assert.NotContains(t, tokens.IDTokenClaims.Claims, oidc_api.ClaimUserGroups, "groups claim must be absent without the groups scope")
}
