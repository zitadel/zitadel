//go:build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
)

func createScimGroup(t *testing.T, orgID, name string) string {
	t.Helper()
	created, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q
	}`, name)))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = Instance.Client.SCIM.Groups.Delete(CTX, orgID, created.ID)
	})
	return created.ID
}

// TestGroups_duplicateDisplayName proves uniqueness violations surface as
// 409 with the RFC 7644 "uniqueness" scimType, including across cases.
func TestGroups_duplicateDisplayName(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	groupName := integration.GroupName()
	createScimGroup(t, orgID, groupName)

	for _, duplicate := range []string{groupName, strings.ToUpper(groupName)} {
		_, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
			"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
			"displayName": %q
		}`, duplicate)))
		require.Errorf(t, err, "duplicate %q must be rejected", duplicate)
		scimErr := new(scim.ScimError)
		require.ErrorAs(t, err, &scimErr)
		assert.Equal(t, "409", scimErr.Status)
		assert.Equal(t, "uniqueness", scimErr.ScimType)
	}
}

// TestGroups_patchUnsupported proves PATCH fails with a proper SCIM error
// pointing the caller to PUT.
func TestGroups_patchUnsupported(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	groupID := createScimGroup(t, orgID, integration.GroupName())

	err := Instance.Client.SCIM.Groups.Update(CTX, orgID, groupID, []byte(`{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		"Operations": [{"op": "replace", "path": "displayName", "value": "new-name"}]
	}`))
	require.Error(t, err)
	scimErr := new(scim.ScimError)
	require.ErrorAs(t, err, &scimErr)
	assert.Equal(t, "501", scimErr.Status)
	require.NotNil(t, scimErr.ZitadelDetail)
	assert.Equal(t, "SCIM-GRP3p", scimErr.ZitadelDetail.ID)
}

