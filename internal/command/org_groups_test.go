package command

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
)

func TestOrgGroupNames(t *testing.T) {
	t.Parallel()
	filterErr := errors.New("filter error")

	tests := []struct {
		name    string
		events  []eventstore.Event
		filter  func(events []eventstore.Event) func(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)
		want    []string
		wantErr error
	}{
		{
			name:    "filter error, propagated",
			filter:  func(_ []eventstore.Event) func(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
				return func(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
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

			filter := func(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
				return tt.events, nil
			}
			if tt.filter != nil {
				filter = tt.filter(tt.events)
			}

			got, err := OrgGroupNames(context.Background(), filter, "org1")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.Equal(t, tt.want, got)
		})
	}
}
