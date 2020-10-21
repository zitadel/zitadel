package eventstore

import (
	"math"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

func testSetColumns(columns Columns) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.Columns(columns)
		return factory
	}
}

func testSetLimit(limit uint64) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.Limit(limit)
		return factory
	}
}

func testSetSequence(sequence uint64) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.SequenceGreater(sequence)
		return factory
	}
}

func testSetAggregateIDs(aggregateIDs ...string) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.AggregateIDs(aggregateIDs...)
		return factory
	}
}

func testSetEventTypes(eventTypes ...EventType) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.EventTypes(eventTypes...)
		return factory
	}
}

func testSetResourceOwner(resourceOwner string) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		factory = factory.ResourceOwner(resourceOwner)
		return factory
	}
}

func testSetSortOrder(asc bool) func(factory *SearchQueryFactory) *SearchQueryFactory {
	return func(factory *SearchQueryFactory) *SearchQueryFactory {
		if asc {
			factory = factory.OrderAsc()
		} else {
			factory = factory.OrderDesc()
		}
		return factory
	}
}

func TestSearchQueryFactorySetters(t *testing.T) {
	type args struct {
		columns        Columns
		aggregateTypes []AggregateType
		setters        []func(*SearchQueryFactory) *SearchQueryFactory
	}
	tests := []struct {
		name string
		args args
		res  *SearchQueryFactory
	}{
		{
			name: "New factory",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user", "org"},
			},
			res: &SearchQueryFactory{
				columns:        repository.Columns(ColumnsEvent),
				aggregateTypes: []AggregateType{"user", "org"},
			},
		},
		{
			name: "set columns",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetColumns(repository.ColumnsMaxSequence)},
			},
			res: &SearchQueryFactory{
				columns: repository.ColumnsMaxSequence,
			},
		},
		{
			name: "set limit",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetLimit(100)},
			},
			res: &SearchQueryFactory{
				limit: 100,
			},
		},
		{
			name: "set sequence",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetSequence(90)},
			},
			res: &SearchQueryFactory{
				eventSequence: 90,
			},
		},
		{
			name: "set aggregateIDs",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetAggregateIDs("1235", "09824")},
			},
			res: &SearchQueryFactory{
				aggregateIDs: []string{"1235", "09824"},
			},
		},
		{
			name: "set eventTypes",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetEventTypes("user.created", "user.updated")},
			},
			res: &SearchQueryFactory{
				eventTypes: []EventType{"user.created", "user.updated"},
			},
		},
		{
			name: "set resource owner",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetResourceOwner("hodor")},
			},
			res: &SearchQueryFactory{
				resourceOwner: "hodor",
			},
		},
		{
			name: "default search query",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{testSetAggregateIDs("1235", "024"), testSetSortOrder(false)},
			},
			res: &SearchQueryFactory{
				aggregateTypes: []AggregateType{"user"},
				aggregateIDs:   []string{"1235", "024"},
				desc:           true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSearchQueryFactory(tt.args.columns, tt.args.aggregateTypes...)
			for _, setter := range tt.args.setters {
				factory = setter(factory)
			}
			if !reflect.DeepEqual(factory, tt.res) {
				t.Errorf("NewSearchQueryFactory() = %v, want %v", factory, tt.res)
			}
		})
	}
}

func TestSearchQueryFactoryBuild(t *testing.T) {
	type args struct {
		columns        Columns
		aggregateTypes []AggregateType
		setters        []func(*SearchQueryFactory) *SearchQueryFactory
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
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{},
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate types",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user", "org"},
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, []AggregateType{"user", "org"}, repository.OperationIn),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequence(100),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    true,
					Limit:   5,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(true),
					testSetSequence(100),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   5,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequence(100),
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetSequence(8),
				},
			},
			res: res{
				isErr: nil,
				query: &repository.SearchQuery{
					Columns: repository.ColumnsEvent,
					Desc:    false,
					Limit:   0,
					Filters: []*repository.Filter{
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldEventType, EventType("user.created"), repository.OperationEquals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and event types",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldEventType, []EventType{"user.created", "user.changed"}, repository.OperationIn),
					},
				},
			},
		},
		{
			name: "filter aggregate type resource owner",
			args: args{
				columns:        ColumnsEvent,
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
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
						repository.NewFilter(repository.FieldAggregateType, AggregateType("user"), repository.OperationEquals),
						repository.NewFilter(repository.FieldResourceOwner, "hodor", repository.OperationEquals),
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
			factory := NewSearchQueryFactory(tt.args.columns, tt.args.aggregateTypes...)
			for _, f := range tt.args.setters {
				factory = f(factory)
			}
			query, err := factory.build()
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error(%T): %v", err, err)
				return
			}
			if err != nil && tt.res.isErr == nil {
				t.Errorf("no error expected: %v", err)
				return
			}

			if !reflect.DeepEqual(query, tt.res.query) {
				t.Errorf("NewSearchQueryFactory() = %+v, want %+v", factory, tt.res)
			}
		})
	}
}
