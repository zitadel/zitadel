package projection

import (
	"database/sql"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/org"
)

type testExecuter struct {
	expectedStmt string
	gottenStmt   string
	shouldExec   bool

	expectedArgs []interface{}
	gottenArgs   []interface{}
	gotExecuted  bool
}

type anyArg struct{}

func (e *testExecuter) Exec(stmt string, args ...interface{}) (sql.Result, error) {
	e.gottenStmt = stmt
	e.gottenArgs = args
	e.gotExecuted = true
	return nil, nil
}

func (e *testExecuter) Validate(t *testing.T) {
	t.Helper()
	if e.shouldExec != e.gotExecuted {
		t.Error("expected to be executed")
		return
	}
	if len(e.gottenArgs) != len(e.expectedArgs) {
		t.Errorf("wrong arg len expected: %d got: %d", len(e.expectedArgs), len(e.gottenArgs))
	} else {
		for i := 0; i < len(e.expectedArgs); i++ {
			if _, ok := e.expectedArgs[i].(anyArg); ok {
				continue
			}
			if e.expectedArgs[i] != e.gottenArgs[i] {
				t.Errorf("wrong argument at index %d: got: %v want: %v", i, e.gottenArgs[i], e.expectedArgs[i])
			}
		}
	}
	if e.gottenStmt != e.expectedStmt {
		t.Errorf("wrong stmt want:\n%s\ngot:\n%s", e.expectedStmt, e.gottenStmt)
	}

}

func TestOrgProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reducePrimaryDomainSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgDomainPrimarySetEventType),
					org.AggregateType,
					[]byte(`{"domain": "domain.new"}`),
				), org.DomainPrimarySetEventMapper),
			},
			reduce: (&OrgProjection{}).reducePrimaryDomainSet,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec:   true,
					expectedStmt: "UPDATE projections.orgs SET (change_date, sequence, domain) = ($1, $2, $3) WHERE (id = $4)",
					expectedArgs: []interface{}{
						anyArg{},
						uint64(15),
						"domain.new",
						"agg-id",
					},
				},
			},
		},
		{
			name: "reduceOrgReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgReactivatedEventType),
					org.AggregateType,
					nil,
				), org.OrgReactivatedEventMapper),
			},
			reduce: (&OrgProjection{}).reduceOrgReactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec:   true,
					expectedStmt: "UPDATE projections.orgs SET (change_date, sequence, org_state) = ($1, $2, $3) WHERE (id = $4)",
					expectedArgs: []interface{}{
						anyArg{},
						uint64(15),
						domain.OrgStateActive,
						"agg-id",
					},
				},
			},
		},
		{
			name: "reduceOrgDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgDeactivatedEventType),
					org.AggregateType,
					nil,
				), org.OrgDeactivatedEventMapper),
			},
			reduce: (&OrgProjection{}).reduceOrgDeactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec:   true,
					expectedStmt: "UPDATE projections.orgs SET (change_date, sequence, org_state) = ($1, $2, $3) WHERE (id = $4)",
					expectedArgs: []interface{}{
						anyArg{},
						uint64(15),
						domain.OrgStateInactive,
						"agg-id",
					},
				},
			},
		},
		{
			name: "reduceOrgChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgChangedEventType),
					org.AggregateType,
					[]byte(`{"name": "new name"}`),
				), org.OrgChangedEventMapper),
			},
			reduce: (&OrgProjection{}).reduceOrgChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec:   true,
					expectedStmt: "UPDATE projections.orgs SET (change_date, sequence, name) = ($1, $2, $3) WHERE (id = $4)",
					expectedArgs: []interface{}{
						anyArg{},
						uint64(15),
						"new name",
						"agg-id",
					},
				},
			},
		},
		{
			name: "reduceOrgChanged no changes",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgChangedEventType),
					org.AggregateType,
					[]byte(`{}`),
				), org.OrgChangedEventMapper),
			},
			reduce: (&OrgProjection{}).reduceOrgChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec: false,
				},
			},
		},
		{
			name: "reduceOrgAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgAddedEventType),
					org.AggregateType,
					[]byte(`{"name": "name"}`),
				), org.OrgAddedEventMapper),
			},
			reduce: (&OrgProjection{}).reduceOrgAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					shouldExec:   true,
					expectedStmt: "INSERT INTO projections.orgs (id, creation_date, change_date, resource_owner, sequence, name, org_state) VALUES ($1, $2, $3, $4, $5, $6, $7)",
					expectedArgs: []interface{}{
						"agg-id",
						anyArg{},
						anyArg{},
						"ro-id",
						uint64(15),
						"name",
						domain.OrgStateActive,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}

func testEvent(
	eventType repository.EventType,
	aggregateType repository.AggregateType,
	data []byte,
) *repository.Event {
	return &repository.Event{
		Sequence:                      15,
		PreviousAggregateSequence:     10,
		PreviousAggregateTypeSequence: 10,
		CreationDate:                  time.Now(),
		Type:                          eventType,
		AggregateType:                 aggregateType,
		Data:                          data,
		Version:                       "v1",
		AggregateID:                   "agg-id",
		ResourceOwner:                 "ro-id",
		ID:                            "event-id",
		EditorService:                 "editor-svc",
		EditorUser:                    "editor-user",
	}
}

func baseEvent(*testing.T) eventstore.EventReader {
	return &eventstore.BaseEvent{}
}

func getEvent(event *repository.Event, mapper func(*repository.Event) (eventstore.EventReader, error)) func(t *testing.T) eventstore.EventReader {
	return func(t *testing.T) eventstore.EventReader {
		e, err := mapper(event)
		if err != nil {
			t.Fatalf("mapper failed: %v", err)
		}
		return e
	}
}

type wantReduce struct {
	aggregateType    eventstore.AggregateType
	sequence         uint64
	previousSequence uint64
	executer         *testExecuter
	err              func(error) bool
}

func assertReduce(t *testing.T, stmt *handler.Statement, err error, want wantReduce) {
	t.Helper()
	if want.err == nil && err != nil {
		t.Errorf("unexpected error of type %T: %v", err, err)
		return
	}
	if want.err != nil && want.err(err) {
		return
	}
	if stmt.AggregateType != want.aggregateType {
		t.Errorf("wront aggregate type: want: %q got: %q", want.aggregateType, stmt.AggregateType)
	}

	if stmt.PreviousSequence != want.previousSequence {
		t.Errorf("wront previous sequence: want: %d got: %d", want.previousSequence, stmt.PreviousSequence)
	}

	if stmt.Sequence != want.sequence {
		t.Errorf("wront sequence: want: %d got: %d", want.sequence, stmt.Sequence)
	}
	if stmt.Execute == nil {
		want.executer.Validate(t)
		return
	}
	err = stmt.Execute(want.executer, orgProjection)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want.executer.Validate(t)
}
