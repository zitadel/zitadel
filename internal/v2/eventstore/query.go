package eventstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/v2/database"
)

type Query struct {
	instance   string
	filters    []*Filter
	reducer    Reducer
	tx         *sql.Tx
	pagination *Pagination
	// TODO: await push
}

func (q *Query) Instance() string {
	return q.instance
}

func (q *Query) Filters() []*Filter {
	return q.filters
}

func (q *Query) Reducer() Reducer {
	return q.reducer
}

func (q *Query) Tx() *sql.Tx {
	return q.tx
}

func (q *Query) Pagination() *Pagination {
	q.ensurePagination()
	return q.pagination
}

func NewQuery(instance string, opts ...QueryOpt) *Query {
	query := &Query{
		instance: instance,
	}

	for _, opt := range opts {
		opt(query)
	}

	return query
}

type QueryOpt func(query *Query)

func QueryReducer(reducer Reducer) QueryOpt {
	return func(query *Query) {
		query.reducer = reducer
	}
}

func QueryTx(tx *sql.Tx) QueryOpt {
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

func AppendFilter(opts ...FilterOpt) QueryOpt {
	return AppendFilters(NewFilter(opts...))
}

var ErrFilterMerge = errors.New("merge failed")

type FilterCreator func() []*Filter

func MergeFilters(creators ...FilterCreator) []*Filter {
	filters := make([]*Filter, 0, len(creators))

	for _, creator := range creators {
		filters = append(filters, creator()...)
	}

	return filters
}

// Merge returns an error if filters diverge
func Merge(filters ...*Filter) []*Filter {
	if len(filters) == 1 {
		return filters
	}
	return filters

	// merged := filters[0]
	// for i := 1; i < len(filters); i++{
	// 	filter := filters[i]
	// 	if merged.desc != filter.desc {
	// 		return nil, FilterMergeErr
	// 	}

	// }

	// return merged, nil
}

type Filter struct {
	parent     *Query
	pagination *Pagination

	aggregateFilters []*AggregateFilter
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
		typ: database.NewTextEqual(typ),
	}

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

type AggregateFilter struct {
	typ    database.TextFilter[string]
	id     database.TextFilter[string]
	events []*EventFilter
}

func (f *AggregateFilter) Type() database.TextFilter[string] {
	return f.typ
}

func (f *AggregateFilter) ID() database.TextFilter[string] {
	return f.id
}

func (f *AggregateFilter) Events() []*EventFilter {
	return f.events
}

type AggregateFilterOpt func(f *AggregateFilter)

func AggregateID(id string) AggregateFilterOpt {
	return func(filter *AggregateFilter) {
		filter.id = database.NewTextEqual(id)
	}
}

func AggregateIDList(cond *database.ListFilter[string]) AggregateFilterOpt {
	return func(filter *AggregateFilter) {
		filter.id = cond
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

func NewEventFilter(opts ...EventFilterOpt) *EventFilter {
	filter := new(EventFilter)

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

type EventFilter struct {
	typ       database.TextFilter[string]
	revision  database.NumberFilter[uint16]
	createdAt database.NumberFilter[time.Time]
	sequence  database.NumberFilter[uint32]
	creator   database.Condition
}

func (f *EventFilter) Type() database.TextFilter[string] {
	return f.typ
}

func (f *EventFilter) Revision() database.NumberFilter[uint16] {
	return f.revision
}

func (f *EventFilter) CreatedAt() database.NumberFilter[time.Time] {
	return f.createdAt
}

func (f *EventFilter) Sequence() database.NumberFilter[uint32] {
	return f.sequence
}

func (f *EventFilter) Creator() database.Condition {
	return f.creator
}

type EventFilterOpt func(f *EventFilter)

func EventType(typ string) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.typ = database.NewTextEqual(typ)
	}
}

func EventRevision(revision database.NumberFilter[uint16]) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.revision = revision
	}
}

func EventCreatedAt(createdAt database.NumberFilter[time.Time]) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.createdAt = createdAt
	}
}

func EventSequence(sequence database.NumberFilter[uint32]) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.sequence = sequence
	}
}

func EventCreator(creator string) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.creator = database.NewTextEqual(creator)
	}
}

func EventCreatorList(cond *database.ListFilter[string]) EventFilterOpt {
	return func(filter *EventFilter) {
		filter.creator = cond
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

// WithPosition prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func PositionGreater(position float64, inPositionOrder uint32) paginationOpt {
	return func(p *Pagination) {
		p.ensurePosition()

		p.position.Position = position
		p.position.InPositionOrder = inPositionOrder
	}
}

// GlobalPositionGreater prepares the condition as follows
// if inPositionOrder is set: position = AND in_tx_order > OR or position >
// if inPositionOrder is NOT set: position >
func GlobalPositionGreater(position *GlobalPosition) paginationOpt {
	return PositionGreater(position.Position, position.InPositionOrder)
}

type Pagination struct {
	pagination *database.Pagination
	position   *GlobalPosition

	desc bool
}

type paginationOpt func(*Pagination)

func (p *Pagination) Pagination() *database.Pagination {
	if p == nil {
		return nil
	}
	return p.pagination
}

func (p *Pagination) Position() *GlobalPosition {
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
	p.position = new(GlobalPosition)
}

func Descending() paginationOpt {
	return func(p *Pagination) {
		p.desc = true
	}
}
