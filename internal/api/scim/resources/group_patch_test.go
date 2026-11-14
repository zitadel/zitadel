package resources

import (
	"encoding/json"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/command"
)

// Payload cases follow Polis npm/test/dsync/{membership.test.ts,data/group-requests.ts}
// at e13ed654, and the Okta / Microsoft Entra SCIM provisioning documentation.
func TestGroupMutations(t *testing.T) {
	tests := []struct {
		name, payload string
		want          []command.GroupMutation
		invalid       bool
	}{
		{name: "Okta pathless rename ignores read-only ID", payload: `[{"op":"replace","value":{"id":"group","displayName":"Engineering"}}]`, want: []command.GroupMutation{{Name: gu.Ptr("Engineering")}}},
		{name: "Entra capitalized add", payload: `[{"op":"Add","path":"members","value":[{"value":"alice","$ref":null},{"value":"bob"}]}]`, want: []command.GroupMutation{{AddMembers: []string{"alice", "bob"}}}},
		{name: "Entra remove only listed members", payload: `[{"op":"Remove","path":"members","value":[{"value":"alice","$ref":null}]}]`, want: []command.GroupMutation{{RemoveMembers: []string{"alice"}}}},
		{name: "Okta filtered removal", payload: `[{"op":"remove","path":"members[value eq \"alice\"]"}]`, want: []command.GroupMutation{{RemoveMembers: []string{"alice"}}}},
		{name: "replace members", payload: `[{"op":"replace","path":"members","value":[]}]`, want: []command.GroupMutation{{Members: gu.Ptr([]string{})}}},
		{name: "remove all members", payload: `[{"op":"remove","path":"members"}]`, want: []command.GroupMutation{{Members: gu.Ptr([]string{})}}},
		{name: "empty Entra removal does not clear membership", payload: `[{"op":"remove","path":"members","value":[]}]`, want: []command.GroupMutation{{RemoveMembers: []string{}}}},
		{name: "ordered batch", payload: `[{"op":"add","path":"members","value":[{"value":"alice"}]},{"op":"remove","path":"members[value eq \"alice\"]"}]`, want: []command.GroupMutation{{AddMembers: []string{"alice"}}, {RemoveMembers: []string{"alice"}}}},
		{name: "schema-qualified case-insensitive attribute", payload: `[{"op":"Replace","path":"urn:ietf:params:scim:schemas:core:2.0:Group:displayName","value":"Sales"}]`, want: []command.GroupMutation{{Name: gu.Ptr("Sales")}}},
		{name: "null member", payload: `[{"op":"add","path":"members","value":[null]}]`, invalid: true},
		{name: "missing member ID", payload: `[{"op":"add","path":"members","value":[{"display":"Alice"}]}]`, invalid: true},
		{name: "nested group", payload: `[{"op":"add","path":"members","value":[{"value":"g","type":"Group"}]}]`, invalid: true},
		{name: "null members", payload: `[{"op":"remove","path":"members","value":null}]`, invalid: true},
		{name: "member object instead of array", payload: `[{"op":"add","path":"members","value":{"value":"alice"}}]`, invalid: true},
		{name: "remove without path", payload: `[{"op":"remove","value":{"members":[]}}]`, invalid: true},
		{name: "empty name", payload: `[{"op":"replace","path":"displayName","value":" "}]`, invalid: true},
		{name: "null name", payload: `[{"op":"replace","path":"displayName","value":null}]`, invalid: true},
		{name: "unsupported filter", payload: `[{"op":"remove","path":"members[value ne \"alice\"]"}]`, invalid: true},
		{name: "unsupported compound filter", payload: `[{"op":"remove","path":"members[value eq \"alice\" or value eq \"bob\"]"}]`, invalid: true},
		{name: "immutable member reference", payload: `[{"op":"replace","path":"members[value eq \"alice\"].value","value":"bob"}]`, invalid: true},
		{name: "unsupported attribute", payload: `[{"op":"replace","path":"unknown","value":"x"}]`, invalid: true},
		{name: "invalid trailing operation returns no mutations", payload: `[{"op":"replace","path":"displayName","value":"Sales"},{"op":"remove","path":"displayName"}]`, invalid: true},
		{name: "set external identity", payload: `[{"op":"Replace","path":"externalId","value":"object-id"}]`, want: []command.GroupMutation{{ExternalID: &command.GroupExternalID{Value: "object-id"}}}},
		{name: "remove external identity", payload: `[{"op":"remove","path":"externalId"}]`, want: []command.GroupMutation{{ExternalID: &command.GroupExternalID{}}}},
		{name: "invalid external identity", payload: `[{"op":"replace","path":"externalId","value":123}]`, invalid: true},
		{name: "null operation", payload: `[null]`, invalid: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ops patch.OperationCollection
			require.NoError(t, json.Unmarshal([]byte(tt.payload), &ops))
			got, err := groupMutations(ops)
			if tt.invalid {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestExcludeGroupMembers(t *testing.T) {
	for _, value := range []string{"members", " Members ,displayName", "urn:ietf:params:scim:schemas:core:2.0:Group:members"} {
		require.True(t, excludeGroupMembers(value))
	}
	for _, value := range []string{"", "members.value", "displayName"} {
		require.False(t, excludeGroupMembers(value))
	}
}

func TestGroupPatchRejectsNullOperationBeforeDispatch(t *testing.T) {
	var request patch.OperationRequest
	require.NoError(t, json.Unmarshal([]byte(`{"Operations":[null]}`), &request))
	require.Error(t, request.Validate())
}
