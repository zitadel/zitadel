package eventstore

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type SearchQueryFactory struct {
	columns        repository.Columns
	limit          uint64
	desc           bool
	aggregateTypes []AggregateType
	aggregateIDs   []string
	eventSequence  uint64
	eventTypes     []EventType
	resourceOwner  string
}

type Columns repository.Columns

const (
	Columns_Event        Columns = repository.Columns_Event
	Columns_Max_Sequence Columns = repository.Columns_Max_Sequence
)

type AggregateType repository.AggregateType
type EventType repository.EventType

func NewSearchQueryFactory(aggregateTypes ...AggregateType) *SearchQueryFactory {
	return &SearchQueryFactory{
		aggregateTypes: aggregateTypes,
	}
}

func (factory *SearchQueryFactory) Columns(columns Columns) *SearchQueryFactory {
	factory.columns = repository.Columns(columns)
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

func (factory *SearchQueryFactory) Build() (*repository.SearchQuery, error) {
	if factory == nil ||
		len(factory.aggregateTypes) < 1 ||
		(factory.columns < 0 || factory.columns >= repository.ColumnsCount) {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-tGAD3", "factory invalid")
	}
	filters := []*repository.Filter{
		factory.aggregateTypeFilter(),
	}

	for _, f := range []func() *repository.Filter{
		factory.aggregateIDFilter,
		factory.eventSequenceFilter,
		factory.eventTypeFilter,
		factory.resourceOwnerFilter,
	} {
		if filter := f(); filter != nil {
			filters = append(filters, filter)
		}
	}

	return &repository.SearchQuery{
		Columns: repository.Columns(factory.columns),
		Limit:   factory.limit,
		Desc:    factory.desc,
		Filters: filters,
	}, nil
}

func (factory *SearchQueryFactory) aggregateIDFilter() *repository.Filter {
	if len(factory.aggregateIDs) < 1 {
		return nil
	}
	if len(factory.aggregateIDs) == 1 {
		return repository.NewFilter(repository.Field_AggregateID, factory.aggregateIDs[0], repository.Operation_Equals)
	}
	return repository.NewFilter(repository.Field_AggregateID, factory.aggregateIDs, repository.Operation_In)
}

func (factory *SearchQueryFactory) eventTypeFilter() *repository.Filter {
	if len(factory.eventTypes) < 1 {
		return nil
	}
	if len(factory.eventTypes) == 1 {
		return repository.NewFilter(repository.Field_EventType, factory.eventTypes[0], repository.Operation_Equals)
	}
	return repository.NewFilter(repository.Field_EventType, factory.eventTypes, repository.Operation_In)
}

func (factory *SearchQueryFactory) aggregateTypeFilter() *repository.Filter {
	if len(factory.aggregateTypes) == 1 {
		return repository.NewFilter(repository.Field_AggregateType, factory.aggregateTypes[0], repository.Operation_Equals)
	}
	return repository.NewFilter(repository.Field_AggregateType, factory.aggregateTypes, repository.Operation_In)
}

func (factory *SearchQueryFactory) eventSequenceFilter() *repository.Filter {
	if factory.eventSequence == 0 {
		return nil
	}
	sortOrder := repository.Operation_Greater
	if factory.desc {
		sortOrder = repository.Operation_Less
	}
	return repository.NewFilter(repository.Field_LatestSequence, factory.eventSequence, sortOrder)
}

func (factory *SearchQueryFactory) resourceOwnerFilter() *repository.Filter {
	if factory.resourceOwner == "" {
		return nil
	}
	return repository.NewFilter(repository.Field_ResourceOwner, factory.resourceOwner, repository.Operation_Equals)
}
