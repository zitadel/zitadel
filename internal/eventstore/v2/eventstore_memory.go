package eventstore

import "context"

type inMemoryEventstore struct {
	events []*Event
}

func (es *inMemoryEventstore) Health(ctx context.Context) error {
	return nil
}

func (es *inMemoryEventstore) PushAggregates(ctx context.Context, aggregates ...*Aggregate) error {
	return nil
}

func (es *inMemoryEventstore) FilterEvents(ctx context.Context, searchQuery *SearchQueryFactory) (events []*Event, err error) {
	return nil, nil
}

func (es *inMemoryEventstore) LatestSequence(ctx context.Context, searchQuery *SearchQueryFactory) (uint64, error) {
	query, err := searchQuery.Build()
	if err != nil {
		return 0, err
	}
query.Filters[0].
	return 0, nil
}
