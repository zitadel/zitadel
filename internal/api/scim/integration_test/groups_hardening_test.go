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
