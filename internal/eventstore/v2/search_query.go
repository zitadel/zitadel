package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
)

type SearchQueryFactory struct {
	columns        Columns
	limit          uint64
	desc           bool
	aggregateTypes []AggregateType
	aggregateIDs   []string
	eventSequence  uint64
	eventTypes     []EventType
	resourceOwner  string
}

type searchQuery struct {
	Columns Columns
	Limit   uint64
	Desc    bool
	Filters []*Filter
}

type Columns int32

const (
	Columns_Event = iota
	Columns_Max_Sequence
	//insert new columns-types above this columnsCount because count is needed for validation
	columnsCount
)

func NewSearchQueryFactory(aggregateTypes ...AggregateType) *SearchQueryFactory {
	return &SearchQueryFactory{
		aggregateTypes: aggregateTypes,
	}
}

func (factory *SearchQueryFactory) Columns(columns Columns) *SearchQueryFactory {
	factory.columns = columns
	return factory
}

func (factory *SearchQueryFactory) Limit(limit uint64) *SearchQueryFactory {
	factory.limit = limit
	return factory
}

func (factory *SearchQueryFactory) SequenceGreater(sequence uint64) *SearchQueryFactory {
	factory.eventSequence = sequence
	return factory
}

func (factory *SearchQueryFactory) AggregateIDs(ids ...string) *SearchQueryFactory {
	factory.aggregateIDs = ids
	return factory
}

func (factory *SearchQueryFactory) EventTypes(types ...EventType) *SearchQueryFactory {
	factory.eventTypes = types
	return factory
}

func (factory *SearchQueryFactory) ResourceOwner(resourceOwner string) *SearchQueryFactory {
	factory.resourceOwner = resourceOwner
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

func (factory *SearchQueryFactory) Build() (*searchQuery, error) {
	if factory == nil ||
		len(factory.aggregateTypes) < 1 ||
		(factory.columns < 0 || factory.columns >= columnsCount) {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-tGAD3", "factory invalid")
	}
	filters := []*Filter{
		factory.aggregateTypeFilter(),
	}

	for _, f := range []func() *Filter{
		factory.aggregateIDFilter,
		factory.eventSequenceFilter,
		factory.eventTypeFilter,
		factory.resourceOwnerFilter,
	} {
		if filter := f(); filter != nil {
			filters = append(filters, filter)
		}
	}

	return &searchQuery{
		Columns: factory.columns,
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (factory *SearchQueryFactory) aggregateIDFilter() *Filter {
	if len(factory.aggregateIDs) < 1 {
		return nil
	}
	if len(factory.aggregateIDs) == 1 {
		return NewFilter(Field_AggregateID, factory.aggregateIDs[0], Operation_Equals)
	}
	return NewFilter(Field_AggregateID, factory.aggregateIDs, Operation_In)
}

func (factory *SearchQueryFactory) eventTypeFilter() *Filter {
	if len(factory.eventTypes) < 1 {
		return nil
	}
	if len(factory.eventTypes) == 1 {
		return NewFilter(Field_EventType, factory.eventTypes[0], Operation_Equals)
	}
	return NewFilter(Field_EventType, factory.eventTypes, Operation_In)
}

func (factory *SearchQueryFactory) aggregateTypeFilter() *Filter {
	if len(factory.aggregateTypes) == 1 {
		return NewFilter(Field_AggregateType, factory.aggregateTypes[0], Operation_Equals)
	}
	return NewFilter(Field_AggregateType, factory.aggregateTypes, Operation_In)
}

func (factory *SearchQueryFactory) eventSequenceFilter() *Filter {
	if factory.eventSequence == 0 {
		return nil
	}
	sortOrder := Operation_Greater
	if factory.desc {
		sortOrder = Operation_Less
	}
	return NewFilter(Field_LatestSequence, factory.eventSequence, sortOrder)
}

func (factory *SearchQueryFactory) resourceOwnerFilter() *Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return NewFilter(Field_ResourceOwner, factory.resourceOwner, Operation_Equals)
}
