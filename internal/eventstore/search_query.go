package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

//SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
	columns       repository.Columns
	limit         uint64
	desc          bool
	resourceOwner string
	queries       []*SearchQuery
}

type SearchQuery struct {
	builder        *SearchQueryBuilder
	aggregateTypes []AggregateType
	aggregateIDs   []string
	eventSequence  uint64
	eventTypes     []EventType
	eventData      map[string]interface{}
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
func NewSearchQueryBuilder(columns Columns) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		columns: repository.Columns(columns),
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

func (factory *SearchQueryBuilder) AddQuery() *SearchQuery {
	query := &SearchQuery{
		builder: factory,
	}
	factory.queries = append(factory.queries, query)

	return query
}

func (factory *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	factory.resourceOwner = resourceOwner
	return factory
}

func (factory *SearchQueryBuilder) OrderDesc() *SearchQueryBuilder {
	factory.desc = true
	return factory
}

func (query SearchQuery) Or() *SearchQuery {
	return query.builder.AddQuery()
}

func (query *SearchQuery) AggregateTypes(types ...AggregateType) *SearchQuery {
	query.aggregateTypes = types
	return query
}

func (query *SearchQuery) SequenceGreater(sequence uint64) *SearchQuery {
	query.eventSequence = sequence
	return query
}

func (query *SearchQuery) AggregateIDs(ids ...string) *SearchQuery {
	query.aggregateIDs = ids
	return query
}

func (query *SearchQuery) EventTypes(types ...EventType) *SearchQuery {
	query.eventTypes = types
	return query
}

func (factory *SearchQueryBuilder) OrderAsc() *SearchQueryBuilder {
	factory.desc = false
	return factory
}

func (query *SearchQuery) EventData(data map[string]interface{}) *SearchQuery {
	query.eventData = data
	return query
}

func (query *SearchQuery) SearchQueryBuilder() *SearchQueryBuilder {
	return query.builder
}

func (factory *SearchQueryBuilder) build() (*repository.SearchQuery, error) {
	if factory == nil ||
		len(factory.queries) < 1 ||
		factory.columns.Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "factory invalid")
	}
	filters := make([][]*repository.Filter, len(factory.queries))

	for i, query := range factory.queries {
		for _, f := range []func() *repository.Filter{
			query.aggregateTypeFilter,
			query.aggregateIDFilter,
			query.eventSequenceFilter,
			query.eventTypeFilter,
			query.eventDataFilter,
			query.builder.resourceOwnerFilter,
		} {
			if filter := f(); filter != nil {
				if err := filter.Validate(); err != nil {
					return nil, err
				}
				filters[i] = append(filters[i], filter)
			}
		}

	}

	return &repository.SearchQuery{
		Columns: factory.columns,
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (query *SearchQuery) aggregateIDFilter() *repository.Filter {
	if len(query.aggregateIDs) < 1 {
		return nil
	}
	if len(query.aggregateIDs) == 1 {
		return repository.NewFilter(repository.FieldAggregateID, query.aggregateIDs[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldAggregateID, query.aggregateIDs, repository.OperationIn)
}

func (query *SearchQuery) eventTypeFilter() *repository.Filter {
	if len(query.eventTypes) < 1 {
		return nil
	}
	if len(query.eventTypes) == 1 {
		return repository.NewFilter(repository.FieldEventType, repository.EventType(query.eventTypes[0]), repository.OperationEquals)
	}
	eventTypes := make([]repository.EventType, len(query.eventTypes))
	for i, eventType := range query.eventTypes {
		eventTypes[i] = repository.EventType(eventType)
	}
	return repository.NewFilter(repository.FieldEventType, eventTypes, repository.OperationIn)
}

func (query *SearchQuery) aggregateTypeFilter() *repository.Filter {
	if len(query.aggregateTypes) == 1 {
		return repository.NewFilter(repository.FieldAggregateType, repository.AggregateType(query.aggregateTypes[0]), repository.OperationEquals)
	}
	aggregateTypes := make([]repository.AggregateType, len(query.aggregateTypes))
	for i, aggregateType := range query.aggregateTypes {
		aggregateTypes[i] = repository.AggregateType(aggregateType)
	}
	return repository.NewFilter(repository.FieldAggregateType, aggregateTypes, repository.OperationIn)
}

func (query *SearchQuery) eventSequenceFilter() *repository.Filter {
	if query.eventSequence == 0 {
		return nil
	}
	sortOrder := repository.OperationGreater
	if query.builder.desc {
		sortOrder = repository.OperationLess
	}
	return repository.NewFilter(repository.FieldSequence, query.eventSequence, sortOrder)
}

func (factory *SearchQueryBuilder) resourceOwnerFilter() *repository.Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldResourceOwner, factory.resourceOwner, repository.OperationEquals)
}

func (query *SearchQuery) eventDataFilter() *repository.Filter {
	if len(query.eventData) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldEventData, query.eventData, repository.OperationJSONContains)
}
