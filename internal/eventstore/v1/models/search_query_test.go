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

func testAddQuery(queryFuncs ...func(*query) *query) func(*SearchQueryFactory) *SearchQueryFactory {
	return func(builder *SearchQueryFactory) *SearchQueryFactory {
		query := builder.AddQuery()
		for _, queryFunc := range queryFuncs {
			queryFunc(query)
		}
		return query.Factory()
	}
}

func testSetSequence(sequence uint64) func(*query) *query {
	return func(q *query) *query {
		q.SequenceGreater(sequence)
		return q
	}
}

func testSetAggregateIDs(aggregateIDs ...string) func(*query) *query {
	return func(q *query) *query {
		q.AggregateIDs(aggregateIDs...)
		return q
	}
}

func testSetAggregateTypes(aggregateTypes ...AggregateType) func(*query) *query {
	return func(q *query) *query {
		q.AggregateTypes(aggregateTypes...)
		return q
	}
}

func testSetEventTypes(eventTypes ...EventType) func(*query) *query {
	return func(q *query) *query {
		q.EventTypes(eventTypes...)
		return q
	}
}

func testSetResourceOwner(resourceOwner string) func(*query) *query {
	return func(q *query) *query {
		q.ResourceOwner(resourceOwner)
		return q
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

func assertFactory(t *testing.T, want, got *SearchQueryFactory) {
	t.Helper()

	if got.columns != want.columns {
		t.Errorf("wrong column: got: %v want: %v", got.columns, want.columns)
	}
	if got.desc != want.desc {
		t.Errorf("wrong desc: got: %v want: %v", got.desc, want.desc)
	}
	if got.limit != want.limit {
		t.Errorf("wrong limit: got: %v want: %v", got.limit, want.limit)
	}
	if len(got.queries) != len(want.queries) {
		t.Errorf("wrong length of queries: got: %v want: %v", len(got.queries), len(want.queries))
	}

	for i, query := range got.queries {
		assertQuery(t, i, want.queries[i], query)
	}
}

func assertQuery(t *testing.T, i int, want, got *query) {
	t.Helper()

	if !reflect.DeepEqual(got.aggregateIDs, want.aggregateIDs) {
		t.Errorf("wrong aggregateIDs in query %d : got: %v want: %v", i, got.aggregateIDs, want.aggregateIDs)
	}
	if !reflect.DeepEqual(got.aggregateTypes, want.aggregateTypes) {
		t.Errorf("wrong aggregateTypes in query %d : got: %v want: %v", i, got.aggregateTypes, want.aggregateTypes)
	}
	if got.sequenceFrom != want.sequenceFrom {
		t.Errorf("wrong sequenceFrom in query %d : got: %v want: %v", i, got.sequenceFrom, want.sequenceFrom)
	}
	if got.sequenceTo != want.sequenceTo {
		t.Errorf("wrong sequenceTo in query %d : got: %v want: %v", i, got.sequenceTo, want.sequenceTo)
	}
	if !reflect.DeepEqual(got.eventTypes, want.eventTypes) {
		t.Errorf("wrong eventTypes in query %d : got: %v want: %v", i, got.eventTypes, want.eventTypes)
	}
}

func TestSearchQueryFactorySetters(t *testing.T) {
	type args struct {
		setters []func(*SearchQueryFactory) *SearchQueryFactory
	}
	tests := []struct {
		name string
		args args
		res  *SearchQueryFactory
	}{
		{
			name: "New factory",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: &SearchQueryFactory{},
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetSequence(90))},
			},
			res: &SearchQueryFactory{
				queries: []*query{
					{
						sequenceFrom: 90,
					},
				},
			},
		},
		{
			name: "set aggregateTypes",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetAggregateTypes("user", "org"))},
			},
			res: &SearchQueryFactory{
				queries: []*query{
					{
						aggregateTypes: []AggregateType{"user", "org"},
					},
				},
			},
		},
		{
			name: "set aggregateIDs",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetAggregateIDs("1235", "09824"))},
			},
			res: &SearchQueryFactory{
				queries: []*query{
					{
						aggregateIDs: []string{"1235", "09824"},
					},
				},
			},
		},
		{
			name: "set eventTypes",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetEventTypes("user.created", "user.updated"))},
			},
			res: &SearchQueryFactory{
				queries: []*query{
					{
						eventTypes: []EventType{"user.created", "user.updated"},
					},
				},
			},
		},
		{
			name: "set resource owner",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetResourceOwner("hodor"))},
			},
			res: &SearchQueryFactory{
				queries: []*query{
					{
						resourceOwner: "hodor",
					},
				},
			},
		},
		{
			name: "default search query",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{testAddQuery(testSetAggregateTypes("user"), testSetAggregateIDs("1235", "024")), testSetSortOrder(false)},
			},
			res: &SearchQueryFactory{
				desc: true,
				queries: []*query{
					{
						aggregateTypes: []AggregateType{"user"},
						aggregateIDs:   []string{"1235", "024"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSearchQueryFactory()
			for _, setter := range tt.args.setters {
				factory = setter(factory)
			}
			assertFactory(t, tt.res, factory)
		})
	}
}

