package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

//SearchQueryFactory represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryFactory struct {
	columns        repository.Columns
	limit          uint64
	desc           bool
	aggregateTypes []AggregateType
	aggregateIDs   []string
	eventSequence  uint64
	eventTypes     []EventType
	eventData      map[string]interface{}
	resourceOwner  string
}

// Columns defines which fields of the event are needed for the query
type Columns repository.Columns

const (
	//ColumnsEvent represents all fields of an event
	ColumnsEvent Columns = repository.ColumnsEvent
	// ColumnsMaxSequence represents the latest sequence of the filtered events
	ColumnsMaxSequence Columns = repository.ColumnsMaxSequence
)

// AggregateType is the object name
type AggregateType repository.AggregateType

// EventType is the description of the change
type EventType repository.EventType

// NewSearchQueryFactory creates a new factory for event filters
// aggregateTypes must contain at least one aggregate type
func NewSearchQueryFactory(columns Columns, aggregateTypes ...AggregateType) *SearchQueryFactory {
	return &SearchQueryFactory{
		columns:        repository.Columns(columns),
		aggregateTypes: aggregateTypes,
	}
}

func (factory *SearchQueryFactory) Columns(columns Columns) *SearchQueryFactory {
	factory.columns = repository.Columns(columns)
	return factory
}

func (factory *SearchQueryFactory) Limit(limit uint64) *SearchQueryFactory {
	factory.limit = limit
	return factory
}

func (factory *SearchQueryFactory) SequenceGreater(sequence uint64) *SearchQueryFactory {
	factory.eventSequence = sequence
	return factory
}

func (factory *SearchQueryFactory) AggregateIDs(ids ...string) *SearchQueryFactory {
	factory.aggregateIDs = ids
	return factory
}

func (factory *SearchQueryFactory) EventTypes(types ...EventType) *SearchQueryFactory {
	factory.eventTypes = types
	return factory
}

func (factory *SearchQueryFactory) ResourceOwner(resourceOwner string) *SearchQueryFactory {
	factory.resourceOwner = resourceOwner
	return factory
}

func (factory *SearchQueryFactory) OrderDesc() *SearchQueryFactory {
	factory.desc = true
	return factory
}

func (factory *SearchQueryFactory) OrderAsc() *SearchQueryFactory {
	factory.desc = false
	return factory
}

func (factory *SearchQueryFactory) EventData(query map[string]interface{}) *SearchQueryFactory {
	factory.eventData = query
	return factory
}

func (factory *SearchQueryFactory) build() (*repository.SearchQuery, error) {
	if factory == nil ||
		len(factory.aggregateTypes) < 1 ||
		factory.columns.Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-tGAD3", "factory invalid")
	}
	filters := []*repository.Filter{
		factory.aggregateTypeFilter(),
	}

	for _, f := range []func() *repository.Filter{
		factory.aggregateIDFilter,
		factory.eventSequenceFilter,
		factory.eventTypeFilter,
		factory.resourceOwnerFilter,
		factory.eventDataFilter,
	} {
		if filter := f(); filter != nil {
			if err := filter.Validate(); err != nil {
				return nil, err
			}
			filters = append(filters, filter)
		}
	}

	return &repository.SearchQuery{
		Columns: repository.Columns(factory.columns),
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (factory *SearchQueryFactory) aggregateIDFilter() *repository.Filter {
	if len(factory.aggregateIDs) < 1 {
		return nil
	}
	if len(factory.aggregateIDs) == 1 {
		return repository.NewFilter(repository.FieldAggregateID, factory.aggregateIDs[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldAggregateID, factory.aggregateIDs, repository.OperationIn)
}

func (factory *SearchQueryFactory) eventTypeFilter() *repository.Filter {
	if len(factory.eventTypes) < 1 {
		return nil
	}
	if len(factory.eventTypes) == 1 {
		return repository.NewFilter(repository.FieldEventType, factory.eventTypes[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldEventType, factory.eventTypes, repository.OperationIn)
}

func (factory *SearchQueryFactory) aggregateTypeFilter() *repository.Filter {
	if len(factory.aggregateTypes) == 1 {
		return repository.NewFilter(repository.FieldAggregateType, factory.aggregateTypes[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldAggregateType, factory.aggregateTypes, repository.OperationIn)
}

func (factory *SearchQueryFactory) eventSequenceFilter() *repository.Filter {
	if factory.eventSequence == 0 {
		return nil
	}
	sortOrder := repository.OperationGreater
	if factory.desc {
		sortOrder = repository.OperationLess
	}
	return repository.NewFilter(repository.FieldSequence, factory.eventSequence, sortOrder)
}

func (factory *SearchQueryFactory) resourceOwnerFilter() *repository.Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldResourceOwner, factory.resourceOwner, repository.OperationEquals)
}

func (factory *SearchQueryFactory) eventDataFilter() *repository.Filter {
	if len(factory.eventData) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldEventData, factory.eventData, repository.OperationJSONContains)
}
