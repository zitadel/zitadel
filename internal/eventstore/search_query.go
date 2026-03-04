package eventstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SearchQueryBuilder represents the builder for your filter
// if invalid data are set the filter will fail
type SearchQueryBuilder struct {
	columns               Columns
	limit                 uint64
	offset                uint32
	desc                  bool
	resourceOwner         string
	instanceID            *string
	instanceIDs           []string
	editorUser            string
	queries               []*SearchQuery
	excludeAggregateIDs   *ExclusionQuery
	tx                    *sql.Tx
	positionAtLeast       decimal.Decimal
	awaitOpenTransactions bool
	creationDateAfter     time.Time
	creationDateBefore    time.Time
	eventSequenceGreater  uint64
}

func (b *SearchQueryBuilder) GetColumns() Columns {
	return b.columns
}

func (b *SearchQueryBuilder) GetLimit() uint64 {
	return b.limit
}

func (b *SearchQueryBuilder) GetOffset() uint32 {
	return b.offset
}

func (b *SearchQueryBuilder) GetDesc() bool {
	return b.desc
}

func (b *SearchQueryBuilder) GetResourceOwner() string {
	return b.resourceOwner
}

func (b *SearchQueryBuilder) GetInstanceID() *string {
	return b.instanceID
}

func (b *SearchQueryBuilder) GetInstanceIDs() []string {
	return b.instanceIDs
}

func (b *SearchQueryBuilder) GetEditorUser() string {
	return b.editorUser
}

func (b *SearchQueryBuilder) GetQueries() []*SearchQuery {
	return b.queries
}

func (b *SearchQueryBuilder) GetExcludeAggregateIDs() *ExclusionQuery {
	return b.excludeAggregateIDs
}

func (b *SearchQueryBuilder) GetTx() *sql.Tx {
	return b.tx
}

func (b SearchQueryBuilder) GetPositionAtLeast() decimal.Decimal {
	return b.positionAtLeast
}

func (b SearchQueryBuilder) GetAwaitOpenTransactions() bool {
	return b.awaitOpenTransactions
}

func (q SearchQueryBuilder) GetEventSequenceGreater() uint64 {
	return q.eventSequenceGreater
}

func (q SearchQueryBuilder) GetCreationDateAfter() time.Time {
	return q.creationDateAfter
}

func (q SearchQueryBuilder) GetCreationDateBefore() time.Time {
	return q.creationDateBefore
}

// ensureInstanceID makes sure that the instance id is always set
func (b *SearchQueryBuilder) ensureInstanceID(ctx context.Context) {
	if b.instanceID == nil && len(b.instanceIDs) == 0 && authz.GetInstance(ctx).InstanceID() != "" {
		b.InstanceID(authz.GetInstance(ctx).InstanceID())
	}
}

type SearchQuery struct {
	builder        *SearchQueryBuilder
	aggregateTypes []AggregateType
	aggregateIDs   []string
	eventTypes     []EventType
	eventData      map[string]interface{}
	positionAfter  decimal.Decimal
}

func (q SearchQuery) GetAggregateTypes() []AggregateType {
	return q.aggregateTypes
}

func (q SearchQuery) GetAggregateIDs() []string {
	return q.aggregateIDs
}

func (q SearchQuery) GetEventTypes() []EventType {
	return q.eventTypes
}

func (q SearchQuery) GetEventData() map[string]interface{} {
	return q.eventData
}

func (q SearchQuery) GetPositionAfter() decimal.Decimal {
	return q.positionAfter
}

type ExclusionQuery struct {
	builder        *SearchQueryBuilder
	aggregateTypes []AggregateType
	eventTypes     []EventType
}

func (q ExclusionQuery) GetAggregateTypes() []AggregateType {
	return q.aggregateTypes
}

func (q ExclusionQuery) GetEventTypes() []EventType {
	return q.eventTypes
}

// Columns defines which fields of the event are needed for the query
type Columns int8

const (
	//ColumnsEvent represents all fields of an event
	ColumnsEvent = iota + 1
	// ColumnsMaxPosition represents the latest sequence of the filtered events
	ColumnsMaxPosition
	// ColumnsInstanceIDs represents the instance ids of the filtered events
	ColumnsInstanceIDs

	columnsCount
)

