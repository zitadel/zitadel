package models

import (
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/errors"
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
				aggregateTypes: []AggregateType{"user", "org"},
			},
			res: &SearchQueryFactory{
				aggregateTypes: []AggregateType{"user", "org"},
			},
		},
		{
			name: "set columns",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testSetColumns(Columns_Max_Sequence)},
			},
			res: &SearchQueryFactory{
				columns: Columns_Max_Sequence,
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
				sequenceFrom: 90,
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
			factory := NewSearchQueryFactory(tt.args.aggregateTypes...)
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
		aggregateTypes []AggregateType
		setters        []func(*SearchQueryFactory) *SearchQueryFactory
	}
	type res struct {
		isErr func(err error) bool
		query *searchQuery
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no aggregate types",
			args: args{
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
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetColumns(columnsCount),
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "filter aggregate type",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
					},
				},
			},
		},
		{
			name: "filter aggregate types",
			args: args{
				aggregateTypes: []AggregateType{"user", "org"},
				setters:        []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, []AggregateType{"user", "org"}, Operation_In),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequence(100),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    true,
					Limit:   5,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_LatestSequence, uint64(100), Operation_Less),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, asc",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(true),
					testSetSequence(100),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   5,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_LatestSequence, uint64(100), Operation_Greater),
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc, max event sequence cols",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetSequence(100),
					testSetColumns(Columns_Max_Sequence),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: Columns_Max_Sequence,
					Desc:    true,
					Limit:   5,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_LatestSequence, uint64(100), Operation_Less),
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate id",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetAggregateIDs("1234"),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_AggregateID, "1234", Operation_Equals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate ids",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetAggregateIDs("1234", "0815"),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_AggregateID, []string{"1234", "0815"}, Operation_In),
					},
				},
			},
		},
		{
			name: "filter aggregate type and sequence greater",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetSequence(8),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_LatestSequence, uint64(8), Operation_Greater),
					},
				},
			},
		},
		{
			name: "filter aggregate type and event type",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetEventTypes("user.created"),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_EventType, EventType("user.created"), Operation_Equals),
					},
				},
			},
		},
		{
			name: "filter aggregate type and event types",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetEventTypes("user.created", "user.changed"),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_EventType, []EventType{"user.created", "user.changed"}, Operation_In),
					},
				},
			},
		},
		{
			name: "filter aggregate type resource owner",
			args: args{
				aggregateTypes: []AggregateType{"user"},
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetResourceOwner("hodor"),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: []*Filter{
						NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						NewFilter(Field_ResourceOwner, "hodor", Operation_Equals),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSearchQueryFactory(tt.args.aggregateTypes...)
			for _, f := range tt.args.setters {
				factory = f(factory)
			}
			query, err := factory.Build()
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error: %v", err)
				return
			}
			if err != nil && tt.res.isErr == nil {
				t.Errorf("no error expected: %v", err)
				return
			}

			if !reflect.DeepEqual(query, tt.res.query) {
				t.Errorf("NewSearchQueryFactory() = %v, want %v", factory, tt.res)
			}
		})
	}
}
