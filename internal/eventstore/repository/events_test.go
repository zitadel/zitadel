package repository

import (
	"context"
	"database/sql"
	"reflect"
	"runtime"
	"testing"

	lib_models "github.com/caos/eventstore-lib/pkg/models"
	"github.com/caos/utils/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type mockEvents struct {
	events []*models.Event
	t      *testing.T
}

func (events *mockEvents) Append(event lib_models.Event) {
	e, ok := event.(*models.Event)
	if !ok {
		events.t.Error("event is not type *models.Event")
		return
	}
	events.events = append(events.events, e)
}

func (events *mockEvents) Len() int {
	return len(events.events)
}
func (events *mockEvents) Get(index int) lib_models.Event {
	if events.Len() < index {
		return nil
	}
	return events.events[index]
}
func (events *mockEvents) GetAll() []lib_models.Event {
	events.t.Fatal("GetAll is not implemented")
	return nil
}
func (events *mockEvents) Insert(position int, event lib_models.Event) {
	events.t.Fatal("Insert is not implemented")
}

func TestSQL_PushEvents(t *testing.T) {
	type fields struct {
		client *dbMock
	}
	type args struct {
		aggregates []lib_models.Aggregate
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
			args:              args{aggregates: []lib_models.Aggregate{}},
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

			args:              args{aggregates: []lib_models.Aggregate{}},
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
					expectInsertEvent(Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModiferUser:      "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
					},
						"asdfölk-234", 45).
					expectInsertEvent(Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc2",
						ModifierTenant:   "tenant2",
						ModiferUser:      "usr2",
						ResourceOwner:    "ro2",
						PreviousSequence: 45,
						Data:             []byte("{}"),
					}, "asdfölk-233", 46).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []lib_models.Aggregate{
					models.NewAggregate("aggID", "aggType", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
						},
						&models.Event{
							ModifierService: "svc2",
							ModifierTenant:  "tenant2",
							ModifierUser:    "usr2",
							ResourceOwner:   "ro2",
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
					expectInsertEvent(Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModiferUser:      "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
					}, "asdfölk-233", 47).
					expectInsertEvent(Event{
						AggregateID:      "aggID2",
						AggregateType:    "aggType2",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModiferUser:      "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 40,
						Data:             []byte("{}"),
					}, "asdfölk-233", 48).
					expectReleaseSavepoint(nil).
					expectCommit(nil),
			},
			args: args{
				aggregates: []lib_models.Aggregate{
					models.NewAggregate("aggID", "aggType", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
						},
					),
					models.NewAggregate("aggID2", "aggType2", 40,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
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
					expectInsertEventError(Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModiferUser:      "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
					}).
					expectReleaseSavepoint(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []lib_models.Aggregate{
					models.NewAggregate("aggID", "aggType", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
						},
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
							ResourceOwner:   "ro",
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
					expectInsertEvent(Event{
						AggregateID:      "aggID",
						AggregateType:    "aggType",
						ModifierService:  "svc",
						ModifierTenant:   "tenant",
						ModiferUser:      "usr",
						ResourceOwner:    "ro",
						PreviousSequence: 34,
						Data:             []byte("{}"),
					}, "asdfölk-233", 47).
					expectReleaseSavepoint(sql.ErrConnDone).
					expectCommit(nil).
					expectRollback(nil),
			},
			args: args{
				aggregates: []lib_models.Aggregate{
					models.NewAggregate("aggID", "aggType", 34,
						&models.Event{
							ModifierService: "svc",
							ModifierTenant:  "tenant",
							ModifierUser:    "usr",
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
				for _, event := range aggregate.Events().GetAll() {
					e := event.(*models.Event)
					if e.Sequence == 0 {
						t.Error("sequence of returned event is not set")
					}
					if e.AggregateType == "" || e.AggregateID == "" {
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
