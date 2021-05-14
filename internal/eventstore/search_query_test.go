package eventstore

import (
	"math"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

func testSetColumns(columns Columns) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.Columns(columns)
		return builder
	}
}

func testSetLimit(limit uint64) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.Limit(limit)
		return builder
	}
}

func testSetSequenceGreater(sequence uint64) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.SequenceGreater(sequence)
		return builder
	}
}

func testSetSequenceLess(sequence uint64) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.SequenceLess(sequence)
		return builder
	}
}

func testSetAggregateIDs(aggregateIDs ...string) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.AggregateIDs(aggregateIDs...)
		return builder
	}
}

func testSetEventTypes(eventTypes ...EventType) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.EventTypes(eventTypes...)
		return builder
	}
}

func testSetResourceOwner(resourceOwner string) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.ResourceOwner(resourceOwner)
		return builder
	}
}

func testSetSortOrder(asc bool) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		if asc {
			builder = builder.OrderAsc()
		} else {
			builder = builder.OrderDesc()
		}
		return builder
	}
}

func TestSearchQuerybuilderSetters(t *testing.T) {
	type args struct {
		columns        Columns
		aggregateTypes []AggregateType
		setters        []func(*SearchQueryBuilder) *SearchQueryBuilder
	}
	tests := []struct {
		name string
		args args
		res  *SearchQueryBuilder
	}{
		{
			name: "New builder",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user", "org"},
			},
			res: &SearchQueryBuilder{
				columns:        repository.Columns(ColumnsEvent),
				aggregateTypes: []AggregateType{"user", "org"},
			},
		},
		{
			name: "set columns",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetColumns(repository.ColumnsMaxSequence)},
			},
			res: &SearchQueryBuilder{
				columns: repository.ColumnsMaxSequence,
			},
		},
		{
			name: "set limit",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetLimit(100)},
			},
			res: &SearchQueryBuilder{
				limit: 100,
			},
		},
		{
			name: "set sequence greater",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetSequenceGreater(90)},
			},
			res: &SearchQueryBuilder{
				eventSequenceGreater: 90,
			},
		},
		{
			name: "set sequence less",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetSequenceLess(90)},
			},
			res: &SearchQueryBuilder{
				eventSequenceLess: 90,
			},
		},
		{
			name: "set aggregateIDs",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetAggregateIDs("1235", "09824")},
			},
			res: &SearchQueryBuilder{
				aggregateIDs: []string{"1235", "09824"},
			},
		},
		{
			name: "set eventTypes",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetEventTypes("user.created", "user.updated")},
			},
			res: &SearchQueryBuilder{
				eventTypes: []EventType{"user.created", "user.updated"},
			},
		},
		{
			name: "set resource owner",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetResourceOwner("hodor")},
			},
			res: &SearchQueryBuilder{
				resourceOwner: "hodor",
			},
		},
		{
			name: "default search query",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters:        []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetAggregateIDs("1235", "024"), testSetSortOrder(false)},
			},
			res: &SearchQueryBuilder{
				aggregateTypes: []AggregateType{"user"},
				aggregateIDs:   []string{"1235", "024"},
				desc:           true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchQueryBuilder(tt.args.columns, tt.args.aggregateTypes...)
			for _, setter := range tt.args.setters {
				builder = setter(builder)
			}
			if !reflect.DeepEqual(builder, tt.res) {
				t.Errorf("NewSearchQuerybuilder() = %v, want %v", builder, tt.res)
			}
		})
	}
}

func TestSearchQuerybuilderBuild(t *testing.T) {
	type args struct {
		columns        Columns
		aggregateTypes []AggregateType
		setters        []func(*SearchQueryBuilder) *SearchQueryBuilder
	}
	type res struct {
		isErr func(err error) bool
		query *repository.SearchQuery
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no aggregate types",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{},
				setters:        []func(*SearchQueryBuilder) *SearchQueryBuilder{},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
				query: nil,
			},
		},
		{
			name: "invalid column (too low)",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetColumns(Columns(-1)),
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid column (too high)",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetColumns(math.MaxInt32),
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "filter aggregate type",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters:        []func(*SearchQueryBuilder) *SearchQueryBuilder{},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate types",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user", "org"},
				setters:        []func(*SearchQueryBuilder) *SearchQueryBuilder{},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, []repository.AggregateType{"user", "org"}, repository.OperationIn),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequenceGreater(100),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   5,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldSequence, uint64(100), repository.OperationLess),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, asc",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetLimit(5),
					testSetSortOrder(true),
					testSetSequenceGreater(100),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   5,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldSequence, uint64(100), repository.OperationGreater),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc, max event sequence cols",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequenceGreater(100),
					testSetColumns(repository.ColumnsMaxSequence),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsMaxSequence,
					Desc:    true,
					Limit:   5,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldSequence, uint64(100), repository.OperationLess),
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate id",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetAggregateIDs("1234"),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldAggregateID, "1234", repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate ids",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetAggregateIDs("1234", "0815"),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldAggregateID, []string{"1234", "0815"}, repository.OperationIn),
					},
				},
			},
		},
		{
			name: "filter aggregate type and sequence greater",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetSequenceGreater(8),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldSequence, uint64(8), repository.OperationGreater),
					},
				},
			},
		},
		{
			name: "filter aggregate type and event type",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetEventTypes("user.created"),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldEventType, repository.EventType("user.created"), repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and event types",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetEventTypes("user.created", "user.changed"),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldEventType, []repository.EventType{"user.created", "user.changed"}, repository.OperationIn),
					},
				},
			},
		},
		{
			name: "filter aggregate type resource owner",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetResourceOwner("hodor"),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldResourceOwner, "hodor", repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and sequence between",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{
					testSetSequenceGreater(8),
					testSetSequenceLess(16),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, repository.AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldSequence, uint64(8), repository.OperationGreater),
						repository.NewFilter(repository.FieldSequence, uint64(16), repository.OperationLess),
					},
				},
			},
		},
		{
			name: "column invalid",
			args: args{
				columns:        Columns(-1),
				aggregateTypes: []AggregateType{"user"},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSearchQueryBuilder(tt.args.columns, tt.args.aggregateTypes...)
			for _, f := range tt.args.setters {
				builder = f(builder)
			}
			query, err := builder.build()
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error(%T): %v", err, err)
				return
			}
			if err != nil && tt.res.isErr == nil {
				t.Errorf("no error expected: %v", err)
				return
			}

			if !reflect.DeepEqual(query, tt.res.query) {
				t.Errorf("NewSearchQuerybuilder() = %+v, want %+v", builder, tt.res.query)
			}
		})
	}
}
