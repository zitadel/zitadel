package saml

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/query"
)

// TestAppendUserGroupsAttribute_ActionTakesPrecedence locks down the conflict
// guard: when a SAML action has already set a "groups" custom attribute, the
// storage must return early without overwriting it. p.query is left nil so a
// regressed guard panics loudly instead of silently double-emitting.
func TestAppendUserGroupsAttribute_ActionTakesPrecedence(t *testing.T) {
	p := &Storage{}
	existing := map[string]*customAttribute{
		"groups": {nameFormat: "action-fmt", attributeValue: []string{"action-only"}},
		"other":  {nameFormat: "fmt", attributeValue: []string{"value"}},
	}

	got, err := p.appendUserGroupsAttribute(context.Background(), "user1", existing)
	require.NoError(t, err)
	require.Contains(t, got, "groups")
	assert.Equal(t, "action-fmt", got["groups"].nameFormat)
	assert.Equal(t, []string{"action-only"}, got["groups"].attributeValue)
	assert.Contains(t, got, "other", "unrelated attributes must be preserved")
}

func TestAppendGroupNamesAttribute(t *testing.T) {
	const expectedFormat = "urn:oasis:names:tc:SAML:2.0:attrname-format:basic"

	t.Run("zero memberships keeps input map", func(t *testing.T) {
		existing := map[string]*customAttribute{
			"other": {nameFormat: "fmt", attributeValue: []string{"value"}},
		}
		got := appendGroupNamesAttribute(existing, nil)
		assert.NotContains(t, got, "groups", "groups attribute must be absent for users without memberships")
		assert.Contains(t, got, "other")
	})

	t.Run("multiple memberships all appended", func(t *testing.T) {
		memberships := []*query.GroupUser{
			{GroupID: "g1", GroupName: "engineering"},
			{GroupID: "g2", GroupName: "ops"},
			{GroupID: "g3", GroupName: "security"},
		}
		got := appendGroupNamesAttribute(nil, memberships)

		groups, ok := got["groups"]
		assert.True(t, ok, "groups attribute must be present")
		assert.Equal(t, expectedFormat, groups.nameFormat)
		assert.Equal(t, []string{"engineering", "ops", "security"}, groups.attributeValue)
	})

	t.Run("empty group names filtered", func(t *testing.T) {
		memberships := []*query.GroupUser{
			{GroupID: "g1", GroupName: ""},
			{GroupID: "g2", GroupName: "ops"},
			{GroupID: "g3", GroupName: ""},
		}
		got := appendGroupNamesAttribute(nil, memberships)

		groups, ok := got["groups"]
		assert.True(t, ok, "groups attribute is appended even when some names are empty")
		assert.Equal(t, []string{"ops"}, groups.attributeValue)
	})

	t.Run("all-empty names yield empty attribute value", func(t *testing.T) {
		memberships := []*query.GroupUser{
			{GroupID: "g1", GroupName: ""},
			{GroupID: "g2", GroupName: ""},
		}
		got := appendGroupNamesAttribute(nil, memberships)

		groups, ok := got["groups"]
		assert.True(t, ok, "memberships > 0 always emits the groups attribute key")
		assert.Empty(t, groups.attributeValue)
	})
}