// TestGroups_replaceUnknownMember proves a PUT referencing a user ID that
// does not exist surfaces as a SCIM error and leaves the group untouched —
// no dangling membership, no silent success.
func TestGroups_replaceUnknownMember(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	groupName := integration.GroupName()
	groupID := createScimGroup(t, orgID, groupName)

	// Replace reads the group from the projection; wait for it to appear
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
		require.NoError(ttt, err)
	}, retryDuration, tick, "timeout waiting for group projection")

	const unknownUserID = "00000000000000000000"
	_, err := Instance.Client.SCIM.Groups.Replace(CTX, orgID, groupID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q,
		"members": [{"value": %q}]
	}`, groupName, unknownUserID)))
	require.Error(t, err)
	scimErr := new(scim.ScimError)
	require.ErrorAs(t, err, &scimErr)
	assert.Equal(t, "400", scimErr.Status)
	require.NotNil(t, scimErr.ZitadelDetail)
	assert.Contains(t, scimErr.ZitadelDetail.Message, "Errors.User.NotFound")

	// the group must still exist with no members
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
		require.NoError(ttt, err)
		assert.Equal(ttt, groupName, group.DisplayName)
		assert.Empty(ttt, group.Members, "no membership must be created on a failed PUT")
	}, retryDuration, tick, "timeout asserting group is unchanged")
}

// TestGroups_replaceMemberAlreadyPresent proves a PUT listing a member that
// is already in the group is a no-op: no error, no duplicate membership in
// the projection — the replaceMembers set diff must filter it out.
func TestGroups_replaceMemberAlreadyPresent(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
	groupName := integration.GroupName()

	created, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q,
		"members": [{"value": %q}]
	}`, groupName, user.GetUserId())))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = Instance.Client.SCIM.Groups.Delete(CTX, orgID, created.ID)
	})

	// replaceMembers reads existingIDs from the projection; wait for it
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, created.ID)
		require.NoError(ttt, err)
		require.Len(ttt, group.Members, 1)
	}, retryDuration, tick, "timeout waiting for initial membership projection")

	// PUT with the same single member — the set diff must produce no work
	_, err = Instance.Client.SCIM.Groups.Replace(CTX, orgID, created.ID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q,
		"members": [{"value": %q}]
	}`, groupName, user.GetUserId())))
	require.NoError(t, err)

	// the projection must still report exactly one member, the same user;
	// even after the projection has had time to apply any stray event
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, created.ID)
		require.NoError(ttt, err)
		assert.Equal(ttt, groupName, group.DisplayName)
		require.Len(ttt, group.Members, 1, "Replace with an already-present member must not duplicate")
		assert.Equal(ttt, user.GetUserId(), group.Members[0].Value)
	}, retryDuration, tick, "timeout asserting idempotent PUT")
}

// TestGroups_listPagination proves count/startIndex paging returns every
// group exactly once with a stable total.
func TestGroups_listPagination(t *testing.T) {
	iamOwnerCtx := Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	org := Instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	orgID := org.GetOrganizationId()

	const numGroups = 5
	createdIDs := make(map[string]bool, numGroups)
	for i := 0; i < numGroups; i++ {
		created, err := Instance.Client.SCIM.Groups.Create(iamOwnerCtx, orgID, []byte(fmt.Sprintf(`{
			"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
			"displayName": %q
		}`, integration.GroupName())))
		require.NoError(t, err)
		createdIDs[created.ID] = true
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		list, err := Instance.Client.SCIM.Groups.List(iamOwnerCtx, orgID, &scim.ListRequest{})
		require.NoError(ttt, err)
		assert.Equal(ttt, numGroups, list.TotalResults)
	}, retryDuration, tick, "timeout waiting for groups projection")

	const pageSize = 2
	seen := make(map[string]int, numGroups)
	for startIndex := 1; startIndex <= numGroups; startIndex += pageSize {
		page, err := Instance.Client.SCIM.Groups.List(iamOwnerCtx, orgID, &scim.ListRequest{
			Count:      gu.Ptr(pageSize),
			StartIndex: gu.Ptr(startIndex),
		})
		require.NoError(t, err)
		assert.Equal(t, numGroups, page.TotalResults)
		assert.Equal(t, startIndex, page.StartIndex)
		assert.LessOrEqual(t, len(page.Resources), pageSize)
		for _, group := range page.Resources {
			seen[group.ID]++
		}
	}
	require.Len(t, seen, numGroups, "paging must return every group")
	for id, count := range seen {
		assert.Equalf(t, 1, count, "group %s must appear exactly once", id)
		assert.Truef(t, createdIDs[id], "unexpected group %s", id)
	}
}

// TestGroups_listFilterUnsupported proves filtering is rejected explicitly
// instead of silently returning unfiltered results.
func TestGroups_listFilterUnsupported(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	_, err := Instance.Client.SCIM.Groups.List(CTX, orgID, &scim.ListRequest{
		Filter: gu.Ptr(`displayName eq "foo"`),
	})
	require.Error(t, err)
	scimErr := new(scim.ScimError)
	require.ErrorAs(t, err, &scimErr)
	assert.Equal(t, "501", scimErr.Status)
}

// TestGroups_crossOrgIsolation proves a SCIM caller authenticated against
// org-A cannot read, replace, or delete a group owned by org-B even when
// the caller has group permissions on org-A. Locks down the ResourceOwner
// guard in getOrgGroup — without it, any caller with group.read on any
// org could fetch groups across the instance by ID.
func TestGroups_crossOrgIsolation(t *testing.T) {
	iamOwnerCtx := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	// group in the secondary org, created with elevated rights
	otherOrgID := SecondaryOrganization.OrganizationId
	otherGroupName := integration.GroupName()
	otherGroup, err := Instance.Client.SCIM.Groups.Create(iamOwnerCtx, otherOrgID, []byte(fmt.Sprintf(`{
		"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
		"displayName": %q
	}`, otherGroupName)))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = Instance.Client.SCIM.Groups.Delete(iamOwnerCtx, otherOrgID, otherGroup.ID)
	})

	// wait until the group is queryable in the secondary org
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := Instance.Client.SCIM.Groups.Get(iamOwnerCtx, otherOrgID, otherGroup.ID)
		require.NoError(ttt, err)
	}, retryDuration, tick, "timeout waiting for secondary-org group projection")

	// CTX is org-owner of the default org. Routing the request through the
	// default org's SCIM endpoint passes authz (the caller has group perms
	// there), so the ResourceOwner check in getOrgGroup is the only thing
	// that prevents the secondary-org group from being returned.
	defaultOrgID := Instance.DefaultOrg.GetId()

	// Pinning the ZitadelDetail.ID to "SCIM-GRP5o" ties each assertion to the
	// guard in getOrgGroup specifically. Without the pin, PUT/DELETE pass even
	// when the guard is dropped, because command-layer permission checks also
	// fail with 404 — masking the real regression.
	const getOrgGroupErrorID = "SCIM-GRP5o"

	t.Run("GET cross-org group returns 404", func(t *testing.T) {
		_, err := Instance.Client.SCIM.Groups.Get(CTX, defaultOrgID, otherGroup.ID)
		require.Error(t, err)
		scimErr := new(scim.ScimError)
		require.ErrorAs(t, err, &scimErr)
		assert.Equal(t, "404", scimErr.Status)
		require.NotNil(t, scimErr.ZitadelDetail)
		assert.Equal(t, getOrgGroupErrorID, scimErr.ZitadelDetail.ID)
	})

	t.Run("PUT cross-org group returns 404", func(t *testing.T) {
		_, err := Instance.Client.SCIM.Groups.Replace(CTX, defaultOrgID, otherGroup.ID, []byte(fmt.Sprintf(`{
			"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
			"displayName": %q
		}`, integration.GroupName())))
		require.Error(t, err)
		scimErr := new(scim.ScimError)
		require.ErrorAs(t, err, &scimErr)
		assert.Equal(t, "404", scimErr.Status)
		require.NotNil(t, scimErr.ZitadelDetail)
		assert.Equal(t, getOrgGroupErrorID, scimErr.ZitadelDetail.ID)
	})

	t.Run("DELETE cross-org group returns 404 and leaves group intact", func(t *testing.T) {
		err := Instance.Client.SCIM.Groups.Delete(CTX, defaultOrgID, otherGroup.ID)
		require.Error(t, err)
		scimErr := new(scim.ScimError)
		require.ErrorAs(t, err, &scimErr)
		assert.Equal(t, "404", scimErr.Status)
		require.NotNil(t, scimErr.ZitadelDetail)
		assert.Equal(t, getOrgGroupErrorID, scimErr.ZitadelDetail.ID)

		// group must still exist in the secondary org with its original name
		got, err := Instance.Client.SCIM.Groups.Get(iamOwnerCtx, otherOrgID, otherGroup.ID)
		require.NoError(t, err)
		assert.Equal(t, otherGroupName, got.DisplayName)
	})
}

// TestGroups_bulk proves groups participate in the /Bulk endpoint:
// create via bulkId, then delete by the returned location.
func TestGroups_bulk(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	groupName := integration.GroupName()

	createBody, err := json.Marshal(&scim.BulkRequest{
		Schemas: []schemas.ScimSchemaType{schemas.IdBulkRequest},
		Operations: []*scim.BulkRequestOperation{{
			Method: http.MethodPost,
			BulkID: "group-1",
			Path:   "/Groups",
			Data: json.RawMessage(fmt.Sprintf(`{
				"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
				"displayName": %q
			}`, groupName)),
		}},
	})
	require.NoError(t, err)

	createResp, err := Instance.Client.SCIM.Bulk(CTX, orgID, createBody)
	require.NoError(t, err)
	require.Len(t, createResp.Operations, 1)
	op := createResp.Operations[0]
	require.Nilf(t, op.Response, "bulk create must succeed, got %v", op.Response)
	assert.Equal(t, "201", op.Status)
	require.NotEmpty(t, op.Location)
	groupID := op.Location[strings.LastIndex(op.Location, "/")+1:]
	t.Cleanup(func() {
		_ = Instance.Client.SCIM.Groups.Delete(CTX, orgID, groupID)
	})

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		group, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
		require.NoError(ttt, err)
		assert.Equal(ttt, groupName, group.DisplayName)
	}, retryDuration, tick, "timeout waiting for bulk-created group")

	deleteBody, err := json.Marshal(&scim.BulkRequest{
		Schemas: []schemas.ScimSchemaType{schemas.IdBulkRequest},
		Operations: []*scim.BulkRequestOperation{{
			Method: http.MethodDelete,
			Path:   "/Groups/" + groupID,
		}},
	})
	require.NoError(t, err)

	deleteResp, err := Instance.Client.SCIM.Bulk(CTX, orgID, deleteBody)
	require.NoError(t, err)
	require.Len(t, deleteResp.Operations, 1)
	assert.Equal(t, "204", deleteResp.Operations[0].Status)

	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
		require.Error(ttt, err)
	}, retryDuration, tick, "timeout waiting for bulk-deleted group")
}
