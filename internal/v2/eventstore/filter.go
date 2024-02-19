package eventstore

import (
	"context"
	"database/sql"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type Filter struct {
	Instances []string

	Limit  uint32
	Offset uint32

	Tx *sql.Tx

	Descending bool

	EventQueries []*EventQuery
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

func FilterAscending() filterOpt {
	return func(f *Filter) *Filter {
		f.Descending = true
		return f
	}
}

func FilterInstances(instance ...string) filterOpt {
	return func(f *Filter) *Filter {
		f.Instances = slices.Compact(append(f.Instances, instance...))
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
func FilterEventQuery(queries ...eventQueryOpt) filterOpt {
	query := new(EventQuery)

	for _, q := range queries {
		q(query)
	}

	return func(f *Filter) *Filter {
		f.EventQueries = append(f.EventQueries, query)
		return f
	}
}
