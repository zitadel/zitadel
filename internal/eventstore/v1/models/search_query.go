package models

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
)

type SearchQueryFactory struct {
	columns Columns
	limit   uint64
	desc    bool
	queries []*query

	InstanceFiltered bool
}

type query struct {
	desc               bool
	aggregateTypes     []AggregateType
	aggregateIDs       []string
	sequenceFrom       uint64
	sequenceTo         uint64
	eventTypes         []EventType
	resourceOwner      string
	instanceID         string
	ignoredInstanceIDs []string
	creationDate       time.Time
	factory            *SearchQueryFactory
}

type searchQuery struct {
	Columns Columns
	Limit   uint64
	Desc    bool
	Filters [][]*Filter
}

type Columns int32

const (
	Columns_Event = iota
	Columns_Max_Sequence
	Columns_InstanceIDs
	// insert new columns-types before this columnsCount because count is needed for validation
	columnsCount
)

// FactoryFromSearchQuery is deprecated because it's for migration purposes. use NewSearchQueryFactory
func FactoryFromSearchQuery(q *SearchQuery) *SearchQueryFactory {
	factory := &SearchQueryFactory{
		columns: q.Columns,
		desc:    q.Desc,
		limit:   q.Limit,
		queries: make([]*query, len(q.Queries)),
	}

	for i, qq := range q.Queries {
		factory.queries[i] = &query{factory: factory}
		for _, filter := range qq.Filters {
			switch filter.field {
			case Field_AggregateType:
				factory.queries[i] = factory.queries[i].aggregateTypesMig(filter.value.([]AggregateType)...)
			case Field_AggregateID:
				if aggregateID, ok := filter.value.(string); ok {
					factory.queries[i] = factory.queries[i].AggregateIDs(aggregateID)
				} else if aggregateIDs, ok := filter.value.([]string); ok {
					factory.queries[i] = factory.queries[i].AggregateIDs(aggregateIDs...)
				}
			case Field_LatestSequence:
				if filter.operation == Operation_Greater {
					factory.queries[i] = factory.queries[i].SequenceGreater(filter.value.(uint64))
				} else {
					factory.queries[i] = factory.queries[i].SequenceLess(filter.value.(uint64))
				}
			case Field_ResourceOwner:
				factory.queries[i] = factory.queries[i].ResourceOwner(filter.value.(string))
			case Field_InstanceID:
				factory.InstanceFiltered = true
				if filter.operation == Operation_Equals {
					factory.queries[i] = factory.queries[i].InstanceID(filter.value.(string))
				} else if filter.operation == Operation_NotIn {
					factory.queries[i] = factory.queries[i].IgnoredInstanceIDs(filter.value.([]string)...)
				}
			case Field_EventType:
				factory.queries[i] = factory.queries[i].EventTypes(filter.value.([]EventType)...)
			case Field_EditorService, Field_EditorUser:
				logging.WithFields("value", filter.value).Panic("field not converted to factory")
			case Field_CreationDate:
				factory.queries[i] = factory.queries[i].CreationDateNewer(filter.value.(time.Time))
			}
		}
	}

	return factory
}

func NewSearchQueryFactory() *SearchQueryFactory {
	return &SearchQueryFactory{}
}

func (factory *SearchQueryFactory) Columns(columns Columns) *SearchQueryFactory {
	factory.columns = columns
	return factory
}

func (factory *SearchQueryFactory) Limit(limit uint64) *SearchQueryFactory {
	factory.limit = limit
	return factory
}

func (factory *SearchQueryFactory) OrderDesc() *SearchQueryFactory {
	factory.desc = true
	return factory
}

func (factory *SearchQueryFactory) OrderAsc() *SearchQueryFactory {
	factory.desc = false
	return factory
}

func (factory *SearchQueryFactory) AddQuery() *query {
	q := &query{factory: factory}
	factory.queries = append(factory.queries, q)
	return q
}

func (q *query) Factory() *SearchQueryFactory {
	return q.factory
}

func (q *query) SequenceGreater(sequence uint64) *query {
	q.sequenceFrom = sequence
	return q
}

func (q *query) SequenceLess(sequence uint64) *query {
	q.sequenceTo = sequence
	return q
}

func (q *query) AggregateTypes(types ...AggregateType) *query {
	q.aggregateTypes = types
	return q
}

