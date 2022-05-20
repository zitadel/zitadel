package eventstore

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

//SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
	columns       repository.Columns
	limit         uint64
	desc          bool
	resourceOwner string
	instanceID    string
	queries       []*SearchQuery
	tx            *sql.Tx
}

type SearchQuery struct {
	builder              *SearchQueryBuilder
	aggregateTypes       []AggregateType
	aggregateIDs         []string
	instanceID           string
	excludedInstanceIDs  []string
	eventSequenceGreater uint64
	eventSequenceLess    uint64
	eventTypes           []EventType
	eventData            map[string]interface{}
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
func NewSearchQueryBuilder(columns Columns) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		columns: repository.Columns(columns),
	}
}

func (builder *SearchQueryBuilder) Matches(event Event, existingLen int) (matches bool) {
	if builder.limit > 0 && uint64(existingLen) >= builder.limit {
		return false
	}
	if builder.resourceOwner != "" && event.Aggregate().ResourceOwner != builder.resourceOwner {
		return false
	}
	if event.Aggregate().InstanceID != "" && builder.instanceID != "" && event.Aggregate().InstanceID != builder.instanceID {
		return false
	}

	if len(builder.queries) == 0 {
		return true
	}
	for _, query := range builder.queries {
		if query.matches(event) {
			return true
		}
	}
	return false
}

//Columns defines which fields are set
func (builder *SearchQueryBuilder) Columns(columns Columns) *SearchQueryBuilder {
	builder.columns = repository.Columns(columns)
	return builder
}

//Limit defines how many events are returned maximally.
func (builder *SearchQueryBuilder) Limit(limit uint64) *SearchQueryBuilder {
	builder.limit = limit
	return builder
}

//ResourceOwner defines the resource owner (org) of the events
func (builder *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	builder.resourceOwner = resourceOwner
	return builder
}

//InstanceID defines the instanceID (system) of the events
func (builder *SearchQueryBuilder) InstanceID(instanceID string) *SearchQueryBuilder {
	builder.instanceID = instanceID
	return builder
}

//OrderDesc changes the sorting order of the returned events to descending
func (builder *SearchQueryBuilder) OrderDesc() *SearchQueryBuilder {
	builder.desc = true
	return builder
}

//OrderAsc changes the sorting order of the returned events to ascending
func (builder *SearchQueryBuilder) OrderAsc() *SearchQueryBuilder {
	builder.desc = false
	return builder
}

//SetTx ensures that the eventstore library uses the existing transaction
func (builder *SearchQueryBuilder) SetTx(tx *sql.Tx) *SearchQueryBuilder {
	builder.tx = tx
	return builder
}

//AddQuery creates a new sub query.
//All fields in the sub query are AND-connected in the storage request.
//Multiple sub queries are OR-connected in the storage request.
func (builder *SearchQueryBuilder) AddQuery() *SearchQuery {
	query := &SearchQuery{
		builder: builder,
	}
	builder.queries = append(builder.queries, query)

	return query
}

//Or creates a new sub query on the search query builder
func (query SearchQuery) Or() *SearchQuery {
	return query.builder.AddQuery()
}

//AggregateTypes filters for events with the given aggregate types
func (query *SearchQuery) AggregateTypes(types ...AggregateType) *SearchQuery {
	query.aggregateTypes = types
	return query
}

//SequenceGreater filters for events with sequence greater the requested sequence
func (query *SearchQuery) SequenceGreater(sequence uint64) *SearchQuery {
	query.eventSequenceGreater = sequence
	return query
}

//SequenceLess filters for events with sequence less the requested sequence
func (query *SearchQuery) SequenceLess(sequence uint64) *SearchQuery {
	query.eventSequenceLess = sequence
	return query
}

//AggregateIDs filters for events with the given aggregate id's
func (query *SearchQuery) AggregateIDs(ids ...string) *SearchQuery {
	query.aggregateIDs = ids
	return query
}

//InstanceID filters for events with the given instanceID
func (query *SearchQuery) InstanceID(instanceID string) *SearchQuery {
	query.instanceID = instanceID
	return query
}

//ExcludedInstanceID filters for events not having the given instanceIDs
func (query *SearchQuery) ExcludedInstanceID(instanceIDs ...string) *SearchQuery {
	query.excludedInstanceIDs = instanceIDs
	return query
}

