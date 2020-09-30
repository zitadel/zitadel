package repository

import "context"

type InMemory struct {
	events []*Event
}

func (repo *InMemory) Health(ctx context.Context) error { return nil }

// PushEvents adds all events of the given aggregates to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (repo *InMemory) Push(ctx context.Context, events ...*Event) error {
	repo.events = append(repo.events, events...)
	return nil
}

// Filter returns all events matching the given search query
func (repo *InMemory) Filter(ctx context.Context, searchQuery *SearchQuery) (events []*Event, err error) {
	indexes := repo.filter(searchQuery)
	events = make([]*Event, len(indexes))
	for i, index := range indexes {
		events[i] = repo.events[index]
	}

	return events, nil
}

func (repo *InMemory) filter(query *SearchQuery) []int {
	foundIndex := make([]int, 0, query.Limit)
events:
	for i, event := range repo.events {
		if query.Limit > 0 && uint64(len(foundIndex)) < query.Limit {
			return foundIndex
		}
		for _, filter := range query.Filters {
			var value interface{}
			switch filter.field {
			case Field_AggregateID:
				value = event.AggregateID
			case Field_EditorService:
				value = event.EditorService
			case Field_EventType:
				value = event.Type
			case Field_AggregateType:
				value = event.AggregateType
			case Field_EditorUser:
				value = event.EditorUser
			case Field_ResourceOwner:
				value = event.ResourceOwner
			case Field_LatestSequence:
				value = event.Sequence
			}
			switch filter.operation {
			case Operation_Equals:
				if filter.value == value {
					foundIndex = append(foundIndex, i)
				}
			case Operation_Greater:
				fallthrough
			case Operation_Less:

				return nil
			case Operation_In:
				values := filter.Value().([]interface{})
				for _, val := range values {
					if val == value {
						foundIndex = append(foundIndex, i)
						continue events
					}
				}
			}
		}
	}

	return foundIndex
}

//LatestSequence returns the latests sequence found by the the search query
func (repo *InMemory) LatestSequence(ctx context.Context, queryFactory *SearchQuery) (uint64, error) {
	return 0, nil
}
