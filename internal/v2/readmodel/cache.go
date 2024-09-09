package readmodel

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

type CachedReadModel[T model] struct {
	cache[T]

	eventStore     *eventstore.Eventstore
	notifications  chan decimal.Decimal
	latestPosition decimal.Decimal
	interestedIn   map[eventstore.AggregateType][]eventstore.EventType
	reduce         v2_es.Reduce
}

// Reduce implements [eventstore.reducer]
func (c *CachedReadModel[T]) Reduce() error {
	return nil
}

// AppendEvents implements [eventstore.reducer]
func (c *CachedReadModel[T]) AppendEvents(events ...eventstore.Event) {
	storageEvents := make([]*v2_es.StorageEvent, 0, len(events))
	for _, event := range events {
		if event.Position() == c.latestPosition {
			continue
		}
		storageEvents = append(storageEvents, eventstore.EventToV2(event))
	}
	if len(storageEvents) == 0 {
		return
	}

	err := c.reduce(storageEvents...)
	logging.OnError(err).Error("could not reduce events")

	c.latestPosition = storageEvents[len(storageEvents)-1].Position.Position
}

func NewCachedReadModel[T model](ctx context.Context, eventStore *eventstore.Eventstore, reduce v2_es.Reduce) *CachedReadModel[T] {
	var t T
	readModel := &CachedReadModel[T]{
		cache:         newMapCache[T](),
		eventStore:    eventStore,
		notifications: make(chan decimal.Decimal),
		interestedIn:  t.InterestedIn(),
		reduce:        reduce,
	}
	go readModel.subscription(ctx)
	readModel.createSubscription()
	return readModel
}

func (c *CachedReadModel[T]) createSubscription() {
	for _, eventTypes := range c.interestedIn {
		c.eventStore.Subscribe(c.notifications, eventTypes...)
	}
}

func (c *CachedReadModel[T]) subscription(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// TODO: unsubscribe, close(c.notifications)
			return
		case position := <-c.notifications:
			// TODO: position as upper bound?
			_ = position
			builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
				AwaitOpenTransactions().
				PositionGreaterEqual(c.latestPosition).
				OrderAsc()
			for aggregateType, eventTypes := range c.interestedIn {
				builder = builder.AddQuery().
					AggregateTypes(aggregateType).
					EventTypes(eventTypes...).
					Builder()
			}
			err := c.eventStore.FilterToReducer(ctx, builder, c)
			// TODO: how to handle retries?
			logging.OnError(err).Error("could not filter to cached read model")
		}
	}
}

type cache[T model] interface {
	get(key string) (T, bool)
	getAll() []T
	set(key string, value T) error
}

var _ cache[model] = (*MapCache[model])(nil)

type MapCache[T model] map[string]T

func newMapCache[T model]() *MapCache[T] {
	m := make(MapCache[T])
	return &m
}

// get implements cache.
func (m *MapCache[T]) get(key string) (T, bool) {
	object, ok := (*m)[key]
	return object, ok
}

// get implements cache.
func (m *MapCache[T]) getAll() []T {
	objects := make([]T, 0, len(*m))

	for _, object := range *m {
		objects = append(objects, object)
	}

	return objects
}

// set implements cache.
func (m *MapCache[T]) set(key string, value T) error {
	(*m)[key] = value
	return nil
}

type model interface {
	InterestedIn() map[eventstore.AggregateType][]eventstore.EventType
	Reduce(events ...*v2_es.StorageEvent) error
}
