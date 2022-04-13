package models

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
)

//SearchQuery is deprecated. Use SearchQueryFactory
type SearchQuery struct {
	Limit   uint64
	Desc    bool
	Filters []*Filter
	Queries []*Query
}

type Query struct {
	searchQuery *SearchQuery
	Filters     []*Filter
}

//NewSearchQuery is deprecated. Use SearchQueryFactory
func NewSearchQuery() *SearchQuery {
	return &SearchQuery{
		Filters: make([]*Filter, 0, 4),
		Queries: make([]*Query, 0),
	}
}

func (q *SearchQuery) AddQuery() *Query {
	query := &Query{
		searchQuery: q,
	}
	q.Queries = append(q.Queries, query)

	return query
}

//SearchQuery returns the SearchQuery of the sub query
func (q *Query) SearchQuery() *SearchQuery {
	return q.searchQuery
}
func (q *Query) setFilter(filter *Filter) *Query {
	for i, f := range q.Filters {
		if f.field == filter.field && f.field != Field_LatestSequence {
			q.Filters[i] = filter
			return q
		}
	}
	q.Filters = append(q.Filters, filter)
	return q
}

func (q *SearchQuery) SetLimit(limit uint64) *SearchQuery {
	q.Limit = limit
	return q
}

func (q *SearchQuery) OrderDesc() *SearchQuery {
	q.Desc = true
	return q
}

func (q *SearchQuery) OrderAsc() *SearchQuery {
	q.Desc = false
	return q
}

func (q *Query) AggregateIDFilter(id string) *Query {
	return q.setFilter(NewFilter(Field_AggregateID, id, Operation_Equals))
}

func (q *Query) AggregateIDsFilter(ids ...string) *Query {
	return q.setFilter(NewFilter(Field_AggregateID, ids, Operation_In))
}

func (q *Query) AggregateTypeFilter(types ...AggregateType) *Query {
	return q.setFilter(NewFilter(Field_AggregateType, types, Operation_In))
}

func (q *Query) EventTypesFilter(types ...EventType) *Query {
	return q.setFilter(NewFilter(Field_EventType, types, Operation_In))
}

func (q *Query) LatestSequenceFilter(sequence uint64) *Query {
	if sequence == 0 {
		return q
	}
	sortOrder := Operation_Greater
	return q.setFilter(NewFilter(Field_LatestSequence, sequence, sortOrder))
}

func (q *Query) SequenceBetween(from, to uint64) *Query {
	q.setFilter(NewFilter(Field_LatestSequence, from, Operation_Greater))
	q.setFilter(NewFilter(Field_LatestSequence, to, Operation_Less))
	return q
}

func (q *Query) ResourceOwnerFilter(resourceOwner string) *Query {
	return q.setFilter(NewFilter(Field_ResourceOwner, resourceOwner, Operation_Equals))
}

func (q *Query) InstanceIDFilter(instanceID string) *Query {
	return q.setFilter(NewFilter(Field_InstanceID, instanceID, Operation_Equals))
}

func (q *Query) IgnoredInstanceIDsFilter(instanceIDs ...string) *Query {
	return q.setFilter(NewFilter(Field_InstanceID, instanceIDs, Operation_NotIn))
}

func (q *Query) CreationDateNewerFilter(time time.Time) *Query {
	return q.setFilter(NewFilter(Field_CreationDate, time, Operation_Greater))
}

func (q *SearchQuery) setFilter(filter *Filter) *SearchQuery {
	for i, f := range q.Filters {
		if f.field == filter.field && f.field != Field_LatestSequence {
			q.Filters[i] = filter
			return q
		}
	}
	q.Filters = append(q.Filters, filter)
	return q
}

func (q *SearchQuery) Validate() error {
	if q == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-J5xQi", "search query is nil")
	}
	if len(q.Queries) == 0 {
		return errors.ThrowPreconditionFailed(nil, "MODEL-pF3DR", "no filters set")
	}
	for _, query := range q.Queries {
		for _, filter := range query.Filters {
			if err := filter.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
