package eventstore

import (
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
)

// SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
	columns         Columns
	limit           uint64
	desc            bool
	resourceOwner   string
	instanceID      string
	editorUser      string
	queries         []*SearchQuery
	tx              *sql.Tx
	allowTimeTravel bool
	positionAfter   float64
}

func (b *SearchQueryBuilder) GetColumns() Columns {
	return b.columns
}

func (b *SearchQueryBuilder) GetLimit() uint64 {
	return b.limit
}

func (b *SearchQueryBuilder) GetDesc() bool {
	return b.desc
}

func (b *SearchQueryBuilder) GetResourceOwner() string {
	return b.resourceOwner
}

func (b *SearchQueryBuilder) GetInstanceID() string {
	return b.instanceID
}

func (b *SearchQueryBuilder) GetEditorUser() string {
	return b.editorUser
}

func (b *SearchQueryBuilder) GetQueries() []*SearchQuery {
	return b.queries
}

func (b *SearchQueryBuilder) GetTx() *sql.Tx {
	return b.tx
}

func (b *SearchQueryBuilder) GetAllowTimeTravel() bool {
	return b.allowTimeTravel
}

func (b SearchQueryBuilder) GetPositionAfter() float64 {
	return b.positionAfter
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
	creationDateAfter    time.Time
}

func (q SearchQuery) GetAggregateTypes() []AggregateType {
	return q.aggregateTypes
}

func (q SearchQuery) GetAggregateIDs() []string {
	return q.aggregateIDs
}

func (q SearchQuery) GetInstanceID() string {
	return q.instanceID
}

func (q SearchQuery) GetExcludedInstanceIDs() []string {
	return q.excludedInstanceIDs
}

func (q SearchQuery) GetEventSequenceGreater() uint64 {
	return q.eventSequenceGreater
}

func (q SearchQuery) GetEventSequenceLess() uint64 {
	return q.eventSequenceLess
}

func (q SearchQuery) GetEventTypes() []EventType {
	return q.eventTypes
}

func (q SearchQuery) GetEventData() map[string]interface{} {
	return q.eventData
}

func (q SearchQuery) GetCreationDateAfter() time.Time {
	return q.creationDateAfter
}

// Columns defines which fields of the event are needed for the query
type Columns int8

const (
	//ColumnsEvent represents all fields of an event
	ColumnsEvent = iota + 1
	// ColumnsMaxSequence represents the latest sequence of the filtered events
	ColumnsMaxSequence
	// ColumnsInstanceIDs represents the instance ids of the filtered events
	ColumnsInstanceIDs

	columnsCount
)

func (c Columns) Validate() error {
	if c <= 0 || c >= columnsCount {
		return errors.ThrowPreconditionFailed(nil, "REPOS-x8R35", "column out of range")
	}
	return nil
}

// NewSearchQueryBuilder creates a new builder for event filters
// aggregateTypes must contain at least one aggregate type
func NewSearchQueryBuilder(columns Columns) *SearchQueryBuilder {
	return &SearchQueryBuilder{
		columns: columns,
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

// Columns defines which fields are set
func (builder *SearchQueryBuilder) Columns(columns Columns) *SearchQueryBuilder {
	builder.columns = columns
	return builder
}

// Limit defines how many events are returned maximally.
func (builder *SearchQueryBuilder) Limit(limit uint64) *SearchQueryBuilder {
	builder.limit = limit
	return builder
}

// ResourceOwner defines the resource owner (org) of the events
func (builder *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	builder.resourceOwner = resourceOwner
	return builder
}

// InstanceID defines the instanceID (system) of the events
func (builder *SearchQueryBuilder) InstanceID(instanceID string) *SearchQueryBuilder {
	builder.instanceID = instanceID
	return builder
}

// OrderDesc changes the sorting order of the returned events to descending
func (builder *SearchQueryBuilder) OrderDesc() *SearchQueryBuilder {
	builder.desc = true
	return builder
}

// OrderAsc changes the sorting order of the returned events to ascending
func (builder *SearchQueryBuilder) OrderAsc() *SearchQueryBuilder {
	builder.desc = false
	return builder
}

// SetTx ensures that the eventstore library uses the existing transaction
func (builder *SearchQueryBuilder) SetTx(tx *sql.Tx) *SearchQueryBuilder {
	builder.tx = tx
	return builder
}

func (builder *SearchQueryBuilder) EditorUser(id string) *SearchQueryBuilder {
	builder.editorUser = id
	return builder
}

// AllowTimeTravel activates the time travel feature of the database if supported
// The queries will be made based on the call time
func (builder *SearchQueryBuilder) AllowTimeTravel() *SearchQueryBuilder {
	builder.allowTimeTravel = true
	return builder
}

// PositionAfter filters for events which happened after the specified time
func (builder *SearchQueryBuilder) PositionAfter(position float64) *SearchQueryBuilder {
	builder.positionAfter = position
	return builder
}

// AddQuery creates a new sub query.
// All fields in the sub query are AND-connected in the storage request.
// Multiple sub queries are OR-connected in the storage request.
func (builder *SearchQueryBuilder) AddQuery() *SearchQuery {
	query := &SearchQuery{
		builder: builder,
	}
	builder.queries = append(builder.queries, query)

	return query
}

// Or creates a new sub query on the search query builder
func (query SearchQuery) Or() *SearchQuery {
	return query.builder.AddQuery()
}

// AggregateTypes filters for events with the given aggregate types
func (query *SearchQuery) AggregateTypes(types ...AggregateType) *SearchQuery {
	query.aggregateTypes = types
	return query
}

// SequenceGreater filters for events with sequence greater the requested sequence
func (query *SearchQuery) SequenceGreater(sequence uint64) *SearchQuery {
	query.eventSequenceGreater = sequence
	return query
}

// SequenceLess filters for events with sequence less the requested sequence
func (query *SearchQuery) SequenceLess(sequence uint64) *SearchQuery {
	query.eventSequenceLess = sequence
	return query
}

// AggregateIDs filters for events with the given aggregate id's
func (query *SearchQuery) AggregateIDs(ids ...string) *SearchQuery {
	query.aggregateIDs = ids
	return query
}

// InstanceID filters for events with the given instanceID
func (query *SearchQuery) InstanceID(instanceID string) *SearchQuery {
	query.instanceID = instanceID
	return query
}

// ExcludedInstanceID filters for events not having the given instanceIDs
func (query *SearchQuery) ExcludedInstanceID(instanceIDs ...string) *SearchQuery {
	query.excludedInstanceIDs = instanceIDs
	return query
}

// CreationDateAfter filters for events which happened after the specified time
func (query *SearchQuery) CreationDateAfter(time time.Time) *SearchQuery {
	query.creationDateAfter = time
	return query
}

// EventTypes filters for events with the given event types
func (query *SearchQuery) EventTypes(types ...EventType) *SearchQuery {
	query.eventTypes = types
	return query
}

// EventData filters for events with the given event data.
// Use this call with care as it will be slower than the other filters.
func (query *SearchQuery) EventData(data map[string]interface{}) *SearchQuery {
	query.eventData = data
	return query
}

// Builder returns the SearchQueryBuilder of the sub query
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
