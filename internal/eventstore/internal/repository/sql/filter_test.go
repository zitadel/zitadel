package sql

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestSQL_Filter(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		events      *mockEvents
		searchQuery *es_models.SearchQueryFactory
	}
	type res struct {
		wantErr   bool
		isErrFunc func(error) bool
		eventsLen int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "only limit filter",
			fields: fields{
				client: mockDB(t).expectFilterEventsLimit("user", 34, 3),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("user").Limit(34),
			},
			res: res{
				eventsLen: 3,
				wantErr:   false,
			},
		},
		{
			name: "only desc filter",
			fields: fields{
				client: mockDB(t).expectFilterEventsDesc("user", 34),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("user").OrderDesc(),
			},
			res: res{
				eventsLen: 34,
				wantErr:   false,
			},
		},
		{
			name: "no events found",
			fields: fields{
				client: mockDB(t).expectFilterEventsError(sql.ErrNoRows),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("nonAggregate"),
			},
			res: res{
				wantErr:   true,
				isErrFunc: errors.IsInternal,
			},
		},
		{
			name: "filter fails because sql internal error",
			fields: fields{
				client: mockDB(t).expectFilterEventsError(sql.ErrConnDone),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("user"),
			},
			res: res{
				wantErr:   true,
				isErrFunc: errors.IsInternal,
			},
		},
		{
			name: "filter by aggregate id",
			fields: fields{
				client: mockDB(t).expectFilterEventsAggregateIDLimit("user", "hop", 5),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("user").Limit(5).AggregateIDs("hop"),
			},
			res: res{
				wantErr:   false,
				isErrFunc: nil,
			},
		},
		{
			name: "filter by aggregate id and aggregate type",
			fields: fields{
				client: mockDB(t).expectFilterEventsAggregateIDTypeLimit("user", "hop", 5),
			},
			args: args{
				events:      &mockEvents{t: t},
				searchQuery: es_models.NewSearchQueryFactory("user").Limit(5).AggregateIDs("hop"),
			},
			res: res{
				wantErr:   false,
				isErrFunc: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client: tt.fields.client.sqlClient,
			}
			events, err := sql.Filter(context.Background(), tt.args.searchQuery.IsPrecondition())
			if (err != nil) != tt.res.wantErr {
				t.Errorf("SQL.Filter() error = %v, wantErr %v", err, tt.res.wantErr)
			}
			if tt.res.eventsLen != 0 && len(events) != tt.res.eventsLen {
				t.Errorf("events has wrong length got: %d want %d", len(events), tt.res.eventsLen)
			}
			if tt.res.wantErr && !tt.res.isErrFunc(err) {
				t.Errorf("got wrong error %v", err)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			tt.fields.client.close()
		})
	}
}

func TestSQL_LatestSequence(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		searchQuery *es_models.SearchQueryFactory
	}
	type res struct {
		wantErr   bool
		isErrFunc func(error) bool
		sequence  uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid query factory",
			args: args{
				searchQuery: nil,
			},
			fields: fields{
				client: mockDB(t),
			},
			res: res{
				wantErr:   true,
				isErrFunc: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no events for aggregate",
			args: args{
				searchQuery: es_models.NewSearchQueryFactory("idiot").Columns(es_models.Columns_Max_Sequence).IsPrecondition(),
			},
			fields: fields{
				client: mockDB(t).expectLatestSequenceFilterError("idiot", sql.ErrNoRows),
			},
			res: res{
				wantErr:  false,
				sequence: 0,
			},
		},
		{
			name: "sql query error",
			args: args{
				searchQuery: es_models.NewSearchQueryFactory("idiot").Columns(es_models.Columns_Max_Sequence).IsPrecondition(),
			},
			fields: fields{
				client: mockDB(t).expectLatestSequenceFilterError("idiot", sql.ErrConnDone),
			},
			res: res{
				wantErr:   true,
				isErrFunc: errors.IsInternal,
				sequence:  0,
			},
		},
		{
			name: "events for aggregate found",
			args: args{
				searchQuery: es_models.NewSearchQueryFactory("user").Columns(es_models.Columns_Max_Sequence).IsPrecondition(),
			},
			fields: fields{
				client: mockDB(t).expectLatestSequenceFilter("user", math.MaxUint64),
			},
			res: res{
				wantErr:  false,
				sequence: math.MaxUint64,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client: tt.fields.client.sqlClient,
			}
			sequence, err := sql.LatestSequence(context.Background(), tt.args.searchQuery)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("SQL.Filter() error = %v, wantErr %v", err, tt.res.wantErr)
			}
			if tt.res.sequence != sequence {
				t.Errorf("events has wrong length got: %d want %d", sequence, tt.res.sequence)
			}
			if tt.res.wantErr && !tt.res.isErrFunc(err) {
				t.Errorf("got wrong error %v", err)
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			tt.fields.client.close()
		})
	}
}
