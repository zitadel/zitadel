package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

//SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
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

// NewSearchQueryBuilder creates a new factory for event filters
// aggregateTypes must contain at least one aggregate type
func NewSearchQueryBuilder(columns Columns, aggregateTypes ...AggregateType) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		columns:        repository.Columns(columns),
		aggregateTypes: aggregateTypes,
	}
}

func (factory *SearchQueryBuilder) Columns(columns Columns) *SearchQueryBuilder {
	factory.columns = repository.Columns(columns)
	return factory
}

func (factory *SearchQueryBuilder) Limit(limit uint64) *SearchQueryBuilder {
	factory.limit = limit
	return factory
}

func (factory *SearchQueryBuilder) SequenceGreater(sequence uint64) *SearchQueryBuilder {
	factory.eventSequence = sequence
	return factory
}

func (factory *SearchQueryBuilder) AggregateIDs(ids ...string) *SearchQueryBuilder {
	factory.aggregateIDs = ids
	return factory
}

func (factory *SearchQueryBuilder) EventTypes(types ...EventType) *SearchQueryBuilder {
	factory.eventTypes = types
	return factory
}

func (factory *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	factory.resourceOwner = resourceOwner
	return factory
}

func (factory *SearchQueryBuilder) OrderDesc() *SearchQueryBuilder {
	factory.desc = true
	return factory
}

func (factory *SearchQueryBuilder) OrderAsc() *SearchQueryBuilder {
	factory.desc = false
	return factory
}

func (factory *SearchQueryBuilder) EventData(query map[string]interface{}) *SearchQueryBuilder {
	factory.eventData = query
	return factory
}

func (factory *SearchQueryBuilder) build() (*repository.SearchQuery, error) {
	if factory == nil ||
		len(factory.aggregateTypes) < 1 ||
		factory.columns.Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "factory invalid")
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
		Columns: factory.columns,
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (factory *SearchQueryBuilder) aggregateIDFilter() *repository.Filter {
	if len(factory.aggregateIDs) < 1 {
		return nil
	}
	if len(factory.aggregateIDs) == 1 {
		return repository.NewFilter(repository.FieldAggregateID, factory.aggregateIDs[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldAggregateID, factory.aggregateIDs, repository.OperationIn)
}

func (factory *SearchQueryBuilder) eventTypeFilter() *repository.Filter {
	if len(factory.eventTypes) < 1 {
		return nil
	}
	if len(factory.eventTypes) == 1 {
		return repository.NewFilter(repository.FieldEventType, repository.EventType(factory.eventTypes[0]), repository.OperationEquals)
	}
	eventTypes := make([]repository.EventType, len(factory.eventTypes))
	for i, eventType := range factory.eventTypes {
		eventTypes[i] = repository.EventType(eventType)
	}
	return repository.NewFilter(repository.FieldEventType, eventTypes, repository.OperationIn)
}

func (factory *SearchQueryBuilder) aggregateTypeFilter() *repository.Filter {
	if len(factory.aggregateTypes) == 1 {
		return repository.NewFilter(repository.FieldAggregateType, repository.AggregateType(factory.aggregateTypes[0]), repository.OperationEquals)
	}
	aggregateTypes := make([]repository.AggregateType, len(factory.aggregateTypes))
	for i, aggregateType := range factory.aggregateTypes {
		aggregateTypes[i] = repository.AggregateType(aggregateType)
	}
	return repository.NewFilter(repository.FieldAggregateType, aggregateTypes, repository.OperationIn)
}

func (factory *SearchQueryBuilder) eventSequenceFilter() *repository.Filter {
	if factory.eventSequence == 0 {
		return nil
	}
	sortOrder := repository.OperationGreater
	if factory.desc {
		sortOrder = repository.OperationLess
	}
	return repository.NewFilter(repository.FieldSequence, factory.eventSequence, sortOrder)
}

func (factory *SearchQueryBuilder) resourceOwnerFilter() *repository.Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldResourceOwner, factory.resourceOwner, repository.OperationEquals)
}

func (factory *SearchQueryBuilder) eventDataFilter() *repository.Filter {
	if len(factory.eventData) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldEventData, factory.eventData, repository.OperationJSONContains)
}
