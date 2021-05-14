package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

//SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
	columns              repository.Columns
	limit                uint64
	desc                 bool
	aggregateTypes       []AggregateType
	aggregateIDs         []string
	eventSequenceGreater uint64
	eventSequenceLess    uint64
	eventTypes           []EventType
	eventData            map[string]interface{}
	resourceOwner        string
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

// NewSearchQueryBuilder creates a new builder for event filters
// aggregateTypes must contain at least one aggregate type
func NewSearchQueryBuilder(columns Columns, aggregateTypes ...AggregateType) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		columns:        repository.Columns(columns),
		aggregateTypes: aggregateTypes,
	}
}

func (builder *SearchQueryBuilder) Columns(columns Columns) *SearchQueryBuilder {
	builder.columns = repository.Columns(columns)
	return builder
}

func (builder *SearchQueryBuilder) Limit(limit uint64) *SearchQueryBuilder {
	builder.limit = limit
	return builder
}

func (builder *SearchQueryBuilder) SequenceGreater(sequence uint64) *SearchQueryBuilder {
	builder.eventSequenceGreater = sequence
	return builder
}

func (builder *SearchQueryBuilder) SequenceLess(sequence uint64) *SearchQueryBuilder {
	builder.eventSequenceLess = sequence
	return builder
}

func (builder *SearchQueryBuilder) AggregateIDs(ids ...string) *SearchQueryBuilder {
	builder.aggregateIDs = ids
	return builder
}

func (builder *SearchQueryBuilder) EventTypes(types ...EventType) *SearchQueryBuilder {
	builder.eventTypes = types
	return builder
}

func (builder *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	builder.resourceOwner = resourceOwner
	return builder
}

func (builder *SearchQueryBuilder) OrderDesc() *SearchQueryBuilder {
	builder.desc = true
	return builder
}

func (builder *SearchQueryBuilder) OrderAsc() *SearchQueryBuilder {
	builder.desc = false
	return builder
}

func (builder *SearchQueryBuilder) EventData(query map[string]interface{}) *SearchQueryBuilder {
	builder.eventData = query
	return builder
}

func (builder *SearchQueryBuilder) build() (*repository.SearchQuery, error) {
	if builder == nil ||
		len(builder.aggregateTypes) < 1 ||
		builder.columns.Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "builder invalid")
	}
	filters := []*repository.Filter{
		builder.aggregateTypeFilter(),
	}

	for _, f := range []func() *repository.Filter{
		builder.aggregateIDFilter,
		builder.eventSequenceGreaterFilter,
		builder.eventSequenceLessFilter,
		builder.eventTypeFilter,
		builder.resourceOwnerFilter,
		builder.eventDataFilter,
	} {
		if filter := f(); filter != nil {
			if err := filter.Validate(); err != nil {
				return nil, err
			}
			filters = append(filters, filter)
		}
	}

	return &repository.SearchQuery{
		Columns: builder.columns,
		Limit:   builder.limit,
		Desc:    builder.desc,
		Filters: filters,
	}, nil
}

func (builder *SearchQueryBuilder) aggregateIDFilter() *repository.Filter {
	if len(builder.aggregateIDs) < 1 {
		return nil
	}
	if len(builder.aggregateIDs) == 1 {
		return repository.NewFilter(repository.FieldAggregateID, builder.aggregateIDs[0], repository.OperationEquals)
	}
	return repository.NewFilter(repository.FieldAggregateID, builder.aggregateIDs, repository.OperationIn)
}

func (builder *SearchQueryBuilder) eventTypeFilter() *repository.Filter {
	if len(builder.eventTypes) < 1 {
		return nil
	}
	if len(builder.eventTypes) == 1 {
		return repository.NewFilter(repository.FieldEventType, repository.EventType(builder.eventTypes[0]), repository.OperationEquals)
	}
	eventTypes := make([]repository.EventType, len(builder.eventTypes))
	for i, eventType := range builder.eventTypes {
		eventTypes[i] = repository.EventType(eventType)
	}
	return repository.NewFilter(repository.FieldEventType, eventTypes, repository.OperationIn)
}

func (builder *SearchQueryBuilder) aggregateTypeFilter() *repository.Filter {
	if len(builder.aggregateTypes) == 1 {
		return repository.NewFilter(repository.FieldAggregateType, repository.AggregateType(builder.aggregateTypes[0]), repository.OperationEquals)
	}
	aggregateTypes := make([]repository.AggregateType, len(builder.aggregateTypes))
	for i, aggregateType := range builder.aggregateTypes {
		aggregateTypes[i] = repository.AggregateType(aggregateType)
	}
	return repository.NewFilter(repository.FieldAggregateType, aggregateTypes, repository.OperationIn)
}

func (builder *SearchQueryBuilder) eventSequenceGreaterFilter() *repository.Filter {
	if builder.eventSequenceGreater == 0 {
		return nil
	}
	sortOrder := repository.OperationGreater
	if builder.desc {
		sortOrder = repository.OperationLess
	}
	return repository.NewFilter(repository.FieldSequence, builder.eventSequenceGreater, sortOrder)
}

func (builder *SearchQueryBuilder) eventSequenceLessFilter() *repository.Filter {
	if builder.eventSequenceLess == 0 {
		return nil
	}
	sortOrder := repository.OperationLess
	if builder.desc {
		sortOrder = repository.OperationGreater
	}
	return repository.NewFilter(repository.FieldSequence, builder.eventSequenceLess, sortOrder)
}

func (builder *SearchQueryBuilder) resourceOwnerFilter() *repository.Filter {
	if builder.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldResourceOwner, builder.resourceOwner, repository.OperationEquals)
}

func (builder *SearchQueryBuilder) eventDataFilter() *repository.Filter {
	if len(builder.eventData) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldEventData, builder.eventData, repository.OperationJSONContains)
}
