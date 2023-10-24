package eventstore

import (
	"reflect"
	"testing"
)

func testSetQuery(queryFuncs ...func(*SearchQueryBuilder) *SearchQueryBuilder) func(*SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		for _, queryFunc := range queryFuncs {
			queryFunc(builder)
		}
		return builder
	}
}

func testSetSequenceGreater(sequence uint64) func(*SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.SequenceGreater(sequence)
		return builder
	}
}

func testAddSubQuery(queryFuncs ...func(*SearchQuery) *SearchQuery) func(*SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		query := builder.AddQuery()
		for _, queryFunc := range queryFuncs {
			queryFunc(query)
		}
		return query.Builder()
	}
}

func testSetColumns(columns Columns) func(factory *SearchQueryBuilder) *SearchQueryBuilder {
	return func(factory *SearchQueryBuilder) *SearchQueryBuilder {
		factory = factory.Columns(columns)
		return factory
	}
}

func testSetLimit(limit uint64) func(builder *SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.Limit(limit)
		return builder
	}
}

func testOr(queryFuncs ...func(*SearchQuery) *SearchQuery) func(*SearchQuery) *SearchQuery {
	return func(query *SearchQuery) *SearchQuery {
		subQuery := query.Or()
		for _, queryFunc := range queryFuncs {
			queryFunc(subQuery)
		}
		return subQuery
	}
}

func testSetAggregateTypes(types ...AggregateType) func(*SearchQuery) *SearchQuery {
	return func(query *SearchQuery) *SearchQuery {
		query = query.AggregateTypes(types...)
		return query
	}
}

func testSetAggregateIDs(aggregateIDs ...string) func(*SearchQuery) *SearchQuery {
	return func(query *SearchQuery) *SearchQuery {
		query = query.AggregateIDs(aggregateIDs...)
		return query
	}
}

func testSetEventTypes(eventTypes ...EventType) func(*SearchQuery) *SearchQuery {
	return func(query *SearchQuery) *SearchQuery {
		query = query.EventTypes(eventTypes...)
		return query
	}
}

func testSetResourceOwner(resourceOwner string) func(*SearchQueryBuilder) *SearchQueryBuilder {
	return func(builder *SearchQueryBuilder) *SearchQueryBuilder {
		builder = builder.ResourceOwner(resourceOwner)
		return builder
	}
}

func testSetSortOrder(asc bool) func(*SearchQueryBuilder) *SearchQueryBuilder {
	return func(query *SearchQueryBuilder) *SearchQueryBuilder {
		if asc {
			query = query.OrderAsc()
		} else {
			query = query.OrderDesc()
		}
		return query
	}
}

func TestSearchQuerybuilderSetters(t *testing.T) {
	type args struct {
		columns Columns
		setters []func(*SearchQueryBuilder) *SearchQueryBuilder
	}
	tests := []struct {
		name string
		args args
		res  *SearchQueryBuilder
	}{
		{
			name: "New builder",
			args: args{
				columns: ColumnsEvent,
			},
			res: &SearchQueryBuilder{
				columns: Columns(ColumnsEvent),
			},
		},
		{
			name: "set columns",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testSetColumns(ColumnsMaxSequence)},
			},
			res: &SearchQueryBuilder{
				columns: ColumnsMaxSequence,
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
				setters: []func(b *SearchQueryBuilder) *SearchQueryBuilder{
					testSetQuery(testSetSequenceGreater(90)),
				},
			},
			res: &SearchQueryBuilder{
				eventSequenceGreater: 90,
			},
		},
		{
			name: "set aggregateIDs",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testAddSubQuery(testSetAggregateIDs("1235", "09824"))},
			},
			res: &SearchQueryBuilder{
				queries: []*SearchQuery{
					{
						aggregateIDs: []string{"1235", "09824"},
					},
				},
			},
		},
		{
			name: "set eventTypes",
			args: args{
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testAddSubQuery(testSetEventTypes("user.created", "user.updated"))},
			},
			res: &SearchQueryBuilder{
				queries: []*SearchQuery{
					{
						eventTypes: []EventType{"user.created", "user.updated"},
					},
				},
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
				setters: []func(*SearchQueryBuilder) *SearchQueryBuilder{testAddSubQuery(testSetAggregateTypes("user"), testSetAggregateIDs("1235", "024")), testSetSortOrder(false)},
			},
			res: &SearchQueryBuilder{
				desc: true,
				queries: []*SearchQuery{
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
			builder := NewSearchQueryBuilder(tt.args.columns)
			for _, setter := range tt.args.setters {
				builder = setter(builder)
			}

			assertBuilder(t, tt.res, builder)
		})
	}
}

func assertBuilder(t *testing.T, want, got *SearchQueryBuilder) {
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
	if got.resourceOwner != want.resourceOwner {
		t.Errorf("wrong : got: %v want: %v", got.resourceOwner, want.resourceOwner)
	}
	if len(got.queries) != len(want.queries) {
		t.Errorf("wrong length of queries: got: %v want: %v", len(got.queries), len(want.queries))
	}

	for i, query := range got.queries {
		assertQuery(t, i, want.queries[i], query)
	}
}

