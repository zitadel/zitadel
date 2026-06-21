package saml

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/query"
)

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
