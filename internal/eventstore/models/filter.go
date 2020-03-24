package models

import (
	"github.com/caos/eventstore-lib/pkg/models"
)

type FilterEventsRequest struct {
	Limit           uint64
	AggregateID     string
	AggregateType   string
	LatestSequence  uint64
	ResourceOwner   string
	ModifierService string
	ModifierUser    string
	ModifierTenant  string
}

const (
	AggregateType models.Field = iota
	AggregateID
	LatestSequence
	ResourceOwner
	ModifierService
	ModifierUser
	ModifierTenant
)

const (
	Equals models.Operation = iota
	Greater
	Less
	In
)

type SearchQuery struct {
	limit   uint64
	desc    bool
	filters []*Filter
}

func NewSearchQuery(limit uint64, desc bool, filters ...*Filter) *SearchQuery {
	return &SearchQuery{
		limit:   limit,
		desc:    desc,
		filters: filters,
	}
}

func (q *SearchQuery) Limit() uint64 {
	return q.limit
}

func (q *SearchQuery) OrderDesc() bool {
	return q.desc
}

func (q *SearchQuery) Filters() []models.Filter {
	filters := make([]models.Filter, len(q.filters))
	for idx, filter := range q.filters {
		filters[idx] = filter
	}

	return filters
}

type Filter struct {
	field     models.Field
	value     interface{}
	operation models.Operation
}

func NewFilter(field models.Field, value interface{}, operation models.Operation) *Filter {
	return &Filter{
		field:     field,
		value:     value,
		operation: operation,
	}
}

func (f *Filter) GetField() models.Field {
	return f.field
}
func (f *Filter) GetOperation() models.Operation {
	return f.operation
}
func (f *Filter) GetValue() interface{} {
	return f.value
}