func assertQuery(t *testing.T, i int, want, got *SearchQuery) {
	t.Helper()

	if !reflect.DeepEqual(got.aggregateIDs, want.aggregateIDs) {
		t.Errorf("wrong aggregateIDs in query %d : got: %v want: %v", i, got.aggregateIDs, want.aggregateIDs)
	}
	if !reflect.DeepEqual(got.aggregateTypes, want.aggregateTypes) {
		t.Errorf("wrong aggregateTypes in query %d : got: %v want: %v", i, got.aggregateTypes, want.aggregateTypes)
	}
	if !reflect.DeepEqual(got.eventData, want.eventData) {
		t.Errorf("wrong eventData in query %d : got: %v want: %v", i, got.eventData, want.eventData)
	}
	// if got.eventSequenceGreater != want.eventSequenceGreater {
	// 	t.Errorf("wrong eventSequenceGreater in query %d : got: %v want: %v", i, got.eventSequenceGreater, want.eventSequenceGreater)
	// }
	if !reflect.DeepEqual(got.eventTypes, want.eventTypes) {
		t.Errorf("wrong eventTypes in query %d : got: %v want: %v", i, got.eventTypes, want.eventTypes)
	}
}

func TestSearchQuery_matches(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name  string
		query *SearchQuery
		event Event
		want  bool
	}{
		{
			name:  "wrong aggregate type",
			query: NewSearchQueryBuilder(ColumnsEvent).AddQuery().AggregateTypes("searched"),
			event: &BaseEvent{
				Agg: &Aggregate{
					Type: "found",
				},
			},
			want: false,
		},
		{
			name:  "wrong aggregate id",
			query: NewSearchQueryBuilder(ColumnsEvent).AddQuery().AggregateIDs("1", "10", "100"),
			event: &BaseEvent{
				Agg: &Aggregate{
					ID: "2",
				},
			},
			want: false,
		},
		{
			name:  "wrong event type",
			query: NewSearchQueryBuilder(ColumnsEvent).AddQuery().EventTypes("event.searched.type"),
			event: &BaseEvent{
				EventType: "event.actual.type",
				Agg:       &Aggregate{},
			},
			want: false,
		},
		{
			name: "matching",
			query: NewSearchQueryBuilder(ColumnsEvent).
				AddQuery().
				AggregateIDs("2").
				AggregateTypes("actual").
				EventTypes("event.actual.type"),
			event: &BaseEvent{
				Seq: 55,
				Agg: &Aggregate{
					ID:   "2",
					Type: "actual",
				},
				EventType: "event.actual.type",
			},
			want: true,
		},
		{
			name:  "matching empty query",
			query: NewSearchQueryBuilder(ColumnsEvent).AddQuery(),
			event: &BaseEvent{
				Seq: 55,
				Agg: &Aggregate{
					ID:   "2",
					Type: "actual",
				},
				EventType: "event.actual.type",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &SearchQuery{
				aggregateTypes: tt.query.aggregateTypes,
				aggregateIDs:   tt.query.aggregateIDs,
				eventTypes:     tt.query.eventTypes,
				eventData:      tt.query.eventData,
			}
			if got := query.matches(tt.event); got != tt.want {
				t.Errorf("SearchQuery.matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchQueryBuilder_Matches(t *testing.T) {
	type args struct {
		event       Event
		existingLen int
	}
	tests := []struct {
		name    string
		builder *SearchQueryBuilder
		args    args
		want    bool
	}{
		{
			name:    "sequence too high",
			builder: NewSearchQueryBuilder(ColumnsEvent).SequenceGreater(60),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						InstanceID: "instance",
					},
					Seq: 60,
				},
			},
			want: false,
		},
		{
			name:    "limit exeeded",
			builder: NewSearchQueryBuilder(ColumnsEvent).Limit(100),
			args: args{
				event:       &BaseEvent{},
				existingLen: 100,
			},
			want: false,
		},
		{
			name:    "wrong resource owner",
			builder: NewSearchQueryBuilder(ColumnsEvent).ResourceOwner("query"),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						ResourceOwner: "ro",
					},
				},
				existingLen: 0,
			},
			want: false,
		},
		{
			name:    "wrong instance",
			builder: NewSearchQueryBuilder(ColumnsEvent).InstanceID("instance"),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						InstanceID: "different instance",
					},
				},
				existingLen: 0,
			},
			want: false,
		},
		{
			name:    "query failed",
			builder: NewSearchQueryBuilder(ColumnsEvent).SequenceGreater(1000),
			args: args{
				event: &BaseEvent{
					Seq: 999,
					Agg: &Aggregate{},
				},
				existingLen: 0,
			},
			want: false,
		},
		{
			name: "matching",
			builder: NewSearchQueryBuilder(ColumnsEvent).
				Limit(1000).
				ResourceOwner("ro").
				InstanceID("instance").
				SequenceGreater(1000),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						ResourceOwner: "ro",
						InstanceID:    "instance",
					},
					Seq: 1001,
				},
				existingLen: 999,
			},
			want: true,
		},
		{
			name:    "matching builder resourceOwner and Instance",
			builder: NewSearchQueryBuilder(ColumnsEvent),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						ResourceOwner: "ro",
						InstanceID:    "instance",
					},
					Seq: 1001,
				},
				existingLen: 999,
			},
			want: true,
		},
		{
			name:    "matching builder resourceOwner only",
			builder: NewSearchQueryBuilder(ColumnsEvent),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						ResourceOwner: "ro",
					},
					Seq: 1001,
				},
				existingLen: 999,
			},
			want: true,
		},
		{
			name:    "matching builder instanceID only",
			builder: NewSearchQueryBuilder(ColumnsEvent),
			args: args{
				event: &BaseEvent{
					Agg: &Aggregate{
						InstanceID: "instance",
					},
					Seq: 1001,
				},
				existingLen: 999,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.Matches(tt.args.event, tt.args.existingLen); got != tt.want {
				t.Errorf("SearchQueryBuilder.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
