//go:build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
)

// The lifecycle and membership cases are adapted from the scenarios in Polis
// npm/test/dsync/{groups,membership}.test.ts (e13ed654), using ZITADEL's HTTP client.
func TestGroups_providerProvisioning(t *testing.T) {
	for _, provider := range []string{"Okta", "Entra"} {
		t.Run(provider, func(t *testing.T) {
			orgID := Instance.DefaultOrg.GetId()
			name := integration.GroupName()
			groupID := createScimGroup(t, orgID, name)
			alice := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone()).GetUserId()
			bob := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone()).GetUserId()
			wait, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			check := func(wantName string, ids ...string) {
				t.Helper()
				require.EventuallyWithT(t, func(ct *assert.CollectT) {
					got, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID)
					require.NoError(ct, err)
					assert.Equal(ct, wantName, got.DisplayName)
					actual := make([]string, 0, len(got.Members))
					for _, m := range got.Members {
						actual = append(actual, m.Value)
					}
					assert.ElementsMatch(ct, ids, actual)
				}, wait, tick)
			}
			update := func(operations string) error {
				return Instance.Client.SCIM.Groups.Update(CTX, orgID, groupID, []byte(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":`+operations+`}`))
			}
			check(name)
			require.EventuallyWithT(t, func(ct *assert.CollectT) {
				groups, err := Instance.Client.SCIM.Groups.List(CTX, orgID, &scim.ListRequest{Filter: gu.Ptr(fmt.Sprintf(`displayName eq %q`, strings.ToUpper(name))), Count: gu.Ptr(100), StartIndex: gu.Ptr(1), ExcludedAttributes: "members"})
				require.NoError(ct, err)
				require.Len(ct, groups.Resources, 1)
				assert.Equal(ct, groupID, groups.Resources[0].ID)
				assert.Equal(ct, 1, groups.TotalResults)
			}, wait, tick)
			add := fmt.Sprintf(`[{"op":"Add","path":"members","value":[{"value":%q},{"value":%q},{"value":%q}]}]`, alice, bob, alice)
			require.NoError(t, update(add))
			require.NoError(t, update(add))
			check(name, alice, bob)
			renamed := integration.GroupName()
			rename := fmt.Sprintf(`[{"op":"replace","value":{"id":%q,"displayName":%q}}]`, groupID, renamed)
			if provider == "Entra" {
				rename = fmt.Sprintf(`[{"op":"Replace","path":"displayName","value":%q}]`, renamed)
			}
			require.NoError(t, update(rename))
			check(renamed, alice, bob)
			var remove string
			if provider == "Entra" {
				remove = fmt.Sprintf(`[{"op":"Remove","path":"members","value":[{"$ref":null,"value":%q}]}]`, alice)
			} else {
				path, _ := json.Marshal(fmt.Sprintf(`members[value eq %q]`, alice))
				remove = fmt.Sprintf(`[{"op":"remove","path":%s}]`, path)
			}
			require.NoError(t, update(remove))
			require.NoError(t, update(remove))
			check(renamed, bob)
			excluded, err := Instance.Client.SCIM.Groups.Get(CTX, orgID, groupID+"?excludedAttributes=members")
			require.NoError(t, err)
			assert.Empty(t, excluded.Members)
			groups, err := Instance.Client.SCIM.Groups.List(CTX, orgID, &scim.ListRequest{Filter: gu.Ptr(fmt.Sprintf(`id eq %q`, groupID))})
			require.NoError(t, err)
			require.Len(t, groups.Resources, 1)
			require.Len(t, groups.Resources[0].Members, 1)
			assert.Equal(t, bob, groups.Resources[0].Members[0].Value)
			// A later invalid member must roll back an earlier rename and removal.
			require.Error(t, update(fmt.Sprintf(`[{"op":"replace","path":"displayName","value":%q},{"op":"remove","path":"members"},{"op":"add","path":"members","value":[{"value":"missing-user"}]}]`, integration.GroupName())))
			check(renamed, bob)
			// A successful subsequent update is a barrier against a transient projection match.
			require.NoError(t, update(fmt.Sprintf(`[{"op":"add","path":"members","value":[{"value":%q}]}]`, alice)))
			check(renamed, alice, bob)
			require.NoError(t, update(`[{"op":"remove","path":"members"}]`))
			check(renamed)
			require.NoError(t, Instance.Client.SCIM.Groups.Delete(CTX, orgID, groupID))
		})
	}
}

func TestGroups_createInvalidMemberIsAtomic(t *testing.T) {
	orgID := Instance.DefaultOrg.GetId()
	name := integration.GroupName()
	_, err := Instance.Client.SCIM.Groups.Create(CTX, orgID, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q,"members":[{"value":"missing-user"}]}`, name)))
	require.Error(t, err)
	// The failed request must not reserve the name or leave an orphan group.
	createScimGroup(t, orgID, name)
}

func TestGroups_patchCrossOrganization(t *testing.T) {
	adminCtx := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	created, createErr := Instance.Client.SCIM.Groups.Create(adminCtx, SecondaryOrganization.OrganizationId, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q}`, integration.GroupName())))
	require.NoError(t, createErr)
	groupID := created.ID
	t.Cleanup(func() {
		_ = Instance.Client.SCIM.Groups.Delete(adminCtx, SecondaryOrganization.OrganizationId, groupID)
	})
	err := Instance.Client.SCIM.Groups.Update(CTX, Instance.DefaultOrg.GetId(), groupID, []byte(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":[{"op":"replace","path":"displayName","value":"stolen"}]}`))
	require.Error(t, err)
	scimErr := new(scim.ScimError)
	require.ErrorAs(t, err, &scimErr)
	assert.Equal(t, "404", scimErr.Status)
}
