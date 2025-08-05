package eventstore

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/v2/database"
)

func TestPaginationOpt(t *testing.T) {
	type args struct {
		opts []paginationOpt
	}
	tests := []struct {
		name string
		args args
		want *Pagination
	}{
		{
			name: "desc",
			args: args{
				opts: []paginationOpt{
					Descending(),
				},
			},
			want: &Pagination{
				desc: true,
			},
		},
		{
			name: "limit",
			args: args{
				opts: []paginationOpt{
					Limit(10),
				},
			},
			want: &Pagination{
				pagination: &database.Pagination{
					Limit: 10,
				},
			},
		},
		{
			name: "offset",
			args: args{
				opts: []paginationOpt{
					Offset(10),
				},
			},
			want: &Pagination{
				pagination: &database.Pagination{
					Offset: 10,
				},
			},
		},
		{
			name: "limit and offset",
			args: args{
				opts: []paginationOpt{
					Limit(10),
					Offset(20),
				},
			},
			want: &Pagination{
				pagination: &database.Pagination{
					Limit:  10,
					Offset: 20,
				},
			},
		},
		{
			name: "global position greater",
			args: args{
				opts: []paginationOpt{
					GlobalPositionGreater(&GlobalPosition{Position: decimal.NewFromInt(10)}),
				},
			},
			want: &Pagination{
				position: &PositionCondition{
					min: &GlobalPosition{
						Position:        decimal.NewFromInt(10),
						InPositionOrder: 0,
					},
				},
			},
		},
		{
			name: "position greater",
			args: args{
				opts: []paginationOpt{
					PositionGreater(decimal.NewFromInt(10), 0),
				},
			},
			want: &Pagination{
				position: &PositionCondition{
					min: &GlobalPosition{
						Position:        decimal.NewFromInt(10),
						InPositionOrder: 0,
					},
				},
				desc: false,
			},
		},
		{
			name: "position less",
			args: args{
				opts: []paginationOpt{
					PositionLess(decimal.NewFromInt(10), 12),
				},
			},
			want: &Pagination{
				position: &PositionCondition{
					max: &GlobalPosition{
						Position:        decimal.NewFromInt(10),
						InPositionOrder: 12,
					},
				},
			},
		},
		{
			name: "global position less",
			args: args{
				opts: []paginationOpt{
					GlobalPositionLess(&GlobalPosition{Position: decimal.NewFromInt(12), InPositionOrder: 24}),
				},
			},
			want: &Pagination{
				position: &PositionCondition{
					max: &GlobalPosition{
						Position:        decimal.NewFromInt(12),
						InPositionOrder: 24,
					},
				},
			},
		},
		{
			name: "position between",
			args: args{
				opts: []paginationOpt{
					PositionBetween(
						&GlobalPosition{decimal.NewFromInt(10), 12},
						&GlobalPosition{decimal.NewFromInt(20), 0},
					),
				},
			},
			want: &Pagination{
				position: &PositionCondition{
					min: &GlobalPosition{
						Position:        decimal.NewFromInt(10),
						InPositionOrder: 12,
					},
					max: &GlobalPosition{
						Position:        decimal.NewFromInt(20),
						InPositionOrder: 0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Pagination)
			for _, opt := range tt.args.opts {
				opt(got)
			}

			if tt.want.Desc() != got.Desc() {
				t.Errorf("unexpected desc %v, want: %v", got.desc, tt.want.desc)
			}
			if !reflect.DeepEqual(tt.want.Pagination(), got.Pagination()) {
				t.Errorf("unexpected pagination %v, want: %v", got.pagination, tt.want.pagination)
			}
			if !reflect.DeepEqual(tt.want.Position(), got.Position()) {
				t.Errorf("unexpected position %v, want: %v", got.position, tt.want.position)
			}
			if !reflect.DeepEqual(tt.want.Position().Max(), got.Position().Max()) {
				t.Errorf("unexpected position.max %v, want: %v", got.Position().max, tt.want.Position().max)
			}
			if !reflect.DeepEqual(tt.want.Position().Min(), got.Position().Min()) {
				t.Errorf("unexpected position.min %v, want: %v", got.Position().min, tt.want.Position().min)
			}
		})
	}
}

func TestEventFilterOpt(t *testing.T) {
	type args struct {
		opts []EventFilterOpt
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *EventFilter
	}{
		{
			name: "EventType",
			args: args{
				opts: []EventFilterOpt{
					SetEventType("test"),
					SetEventType("test2"),
				},
			},
			want: &EventFilter{
				types: []string{"test2"},
			},
		},
		{
			name: "EventTypes",
			args: args{
				opts: []EventFilterOpt{
					SetEventTypes("a", "s"),
					SetEventTypes("d", "f"),
				},
			},
			want: &EventFilter{
				types: []string{"d", "f"},
			},
		},
		{
			name: "AppendEventTypes",
			args: args{
				opts: []EventFilterOpt{
					AppendEventTypes("a", "s"),
					AppendEventTypes("d", "f"),
				},
			},
			want: &EventFilter{
				types: []string{"a", "s", "d", "f"},
			},
		},
		{
			name: "EventRevisionEquals",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionEquals(12),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberEquals[uint16](12),
					value:     toPtr(uint16(12)),
				},
			},
		},
		{
			name: "EventRevisionAtLeast",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionAtLeast(12),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberAtLeast[uint16](12),
					value:     toPtr(uint16(12)),
				},
			},
		},
		{
			name: "EventRevisionGreater",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionGreater(12),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberGreater[uint16](12),
					value:     toPtr(uint16(12)),
				},
			},
		},
		{
			name: "EventRevisionAtMost",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionAtMost(12),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberAtMost[uint16](12),
					value:     toPtr(uint16(12)),
				},
			},
		},
		{
			name: "EventRevisionLess",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionLess(12),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberLess[uint16](12),
					value:     toPtr(uint16(12)),
				},
			},
		},
		{
			name: "EventRevisionBetween",
			args: args{
				opts: []EventFilterOpt{
					EventRevisionBetween(12, 20),
				},
			},
			want: &EventFilter{
				revision: &filter[uint16]{
					condition: database.NewNumberBetween[uint16](12, 20),
					min:       toPtr(uint16(12)),
					max:       toPtr(uint16(20)),
				},
			},
		},
		{
			name: "EventCreatedAtEquals",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtEquals(now),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberEquals(now),
					value:     toPtr(now),
				},
			},
		},
		{
			name: "EventCreatedAtAtLeast",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtAtLeast(now),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberAtLeast(now),
					value:     toPtr(now),
				},
			},
		},
		{
			name: "EventCreatedAtGreater",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtGreater(now),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberGreater(now),
					value:     toPtr(now),
				},
			},
		},
		{
			name: "EventCreatedAtAtMost",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtAtMost(now),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberAtMost(now),
					value:     toPtr(now),
				},
			},
		},
		{
			name: "EventCreatedAtLess",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtLess(now),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberLess(now),
					value:     toPtr(now),
				},
			},
		},
		{
			name: "EventCreatedAtBetween",
			args: args{
				opts: []EventFilterOpt{
					EventCreatedAtBetween(now, now.Add(1*time.Second)),
				},
			},
			want: &EventFilter{
				createdAt: &filter[time.Time]{
					condition: database.NewNumberBetween(now, now.Add(1*time.Second)),
					min:       toPtr(now),
					max:       toPtr(now.Add(1 * time.Second)),
				},
			},
		},
		{
			name: "EventSequenceEquals",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceEquals(12),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberEquals[uint32](12),
					value:     toPtr(uint32(12)),
				},
			},
		},
		{
			name: "EventSequenceAtLeast",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceAtLeast(12),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberAtLeast[uint32](12),
					value:     toPtr(uint32(12)),
				},
			},
		},
		{
			name: "EventSequenceGreater",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceGreater(12),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberGreater[uint32](12),
					value:     toPtr(uint32(12)),
				},
			},
		},
		{
			name: "EventSequenceAtMost",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceAtMost(12),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberAtMost[uint32](12),
					value:     toPtr(uint32(12)),
				},
			},
		},
		{
			name: "EventSequenceLess",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceLess(12),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberLess[uint32](12),
					value:     toPtr(uint32(12)),
				},
			},
		},
		{
			name: "EventSequenceBetween",
			args: args{
				opts: []EventFilterOpt{
					EventSequenceBetween(12, 24),
				},
			},
			want: &EventFilter{
				sequence: &filter[uint32]{
					condition: database.NewNumberBetween[uint32](12, 24),
					min:       toPtr(uint32(12)),
					max:       toPtr(uint32(24)),
				},
			},
		},
		{
			name: "EventCreatorsEqual",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsEqual("cr", "ea", "tor"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewListEquals("cr", "ea", "tor"),
					value:     toPtr([]string{"cr", "ea", "tor"}),
				},
			},
		},
		{
			name: "EventCreatorsEqual no params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsEqual(),
				},
			},
			want: &EventFilter{},
		},
		{
			name: "EventCreatorsEqual one params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsEqual("asdf"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewTextEqual("asdf"),
					value:     toPtr([]string{"asdf"}),
				},
			},
		},
		{
			name: "EventCreatorsContains",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsContains("cr", "ea", "tor"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewListContains("cr", "ea", "tor"),
					value:     toPtr([]string{"cr", "ea", "tor"}),
				},
			},
		},
		{
			name: "EventCreatorsContains no params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsContains(),
				},
			},
			want: &EventFilter{},
		},
		{
			name: "EventCreatorsContains one params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsContains("asdf"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewTextEqual("asdf"),
					value:     toPtr([]string{"asdf"}),
				},
			},
		},
		{
			name: "EventCreatorsNotContains",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsNotContains("cr", "ea", "tor"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewListNotContains("cr", "ea", "tor"),
					value:     toPtr([]string{"cr", "ea", "tor"}),
				},
			},
		},
		{
			name: "EventCreatorsNotContains no params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsNotContains(),
				},
			},
			want: &EventFilter{},
		},
		{
			name: "EventCreatorsNotContains one params",
			args: args{
				opts: []EventFilterOpt{
					EventCreatorsNotContains("asdf"),
				},
			},
			want: &EventFilter{
				creators: &filter[[]string]{
					condition: database.NewTextUnequal("asdf"),
					value:     toPtr([]string{"asdf"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEventFilter(tt.args.opts...)

			if !reflect.DeepEqual(tt.want.Types(), got.Types()) {
				t.Errorf("unexpected types %v, want: %v", got.types, tt.want.types)
			}
			if !reflect.DeepEqual(tt.want.Revision(), got.Revision()) {
				t.Errorf("unexpected revision %v, want: %v", got.revision, tt.want.revision)
			}
			if !reflect.DeepEqual(tt.want.CreatedAt(), got.CreatedAt()) {
				t.Errorf("unexpected createdAt %v, want: %v", got.createdAt, tt.want.createdAt)
			}
			if !reflect.DeepEqual(tt.want.Sequence(), got.Sequence()) {
				t.Errorf("unexpected sequence %v, want: %v", got.sequence, tt.want.sequence)
			}
			if !reflect.DeepEqual(tt.want.Creators(), got.Creators()) {
				t.Errorf("unexpected creators %v, want: %v", got.creators, tt.want.creators)
			}
		})
	}
}

func TestAggregateFilter(t *testing.T) {
	type args struct {
		opts []AggregateFilterOpt
	}
	tests := []struct {
		name string
		args args
		want *AggregateFilter
	}{
		{
			name: "AggregateID",
			args: args{
				opts: []AggregateFilterOpt{
					SetAggregateID("asdf"),
				},
			},
			want: &AggregateFilter{
				ids: []string{"asdf"},
			},
		},
		{
			name: "AggregateIDs",
			args: args{
				opts: []AggregateFilterOpt{
					AggregateIDs("a", "s"),
					AggregateIDs("d", "f"),
				},
			},
			want: &AggregateFilter{
				ids: []string{"d", "f"},
			},
		},
		{
			name: "AggregateIDs",
			args: args{
				opts: []AggregateFilterOpt{
					AppendAggregateIDs("a", "s"),
					AppendAggregateIDs("d", "f"),
				},
			},
			want: &AggregateFilter{
				ids: []string{"a", "s", "d", "f"},
			},
		},
		{
			name: "AppendEvent",
			args: args{
				opts: []AggregateFilterOpt{
					AppendEvent(AppendEventTypes("asdf")),
					AppendEvent(AppendEventTypes("asdf")),
				},
			},
			want: &AggregateFilter{
				events: make([]*EventFilter, 2),
			},
		},
		{
			name: "AppendEvents",
			args: args{
				opts: []AggregateFilterOpt{
					AppendEvents(NewEventFilter()),
					AppendEvents(NewEventFilter()),
				},
			},
			want: &AggregateFilter{
				events: make([]*EventFilter, 2),
			},
		},
		{
			name: "Events",
			args: args{
				opts: []AggregateFilterOpt{
					SetEvents(NewEventFilter()),
					SetEvents(NewEventFilter()),
				},
			},
			want: &AggregateFilter{
				events: make([]*EventFilter, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAggregateFilter("", tt.args.opts...)

			if tt.want.typ != got.typ {
				t.Errorf("unexpected typ %v, want: %v", got.typ, tt.want.typ)
			}
			if !reflect.DeepEqual(tt.want.Type(), got.Type()) {
				t.Errorf("unexpected typ %v, want: %v", got.typ, tt.want.typ)
			}
			if !reflect.DeepEqual(tt.want.IDs(), got.IDs()) {
				t.Errorf("unexpected ids %v, want: %v", got.ids, tt.want.ids)
			}
			if len(tt.want.Events()) != len(got.Events()) {
				t.Errorf("unexpected length of events %v, want: %v", len(got.events), len(tt.want.events))
			}
		})
	}
}

func TestFilterOpt(t *testing.T) {
	type args struct {
		opts []FilterOpt
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{
			name: "limit 1",
			args: args{
				opts: []FilterOpt{
					FilterPagination(Limit(10)),
					FilterPagination(Limit(1)),
				},
			},
			want: &Filter{
				pagination: &Pagination{
					pagination: &database.Pagination{
						Limit: 1,
					},
				},
			},
		},
		{
			name: "AppendAggregateFilter",
			args: args{
				opts: []FilterOpt{
					AppendAggregateFilter("typ"),
					AppendAggregateFilter("typ2"),
				},
			},
			want: &Filter{
				aggregateFilters: make([]*AggregateFilter, 2),
			},
		},
		{
			name: "AppendAggregateFilters",
			args: args{
				opts: []FilterOpt{
					AppendAggregateFilters(NewAggregateFilter("typ")),
					AppendAggregateFilters(NewAggregateFilter("typ2")),
				},
			},
			want: &Filter{
				aggregateFilters: make([]*AggregateFilter, 2),
			},
		},
		{
			name: "AggregateFilters",
			args: args{
				opts: []FilterOpt{
					SetAggregateFilters(NewAggregateFilter("typ")),
					SetAggregateFilters(NewAggregateFilter("typ2")),
				},
			},
			want: &Filter{
				aggregateFilters: make([]*AggregateFilter, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFilter(tt.args.opts...)
			parent := NewQuery("instance", nil)
			got.parent = parent
			tt.want.parent = parent

			if !reflect.DeepEqual(tt.want.Pagination(), got.Pagination()) {
				t.Errorf("unexpected pagination %v, want: %v", got.pagination, tt.want.pagination)
			}
			if len(tt.want.AggregateFilters()) != len(got.AggregateFilters()) {
				t.Errorf("unexpected length of aggregateFilters %v, want: %v", len(got.aggregateFilters), len(tt.want.aggregateFilters))
			}
		})
	}
}

func TestQueryOpt(t *testing.T) {
	type args struct {
		opts []QueryOpt
	}
	var tx sql.Tx
	tests := []struct {
		name string
		args args
		want *Query
	}{
		{
			name: "limit 1",
			args: args{
				opts: []QueryOpt{
					QueryPagination(Limit(10)),
					QueryPagination(Limit(1)),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
				pagination: &Pagination{
					pagination: &database.Pagination{
						Limit: 1,
					},
				},
			},
		},
		{
			name: "with tx",
			args: args{
				opts: []QueryOpt{
					SetQueryTx(&tx),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
				tx: &tx,
			},
		},
		{
			name: "instance",
			args: args{
				opts: []QueryOpt{
					SetInstance("instance2"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance2"),
					value:     toPtr([]string{"instance2"}),
				},
			},
		},
		{
			name: "InstanceEqual no param",
			args: args{
				opts: []QueryOpt{
					InstancesEqual(),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
			},
		},
		{
			name: "InstanceEqual 1 param",
			args: args{
				opts: []QueryOpt{
					InstancesEqual("instance2"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance2"),
					value:     toPtr([]string{"instance2"}),
				},
			},
		},
		{
			name: "InstanceEqual 2 params",
			args: args{
				opts: []QueryOpt{
					InstancesEqual("instance2"),
					InstancesEqual("inst", "ancestor"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewListEquals("inst", "ancestor"),
					value:     toPtr([]string{"inst", "ancestor"}),
				},
			},
		},
		{
			name: "InstancesContains no param",
			args: args{
				opts: []QueryOpt{
					InstancesContains(),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
			},
		},
		{
			name: "InstancesContains 1 param",
			args: args{
				opts: []QueryOpt{
					InstancesContains("instance2"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance2"),
					value:     toPtr([]string{"instance2"}),
				},
			},
		},
		{
			name: "InstancesContains 2 params",
			args: args{
				opts: []QueryOpt{
					InstancesContains("instance2"),
					InstancesContains("inst", "ancestor"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewListContains("inst", "ancestor"),
					value:     toPtr([]string{"inst", "ancestor"}),
				},
			},
		},
		{
			name: "InstancesNotContains no param",
			args: args{
				opts: []QueryOpt{
					InstancesNotContains(),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
			},
		},
		{
			name: "InstancesNotContains 1 param",
			args: args{
				opts: []QueryOpt{
					InstancesNotContains("instance2"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextUnequal("instance2"),
					value:     toPtr([]string{"instance2"}),
				},
			},
		},
		{
			name: "InstancesNotContains 2 params",
			args: args{
				opts: []QueryOpt{
					InstancesNotContains("instance2"),
					InstancesNotContains("inst", "ancestor"),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewListNotContains("inst", "ancestor"),
					value:     toPtr([]string{"inst", "ancestor"}),
				},
			},
		},
		{
			name: "AppendFilters",
			args: args{
				opts: []QueryOpt{
					AppendFilters(NewFilter()),
					AppendFilters(NewFilter()),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
				filters: make([]*Filter, 2),
			},
		},
		{
			name: "AppendFilter",
			args: args{
				opts: []QueryOpt{
					AppendFilter(),
					AppendFilter(),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
				filters: make([]*Filter, 2),
			},
		},
		{
			name: "Filter",
			args: args{
				opts: []QueryOpt{
					SetFilters(NewFilter()),
					SetFilters(NewFilter()),
				},
			},
			want: &Query{
				instances: &filter[[]string]{
					condition: database.NewTextEqual("instance"),
					value:     toPtr([]string{"instance"}),
				},
				filters: make([]*Filter, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQuery("instance", nil, tt.args.opts...)

			if !reflect.DeepEqual(tt.want.Instance(), got.Instance()) {
				t.Errorf("unexpected instances %v, want: %v", got.instances, tt.want.instances)
			}
			if len(tt.want.Filters()) != len(got.Filters()) {
				t.Errorf("unexpected length of filters %v, want: %v", len(got.filters), len(tt.want.filters))
			}
			if !reflect.DeepEqual(tt.want.Tx(), got.Tx()) {
				t.Errorf("unexpected tx %v, want: %v", got.tx, tt.want.tx)
			}
			if !reflect.DeepEqual(tt.want.Pagination(), got.Pagination()) {
				t.Errorf("unexpected pagination %v, want: %v", got.pagination, tt.want.pagination)
			}
		})
	}
}

func toPtr[T any](value T) *T {
	return &value
}