//EventTypes filters for events with the given event types
func (query *SearchQuery) EventTypes(types ...EventType) *SearchQuery {
	query.eventTypes = types
	return query
}

//EventData filters for events with the given event data.
//Use this call with care as it will be slower than the other filters.
func (query *SearchQuery) EventData(data map[string]interface{}) *SearchQuery {
	query.eventData = data
	return query
}

//Builder returns the SearchQueryBuilder of the sub query
func (query *SearchQuery) Builder() *SearchQueryBuilder {
	return query.builder
}

func (query *SearchQuery) matches(event Event) bool {
	if query.eventSequenceLess > 0 && event.Sequence() >= query.eventSequenceLess {
		return false
	}
	if query.eventSequenceGreater > 0 && event.Sequence() <= query.eventSequenceGreater {
		return false
	}
	if ok := isAggreagteTypes(event.Aggregate(), query.aggregateTypes...); len(query.aggregateTypes) > 0 && !ok {
		return false
	}
	if ok := isAggregateIDs(event.Aggregate(), query.aggregateIDs...); len(query.aggregateIDs) > 0 && !ok {
		return false
	}
	if event.Aggregate().InstanceID != "" && query.instanceID != "" && event.Aggregate().InstanceID != query.instanceID {
		return false
	}
	if ok := isEventTypes(event, query.eventTypes...); len(query.eventTypes) > 0 && !ok {
		return false
	}
	return true
}

func (builder *SearchQueryBuilder) build(instanceID string) (*repository.SearchQuery, error) {
	if builder == nil ||
		len(builder.queries) < 1 ||
		builder.columns.Validate() != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-4m9gs", "builder invalid")
	}
	builder.instanceID = instanceID
	filters := make([][]*repository.Filter, len(builder.queries))

	for i, query := range builder.queries {
		for _, f := range []func() *repository.Filter{
			query.aggregateTypeFilter,
			query.aggregateIDFilter,
			query.eventTypeFilter,
			query.eventDataFilter,
			query.eventSequenceGreaterFilter,
			query.eventSequenceLessFilter,
			query.instanceIDFilter,
			query.excludedInstanceIDFilter,
			query.builder.resourceOwnerFilter,
			query.builder.instanceIDFilter,
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
		Columns: builder.columns,
		Limit:   builder.limit,
		Desc:    builder.desc,
		Filters: filters,
		Tx:      builder.tx,
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

func (query *SearchQuery) eventSequenceGreaterFilter() *repository.Filter {
	if query.eventSequenceGreater == 0 {
		return nil
	}
	sortOrder := repository.OperationGreater
	if query.builder.desc {
		sortOrder = repository.OperationLess
	}
	return repository.NewFilter(repository.FieldSequence, query.eventSequenceGreater, sortOrder)
}

func (query *SearchQuery) eventSequenceLessFilter() *repository.Filter {
	if query.eventSequenceLess == 0 {
		return nil
	}
	sortOrder := repository.OperationLess
	if query.builder.desc {
		sortOrder = repository.OperationGreater
	}
	return repository.NewFilter(repository.FieldSequence, query.eventSequenceLess, sortOrder)
}

func (query *SearchQuery) instanceIDFilter() *repository.Filter {
	if query.instanceID == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldInstanceID, query.instanceID, repository.OperationEquals)
}

func (query *SearchQuery) excludedInstanceIDFilter() *repository.Filter {
	if len(query.excludedInstanceIDs) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldInstanceID, query.excludedInstanceIDs, repository.OperationNotIn)
}

func (builder *SearchQueryBuilder) resourceOwnerFilter() *repository.Filter {
	if builder.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldResourceOwner, builder.resourceOwner, repository.OperationEquals)
}

func (builder *SearchQueryBuilder) instanceIDFilter() *repository.Filter {
	if builder.instanceID == "" {
		return nil
	}
	return repository.NewFilter(repository.FieldInstanceID, builder.instanceID, repository.OperationEquals)
}

func (query *SearchQuery) eventDataFilter() *repository.Filter {
	if len(query.eventData) == 0 {
		return nil
	}
	return repository.NewFilter(repository.FieldEventData, query.eventData, repository.OperationJSONContains)
}