func (c Columns) Validate() error {
	if c <= 0 || c >= columnsCount {
		return zerrors.ThrowPreconditionFailed(nil, "REPOS-x8R35", "column out of range")
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

func (builder *SearchQueryBuilder) Matches(commands ...Command) []Command {
	matches := make([]Command, 0, len(commands))
	for i, command := range commands {
		if builder.limit > 0 && builder.limit <= uint64(len(matches)) {
			break
		}
		if builder.offset > 0 && uint32(i) < builder.offset {
			continue
		}

		if builder.matchCommand(command) {
			matches = append(matches, command)
		}
	}

	return matches
}

type sequencer interface {
	Sequence() uint64
}

func (builder *SearchQueryBuilder) matchCommand(command Command) bool {
	if builder.resourceOwner != "" && command.Aggregate().ResourceOwner != builder.resourceOwner {
		return false
	}
	if command.Aggregate().InstanceID != "" && builder.instanceID != nil && *builder.instanceID != "" && command.Aggregate().InstanceID != *builder.instanceID {
		return false
	}
	if seq, ok := command.(sequencer); ok {
		if builder.eventSequenceGreater > 0 && seq.Sequence() <= builder.eventSequenceGreater {
			return false
		}
	}

	if len(builder.queries) == 0 {
		return true
	}
	for _, query := range builder.queries {
		if query.matches(command) {
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

// Limit defines how many events are returned maximally.
func (builder *SearchQueryBuilder) Offset(offset uint32) *SearchQueryBuilder {
	builder.offset = offset
	return builder
}

// ResourceOwner defines the resource owner (org or instance) of the events
func (builder *SearchQueryBuilder) ResourceOwner(resourceOwner string) *SearchQueryBuilder {
	builder.resourceOwner = resourceOwner
	return builder
}

// InstanceID defines the instanceID (system) of the events
func (builder *SearchQueryBuilder) InstanceID(instanceID string) *SearchQueryBuilder {
	builder.instanceID = &instanceID
	return builder
}

// InstanceIDs defines the instanceIDs (system) of the events
func (builder *SearchQueryBuilder) InstanceIDs(instanceIDs []string) *SearchQueryBuilder {
	builder.instanceIDs = instanceIDs
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

// PositionAtLeast filters for events which happened after the specified time
func (builder *SearchQueryBuilder) PositionAtLeast(position decimal.Decimal) *SearchQueryBuilder {
	builder.positionAtLeast = position
	return builder
}

// AwaitOpenTransactions filters for events which are older than the oldest transaction of the database
func (builder *SearchQueryBuilder) AwaitOpenTransactions() *SearchQueryBuilder {
	builder.awaitOpenTransactions = true
	return builder
}

// SequenceGreater filters for events with sequence greater the requested sequence
func (builder *SearchQueryBuilder) SequenceGreater(sequence uint64) *SearchQueryBuilder {
	builder.eventSequenceGreater = sequence
	return builder
}

// CreationDateAfter filters for events which happened after the specified time
func (builder *SearchQueryBuilder) CreationDateAfter(creationDate time.Time) *SearchQueryBuilder {
	if creationDate.IsZero() || creationDate.Unix() == 0 {
		return builder
	}
	builder.creationDateAfter = creationDate
	return builder
}

// CreationDateBefore filters for events which happened before the specified time
func (builder *SearchQueryBuilder) CreationDateBefore(creationDate time.Time) *SearchQueryBuilder {
	if creationDate.IsZero() || creationDate.Unix() == 0 {
		return builder
	}
	builder.creationDateBefore = creationDate
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

// ExcludeAggregateIDs excludes events from the aggregate IDs returned by the [ExclusionQuery].
// There can be only 1 exclusion query. Subsequent calls overwrite previous definitions.
func (builder *SearchQueryBuilder) ExcludeAggregateIDs() *ExclusionQuery {
	query := &ExclusionQuery{
		builder: builder,
	}
	builder.excludeAggregateIDs = query
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

// AggregateIDs filters for events with the given aggregate id's
func (query *SearchQuery) AggregateIDs(ids ...string) *SearchQuery {
	query.aggregateIDs = ids
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

func (query *SearchQuery) PositionAfter(position decimal.Decimal) *SearchQuery {
	query.positionAfter = position
	return query
}

// Builder returns the SearchQueryBuilder of the sub query
func (query *SearchQuery) Builder() *SearchQueryBuilder {
	return query.builder
}

func (query *SearchQuery) matches(command Command) bool {
	if ok := isAggregateTypes(command.Aggregate(), query.aggregateTypes...); len(query.aggregateTypes) > 0 && !ok {
		return false
	}
	if ok := isAggregateIDs(command.Aggregate(), query.aggregateIDs...); len(query.aggregateIDs) > 0 && !ok {
		return false
	}
	if ok := isEventTypes(command, query.eventTypes...); len(query.eventTypes) > 0 && !ok {
		return false
	}
	return true
}

// AggregateTypes filters for events with the given aggregate types
func (query *ExclusionQuery) AggregateTypes(types ...AggregateType) *ExclusionQuery {
	query.aggregateTypes = types
	return query
}

// EventTypes filters for events with the given event types
func (query *ExclusionQuery) EventTypes(types ...EventType) *ExclusionQuery {
	query.eventTypes = types
	return query
}

// Builder returns the SearchQueryBuilder of the sub query
func (query *ExclusionQuery) Builder() *SearchQueryBuilder {
	return query.builder
}
