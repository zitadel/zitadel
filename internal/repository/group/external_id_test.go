package group

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestGroupExternalIDConstraints(t *testing.T) {
	aggregate := &NewAggregate("group", "org").Aggregate
	for _, tt := range []struct {
		name, old, value string
		actions          []eventstore.UniqueConstraintAction
	}{
		{"add", "", "abc", []eventstore.UniqueConstraintAction{eventstore.UniqueConstraintAdd}},
		{"replace", "abc", "xyz", []eventstore.UniqueConstraintAction{eventstore.UniqueConstraintRemove, eventstore.UniqueConstraintAdd}},
		{"remove", "abc", "", []eventstore.UniqueConstraintAction{eventstore.UniqueConstraintRemove}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			constraints := NewGroupExternalIDSetEvent(context.Background(), aggregate, "entra", tt.old, tt.value).UniqueConstraints()
			require.Len(t, constraints, len(tt.actions))
			for i, constraint := range constraints {
				require.Equal(t, tt.actions[i], constraint.Action)
				require.False(t, constraint.IsGlobal)
				require.Equal(t, "claimant:group:group", constraint.Claimant)
				if constraint.Action == eventstore.UniqueConstraintAdd {
					require.Equal(t, []string{"org:org", "claimant:group:group"}, constraint.Owners)
				}
			}
		})
	}
	key := externalIDConstraintKey("org", "entra", "ABC")
	require.Equal(t, key, strings.ToLower(key))
	for _, tuple := range [][3]string{{"org", "entra", "abc"}, {"org", "okta", "ABC"}, {"other", "entra", "ABC"}, {"org:entra", "", "ABC"}} {
		require.NotEqual(t, key, externalIDConstraintKey(tuple[0], tuple[1], tuple[2]))
	}
	removed := NewGroupRemovedEvent(context.Background(), aggregate, "name", map[string]string{"entra": "ABC"}).UniqueConstraints()
	require.Len(t, removed, 2)
	require.Equal(t, NewRemoveGroupExternalIDUniqueConstraint("group", "org", "entra", "ABC"), removed[1])
}

func TestGroupExternalIDSetEventPersistsOldExternalID(t *testing.T) {
	aggregate := &NewAggregate("group", "org").Aggregate
	event := NewGroupExternalIDSetEvent(context.Background(), aggregate, "entra", "old", "new")
	payload, err := json.Marshal(event.Payload())
	require.NoError(t, err)

	stored := new(GroupExternalIDSetEvent)
	require.NoError(t, json.Unmarshal(payload, stored))
	stored.SetBaseEvent(&event.BaseEvent)

	require.Equal(t, "old", stored.OldExternalID)
	require.Equal(t, event.UniqueConstraints(), stored.UniqueConstraints())
	require.Equal(t, NewRemoveGroupExternalIDUniqueConstraint("group", "org", "entra", "old"), stored.UniqueConstraints()[0])
	fields := stored.Fields()
	require.Equal(t, event.Fields(), fields)
	require.Len(t, fields, 2)
	require.Equal(t, removeGroupField(aggregate, externalIDConstraintKey("org", "entra", "old"), FieldExternalIDKey), fields[0])
}

func TestGroupRemovedEventSkipUniqueConstraints(t *testing.T) {
	aggregate := &NewAggregate("group", "org").Aggregate
	removed := NewGroupRemovedEvent(context.Background(), aggregate, "name", map[string]string{"entra": "ABC"}).
		WithMembers([]string{"user"}).WithExternalIDKeys([]string{"key"})
	require.Len(t, removed.UniqueConstraints(), 4)
	eventstore.SkipCommandUniqueConstraints(removed)
	require.Empty(t, removed.UniqueConstraints())
}
