package group

import (
	"context"
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
				require.Equal(t, eventstore.AggregateType(AggregateType), constraint.ClaimantType)
				require.Equal(t, aggregate.ID, constraint.ClaimantID)
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
