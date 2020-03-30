package sql

import (
	"context"
	"database/sql"
	"reflect"
	"runtime"
	"testing"

	"github.com/caos/utils/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type mockEvents struct {
	events []*models.Event
	t      *testing.T
}

func TestSQL_PushEvents(t *testing.T) {
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
					expectPrepareInsert().
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args:              args{aggregates: []*models.Aggregate{}},
			shouldCheckEvents: false,
			isError:           noErr,
		},
		{
			name: "no aggregates release fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert().
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil),
			},

			args:              args{aggregates: []*models.Aggregate{}},
			isError:           errors.IsInternal,
			shouldCheckEvents: false,
		},
		{
			name: "one aggregate two events success",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectSavepoint().
					expectPrepareInsert().
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModifierUser:     "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Typ:              "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					},
						"asdfölk-234", 45).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc2",
						ModifierTenant:   "tenant2",
						ModifierUser:     "usr2",
						ResourceOwner:    "ro2",
						PreviousSequence: 45,
						Typ:              "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 46).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					models.MustNewAggregate("aggID", "aggType", models.MustVersion(0, 0, 1), 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
							Typ:             "eventTyp",
						},
						&models.Event{
							ModifierService: "svc2",
							ModifierTenant:  "tenant2",
							ModifierUser:    "usr2",
							ResourceOwner:   "ro2",
							Typ:             "eventTyp",
						},
					),
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
					expectPrepareInsert().
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModifierUser:     "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
						Typ:              "eventTyp",
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 47).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID2",
						AggregateType:    "aggType2",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModifierUser:     "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 40,
						Data:             []byte("{}"),
						Typ:              "eventTyp",
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 48).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					models.MustNewAggregate("aggID", "aggType", "v0.0.1", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
							Typ:             "eventTyp",
						},
					),
					models.MustNewAggregate("aggID2", "aggType2", "v0.0.1", 40,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
							Typ:             "eventTyp",
						},
					),
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
					expectInsertEventError(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModifierUser:     "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
						Typ:              "eventTyp",
						AggregateVersion: "v0.0.1",
					}).
					expectReleaseSavepoint(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					models.MustNewAggregate("aggID", "aggType", "v0.0.1", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
							Typ:             "eventTyp",
						},
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
							Typ:             "eventTyp",
						},
					),
				},
			},
			isError:           errors.IsInternal,
			shouldCheckEvents: false,
		},
		{
			name: "one event, release savepoint fails",
			fields: fields{
				client: mockDB(t).
					expectBegin(nil).
					expectPrepareInsert().
					expectSavepoint().
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModifierUser:     "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Typ:              "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 47).
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					models.MustNewAggregate("aggID", "aggType", "v0.0.1", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							Typ:             "eventTyp",
							ResourceOwner:   "ro",
						},
					),
				},
			},
			isError:           errors.IsInternal,
			shouldCheckEvents: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := &SQL{
				client: tt.fields.client.sqlClient,
			}
			err := sql.PushEvents(context.Background(), tt.args.aggregates...)
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
