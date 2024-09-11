package projection

import (
	"maps"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type Reducers map[string]map[string]eventstore.ReduceEvent

// todo: does not work
func MergeReducers(reducers ...Reducers) Reducers {
	mergedReduces := make(map[string]map[string][]eventstore.ReduceEvent, len(reducers))

	for _, r := range reducers {
		for aggregateType, eventReducers := range r {
			if _, ok := mergedReduces[aggregateType]; !ok {
				mergedReduces[aggregateType] = map[string][]eventstore.ReduceEvent{}
			}
			for eventType, reducer := range eventReducers {
				mergedReduces[aggregateType][eventType] = append(mergedReduces[aggregateType][eventType], reducer)
			}
		}
	}

	merged := make(Reducers, len(mergedReduces))
	for aggregateType, eventReducers := range mergedReduces {
		if _, ok := merged[aggregateType]; !ok {
			merged[aggregateType] = make(map[string]eventstore.ReduceEvent, len(eventReducers))
		}
		for eventType, reducers := range eventReducers {
			merged[aggregateType][eventType] = func(event *eventstore.StorageEvent) error {
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

func OverwriteReduces(reducers Reducers, overwrite eventstore.ReduceEvent) Reducers {
	overwritten := maps.Clone(reducers)
	for aggregateType, eventReducers := range overwritten {
		for eventType := range eventReducers {
			overwritten[aggregateType][eventType] = overwrite
		}
	}
	return overwritten
}
