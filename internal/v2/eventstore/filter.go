package eventstore

import (
	"context"
	"database/sql"
	"slices"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/database"
)

// TODO: improve instances, position, eventFilters
func MergeFilters(filters ...*Filter) *Filter {
	merged := new(Filter)

	for _, filter := range filters {
		if merged.Instances == nil {
			merged.Instances = filter.Instances
		}
		if merged.Position == nil {
			merged.Position = filter.Position
		}

		if merged.Limit < filter.Limit {
			merged.Limit = filter.Limit
		}
		if merged.Offset == 0 || merged.Offset > filter.Offset {
			merged.Offset = filter.Offset
		}

		if merged.Tx == nil {
			merged.Tx = filter.Tx
		}

		if merged.Descending != filter.Descending {
			logging.Panic("Filter: sort order of filter was not equal")
		}

		merged.EventFilters = mergeEventFilters(merged.EventFilters, filter.EventFilters)
	}

	return merged
}

// TODO: are there possibility to merge specific filters?
func mergeEventFilters(existing []*EventFilter, additional []*EventFilter) []*EventFilter {
	if len(existing) == 0 {
		return additional
	}

	merged := slices.Clone(existing)
	merged = append(merged, additional...)

	return merged
}

type Filter struct {
	Instances database.Filter
	Position  database.Filter

	Limit  uint32
	Offset uint32

	Tx *sql.Tx

	Descending bool

	EventFilters []*EventFilter
}

type filterOpt func(f *Filter) *Filter

func NewFilter(ctx context.Context, opts ...filterOpt) *Filter {
	f := new(Filter)

	opts = append(
		[]filterOpt{
			FilterInstances(authz.GetInstance(ctx).InstanceID()),
		},
		opts...,
	)

	for _, opt := range opts {
		f = opt(f)
	}

	return f
}

func FilterLimit(limit uint32) filterOpt {
	return func(f *Filter) *Filter {
		f.Limit = limit
		return f
	}
}

func FilterOffset(offset uint32) filterOpt {
	return func(f *Filter) *Filter {
		f.Offset = offset
		return f
	}
}

func FilterDescending() filterOpt {
	return func(f *Filter) *Filter {
		f.Descending = true
		return f
	}
}

func FilterInstances(instances ...string) filterOpt {
	return func(f *Filter) *Filter {
		// f.Instances = slices.Compact(append(f.Instances, instance...))
		if len(instances) == 0 {
			return f
		}
		if len(instances) == 1 {
			f.Instances = database.NewTextEqual(instances[0])
			return f
		}

		f.Instances = database.NewListContains(instances)
		return f
	}
}

func FilterInTx(tx *sql.Tx) filterOpt {
	return func(f *Filter) *Filter {
		f.Tx = tx
		return f
	}
}

// FilterEventQuery creates a new sub clause for the given options
// sub clauses allow filters on specific events
// sub clauses are OR connected in the storage query
func FilterEventQuery(opts ...eventFilterOpt) filterOpt {
	return func(f *Filter) *Filter {
		f.EventFilters = append(f.EventFilters, NewEventFilter(opts...))
		return f
	}
}

func FilterPositionEquals(pos float64) filterOpt {
	return func(f *Filter) *Filter {
		if pos == 0 {
			return f
		}
		f.Position = database.NewNumberEquals(pos)
		return f
	}
}

func FilterPositionAtLeast(pos float64) filterOpt {
	return func(f *Filter) *Filter {
		if pos == 0 {
			return f
		}
		f.Position = database.NewNumberAtLeast(pos)
		return f
	}
}

func FilterPositionGreater(pos float64) filterOpt {
	return func(f *Filter) *Filter {
		if pos == 0 {
			return f
		}
		f.Position = database.NewNumberGreater(pos)
		return f
	}
}

func FilterPositionAtMost(pos float64) filterOpt {
	return func(f *Filter) *Filter {
		if pos == 0 {
			return f
		}
		f.Position = database.NewNumberAtMost(pos)
		return f
	}
}

func FilterPositionLess(pos float64) filterOpt {
	return func(f *Filter) *Filter {
		if pos == 0 {
			return f
		}
		f.Position = database.NewNumberLess(pos)
		return f
	}
}

func FilterPositionBetween(min, max float64) filterOpt {
	return func(f *Filter) *Filter {
		if min == 0 && max == 0 {
			return f
		}
		f.Position = database.NewNumberBetween(min, max)
		return f
	}
}
