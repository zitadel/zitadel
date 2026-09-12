//go:build integration

package oidc_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

// TestServer_UserInfo_GroupRevocation asserts the core revocation promise of the
// query-time design: removing the membership or the grant immediately stops the
// derived role and group claims, without new tokens being issued.
func TestServer_UserInfo_GroupRevocation(t *testing.T) {
	const roleViaGroup = "role-revoked-via-group"

	clientID, projectID := createClient(t, Instance)
	_, err := Instance.Client.Mgmt.BulkAddProjectRoles(CTX, &management.BulkAddProjectRolesRequest{
		ProjectId: projectID,
		Roles: []*management.BulkAddProjectRolesRequest_Role{
			{Key: roleViaGroup, DisplayName: roleViaGroup},
		},
	})
	require.NoError(t, err)

	groupName := integration.GroupName()
	group := Instance.CreateGroup(CTXIAM, t, Instance.DefaultOrg.GetId(), groupName)
	Instance.AddUsersToGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	t.Cleanup(func() {
		// the User is shared in this package; other tests assert its exact memberships
		Instance.RemoveUsersFromGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	})

	grant, err := Instance.Client.GroupV2.CreateGroupGrant(CTXIAM, &group_v2.CreateGroupGrantRequest{
		GroupId:   group.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{roleViaGroup},
	})
	require.NoError(t, err)

	scope := []string{
		oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
		oidc_api.ScopeProjectRolePrefix + roleViaGroup,
		oidc_api.ScopeUserGroups,
	}
	tokens := getTokens(t, clientID, scope)

	provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)

	roleClaimKey := fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, projectID)
	containsGroupClaim := func(claims map[string]any) bool {
		groupsClaim, ok := claims[oidc_api.ClaimUserGroups].([]any)
		if !ok {
			return false
		}
		for _, entry := range groupsClaim {
			if m, ok := entry.(map[string]any); ok && m["display"] == groupName {
				return true
			}
		}
		return false
	}
	containsRole := func(claims map[string]any) bool {
		roleClaim, ok := claims[roleClaimKey].(map[string]any)
		if !ok {
			return false
		}
		_, ok = roleClaim[roleViaGroup]
		return ok
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)

	// both the role and the group claim are present
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
		require.NoError(ttt, err)
		assert.True(ttt, containsRole(userinfo.Claims), "group-granted role expected, claims: %v", userinfo.Claims)
		assert.True(ttt, containsGroupClaim(userinfo.Claims), "group claim expected, claims: %v", userinfo.Claims)
	}, retryDuration, tick, "timeout waiting for derived claims")

	// deleting the grant revokes the role, the membership claim stays
	_, err = Instance.Client.GroupV2.DeleteGroupGrant(CTXIAM, &group_v2.DeleteGroupGrantRequest{Id: grant.GetId()})
	require.NoError(t, err)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
		require.NoError(ttt, err)
		assert.False(ttt, containsRole(userinfo.Claims), "role should be revoked, claims: %v", userinfo.Claims)
		assert.True(ttt, containsGroupClaim(userinfo.Claims), "group claim should remain, claims: %v", userinfo.Claims)
	}, retryDuration, tick, "timeout waiting for grant revocation")

	// removing the user from the group revokes the group claim
	Instance.RemoveUsersFromGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
		require.NoError(ttt, err)
		assert.False(ttt, containsGroupClaim(userinfo.Claims), "group claim should be revoked, claims: %v", userinfo.Claims)
	}, retryDuration, tick, "timeout waiting for membership revocation")
}
