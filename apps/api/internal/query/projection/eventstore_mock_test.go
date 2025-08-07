package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
)

var _ handler.EventStore = (*mockEventStore)(nil)

type mockEventStore struct {
	instanceIDsResponse [][]string
	instanceIDCounter   int
	filterResponse      [][]eventstore.Event
	filterCounter       int
	pushResponse        [][]eventstore.Event
	pushCounter         int
}

func newMockEventStore() *mockEventStore {
	return new(mockEventStore)
}

func (m *mockEventStore) appendFilterResponse(events []eventstore.Event) *mockEventStore {
	m.filterResponse = append(m.filterResponse, events)
	return m
}

func (m *mockEventStore) InstanceIDs(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]string, error) {
	m.instanceIDCounter++
	return m.instanceIDsResponse[m.instanceIDCounter-1], nil
}

func (m *mockEventStore) Filter(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	m.filterCounter++
	return m.filterResponse[m.filterCounter-1], nil
}

func (m *mockEventStore) FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error {
	m.filterCounter++
	events := m.filterResponse[m.filterCounter-1]
	reducer.AppendEvents(events...)
	return reducer.Reduce()
}

func (m *mockEventStore) Push(ctx context.Context, cmds ...eventstore.Command) ([]eventstore.Event, error) {
	m.pushCounter++
	return m.pushResponse[m.pushCounter-1], nil
}

func (m *mockEventStore) FillFields(ctx context.Context, events ...eventstore.FillFieldsEvent) error {
	return nil
}
