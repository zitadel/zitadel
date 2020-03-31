package sql

import (
	"context"
	"database/sql"
	"reflect"
	"runtime"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
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
						EditorService:    "svc",
						EditorOrg:        "tenant",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Type:             "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					},
						"asdfölk-234", 45).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc2",
						EditorOrg:        "tenant2",
						EditorUser:       "usr2",
						ResourceOwner:    "ro2",
						PreviousSequence: 45,
						Type:             "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 46).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					&models.Aggregate{
						Events: []*models.Event{
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc2",
								EditorOrg:        "tenant2",
								EditorUser:       "usr2",
								ResourceOwner:    "ro2",
								Type:             "eventTyp",
								PreviousSequence: 0,
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
					expectPrepareInsert().
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorOrg:        "tenant",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 47).
					expectInsertEvent(&models.Event{
						AggregateID:      "aggID2",
						AggregateType:    "aggType2",
						EditorService:    "svc",
						EditorOrg:        "tenant",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 40,
						Data:             []byte("{}"),
						Type:             "eventTyp",
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 48).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					&models.Aggregate{
						Events: []*models.Event{
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
						},
					},
					&models.Aggregate{
						Events: []*models.Event{
							&models.Event{
								AggregateID:      "aggID2",
								AggregateType:    "aggType2",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 40,
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
					expectInsertEventError(&models.Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						EditorService:    "svc",
						EditorOrg:        "tenant",
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
					&models.Aggregate{
						Events: []*models.Event{
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 0,
							},
						},
					},
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
						EditorService:    "svc",
						EditorOrg:        "tenant",
						EditorUser:       "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Type:             "eventTyp",
						Data:             []byte("{}"),
						AggregateVersion: "v0.0.1",
					}, "asdfölk-233", 47).
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []*models.Aggregate{
					&models.Aggregate{
						Events: []*models.Event{
							&models.Event{
								AggregateID:      "aggID",
								AggregateType:    "aggType",
								AggregateVersion: "v0.0.1",
								EditorService:    "svc",
								EditorOrg:        "tenant",
								EditorUser:       "usr",
								ResourceOwner:    "ro",
								Type:             "eventTyp",
								PreviousSequence: 34,
							},
						},
					},
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
