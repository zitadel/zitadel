package command

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	repo "github.com/zitadel/zitadel/internal/repository/group"
)

func expectGroupMembershipsFixture(userID, orgID string, ids ...string) expect {
	var events []eventstore.Event
	for _, id := range ids {
		aggregate := &repo.NewAggregate(id, orgID).Aggregate
		events = append(events, repo.NewGroupUserAddedEvent(context.Background(), aggregate, userID))
	}
	return expectGroupFields(events...)
}

func TestRemoveUserGroupMemberships_authoritativeState(t *testing.T) {
	ctx := context.Background()
	aggregate := &repo.NewAggregate("g", "org").Aggregate
	c := &Commands{eventstore: expectEventstore(expectGroupFields(
		repo.NewGroupUserAddedEvent(ctx, aggregate, "removed"),
		repo.NewGroupUserRemovedEvent(ctx, aggregate, "removed"),
		repo.NewGroupUserAddedEvent(ctx, aggregate, "current"),
	))(t)}
	events, err := c.removeUserGroupMemberships(ctx, "current", "org")
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "current", events[0].(*repo.GroupUserRemovedEvent).UserID)
	require.Equal(t, "org", events[0].Aggregate().ResourceOwner)
	// No projected list is supplied; the committed addition alone determines the cascade.
}

type groupFieldValue struct{ value any }

func (v groupFieldValue) Unmarshal(ptr any) error {
	data, err := json.Marshal(v.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, ptr)
}

func groupFieldsFixture(events ...eventstore.Event) []*eventstore.SearchResult {
	fields := make(map[[3]string]*eventstore.SearchResult)
	for _, event := range events {
		for _, operation := range event.(eventstore.FillFieldsEvent).Fields() {
			if field := operation.Set; field != nil {
				fields[[3]string{field.Aggregate.ID, field.Object.ID, field.FieldName}] = &eventstore.SearchResult{Aggregate: *field.Aggregate, Object: field.Object, FieldName: field.FieldName, Value: groupFieldValue{field.Value.Value}}
			} else {
				for key, field := range fields {
					if field.Aggregate.ID != event.Aggregate().ID {
						continue
					}
					if id, ok := operation.Remove[eventstore.FieldTypeObjectID]; ok && id != field.Object.ID {
						continue
					}
					if name, ok := operation.Remove[eventstore.FieldTypeFieldName]; ok && name != field.FieldName {
						continue
					}
					delete(fields, key)
				}
			}
		}
	}
	results := make([]*eventstore.SearchResult, 0, len(fields))
	for _, field := range fields {
		results = append(results, field)
	}
	return results
}

func expectGroupFields(events ...eventstore.Event) expect {
	return func(m *mock.MockRepository) { m.ExpectSearch(groupFieldsFixture(events...)...) }
}

func expectGroupSearchError(err error) expect {
	return func(m *mock.MockRepository) {
		m.MockSearcher.EXPECT().Search(gomock.Any(), gomock.Any()).Return(nil, err)
	}
}
