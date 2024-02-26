package eventstore

import (
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/v2/database"
)

var FilterMergeErr = errors.New("merge failed")

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
	desc       bool
	pagination *pagination

	aggregateFilters []*AggregateFilter
}

func (f *Filter) Desc() bool {
	return f.desc
}

func (f *Filter) Pagination() *pagination {
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

func Descending() FilterOpt {
	return func(f *Filter) {
		f.desc = true
	}
}

func AppendAggregateFilter(typ string, opts ...AggregateFilterOpt) FilterOpt {
	return AppendFilters(NewAggregateFilter(typ, opts...))
}

func AppendFilters(filters ...*AggregateFilter) FilterOpt {
	return func(mf *Filter) {
		mf.aggregateFilters = append(mf.aggregateFilters, filters...)
	}
}

func WithLimit(limit uint32) FilterOpt {
	return func(f *Filter) {
		f.ensurePagination()
		f.pagination.ensurePagination()

		f.pagination.pagination.Limit = limit
	}
}

func WithOffset(offset uint32) FilterOpt {
	return func(f *Filter) {
		f.ensurePagination()
		f.pagination.ensurePagination()

		f.pagination.pagination.Offset = offset
	}
}

// WithPosition prepares the condition as follows if inTxOrder is set:
// * position > or >=: ((position AND in_tx_order) OR position >/>=)
// * position < or <=: ((position AND in_tx_order) OR position </<=)
// TODO: between
func WithPosition(position database.NumberFilter[float64], inTxOrder database.NumberFilter[uint32]) FilterOpt {
	return func(f *Filter) {
		f.ensurePagination()
		f.pagination.ensurePosition()

		f.pagination.position.position = position
		f.pagination.position.inTxOrder = inTxOrder
	}
}

func WithPositionAtLeast(position float64, inTxOrder uint32) FilterOpt {
	return func(f *Filter) {
		f.ensurePagination()
		f.pagination.ensurePosition()

		f.pagination.position.position = database.NewNumberAtLeast(position)
		if inTxOrder > 0 {
			f.pagination.position.inTxOrder = database.NewNumberAtLeast(inTxOrder)
		}
	}
}

func (f *Filter) ensurePagination() {
	if f.pagination != nil {
		return
	}
	f.pagination = new(pagination)
}

func (p *pagination) ensurePosition() {
	if p.position != nil {
		return
	}
	p.position = new(positionFilter)
}

func (p *pagination) ensurePagination() {
	if p.pagination != nil {
		return
	}
	p.pagination = new(database.Pagination)
}

type pagination struct {
	pagination *database.Pagination
	position   *positionFilter
}

func (p *pagination) Pagination() *database.Pagination {
	return p.pagination
}

func (p *pagination) Position() *positionFilter {
	return p.position
}

type positionFilter struct {
	position  database.NumberFilter[float64]
	inTxOrder database.NumberFilter[uint32]
}

func (f *positionFilter) Position() database.NumberFilter[float64] {
	return f.position
}

func (f *positionFilter) InTxOrder() database.NumberFilter[uint32] {
	return f.inTxOrder
}

func NewAggregateFilter(typ string, opts ...AggregateFilterOpt) *AggregateFilter {
	f := &AggregateFilter{
		typ: database.NewTextEqual(typ),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

type AggregateFilter struct {
	typ    *database.TextFilter[string]
	id     database.Condition
	events []*EventFilter
}

func (f *AggregateFilter) Type() *database.TextFilter[string] {
	return f.typ
}

func (f *AggregateFilter) ID() database.Condition {
	return f.id
}

func (f *AggregateFilter) Events() []*EventFilter {
	return f.events
}

type AggregateFilterOpt func(f *AggregateFilter)

func WithAggregateID(id string) AggregateFilterOpt {
	return func(f *AggregateFilter) {
		f.id = database.NewTextEqual(id)
	}
}

func WithAggregateIDList(cond *database.ListFilter[string]) AggregateFilterOpt {
	return func(f *AggregateFilter) {
		f.id = cond
	}
}

func AppendEvent(opts ...EventFilterOpt) AggregateFilterOpt {
	return AppendEvents(NewEventFilter(opts...))
}

func AppendEvents(events ...*EventFilter) AggregateFilterOpt {
	return func(f *AggregateFilter) {
		f.events = append(f.events, events...)
	}
}

func NewEventFilter(opts ...EventFilterOpt) *EventFilter {
	f := new(EventFilter)

	for _, opt := range opts {
		opt(f)
	}

	return f
}

type EventFilter struct {
	typ       *database.TextFilter[string]
	revision  database.NumberFilter[uint16]
	createdAt database.NumberFilter[time.Time]
	sequence  database.NumberFilter[uint32]
	creator   database.Condition
}

func (f *EventFilter) Type() *database.TextFilter[string] {
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

func WithEventType(typ string) EventFilterOpt {
	return func(f *EventFilter) {
		f.typ = database.NewTextEqual(typ)
	}
}

func WithRevision(revision database.NumberFilter[uint16]) EventFilterOpt {
	return func(f *EventFilter) {
		f.revision = revision
	}
}

func WithCreatedAt(createdAt database.NumberFilter[time.Time]) EventFilterOpt {
	return func(f *EventFilter) {
		f.createdAt = createdAt
	}
}

func WithSequence(sequence database.NumberFilter[uint32]) EventFilterOpt {
	return func(f *EventFilter) {
		f.sequence = sequence
	}
}

func WithCreator(creator string) EventFilterOpt {
	return func(f *EventFilter) {
		f.creator = database.NewTextEqual(creator)
	}
}

func WithCreatorList(cond *database.ListFilter[string]) EventFilterOpt {
	return func(f *EventFilter) {
		f.creator = cond
	}
}
