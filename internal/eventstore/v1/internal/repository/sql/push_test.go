package sql

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"runtime"
	"testing"

	z_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type mockEvents struct {
	events []*models.Event
	t      *testing.T
}

func TestSQL_PushAggregates(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		aggregates []*models.Aggregate
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		isError           func(error) bool
		shouldCheckEvents bool
	}{
		{
			name: "no aggregates",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args:              args{aggregates: []*models.Aggregate{}},
			shouldCheckEvents: false,
			isError:           noErr,
		},
		{
			name: "prepare fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(sql.ErrConnDone).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args:              args{aggregates: []*models.Aggregate{}},
			shouldCheckEvents: false,
			isError:           func(err error) bool { return errors.Is(err, sql.ErrConnDone) },
		},
		{
			name: "no aggregates release fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil),
			},

			args:              args{aggregates: []*models.Aggregate{}},
			isError:           z_errors.IsInternal,
			shouldCheckEvents: false,
		},
		{
			name: "aggregate precondtion fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectFilterEventsError(z_errors.CreateCaosError(nil, "SQL-IzJOf", "err")).
					expectRollback(nil),
			},

			args:              args{aggregates: []*models.Aggregate{aggregateWithPrecondition(&models.Aggregate{}, models.NewSearchQuery().SetLimit(1), nil)}},
			isError:           z_errors.IsPreconditionFailed,
			shouldCheckEvents: false,
		},
		{
			name: "one aggregate two events success",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, 45).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc2",
						EditorUser:       "usr2",
						ResourceOwner:    "ro2",
						PreviousSequence: 45,
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, 46).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					{
						PreviousSequence: 34,
						Events: []*models.Event{
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
							},
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc2",
								EditorUser:       "usr2",
								ResourceOwner:    "ro2",
								Type:             "eventTyp",
							},
						},
					},
				},
			},
			shouldCheckEvents: true,
			isError:           noErr,
		},
		{
			name: "two aggregates one event per aggregate success",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, 47).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID2",
						AggregateType:    "aggType2",
						EditorService:    "svc",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 40,
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, 48).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					{
						PreviousSequence: 34,
						Events: []*models.Event{
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
							},
						},
					},
					{
						PreviousSequence: 40,
						Events: []*models.Event{
							{
								AggregateID:      "aggID2",
								AggregateType:    "aggType2",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
							},
						},
					},
				},
			},
			shouldCheckEvents: true,
			isError:           noErr,
		},
		{
			name: "first event fails no action with second event",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert(nil).
					expectInsertEventError(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}).
					expectReleaseSavepoint(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					{
						Events: []*models.Event{
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 0,
							},
						},
					},
				},
			},
			isError:           z_errors.IsInternal,
			shouldCheckEvents: false,
		},
		{
			name: "one event, release savepoint fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectPrepareInsert(nil).
					expectSavepoint().
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Type:             "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					}, 47).
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					{
						Events: []*models.Event{
							{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
						},
					},
				},
			},
			isError:           z_errors.IsInternal,
			shouldCheckEvents: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client: tt.fields.client.sqlClient,
			}
			err := sql.PushAggregates(context.Background(), tt.args.aggregates...)
			if err != nil && !tt.isError(err) {
				t.Errorf("wrong error type = %v, errFunc %s", err, functionName(tt.isError))
			}
			if !tt.shouldCheckEvents {
				return
			}
			for _, aggregate := range tt.args.aggregates {
				for _, event := range aggregate.Events {
					if event.Sequence == 0 {
						t.Error("sequence of returned event is not set")
					}
					if event.AggregateType == "" || event.AggregateID == "" {
						t.Error("aggregate of event is not set")
					}
				}
			}
			if err := tt.fields.client.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all database expectaions met: %s", err)
			}
		})
	}
}

func noErr(err error) bool {
	return err == nil
}

func functionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Test_precondtion(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		aggregate *models.Aggregate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isErr  func(error) bool
	}{
		{
			name: "no precondition",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil),
			},
			args: args{
				aggregate: &models.Aggregate{},
			},
		},
		{
			name: "precondition fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).expectFilterEventsLimit("test", 5, 0),
			},
			args: args{
				aggregate: aggregateWithPrecondition(&models.Aggregate{}, models.NewSearchQuery().SetLimit(5).AggregateTypeFilter("test"), validationFunc(z_errors.ThrowPreconditionFailed(nil, "SQL-LBIKm", "err"))),
			},
			isErr: z_errors.IsPreconditionFailed,
		},
		{
			name: "precondition with filter error",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).expectFilterEventsError(z_errors.ThrowInternal(nil, "SQL-ac9EW", "err")),
			},
			args: args{
				aggregate: aggregateWithPrecondition(&models.Aggregate{}, models.NewSearchQuery().SetLimit(5).AggregateTypeFilter("test"), validationFunc(z_errors.CreateCaosError(nil, "SQL-LBIKm", "err"))),
			},
			isErr: z_errors.IsPreconditionFailed,
		},
		{
			name: "precondition no events",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).expectFilterEventsLimit("test", 5, 0),
			},
			args: args{
				aggregate: aggregateWithPrecondition(&models.Aggregate{}, models.NewSearchQuery().SetLimit(5).AggregateTypeFilter("test"), validationFunc(nil)),
			},
		},
		{
			name: "precondition with events",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).expectFilterEventsLimit("test", 5, 3),
			},
			args: args{
				aggregate: aggregateWithPrecondition(&models.Aggregate{}, models.NewSearchQuery().SetLimit(5).AggregateTypeFilter("test"), validationFunc(nil)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tt.fields.client.sqlClient.Begin()
			if err != nil {
				t.Errorf("unable to start tx %v", err)
				t.FailNow()
			}
			err = precondtion(tx, tt.args.aggregate)
			if (err != nil) && (tt.isErr == nil) {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.isErr != nil && !tt.isErr(err) {
				t.Errorf("precondtion() wrong error %T, %v", err, err)
			}
		})
	}
}

func aggregateWithPrecondition(aggregate *models.Aggregate, query *models.SearchQuery, precondition func(...*models.Event) error) *models.Aggregate {
	aggregate.SetPrecondition(query, precondition)
	return aggregate
}

func validationFunc(err error) func(events ...*models.Event) error {
	return func(events ...*models.Event) error {
		return err
	}
}
