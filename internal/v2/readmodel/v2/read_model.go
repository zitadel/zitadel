package readmodel

import (
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/projection"
)

type readModel struct {
	reducers projection.Reducers

	reduceErr error
}

// AppendEvents implements eventstore.reducer.
func (rm *readModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		storageEvent := eventstore.EventToV2(event)
		reducer, ok := rm.reducers[storageEvent.Aggregate.Type][storageEvent.Type]
		if !ok {
			continue
		}
		if err := reducer(storageEvent); err != nil {
			rm.reduceErr = err
			return
		}
	}
}

// Reduce implements eventstore.reducer.
func (rm *readModel) Reduce() error {
	if rm.reduceErr != nil {
		err := rm.reduceErr
		rm.reduceErr = nil
		return err
	}
	return nil
}

type manager interface {
	EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder
	// Reducers returns the aggregate types and event types to reduce.
	Reducers() projection.Reducers
	// Name returns the name of the read model.
	Name() string
}
