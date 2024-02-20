package eventstore

import (
	"github.com/zitadel/zitadel/internal/v2/database"
)

type EventFilter struct {
	AggregateTypes database.Filter
	AggregateIDs   database.Filter
	EventTypes     database.Filter
	Sequence       database.Filter
}

func (filter *EventFilter) Filters() []database.Filter {
	filters := make([]database.Filter, 0, 4)

	if filter.AggregateTypes != nil {
		filters = append(filters, filter.AggregateTypes)
	}
	if filter.AggregateIDs != nil {
		filters = append(filters, filter.AggregateIDs)
	}
	if filter.EventTypes != nil {
		filters = append(filters, filter.EventTypes)
	}
	if filter.Sequence != nil {
		filters = append(filters, filter.Sequence)
	}

	return filters
}

func NewEventFilter(opts ...eventFilterOpt) *EventFilter {
	f := new(EventFilter)

	for _, opt := range opts {
		f = opt(f)
	}

	return f
}

type eventFilterOpt func(filter *EventFilter) *EventFilter

func FilterAggregateTypes(types ...string) eventFilterOpt {
	return func(filter *EventFilter) *EventFilter {
		if len(types) == 0 {
			return filter
		}
		if len(types) == 1 {
			filter.AggregateTypes = database.NewTextEqual(types[0])
			return filter
		}

		filter.AggregateTypes = database.NewListContains(types)
		return filter
	}
}

func FilterAggregateIDs(ids ...string) eventFilterOpt {
	return func(filter *EventFilter) *EventFilter {
		if len(ids) == 0 {
			return filter
		}
		if len(ids) == 1 {
			filter.AggregateIDs = database.NewTextEqual(ids[0])
			return filter
		}

		filter.AggregateIDs = database.NewListContains(ids)
		return filter
	}
}

func FilterEventTypes(types ...string) eventFilterOpt {
	return func(filter *EventFilter) *EventFilter {
		if len(types) == 0 {
			return filter
		}
		if len(types) == 1 {
			filter.EventTypes = database.NewTextEqual(types[0])
			return filter
		}

		filter.EventTypes = database.NewListContains(types)
		return filter
	}
}

func FilterSequenceEquals(pos float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if pos == 0 {
			return f
		}
		f.Sequence = database.NewNumberEquals(pos)
		return f
	}
}

func FilterSequenceAtLeast(pos float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if pos == 0 {
			return f
		}
		f.Sequence = database.NewNumberAtLeast(pos)
		return f
	}
}

func FilterSequenceGreater(pos float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if pos == 0 {
			return f
		}
		f.Sequence = database.NewNumberGreater(pos)
		return f
	}
}

func FilterSequenceAtMost(pos float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if pos == 0 {
			return f
		}
		f.Sequence = database.NewNumberAtMost(pos)
		return f
	}
}

func FilterSequenceLess(pos float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if pos == 0 {
			return f
		}
		f.Sequence = database.NewNumberLess(pos)
		return f
	}
}

func FilterSequenceBetween(min, max float64) eventFilterOpt {
	return func(f *EventFilter) *EventFilter {
		if min == 0 && max == 0 {
			return f
		}
		f.Sequence = database.NewNumberBetween(min, max)
		return f
	}
}
