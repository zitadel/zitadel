package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/zitadel/zitadel/internal/v2/database"
)

type Querier interface {
	healthier
	Query(ctx context.Context, query *Query) (eventCount int, err error)
}

type Query struct {
	instances  *filter[[]string]
	filters    []*Filter
	tx         *sql.Tx
	pagination *Pagination
	reducer    Reducer
	// TODO: await push
}

func (q *Query) Instance() database.Condition {
	return q.instances.condition
}

func (q *Query) Filters() []*Filter {
	return q.filters
}

func (q *Query) Tx() *sql.Tx {
	return q.tx
}

func (q *Query) Pagination() *Pagination {
	q.ensurePagination()
	return q.pagination
}

func (q *Query) Reduce(events ...*Event[StoragePayload]) error {
	return q.reducer.Reduce(events...)
}

func NewQuery(instance string, reducer Reducer, opts ...QueryOpt) *Query {
	query := &Query{
		reducer: reducer,
	}

	for _, opt := range append([]QueryOpt{SetInstance(instance)}, opts...) {
		opt(query)
	}

	return query
}

type QueryOpt func(q *Query)

func SetInstance(instance string) QueryOpt {
	return InstancesEqual(instance)
}

func InstancesEqual(instances ...string) QueryOpt {
	return func(q *Query) {
		var cond database.Condition
		switch len(instances) {
		case 0:
			return
		case 1:
			cond = database.NewTextEqual(instances[0])
		default:
			cond = database.NewListEquals(instances...)
		}
		q.instances = &filter[[]string]{
			condition: cond,
			value:     &instances,
		}
	}
}

func InstancesContains(instances ...string) QueryOpt {
	return func(f *Query) {
		var cond database.Condition
		switch len(instances) {
		case 0:
			return
		case 1:
			cond = database.NewTextEqual(instances[0])
		default:
			cond = database.NewListContains(instances...)
		}

		f.instances = &filter[[]string]{
			condition: cond,
			value:     &instances,
		}
	}
}

func InstancesNotContains(instances ...string) QueryOpt {
	return func(f *Query) {
		var cond database.Condition
		switch len(instances) {
		case 0:
			return
		case 1:
			cond = database.NewTextUnequal(instances[0])
		default:
			cond = database.NewListNotContains(instances...)
		}
		f.instances = &filter[[]string]{
			condition: cond,
			value:     &instances,
		}
	}
}

func SetQueryTx(tx *sql.Tx) QueryOpt {
	return func(query *Query) {
		query.tx = tx
	}
}

func QueryPagination(opts ...paginationOpt) QueryOpt {
	return func(query *Query) {
		query.ensurePagination()

		for _, opt := range opts {
			opt(query.pagination)
		}
	}
}

func (q *Query) ensurePagination() {
	if q.pagination != nil {
		return
	}
	q.pagination = new(Pagination)
}

func AppendFilters(filters ...*Filter) QueryOpt {
	return func(query *Query) {
		for _, filter := range filters {
			filter.parent = query
		}
		query.filters = append(query.filters, filters...)
	}
}

func SetFilters(filters ...*Filter) QueryOpt {
	return func(query *Query) {
		for _, filter := range filters {
			filter.parent = query
		}
		query.filters = filters
	}
}

func AppendFilter(opts ...FilterOpt) QueryOpt {
	return AppendFilters(NewFilter(opts...))
}

var ErrFilterMerge = errors.New("merge failed")

type FilterCreator func() []*Filter

func MergeFilters(filters ...[]*Filter) []*Filter {
	// TODO: improve merge by checking fields of filters and merge filters if possible
	// this will reduce cost of queries which do multiple filters
	return slices.Concat(filters...)
}

type Filter struct {
	parent     *Query
	pagination *Pagination

	aggregateFilters []*AggregateFilter
}

func (f *Filter) Parent() *Query {
	return f.parent
}

func (f *Filter) Pagination() *Pagination {
	if f.pagination == nil {
		return f.parent.Pagination()
	}
	return f.pagination
}

func (f *Filter) AggregateFilters() []*AggregateFilter {
	return f.aggregateFilters
}

func NewFilter(opts ...FilterOpt) *Filter {
	f := new(Filter)

	for _, opt := range opts {
		opt(f)
	}

	return f
}

type FilterOpt func(f *Filter)

func AppendAggregateFilter(typ string, opts ...AggregateFilterOpt) FilterOpt {
	return AppendAggregateFilters(NewAggregateFilter(typ, opts...))
}

func AppendAggregateFilters(filters ...*AggregateFilter) FilterOpt {
	return func(mf *Filter) {
		mf.aggregateFilters = append(mf.aggregateFilters, filters...)
	}
}

func SetAggregateFilters(filters ...*AggregateFilter) FilterOpt {
	return func(mf *Filter) {
		mf.aggregateFilters = filters
	}
}