func (q *query) AggregateIDs(ids ...string) *query {
	q.aggregateIDs = ids
	return q
}

func (q *query) aggregateTypesMig(types ...AggregateType) *query {
	q.aggregateTypes = types
	return q
}

func (q *query) EventTypes(types ...EventType) *query {
	q.eventTypes = types
	return q
}

func (q *query) ResourceOwner(resourceOwner string) *query {
	q.resourceOwner = resourceOwner
	return q
}

func (q *query) InstanceID(instanceID string) *query {
	q.instanceID = instanceID
	return q
}

func (q *query) IgnoredInstanceIDs(instanceIDs ...string) *query {
	q.ignoredInstanceIDs = instanceIDs
	return q
}

func (q *query) CreationDateNewer(time time.Time) *query {
	q.creationDate = time
	return q
}

func (factory *SearchQueryFactory) Build() (*searchQuery, error) {
	if factory == nil ||
		len(factory.queries) < 1 ||
		(factory.columns < 0 || factory.columns >= columnsCount) {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-tGAD3", "factory invalid")
	}
	filters := make([][]*Filter, len(factory.queries))

	for i, query := range factory.queries {
		for _, f := range []func() *Filter{
			query.aggregateTypeFilter,
			query.aggregateIDFilter,
			query.sequenceFromFilter,
			query.sequenceToFilter,
			query.eventTypeFilter,
			query.resourceOwnerFilter,
			query.instanceIDFilter,
			query.ignoredInstanceIDsFilter,
			query.creationDateNewerFilter,
		} {
			if filter := f(); filter != nil {
				filters[i] = append(filters[i], filter)
			}
		}
	}

	return &searchQuery{
		Columns: factory.columns,
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (q *query) aggregateIDFilter() *Filter {
	if len(q.aggregateIDs) < 1 {
		return nil
	}
	if len(q.aggregateIDs) == 1 {
		return NewFilter(Field_AggregateID, q.aggregateIDs[0], Operation_Equals)
	}
	return NewFilter(Field_AggregateID, q.aggregateIDs, Operation_In)
}

func (q *query) eventTypeFilter() *Filter {
	if len(q.eventTypes) < 1 {
		return nil
	}
	if len(q.eventTypes) == 1 {
		return NewFilter(Field_EventType, q.eventTypes[0], Operation_Equals)
	}
	return NewFilter(Field_EventType, q.eventTypes, Operation_In)
}

func (q *query) aggregateTypeFilter() *Filter {
	if len(q.aggregateTypes) < 1 {
		return nil
	}
	if len(q.aggregateTypes) == 1 {
		return NewFilter(Field_AggregateType, q.aggregateTypes[0], Operation_Equals)
	}
	return NewFilter(Field_AggregateType, q.aggregateTypes, Operation_In)
}

func (q *query) sequenceFromFilter() *Filter {
	if q.sequenceFrom == 0 {
		return nil
	}
	sortOrder := Operation_Greater
	if q.factory.desc {
		sortOrder = Operation_Less
	}
	return NewFilter(Field_LatestSequence, q.sequenceFrom, sortOrder)
}

func (q *query) sequenceToFilter() *Filter {
	if q.sequenceTo == 0 {
		return nil
	}
	sortOrder := Operation_Less
	if q.factory.desc {
		sortOrder = Operation_Greater
	}
	return NewFilter(Field_LatestSequence, q.sequenceTo, sortOrder)
}

func (q *query) resourceOwnerFilter() *Filter {
	if q.resourceOwner == "" {
		return nil
	}
	return NewFilter(Field_ResourceOwner, q.resourceOwner, Operation_Equals)
}

func (q *query) instanceIDFilter() *Filter {
	if q.instanceID == "" {
		return nil
	}
	return NewFilter(Field_InstanceID, q.instanceID, Operation_Equals)
}

func (q *query) ignoredInstanceIDsFilter() *Filter {
	if len(q.ignoredInstanceIDs) == 0 {
		return nil
	}
	return NewFilter(Field_InstanceID, q.ignoredInstanceIDs, Operation_NotIn)
}

func (q *query) creationDateNewerFilter() *Filter {
	if q.creationDate.IsZero() {
		return nil
	}
	return NewFilter(Field_CreationDate, q.creationDate, Operation_Greater)
}
