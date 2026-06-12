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

// TestServer_UserInfo_GroupGrantRoles asserts that project roles granted to a group
// surface as role claims for its members, merged with their personal grants.
func TestServer_UserInfo_GroupGrantRoles(t *testing.T) {
	const roleViaGroup = "role-via-group"

	clientID, projectID := createClient(t, Instance)

	// the role exists on the project, but is granted to a group instead of the user
	_, err := Instance.Client.Mgmt.BulkAddProjectRoles(CTX, &management.BulkAddProjectRolesRequest{
		ProjectId: projectID,
		Roles: []*management.BulkAddProjectRolesRequest_Role{
			{Key: roleViaGroup, DisplayName: roleViaGroup},
		},
	})
	require.NoError(t, err)

	group := Instance.CreateGroup(CTXIAM, t, Instance.DefaultOrg.GetId(), integration.GroupName())
	Instance.AddUsersToGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	t.Cleanup(func() {
		// the User is shared in this package; other tests assert its exact memberships
		Instance.RemoveUsersFromGroup(CTXIAM, t, group.GetId(), []string{User.GetUserId()})
	})

	_, err = Instance.Client.GroupV2.CreateGroupGrant(CTXIAM, &group_v2.CreateGroupGrantRequest{
		GroupId:   group.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{roleViaGroup},
	})
	require.NoError(t, err)

	scope := []string{
		oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
		oidc_api.ScopeProjectRolePrefix + roleViaGroup,
	}
	tokens := getTokens(t, clientID, scope)

	provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
		require.NoError(ttt, err)

		roleClaim, ok := userinfo.Claims[fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, projectID)].(map[string]any)
		require.Truef(ttt, ok, "project role claim not found, claims: %v", userinfo.Claims)

		role, ok := roleClaim[roleViaGroup].(map[string]any)
		require.Truef(ttt, ok, "group-granted role not asserted, roles: %v", roleClaim)
		assert.Contains(ttt, role, Instance.DefaultOrg.GetId())
	}, retryDuration, tick, "timeout waiting for group-granted role in userinfo")
}
