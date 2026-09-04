//go:build integration

package group_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// TestServer_GroupScale_100Members proves counts and paging stay correct on a
// group two orders of magnitude larger than the other tests use.
func TestServer_GroupScale_100Members(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	orgID := orgResp.GetOrganizationId()
	group := instance.CreateGroup(iamOwnerCtx, t, orgID, integration.GroupName())

	const numUsers = 100
	userIDs := make([]string, numUsers)
	var eg errgroup.Group
	eg.SetLimit(10)
	for i := 0; i < numUsers; i++ {
		eg.Go(func() error {
			userIDs[i] = instance.CreateHumanUserVerified(iamOwnerCtx, orgID, integration.Email(), integration.Phone()).GetUserId()
			return nil
		})
	}
	require.NoError(t, eg.Wait())

	const batchSize = 25
	for start := 0; start < numUsers; start += batchSize {
		instance.AddUsersToGroup(iamOwnerCtx, t, group.GetId(), userIDs[start:start+batchSize])
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		groupResp, err := instance.Client.GroupV2.GetGroup(iamOwnerCtx, &group_v2.GetGroupRequest{Id: group.GetId()})
		require.NoError(ttt, err)
		assert.Equal(ttt, uint64(numUsers), groupResp.GetGroup().GetUserCount())
	}, retryDuration, tick, "timeout waiting for full membership")

	const pageSize = 30
	seen := make(map[string]int, numUsers)
	var pages int
	for offset := uint64(0); ; offset += pageSize {
		start := time.Now()
		page, err := instance.Client.GroupV2.ListGroupUsers(iamOwnerCtx, &group_v2.ListGroupUsersRequest{
			Filters: []*group_v2.GroupUsersSearchFilter{{
				Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{group.GetId()}},
				},
			}},
			Pagination: &filter.PaginationRequest{
				Offset: offset,
				Limit:  pageSize,
				Asc:    true,
			},
		})
		require.NoError(t, err)
		t.Logf("page %d: %d users in %s", pages, len(page.GetGroupUsers()), time.Since(start))
		assert.Equal(t, uint64(numUsers), page.GetPagination().GetTotalResult())
		require.LessOrEqual(t, len(page.GetGroupUsers()), pageSize)
		for _, member := range page.GetGroupUsers() {
			seen[member.GetUser().GetId()]++
		}
		pages++
		if len(page.GetGroupUsers()) < pageSize {
			break
		}
	}
	assert.Equal(t, 4, pages)
	require.Len(t, seen, numUsers, "paging must return every member")
	for userID, count := range seen {
		assert.Equalf(t, 1, count, "member %s must appear exactly once", userID)
	}

	// removal at scale keeps the count consistent
	_, err := instance.Client.GroupV2.RemoveUsersFromGroup(iamOwnerCtx, &group_v2.RemoveUsersFromGroupRequest{
		Id:      group.GetId(),
		UserIds: userIDs[:batchSize],
	})
	require.NoError(t, err)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		groupResp, err := instance.Client.GroupV2.GetGroup(iamOwnerCtx, &group_v2.GetGroupRequest{Id: group.GetId()})
		require.NoError(ttt, err)
		assert.Equal(ttt, uint64(numUsers-batchSize), groupResp.GetGroup().GetUserCount())
	}, retryDuration, tick, "timeout waiting for membership removal")
}
