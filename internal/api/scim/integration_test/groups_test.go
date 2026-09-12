//go:build integration

package integration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
)

func TestGroups_lifecycle(t *testing.T) {
	groupName := integration.GroupName()
	renamedGroupName := integration.GroupName()
	orgID := Instance.DefaultOrg.GetId()

	user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())

	// create
	created, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q,
		"members": [{"value": %q}]
	}`, groupName, user.GetUserId())))
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	assert.Equal(t, groupName, created.DisplayName)

	// get returns the member
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, created.ID)
		require.NoError(ttt, err)
		assert.Equal(ttt, groupName, group.DisplayName)
		require.Len(ttt, group.Members, 1)
		assert.Equal(ttt, user.GetUserId(), group.Members[0].Value)
	}, retryDuration, tick, "timeout waiting for group with member")

	// list contains the group
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		list, err := Instance.Client.SCIM.Groups.List(CTX, orgID, &scim.ListRequest{})
		require.NoError(ttt, err)
		found := false
		for _, group := range list.Resources {
			if group.ID == created.ID {
				found = true
			}
		}
		assert.True(ttt, found, "created group should be listed")
	}, retryDuration, tick, "timeout waiting for group listing")

	// replace renames and removes the member
	_, err = Instance.Client.SCIM.Groups.Replace(CTX, orgID, created.ID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q,
		"members": []
	}`, renamedGroupName)))
	require.NoError(t, err)
	// the projection applies the change eventually
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, created.ID)
		require.NoError(ttt, err)
		assert.Equal(ttt, renamedGroupName, group.DisplayName)
		assert.Empty(ttt, group.Members)
	}, retryDuration, tick, "timeout waiting for replaced group")

	// delete; a second delete returns not found per RFC 7644
	require.NoError(t, Instance.Client.SCIM.Groups.Delete(CTX, orgID, created.ID))
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		err := Instance.Client.SCIM.Groups.Delete(CTX, orgID, created.ID)
		require.Error(ttt, err)
	}, retryDuration, tick, "timeout waiting for group deletion")
}
