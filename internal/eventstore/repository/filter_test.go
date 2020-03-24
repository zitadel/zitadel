package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/caos/eventstore-lib/pkg/models"
	"github.com/caos/utils/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/jinzhu/gorm"
)

type mockSearchQuery struct {
	t       *testing.T
	limit   uint64
	desc    bool
	filters []*mockFilter
}

func (f *mockSearchQuery) Limit() uint64 {
	return f.limit
}

func (f *mockSearchQuery) OrderDesc() bool {
	return f.desc
}

func (f *mockSearchQuery) Filters() []models.Filter {
	filters := make([]models.Filter, len(f.filters))
	for i, filter := range f.filters {
		filters[i] = filter
	}
	return filters
}

type mockFilter struct {
	key       models.Field
	operation models.Operation
	value     interface{}
}

func (f *mockFilter) GetField() models.Field {
	return f.key
}
func (f *mockFilter) GetOperation() models.Operation {
	return f.operation
}
func (f *mockFilter) GetValue() interface{} {
	return f.value
}

func TestSQL_Filter(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		events      models.Events
		searchQuery models.SearchQuery
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		eventsLen int
		wantErr   bool
		isErrFunc isExpectedError
	}{
		{
			name: "only limit filter",
			fields: fields{
				client: mockDB(t).expectFilterEventsLimit(34, 3),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: &mockSearchQuery{limit: 34},
			},
			eventsLen: 3,
			wantErr:   false,
		},
		{
			name: "only desc filter",
			fields: fields{
				client: mockDB(t).expectFilterEventsDesc(34),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: &mockSearchQuery{desc: true},
			},
			eventsLen: 34,
			wantErr:   false,
		},
		{
			name: "no events found",
			fields: fields{
				client: mockDB(t).expectFilterEventsError(gorm.ErrRecordNotFound),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: &mockSearchQuery{},
			},
			wantErr:   true,
			isErrFunc: errors.IsInternal,
		},
		{
			name: "filter fails because sql internal error",
			fields: fields{
				client: mockDB(t).expectFilterEventsError(sql.ErrConnDone),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: &mockSearchQuery{},
			},
			wantErr:   true,
			isErrFunc: errors.IsInternal,
		},
		{
			name: "filter by aggregate id",
			fields: fields{
				client: mockDB(t).expectFilterEventsAggregateIDLimit("hop", 5),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: &mockSearchQuery{limit: 5, filters: []*mockFilter{&mockFilter{key: es_models.AggregateID, operation: es_models.Equals, value: "hop"}}},
			},
			wantErr:   false,
			isErrFunc: nil,
		},
		{
			name: "filter by aggregate id and aggregate type",
			fields: fields{
				client: mockDB(t).expectFilterEventsAggregateIDTypeLimit("hop", "user", 5),
			},
			args: args{
				events: &mockEvents{t: t},
				searchQuery: &mockSearchQuery{limit: 5, filters: []*mockFilter{
					&mockFilter{key: es_models.AggregateID, operation: es_models.Equals, value: "hop"},
					&mockFilter{key: es_models.AggregateType, operation: es_models.In, value: []string{"user"}},
				}},
			},
			wantErr:   false,
			isErrFunc: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client:    tt.fields.client.db,
				sqlClient: tt.fields.client.sqlClient,
			}
			err := sql.Filter(context.Background(), tt.args.events, tt.args.searchQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.UnlockAggregates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.eventsLen != 0 && tt.args.events.Len() != tt.eventsLen {
				t.Errorf("events has wrong length got: %d want %d", tt.args.events.Len(), tt.eventsLen)
			}
			if tt.wantErr && !tt.isErrFunc(err) {
				t.Errorf("got wrong error %v", err)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			tt.fields.client.close()
		})
	}
}

func Test_getCondition(t *testing.T) {
	type args struct {
		filter models.Filter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single value",
			args: args{
				filter: &mockFilter{
					key:       es_models.LatestSequence,
					operation: es_models.Greater,
					value:     34,
				},
			},
			want: "event_sequence > ?",
		},
		{
			name: "list value",
			args: args{
				filter: &mockFilter{
					key:       es_models.AggregateType,
					operation: es_models.In,
					value:     []string{"a", "b"},
				},
			},
			want: "aggregate_type IN (?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCondition(tt.args.filter); got != tt.want {
				t.Errorf("getCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
