package eventstore

import (
	"context"
)

type Eventstore interface {
	Health(ctx context.Context) error
	PushAggregates(ctx context.Context, aggregates ...*Aggregate) error
	FilterEvents(ctx context.Context, searchQuery *SearchQueryFactory) (events []*Event, err error)
	LatestSequence(ctx context.Context, searchQuery *SearchQueryFactory) (uint64, error)
}
