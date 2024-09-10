package readmodel

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

// listManager manages lists of objects.
type listManager interface {
	manager
}

type listReadModel struct {
	LatestPosition v2_es.GlobalPosition

	manager   listManager
	es        *eventstore.Eventstore
	reduceErr error
}

// newListReadModel decorates the manager.
// It manages
// - the subscription to the eventstore
// - the reduction of events
func newListReadModel(ctx context.Context, manager listManager, es *eventstore.Eventstore) *listReadModel {
	readModel := &listReadModel{
		manager: manager,
		es:      es,
	}

	go readModel.subscribe(ctx)

	return readModel
}

func (rm *listReadModel) subscribe(ctx context.Context) {
	positions := make(chan decimal.Decimal)

	for _, eventTypes := range rm.manager.Reducers() {
		for eventType := range eventTypes {
			rm.es.Subscribe(positions, eventstore.EventType(eventType))
		}
	}

	for {
		select {
		case <-ctx.Done():
			// TODO: unsubscribe, close(positions)
			return
		case position := <-positions:
			// TODO: implement batching
			err := rm.es.FilterToReducer(
				ctx,
				rm.manager.EventstoreV3Query(position).
					OrderAsc().
					AwaitOpenTransactions().
					PositionGreaterEqual(rm.LatestPosition.Position),
				rm,
			)
			rm.log().OnError(err).Error("unable to filter to read model")
		}
	}
}

// AppendEvents implements [eventstore.reducer].
func (rm *listReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		storageEvent := eventstore.EventToV2(event)
		reduce := rm.manager.Reducers()[storageEvent.Aggregate.Type][storageEvent.Type]
		if reduce == nil {
			rm.logEvent(storageEvent).Debug("no reducer found")
			rm.LatestPosition = storageEvent.Position
			continue
		}
		err := reduce(storageEvent)
		if err != nil {
			rm.reduceErr = err
			rm.logEvent(storageEvent).WithError(err).Error("could not reduce events")
			return
		}
		rm.LatestPosition = storageEvent.Position
	}
}

// Reduce implements [eventstore.reducer].
func (a *listReadModel) Reduce() error {
	err := a.reduceErr
	a.reduceErr = nil
	return err
}

func (rm *listReadModel) logEvent(event *v2_es.StorageEvent) *logging.Entry {
	return rm.log().
		WithField("position", event.Position.Position.String()).
		WithField("in_position_order", event.Position.InPositionOrder)
}

func (rm *listReadModel) log() *logging.Entry {
	return logging.WithFields(
		"name", rm.manager.Name(),
	)
}
