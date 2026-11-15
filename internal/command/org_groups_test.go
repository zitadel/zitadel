package command

import (
	"context"
	"errors"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

func TestOrgGroupRemovals(t *testing.T) {
	t.Parallel()
	filterErr := errors.New("filter error")

	tests := []struct {
		name    string
		events  []eventstore.Event
		filter  func(events []eventstore.Event) func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error)
		want    []string
		wantErr error
	}{
		{
			name: "filter error, propagated",
			filter: func(_ []eventstore.Event) func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error) {
				return func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error) {
					return nil, filterErr
				}
			},
			wantErr: filterErr,
		},
		{
			name: "no groups, empty",
			want: []string{},
		},
		{
			name: "added and renamed, latest name returned",
			events: []eventstore.Event{
				group.NewGroupAddedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					"",
				),
				group.NewGroupChangedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					[]group.GroupChanges{group.ChangeName(gu.Ptr("renamed"))},
				),
			},
			want: []string{"renamed"},
		},
		{
			name: "description-only change keeps the added name",
			events: []eventstore.Event{
				group.NewGroupAddedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					"",
				),
				group.NewGroupChangedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					[]group.GroupChanges{group.ChangeDescription(gu.Ptr("new description"))},
				),
			},
			want: []string{"original"},
		},
		{
			name: "removed group is excluded",
			events: []eventstore.Event{
				group.NewGroupAddedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"deleted",
					"",
				),
				group.NewGroupRemovedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"deleted",
				),
				group.NewGroupAddedEvent(context.Background(),
					&group.NewAggregate("group2", "org1").Aggregate,
					"kept",
					"",
				),
			},
			want: []string{"kept"},
		},
		{
			name: "renamed then removed is excluded",
			events: []eventstore.Event{
				group.NewGroupAddedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					"",
				),
				group.NewGroupChangedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"original",
					[]group.GroupChanges{group.ChangeName(gu.Ptr("renamed"))},
				),
				group.NewGroupRemovedEvent(context.Background(),
					&group.NewAggregate("group1", "org1").Aggregate,
					"renamed",
				),
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filter := func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error) {
				return groupFieldsFixture(tt.events...), nil
			}
			if tt.filter != nil {
				filter = tt.filter(tt.events)
			}

			events, err := orgGroupRemovals(context.Background(), filter, "org1")
			var got []*eventstore.UniqueConstraint
			for _, event := range events {
				got = append(got, event.UniqueConstraints()...)
			}
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			want := make([]*eventstore.UniqueConstraint, 0, len(tt.want))
			for _, name := range tt.want {
				want = append(want, group.NewRemoveGroupNameUniqueConstraint(name, "org1"))
			}
			assert.ElementsMatch(t, want, got)
		})
	}
}

func TestOrgGroupRemovals_externalIDs(t *testing.T) {
	ctx := context.Background()
	aggregate := &group.NewAggregate("group1", "org1").Aggregate
	filter := func(context.Context, ...map[eventstore.FieldType]any) ([]*eventstore.SearchResult, error) {
		return groupFieldsFixture(
			group.NewGroupAddedEvent(ctx, aggregate, "group", ""),
			group.NewGroupExternalIDSetEvent(ctx, aggregate, "entra", "", "old"),
			group.NewGroupExternalIDSetEvent(ctx, aggregate, "entra", "old", "new"),
			group.NewGroupExternalIDSetEvent(ctx, aggregate, "okta", "", "removed"),
			group.NewGroupExternalIDSetEvent(ctx, aggregate, "okta", "removed", ""),
		), nil
	}
	events, err := orgGroupRemovals(ctx, filter, "org1")
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.ElementsMatch(t, []*eventstore.UniqueConstraint{group.NewRemoveGroupNameUniqueConstraint("group", "org1"), group.NewRemoveGroupExternalIDUniqueConstraint("group1", "org1", "entra", "new")}, events[0].UniqueConstraints())
}