func FilterPagination(opts ...paginationOpt) FilterOpt {
	return func(filter *Filter) {
		filter.ensurePagination()

		for _, opt := range opts {
			opt(filter.pagination)
		}
	}
}

func (f *Filter) ensurePagination() {
	if f.pagination != nil {
		return
	}
	f.pagination = new(Pagination)
}

func NewAggregateFilter(typ string, opts ...AggregateFilterOpt) *AggregateFilter {
	filter := &AggregateFilter{
		typ: typ,
	}

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

type AggregateFilter struct {
	typ    string
	ids    []string
	events []*EventFilter
}

func (f *AggregateFilter) Type() *database.TextFilter[string] {
	return database.NewTextEqual(f.typ)
}

func (f *AggregateFilter) IDs() database.Condition {
	if len(f.ids) == 0 {
		return nil
	}
	if len(f.ids) == 1 {
		return database.NewTextEqual(f.ids[0])
	}

	return database.NewListContains(f.ids...)
}

func (f *AggregateFilter) Events() []*EventFilter {
	return f.events
}

type AggregateFilterOpt func(f *AggregateFilter)

func SetAggregateID(id string) AggregateFilterOpt {
	return func(filter *AggregateFilter) {
		filter.ids = []string{id}
	}
}

func AppendAggregateIDs(ids ...string) AggregateFilterOpt {
	return func(f *AggregateFilter) {
		f.ids = append(f.ids, ids...)
	}
}

// AggregateIDs sets the given ids as search param
func AggregateIDs(ids ...string) AggregateFilterOpt {
	return func(f *AggregateFilter) {
		f.ids = ids
	}
}

func AppendEvent(opts ...EventFilterOpt) AggregateFilterOpt {
	return AppendEvents(NewEventFilter(opts...))
}

func AppendEvents(events ...*EventFilter) AggregateFilterOpt {
	return func(filter *AggregateFilter) {
		filter.events = append(filter.events, events...)
	}
}

func SetEvents(events ...*EventFilter) AggregateFilterOpt {
	return func(filter *AggregateFilter) {
		filter.events = events
	}
}

func NewEventFilter(opts ...EventFilterOpt) *EventFilter {
	filter := new(EventFilter)

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

type EventFilter struct {
	types     []string
	revision  *filter[uint16]
	createdAt *filter[time.Time]
	sequence  *filter[uint32]
	creators  *filter[[]string]
}

type filter[T any] struct {
	condition database.Condition
	// the following fields are considered as one of
	// you can either have value and max or value
	min, max *T
	value    *T
}

func (f *EventFilter) Types() database.Condition {
	switch len(f.types) {
	case 0:
		return nil
	case 1:
		return database.NewTextEqual(f.types[0])
	default:
		return database.NewListContains(f.types...)
	}
}

func (f *EventFilter) Revision() database.Condition {
	if f.revision == nil {
		return nil
	}
	return f.revision.condition
}

func (f *EventFilter) CreatedAt() database.Condition {
	if f.createdAt == nil {
		return nil
	}
	return f.createdAt.condition
}

func (f *EventFilter) Sequence() database.Condition {
	if f.sequence == nil {
		return nil
	}
	return f.sequence.condition
}

func (f *EventFilter) Creators() database.Condition {
	if f.creators == nil {
		return nil
	}
	return f.creators.condition
}

type EventFilterOpt func(f *EventFilter)

func SetEventType(typ string) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.types = []string{typ}
	}
}

// SetEventTypes overwrites the currently set types
func SetEventTypes(types ...string) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.types = types
	}
}

// AppendEventTypes appends the types the currently set types
func AppendEventTypes(types ...string) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.types = append(filter.types, types...)
	}
}

func EventRevisionEquals(revision uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberEquals(revision),
			value:     &revision,
		}
	}
}

func EventRevisionAtLeast(revision uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberAtLeast(revision),
			value:     &revision,
		}
	}
}

func EventRevisionGreater(revision uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberGreater(revision),
			value:     &revision,
		}
	}
}

func EventRevisionAtMost(revision uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberAtMost(revision),
			value:     &revision,
		}
	}
}

func EventRevisionLess(revision uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberLess(revision),
			value:     &revision,
		}
	}
}

func EventRevisionBetween(min, max uint16) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = &filter[uint16]{
			condition: database.NewNumberBetween(min, max),
			min:       &min,
			max:       &max,
		}
	}
}

func EventCreatedAtEquals(createdAt time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberEquals(createdAt),
			value:     &createdAt,
		}
	}
}

func EventCreatedAtAtLeast(createdAt time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberAtLeast(createdAt),
			value:     &createdAt,
		}
	}
}

func EventCreatedAtGreater(createdAt time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberGreater(createdAt),
			value:     &createdAt,
		}
	}
}

func EventCreatedAtAtMost(createdAt time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberAtMost(createdAt),
			value:     &createdAt,
		}
	}
}

func EventCreatedAtLess(createdAt time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberLess(createdAt),
			value:     &createdAt,
		}
	}
}

