package readmodel

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

type CachedReadModel[M objectManager, T cacheModel] struct {
	cache[T]

	eventStore     *eventstore.Eventstore
	notifications  chan decimal.Decimal
	latestPosition v2_es.GlobalPosition
	// interestedIn   map[eventstore.AggregateType][]eventstore.EventType
	// reduce         v2_es.Reduce

	manager M
}

// Reduce implements [eventstore.reducer]
func (c *CachedReadModel[M, T]) Reduce() error {
	return nil
}

// AppendEvents implements [eventstore.reducer]
func (c *CachedReadModel[M, T]) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		storageEvent := eventstore.EventToV2(event)
		if storageEvent.Position.IsLessOrEqual(c.latestPosition) {
			continue
		}

		reduce := c.manager.Reducers()[string(event.Aggregate().Type)][string(event.Type())]
		if reduce == nil {
			continue
		}
		err := reduce(storageEvent)
		logging.WithFields("position", event.Position().String(), "in_tx_order", event.InTxOrder()).OnError(err).Error("could not reduce events")

		c.latestPosition = storageEvent.Position
	}
}

func NewCachedReadModel[M objectManager, T cacheModel](ctx context.Context, manager M, eventStore *eventstore.Eventstore) *CachedReadModel[M, T] {
	readModel := &CachedReadModel[M, T]{
		cache:         newMapCache[T](),
		eventStore:    eventStore,
		notifications: make(chan decimal.Decimal),
		manager:       manager,
	}
	go readModel.subscription(ctx)
	readModel.createSubscription()
	return readModel
}

func (c *CachedReadModel[M, T]) createSubscription() {
	for _, eventTypes := range c.manager.Reducers() {
		for eventType := range eventTypes {
			c.eventStore.Subscribe(c.notifications, eventstore.EventType(eventType))
		}
	}
}

func (c *CachedReadModel[M, T]) subscription(ctx context.Context) {
	eventFilters := make(map[eventstore.AggregateType][]eventstore.EventType, len(c.manager.Reducers()))
	for aggregateType, eventTypes := range c.manager.Reducers() {
		for eventType := range eventTypes {
			eventFilters[eventstore.AggregateType(aggregateType)] = append(eventFilters[eventstore.AggregateType(aggregateType)], eventstore.EventType(eventType))
		}
	}

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
				PositionGreaterEqual(c.latestPosition.Position).
				OrderAsc()
			for aggregateType, eventTypes := range eventFilters {
				builder = builder.AddQuery().
					AggregateTypes(eventstore.AggregateType(aggregateType)).
					EventTypes(eventTypes...).
					Builder()
			}
			err := c.eventStore.FilterToReducer(ctx, builder, c)
			// TODO: how to handle retries?
			logging.OnError(err).Error("could not filter to cached read model")
		}
	}
}

type cache[T cacheModel] interface {
	get(key string) (T, bool)
	getAll() []T
	set(key string, value T) error
}

var _ cache[cacheModel] = (*MapCache[cacheModel])(nil)

type MapCache[T cacheModel] map[string]T

func newMapCache[T cacheModel]() *MapCache[T] {
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

type objectManager interface {
	Reducers() map[string]map[string]v2_es.ReduceEvent
}

type cacheModel interface {
	// Reducers() map[string]map[string]v2_es.ReduceEvent
}
