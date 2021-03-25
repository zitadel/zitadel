package models

import (
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

type SearchQueryFactory struct {
	columns        Columns
	limit          uint64
	desc           bool
	aggregateTypes []AggregateType
	aggregateIDs   []string
	sequenceFrom   uint64
	sequenceTo     uint64
	eventTypes     []EventType
	resourceOwner  string
	creationDate   time.Time
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
	//insert new columns-types before this columnsCount because count is needed for validation
	columnsCount
)

//FactoryFromSearchQuery is deprecated because it's for migration purposes. use NewSearchQueryFactory
func FactoryFromSearchQuery(query *SearchQuery) *SearchQueryFactory {
	factory := &SearchQueryFactory{
		columns: Columns_Event,
		desc:    query.Desc,
		limit:   query.Limit,
	}

	for _, filter := range query.Filters {
		switch filter.field {
		case Field_AggregateType:
			factory = factory.aggregateTypesMig(filter.value.([]AggregateType)...)
		case Field_AggregateID:
			if aggregateID, ok := filter.value.(string); ok {
				factory = factory.AggregateIDs(aggregateID)
			} else if aggregateIDs, ok := filter.value.([]string); ok {
				factory = factory.AggregateIDs(aggregateIDs...)
			}
		case Field_LatestSequence:
			if filter.operation == Operation_Greater {
				factory = factory.SequenceGreater(filter.value.(uint64))
			} else {
				factory = factory.SequenceLess(filter.value.(uint64))
			}
		case Field_ResourceOwner:
			factory = factory.ResourceOwner(filter.value.(string))
		case Field_EventType:
			factory = factory.EventTypes(filter.value.([]EventType)...)
		case Field_EditorService, Field_EditorUser:
			logging.Log("MODEL-Mr0VN").WithField("value", filter.value).Panic("field not converted to factory")
		case Field_CreationDate:
			factory = factory.CreationDateNewer(filter.value.(time.Time))
		}
	}

	return factory
}

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
	factory.sequenceFrom = sequence
	return factory
}

func (factory *SearchQueryFactory) SequenceLess(sequence uint64) *SearchQueryFactory {
	factory.sequenceTo = sequence
	return factory
}

func (factory *SearchQueryFactory) AggregateIDs(ids ...string) *SearchQueryFactory {
	factory.aggregateIDs = ids
	return factory
}

func (factory *SearchQueryFactory) aggregateTypesMig(types ...AggregateType) *SearchQueryFactory {
	factory.aggregateTypes = types
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

func (factory *SearchQueryFactory) CreationDateNewer(time time.Time) *SearchQueryFactory {
	factory.creationDate = time
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
		factory.sequenceFromFilter,
		factory.sequenceToFilter,
		factory.eventTypeFilter,
		factory.resourceOwnerFilter,
		factory.creationDateNewerFilter,
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

func (factory *SearchQueryFactory) sequenceFromFilter() *Filter {
	if factory.sequenceFrom == 0 {
		return nil
	}
	sortOrder := Operation_Greater
	if factory.desc {
		sortOrder = Operation_Less
	}
	return NewFilter(Field_LatestSequence, factory.sequenceFrom, sortOrder)
}

func (factory *SearchQueryFactory) sequenceToFilter() *Filter {
	if factory.sequenceTo == 0 {
		return nil
	}
	sortOrder := Operation_Less
	if factory.desc {
		sortOrder = Operation_Greater
	}
	return NewFilter(Field_LatestSequence, factory.sequenceTo, sortOrder)
}

func (factory *SearchQueryFactory) resourceOwnerFilter() *Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return NewFilter(Field_ResourceOwner, factory.resourceOwner, Operation_Equals)
}

func (factory *SearchQueryFactory) creationDateNewerFilter() *Filter {
	if factory.creationDate.IsZero() {
		return nil
	}
	return NewFilter(Field_CreationDate, factory.creationDate, Operation_Greater)
}
