package models

import "github.com/caos/zitadel/internal/errors"

type SearchQuery struct {
	limit   uint64
	desc    bool
	filters []*Filter
}

func NewSearchQuery(limit uint64, desc bool) *SearchQuery {
	return &SearchQuery{
		limit:   limit,
		desc:    desc,
		filters: make([]*Filter, 0, 4),
	}
}

func (q *SearchQuery) Limit() uint64 {
	return q.limit
}

func (q *SearchQuery) OrderDesc() bool {
	return q.desc
}

func (q *SearchQuery) Filters() []*Filter {
	return q.filters
}

func (q *SearchQuery) Validate() error {
	if q == nil {
		return errors.ThrowPreconditionFailed(nil, "MODEL-J5xQi", "search query is nil")
	}
	for _, filter := range q.filters {
		if err := filter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (q *SearchQuery) AggregateIDFilter(id string) *SearchQuery {
	return q.setFilter(NewFilter(AggregateID, id, Equals))
}

func (q *SearchQuery) AggregateTypeFilter(types ...string) *SearchQuery {
	return q.setFilter(NewFilter(AggregateType, types, In))
}

func (q *SearchQuery) LatestSequenceFilter(sequence uint64) *SearchQuery {
	sortOrder := Greater
	if q.desc {
		sortOrder = Less
	}
	return q.setFilter(NewFilter(LatestSequence, sequence, sortOrder))
}

func (q *SearchQuery) ResourceOwnerFilter(resourceOwner string) *SearchQuery {
	q.setFilter(NewFilter(ResourceOwner, resourceOwner, Equals))
	return q
}

func (q *SearchQuery) setFilter(filter *Filter) *SearchQuery {
	for _, f := range q.filters {
		if f.field == filter.field {
			f = filter
			return q
		}
	}
	q.filters = append(q.filters, filter)
	return q
}