func EventCreatedAtBetween(min, max time.Time) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = &filter[time.Time]{
			condition: database.NewNumberBetween(min, max),
			min:       &min,
			max:       &max,
		}
	}
}

func EventSequenceEquals(sequence uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberEquals(sequence),
			value:     &sequence,
		}
	}
}

func EventSequenceAtLeast(sequence uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberAtLeast(sequence),
			value:     &sequence,
		}
	}
}

func EventSequenceGreater(sequence uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberGreater(sequence),
			value:     &sequence,
		}
	}
}

func EventSequenceAtMost(sequence uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberAtMost(sequence),
			value:     &sequence,
		}
	}
}

func EventSequenceLess(sequence uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberLess(sequence),
			value:     &sequence,
		}
	}
}

func EventSequenceBetween(min, max uint32) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = &filter[uint32]{
			condition: database.NewNumberBetween(min, max),
			min:       &min,
			max:       &max,
		}
	}
}

func EventCreatorsEqual(creators ...string) EventFilterOpt {
	return func(f *EventFilter) {
		var cond database.Condition
		switch len(creators) {
		case 0:
			return
		case 1:
			cond = database.NewTextEqual(creators[0])
		default:
			cond = database.NewListEquals(creators...)
		}
		f.creators = &filter[[]string]{
			condition: cond,
			value:     &creators,
		}
	}
}

func EventCreatorsContains(creators ...string) EventFilterOpt {
	return func(f *EventFilter) {
		var cond database.Condition
		switch len(creators) {
		case 0:
			return
		case 1:
			cond = database.NewTextEqual(creators[0])
		default:
			cond = database.NewListContains(creators...)
		}

		f.creators = &filter[[]string]{
			condition: cond,
			value:     &creators,
		}
	}
}

func EventCreatorsNotContains(creators ...string) EventFilterOpt {
	return func(f *EventFilter) {
		var cond database.Condition
		switch len(creators) {
		case 0:
			return
		case 1:
			cond = database.NewTextUnequal(creators[0])
		default:
			cond = database.NewListNotContains(creators...)
		}
		f.creators = &filter[[]string]{
			condition: cond,
			value:     &creators,
		}
	}
}

func Limit(limit uint32) paginationOpt {
	return func(p *Pagination) {
		p.ensurePagination()

		p.pagination.Limit = limit
	}
}

func Offset(offset uint32) paginationOpt {
	return func(p *Pagination) {
		p.ensurePagination()

		p.pagination.Offset = offset
	}
}

type PositionCondition struct {
	min, max *GlobalPosition
}

func (pc *PositionCondition) Max() *GlobalPosition {
	if pc == nil || pc.max == nil {
		return nil
	}
	max := *pc.max
	return &max
}

func (pc *PositionCondition) Min() *GlobalPosition {
	if pc == nil || pc.min == nil {
		return nil
	}
	min := *pc.min
	return &min
}

// PositionGreater prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func PositionGreater(position float64, inPositionOrder uint32) paginationOpt {
	return func(p *Pagination) {
		p.ensurePosition()
		p.position.min = &GlobalPosition{
			Position:        position,
			InPositionOrder: inPositionOrder,
		}
	}
}

// GlobalPositionGreater prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func GlobalPositionGreater(position *GlobalPosition) paginationOpt {
	return PositionGreater(position.Position, position.InPositionOrder)
}

// PositionLess prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func PositionLess(position float64, inPositionOrder uint32) paginationOpt {
	return func(p *Pagination) {
		p.ensurePosition()
		p.position.max = &GlobalPosition{
			Position:        position,
			InPositionOrder: inPositionOrder,
		}
	}
}

func PositionBetween(min, max *GlobalPosition) paginationOpt {
	return func(p *Pagination) {
		GlobalPositionGreater(min)(p)
		GlobalPositionLess(max)(p)
	}
}

// GlobalPositionLess prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func GlobalPositionLess(position *GlobalPosition) paginationOpt {
	return PositionLess(position.Position, position.InPositionOrder)
}

type Pagination struct {
	pagination *database.Pagination
	position   *PositionCondition

	desc bool
}

type paginationOpt func(*Pagination)

func (p *Pagination) Pagination() *database.Pagination {
	if p == nil {
		return nil
	}
	return p.pagination
}

func (p *Pagination) Position() *PositionCondition {
	if p == nil {
		return nil
	}
	return p.position
}

func (p *Pagination) Desc() bool {
	if p == nil {
		return false
	}

	return p.desc
}

func (p *Pagination) ensurePagination() {
	if p.pagination != nil {
		return
	}
	p.pagination = new(database.Pagination)
}

func (p *Pagination) ensurePosition() {
	if p.position != nil {
		return
	}
	p.position = new(PositionCondition)
}

func Descending() paginationOpt {
	return func(p *Pagination) {
		p.desc = true
	}
}