func TestSearchQueryFactoryBuild(t *testing.T) {
	type args struct {
		setters []func(*SearchQueryFactory) *SearchQueryFactory
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
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
				query: nil,
			},
		},
		{
			name: "invalid column (too low)",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetColumns(Columns(-1)),
					testAddQuery(testSetAggregateTypes("user")),
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid column (too high)",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetColumns(columnsCount),
					testAddQuery(testSetAggregateTypes("user")),
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "filter aggregate type",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(testSetAggregateTypes("user")),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate types",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(testSetAggregateTypes("user", "org")),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, []AggregateType{"user", "org"}, Operation_In),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testAddQuery(
						testSetAggregateTypes("user"),
						testSetSequence(100),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    true,
					Limit:   5,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_LatestSequence, uint64(100), Operation_Less),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, asc",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(true),
					testAddQuery(
						testSetSequence(100),
						testSetAggregateTypes("user"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   5,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_LatestSequence, uint64(100), Operation_Greater),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type, limit, desc, max event sequence cols",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testSetLimit(5),
					testSetSortOrder(false),
					testSetColumns(Columns_Max_Sequence),
					testAddQuery(
						testSetSequence(100),
						testSetAggregateTypes("user"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: Columns_Max_Sequence,
					Desc:    true,
					Limit:   5,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_LatestSequence, uint64(100), Operation_Less),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate id",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetAggregateIDs("1234"),
						testSetAggregateTypes("user"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_AggregateID, "1234", Operation_Equals),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type and aggregate ids",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetAggregateIDs("1234", "0815"),
						testSetAggregateTypes("user"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_AggregateID, []string{"1234", "0815"}, Operation_In),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type and sequence greater",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetSequence(8),
						testSetAggregateTypes("user"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_LatestSequence, uint64(8), Operation_Greater),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type and event type",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetAggregateTypes("user"),
						testSetEventTypes("user.created"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_EventType, EventType("user.created"), Operation_Equals),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type and event types",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetAggregateTypes("user"),
						testSetEventTypes("user.created", "user.changed"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_EventType, []EventType{"user.created", "user.changed"}, Operation_In),
						},
					},
				},
			},
		},
		{
			name: "filter aggregate type resource owner",
			args: args{
				setters: []func(*SearchQueryFactory) *SearchQueryFactory{
					testAddQuery(
						testSetAggregateTypes("user"),
						testSetResourceOwner("hodor"),
					),
				},
			},
			res: res{
				isErr: nil,
				query: &searchQuery{
					Columns: 0,
					Desc:    false,
					Limit:   0,
					Filters: [][]*Filter{
						{
							NewFilter(Field_AggregateType, AggregateType("user"), Operation_Equals),
							NewFilter(Field_ResourceOwner, "hodor", Operation_Equals),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSearchQueryFactory()
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
