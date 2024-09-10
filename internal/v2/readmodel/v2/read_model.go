package readmodel

import (
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
)

type readModel struct {
	reducers map[string]map[string]v2_es.ReduceEvent

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
	Reducers() map[string]map[string]v2_es.ReduceEvent
	// Name returns the name of the read model.
	Name() string
}

// todo: does not work
func mergeReducers(reducers ...map[string]map[string]v2_es.ReduceEvent) map[string]map[string]v2_es.ReduceEvent {
	mergedReduces := make(map[string]map[string][]v2_es.ReduceEvent, len(reducers))

	for _, r := range reducers {
		for aggregateType, eventReducers := range r {
			if _, ok := mergedReduces[aggregateType]; !ok {
				mergedReduces[aggregateType] = map[string][]v2_es.ReduceEvent{}
			}
			for eventType, reducer := range eventReducers {
				mergedReduces[aggregateType][eventType] = append(mergedReduces[aggregateType][eventType], reducer)
			}
		}
	}

	merged := make(map[string]map[string]v2_es.ReduceEvent, len(mergedReduces))
	for aggregateType, eventReducers := range mergedReduces {
		if _, ok := merged[aggregateType]; !ok {
			merged[aggregateType] = make(map[string]v2_es.ReduceEvent, len(eventReducers))
		}
		for eventType, reducers := range eventReducers {
			merged[aggregateType][eventType] = func(event *v2_es.StorageEvent) error {
				for _, reduce := range reducers {
					if err := reduce(event); err != nil {
						return err
					}
				}
				return nil
			}
		}
	}

	return merged
}
