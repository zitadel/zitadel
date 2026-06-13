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

// TestGroups_nestedGroupsRejected proves the flat group model is explicit:
// members of type Group (allowed by RFC 7643) fail loudly instead of being
// misread as user IDs.
func TestGroups_nestedGroupsRejected(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	otherGroupID := createScimGroup(t, orgID, integration.GroupName())

	tests := []struct {
		name         string
		memberType   string
		wantErrID    string
		wantScimType string
	}{
		{
			name:         "type Group",
			memberType:   "Group",
			wantErrID:    "SCIM-GRP6n",
			wantScimType: "invalidValue",
		},
		{
			name:         "unknown type",
			memberType:   "Robot",
			wantErrID:    "SCIM-GRP7t",
			wantScimType: "invalidValue",
		},
	}
	for _, tt := range tests {
		t.Run("create with member "+tt.name, func(t *testing.T) {
			_, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
				"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
				"displayName": %q,
				"members": [{"value": %q, "type": %q}]
			}`, integration.GroupName(), otherGroupID, tt.memberType)))
			require.Error(t, err)
			scimErr := new(scim.ScimError)
			require.ErrorAs(t, err, &scimErr)
			assert.Equal(t, "400", scimErr.Status)
			assert.Equal(t, tt.wantScimType, scimErr.ScimType)
			require.NotNil(t, scimErr.ZitadelDetail)
			assert.Equal(t, tt.wantErrID, scimErr.ZitadelDetail.ID)
		})
	}

	t.Run("replace with member type Group", func(t *testing.T) {
		groupName := integration.GroupName()
		groupID := createScimGroup(t, orgID, groupName)
		_, err := Instance.Client.SCIM.Groups.Replace(CTX, orgID, groupID, []byte(fmt.Sprintf(`{
			"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
			"displayName": %q,
			"members": [{"value": %q, "type": "Group"}]
		}`, groupName, otherGroupID)))
		require.Error(t, err)
		scimErr := new(scim.ScimError)
		require.ErrorAs(t, err, &scimErr)
		assert.Equal(t, "400", scimErr.Status)
		assert.Equal(t, "invalidValue", scimErr.ScimType)
	})

	t.Run("explicit type User is accepted", func(t *testing.T) {
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		created, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{
			"schemas": ["urn:ietf:params:scim:schemas:core:2.0:Group"],
			"displayName": %q,
			"members": [{"value": %q, "type": "User"}]
		}`, integration.GroupName(), user.GetUserId())))
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = Instance.Client.SCIM.Groups.Delete(CTX, orgID, created.ID)
		})
	})
}

// TestUsers_groupsAttribute proves User resources expose the read-only groups
// attribute of RFC 7643 section 4.1.2, with all memberships marked direct.
func TestUsers_groupsAttribute(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	groupName := integration.GroupName()
	groupID := createScimGroup(t, orgID, groupName)
	user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
	Instance.AddUsersToGroup(CTX, t, groupID, []string{user.GetUserId()})

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		fetched, err := Instance.Client.SCIM.Users.Get(CTX, orgID, user.GetUserId())
		require.NoError(ttt, err)
		require.Len(ttt, fetched.Groups, 1)
		group := fetched.Groups[0]
		assert.Equal(ttt, groupID, group.Value)
		assert.Equal(ttt, groupName, group.Display)
		assert.Equal(ttt, "direct", group.Type)
		assert.Contains(ttt, group.Ref, "/Groups/"+groupID)
	}, retryDuration, tick, "timeout waiting for groups attribute")

	// memberless users carry no groups attribute
	loner := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
	fetched, err := Instance.Client.SCIM.Users.Get(CTX, orgID, loner.GetUserId())
	require.NoError(t, err)
	assert.Empty(t, fetched.Groups)
}
