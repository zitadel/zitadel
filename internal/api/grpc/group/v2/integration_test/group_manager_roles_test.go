//go:build integration

package group_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

// TestServer_GroupManagerRoles_AuthorizationInAction proves that manager roles
// granted to a group take effect in real permission checks for its members,
// and stop applying as soon as the member leaves the group.
func TestServer_GroupManagerRoles_AuthorizationInAction(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	orgID := orgResp.GetOrganizationId()

	// a plain service account without any memberships, acting with its own token
	machine := instance.CreateUserTypeMachine(iamOwnerCtx, orgID)
	pat := instance.CreatePersonalAccessToken(iamOwnerCtx, machine.GetId())
	userCtx := integration.WithAuthorizationToken(CTX, pat.GetToken())

	// baseline: the user has no memberships and may not read org members
	baselineMemberships, err := instance.Client.Auth.ListMyMemberships(userCtx, &auth.ListMyMembershipsRequest{})
	require.NoError(t, err)
	require.Empty(t, baselineMemberships.GetResult())
	_, err = instance.Client.Mgmt.ListOrgMembers(userCtx, &management.ListOrgMembersRequest{})
	require.Error(t, err, "baseline must be forbidden")

	// grant ORG_OWNER_VIEWER to a group and add the user
	group := instance.CreateGroup(iamOwnerCtx, t, orgID, integration.GroupName())
	_, err = instance.Client.GroupV2.SetGroupManagerRoles(iamOwnerCtx, &group_v2.SetGroupManagerRolesRequest{
		GroupId: group.GetId(),
		Roles:   []string{"ORG_OWNER_VIEWER"},
	})
	require.NoError(t, err)
	instance.AddUsersToGroup(iamOwnerCtx, t, group.GetId(), []string{machine.GetId()})

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)

	// the membership resolves with the granted roles
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		memberships, err := instance.Client.Auth.ListMyMemberships(userCtx, &auth.ListMyMembershipsRequest{})
		require.NoError(ttt, err)
		require.Len(ttt, memberships.GetResult(), 1)
		membership := memberships.GetResult()[0]
		assert.Equal(ttt, orgID, membership.GetOrgId())
		assert.Contains(ttt, membership.GetRoles(), "ORG_OWNER_VIEWER")
	}, retryDuration, tick, "timeout waiting for group-supplied membership")

	// the user can now act with the role
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		members, err := instance.Client.Mgmt.ListOrgMembers(userCtx, &management.ListOrgMembersRequest{})
		require.NoError(ttt, err)
		assert.NotNil(ttt, members)
	}, retryDuration, tick, "timeout waiting for group-supplied permission")

	// leaving the group revokes the access
	instance.RemoveUsersFromGroup(iamOwnerCtx, t, group.GetId(), []string{machine.GetId()})
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		memberships, err := instance.Client.Auth.ListMyMemberships(userCtx, &auth.ListMyMembershipsRequest{})
		require.NoError(ttt, err)
		assert.Empty(ttt, memberships.GetResult())

		_, err = instance.Client.Mgmt.ListOrgMembers(userCtx, &management.ListOrgMembersRequest{})
		assert.Error(ttt, err, "access must be revoked after leaving the group")
	}, retryDuration, tick, "timeout waiting for revocation")
}
