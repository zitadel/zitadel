//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
)

func TestGroups_externalID(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	externalID := integration.GroupName()
	name := integration.GroupName()
	create := func(ctx context.Context, org, name, id string) (string, error) {
		result, err := Instance.Client.SCIM.Groups.Create(ctx, org, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q,"externalId":%q}`, name, id)))
		if err != nil {
			return "", err
		}
		assert.Equal(t, id, result.ExternalID)
		return result.ID, nil
	}
	groupID, err := create(CTX, orgID, name, externalID)
	require.NoError(t, err)
	wait, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	lookup := func(ctx context.Context, org, id, wantID, wantName string) {
		t.Helper()
		require.EventuallyWithT(t, func(ct *assert.CollectT) {
			result, err := Instance.Client.SCIM.Groups.List(ctx, org, &scim.ListRequest{Filter: gu.Ptr(fmt.Sprintf(`externalId eq %q`, id)), ExcludedAttributes: "members"})
			require.NoError(ct, err)
			if wantID == "" {
				assert.Empty(ct, result.Resources)
				return
			}
			require.Len(ct, result.Resources, 1)
			assert.Equal(ct, wantID, result.Resources[0].ID)
			assert.Equal(ct, wantName, result.Resources[0].DisplayName)
			assert.Equal(ct, id, result.Resources[0].ExternalID)
		}, wait, tick)
	}
	lookup(CTX, orgID, externalID, groupID, name)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		got, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
		require.NoError(ct, err)
		assert.Equal(ct, externalID, got.ExternalID)
	}, wait, tick)
	_, err = create(CTX, orgID, integration.GroupName(), externalID)
	require.Error(t, err)
	scimErr := new(scim.ScimError)
	require.ErrorAs(t, err, &scimErr)
	assert.Equal(t, "409", scimErr.Status)
	patch := func(ops string) error {
		return Instance.Client.SCIM.Groups.Update(CTX, orgID, groupID, []byte(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":`+ops+`}`))
	}
	renamed := integration.GroupName()
	require.NoError(t, patch(fmt.Sprintf(`[{"op":"replace","path":"displayName","value":%q}]`, renamed)))
	lookup(CTX, orgID, externalID, groupID, renamed)
	replacement := externalID + "-replacement"
	require.NoError(t, patch(fmt.Sprintf(`[{"op":"replace","path":"externalId","value":%q}]`, replacement)))
	lookup(CTX, orgID, replacement, groupID, renamed)
	lookup(CTX, orgID, externalID, "", "")
	// Invalid membership must roll back both identity and display-name changes.
	require.Error(t, patch(fmt.Sprintf(`[{"op":"replace","value":{"externalId":%q,"displayName":"invalid rename"}},{"op":"add","path":"members","value":[{"value":"missing-user"}]}]`, externalID)))
	lookup(CTX, orgID, replacement, groupID, renamed)
	// Removing a binding releases its uniqueness constraint for another group.
	require.NoError(t, patch(`[{"op":"remove","path":"externalId"}]`))
	lookup(CTX, orgID, replacement, "", "")
	reused, err := create(CTX, orgID, integration.GroupName(), replacement)
	require.NoError(t, err)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		_, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, reused)
		require.NoError(ct, err)
	}, wait, tick)
	require.NoError(t, Instance.Client.SCIM.Groups.Delete(CTX, orgID, reused))
	_, err = create(CTX, orgID, integration.GroupName(), replacement)
	require.NoError(t, err)
	// PUT changes the source binding atomically, and omitting externalId clears it.
	replaced, err := Instance.Client.SCIM.Groups.Replace(CTX, orgID, groupID, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q,"externalId":%q}`, renamed, externalID)))
	require.NoError(t, err)
	assert.Equal(t, externalID, replaced.ExternalID)
	lookup(CTX, orgID, externalID, groupID, renamed)
	err = patch(fmt.Sprintf(`[{"op":"replace","value":{"externalId":%q,"displayName":%q}}]`, replacement, integration.GroupName()))
	require.Error(t, err)
	conflict := new(scim.ScimError)
	require.ErrorAs(t, err, &conflict)
	assert.Equal(t, "409", conflict.Status)
	lookup(CTX, orgID, externalID, groupID, renamed)
	_, err = Instance.Client.SCIM.Groups.Replace(CTX, orgID, groupID, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q}`, renamed)))
	require.NoError(t, err)
	lookup(CTX, orgID, externalID, "", "")
	// External IDs are scoped to an organization and are case-exact.
	admin := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	_, err = create(admin, SecondaryOrganization.OrganizationId, integration.GroupName(), replacement)
	require.NoError(t, err)
	_, err = create(CTX, orgID, integration.GroupName(), "Case-"+externalID)
	require.NoError(t, err)
	_, err = create(CTX, orgID, integration.GroupName(), "case-"+externalID)
	require.NoError(t, err)
}

func TestGroups_externalIDProvisioningSource(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	externalID := integration.GroupName()
	contexts := make([]context.Context, 2)
	for i, domain := range []string{"entra-" + externalID, "okta-" + externalID} {
		userID, pat, err := Instance.CreateMachineUserPATWithMembership(CTX, "ORG_OWNER")
		require.NoError(t, err)
		setAndEnsureMetadata(t, userID, "urn:zitadel:scim:provisioningDomain", domain)
		contexts[i] = integration.WithAuthorizationToken(CTX, pat)
	}
	ids := make([]string, 2)
	for i, ctx := range contexts {
		created, err := Instance.Client.SCIM.Groups.Create(ctx, orgID, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q,"externalId":%q}`, integration.GroupName(), externalID)))
		require.NoError(t, err)
		ids[i] = created.ID
	}
	wait, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	for i, ctx := range contexts {
		require.EventuallyWithT(t, func(ct *assert.CollectT) {
			result, err := Instance.Client.SCIM.Groups.List(ctx, orgID, &scim.ListRequest{Filter: gu.Ptr(fmt.Sprintf(`externalId eq %q`, externalID))})
			require.NoError(ct, err)
			require.Len(ct, result.Resources, 1)
			assert.Equal(ct, ids[i], result.Resources[0].ID)
			other, err := Instance.Client.SCIM.Groups.Get(ctx, orgID, ids[1-i])
			require.NoError(ct, err)
			assert.Empty(ct, other.ExternalID)
		}, wait, tick)
	}
}
