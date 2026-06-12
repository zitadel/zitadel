//go:build integration

package group_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

// TestServer_AddUsersToGroup_Concurrent proves the OD-9 design: concurrent,
// overlapping member additions neither wedge the projection nor drop members.
func TestServer_AddUsersToGroup_Concurrent(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	group := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.GroupName())

	const numUsers = 8
	userIDs := make([]string, numUsers)
	for i := 0; i < numUsers; i++ {
		userIDs[i] = instance.CreateHumanUserVerified(iamOwnerCtx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone()).GetUserId()
	}

	// every call overlaps with all previous ones, racing duplicate `users.added`
	// events; same-aggregate sequence contention may reject a call after the
	// eventstore exhausts its retries (issue 7202), so the client retries like a
	// real caller would
	var wg sync.WaitGroup
	errs := make([]error, numUsers)
	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for attempt := 0; attempt < 5; attempt++ {
				_, errs[i] = instance.Client.GroupV2.AddUsersToGroup(iamOwnerCtx, &group_v2.AddUsersToGroupRequest{
					Id:      group.GetId(),
					UserIds: userIDs[:i+1],
				})
				if errs[i] == nil {
					return
				}
				time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		require.NoErrorf(t, err, "concurrent add %d", i)
	}

	// the projection must contain every member exactly once
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		resp, err := instance.Client.GroupV2.ListGroupUsers(iamOwnerCtx, &group_v2.ListGroupUsersRequest{
			Filters: []*group_v2.GroupUsersSearchFilter{{
				Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{group.GetId()}},
				},
			}},
		})
		require.NoError(ttt, err)
		gotUsers := make(map[string]int, len(resp.GetGroupUsers()))
		for _, member := range resp.GetGroupUsers() {
			gotUsers[member.GetUser().GetId()]++
		}
		require.Lenf(ttt, gotUsers, numUsers, "all members expected, got %v", gotUsers)
		for userID, count := range gotUsers {
			assert.Equalf(ttt, 1, count, "member %s duplicated", userID)
		}

		groupResp, err := instance.Client.GroupV2.GetGroup(iamOwnerCtx, &group_v2.GetGroupRequest{Id: group.GetId()})
		require.NoError(ttt, err)
		assert.Equal(ttt, uint64(numUsers), groupResp.GetGroup().GetUserCount())
	}, retryDuration, tick, "timeout waiting for complete membership")
}

// TestServer_UpdateGroup_caseInsensitiveRename proves renames cannot collide
// with an existing name in a different case.
func TestServer_UpdateGroup_caseInsensitiveRename(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	existingName := integration.GroupName()
	instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), existingName)
	other := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.GroupName())

	_, err := instance.Client.GroupV2.UpdateGroup(iamOwnerCtx, &group_v2.UpdateGroupRequest{
		Id:   other.GetId(),
		Name: gu.Ptr(strings.ToUpper(existingName)),
	})
	require.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

// TestServer_DeleteUser_removesMemberships proves a deleted user disappears
// from group listings and counts.
func TestServer_DeleteUser_removesMemberships(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	group := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.GroupName())
	member := instance.CreateHumanUserVerified(iamOwnerCtx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
	instance.AddUsersToGroup(iamOwnerCtx, t, group.GetId(), []string{member.GetUserId()})

	listMembers := func() (int, error) {
		resp, err := instance.Client.GroupV2.ListGroupUsers(iamOwnerCtx, &group_v2.ListGroupUsersRequest{
			Filters: []*group_v2.GroupUsersSearchFilter{{
				Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{group.GetId()}},
				},
			}},
		})
		if err != nil {
			return 0, err
		}
		return len(resp.GetGroupUsers()), nil
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		count, err := listMembers()
		require.NoError(ttt, err)
		assert.Equal(ttt, 1, count)
	}, retryDuration, tick, "timeout waiting for membership")

	_, err := instance.Client.UserV2.DeleteUser(iamOwnerCtx, &user.DeleteUserRequest{UserId: member.GetUserId()})
	require.NoError(t, err)

	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		count, err := listMembers()
		require.NoError(ttt, err)
		assert.Equal(ttt, 0, count)

		groupResp, err := instance.Client.GroupV2.GetGroup(iamOwnerCtx, &group_v2.GetGroupRequest{Id: group.GetId()})
		require.NoError(ttt, err)
		assert.Equal(ttt, uint64(0), groupResp.GetGroup().GetUserCount())
	}, retryDuration, tick, "timeout waiting for membership cleanup")
}

// TestServer_DeleteOrganization_withGroups proves an organization with groups,
// members, and grants can be removed without leaving group state behind.
func TestServer_DeleteOrganization_withGroups(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	orgID := orgResp.GetOrganizationId()
	group := instance.CreateGroup(iamOwnerCtx, t, orgID, integration.GroupName())
	member := instance.CreateHumanUserVerified(iamOwnerCtx, orgID, integration.Email(), integration.Phone())
	instance.AddUsersToGroup(iamOwnerCtx, t, group.GetId(), []string{member.GetUserId()})
	projectID := createProjectWithRole(iamOwnerCtx, t, orgID, "role1")
	_, err := instance.Client.GroupV2.CreateGroupGrant(iamOwnerCtx, &group_v2.CreateGroupGrantRequest{
		GroupId:   group.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{"role1"},
	})
	require.NoError(t, err)

	_, err = instance.Client.OrgV2.DeleteOrganization(iamOwnerCtx, &org.DeleteOrganizationRequest{
		OrganizationId: orgID,
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		resp, err := instance.Client.GroupV2.ListGroups(iamOwnerCtx, &group_v2.ListGroupsRequest{
			Filters: []*group_v2.GroupsSearchFilter{{
				Filter: &group_v2.GroupsSearchFilter_OrganizationId{
					OrganizationId: &filter.IDFilter{Id: orgID},
				},
			}},
		})
		require.NoError(ttt, err)
		assert.Empty(ttt, resp.GetGroups(), "groups of the removed organization must disappear")
	}, retryDuration, tick, "timeout waiting for org group cleanup")
}
