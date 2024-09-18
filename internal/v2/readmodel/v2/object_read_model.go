package readmodel

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

// objectManager manages single objects.
type objectManager interface {
	manager
}

type objectReadModel struct {
	manager   objectManager
	es        *eventstore.Eventstore
	reduceErr error
}

// newObjectReadModel decorates the manager.
// It manages
// - the subscription to the eventstore
// - the reduction of events
func newObjectReadModel(ctx context.Context, manager objectManager, es *eventstore.Eventstore) *objectReadModel {
	readModel := &objectReadModel{
		manager: manager,
		es:      es,
	}

	go readModel.subscribe(ctx)

	return readModel
}

func (rm *objectReadModel) subscribe(ctx context.Context) {
	notifications := rm.es.Subscribe(rm.manager.Reducers().EventTypes()...)
	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-notifications:
			if !ok {
				return
			}
			// TODO: implement batching
			err := rm.es.FilterToReducer(
				ctx,
				rm.manager.EventstoreV3Query(n.Position).
					OrderAsc().
					AwaitOpenTransactions(),
				rm,
			)
			rm.log().OnError(err).Error("unable to filter to read model")
		}
	}
}

// AppendEvents implements [eventstore.reducer].
func (rm *objectReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		storageEvent := eventstore.EventToV2(event)
		reduce := rm.manager.Reducers()[storageEvent.Aggregate.Type][storageEvent.Type]
		if reduce == nil {
			rm.logEvent(storageEvent).Debug("no reducer found")
			continue
		}
		err := reduce(storageEvent)
		if err != nil {
			rm.reduceErr = err
			rm.logEvent(storageEvent).WithError(err).Error("could not reduce events")
			return
		}
	}
}

// Reduce implements [eventstore.reducer].
func (a *objectReadModel) Reduce() error {
	err := a.reduceErr
	a.reduceErr = nil
	return err
}

func (rm *objectReadModel) init(ctx context.Context) {
	err := rm.es.FilterToReducer(ctx, rm.manager.EventstoreV3Query(decimal.Zero), rm)
	rm.log().OnError(err).Error("unable to init read model")
}

func (rm *objectReadModel) logEvent(event *v2_es.StorageEvent) *logging.Entry {
	return rm.log().
		WithField("position", event.Position.Position.String()).
		WithField("in_position_order", event.Position.InPositionOrder)
}

func (rm *objectReadModel) log() *logging.Entry {
	return logging.WithFields(
		"name", rm.manager.Name(),
	)
}
