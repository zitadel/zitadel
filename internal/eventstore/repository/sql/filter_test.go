package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/caos/utils/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/jinzhu/gorm"
)

func TestSQL_Filter(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		events      *mockEvents
		searchQuery *es_models.SearchQuery
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		eventsLen int
		wantErr   bool
		isErrFunc func(error) bool
	}{
		{
			name: "only limit filter",
			fields: fields{
				client: mockDB(t).expectFilterEventsLimit(34, 3),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQuery(34, false),
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
				searchQuery: es_models.NewSearchQuery(0, true),
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
				searchQuery: &es_models.SearchQuery{},
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
				searchQuery: &es_models.SearchQuery{},
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
				searchQuery: es_models.NewSearchQuery(5, false).AggregateIDFilter("hop"),
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
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQuery(5, false).AggregateIDFilter("hop").AggregateTypeFilter("user"),
			},
			wantErr:   false,
			isErrFunc: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client: tt.fields.client.sqlClient,
			}
			events, err := sql.Filter(context.Background(), tt.args.searchQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQL.UnlockAggregates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.eventsLen != 0 && len(events) != tt.eventsLen {
				t.Errorf("events has wrong length got: %d want %d", len(events), tt.eventsLen)
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
		filter *es_models.Filter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "single value",
			args: args{
				filter: es_models.NewFilter(es_models.LatestSequence, 34, es_models.Greater),
			},
			want: "event_sequence > ?",
		},
		{
			name: "list value",
			args: args{
				filter: es_models.NewFilter(es_models.AggregateType, []string{"a", "b"}, es_models.In),
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
