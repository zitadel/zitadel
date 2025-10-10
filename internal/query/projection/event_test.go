package projection

import (
	"database/sql"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func testEvent(
	eventType eventstore.EventType,
	aggregateType eventstore.AggregateType,
	data []byte,
	opts ...eventOption,
) *repository.Event {
	return timedTestEvent(eventType, aggregateType, data, time.Now(), opts...)
}

func toSystemEvent(event *repository.Event) *repository.Event {
	event.EditorUser = "SYSTEM"
	return event
}

func timedTestEvent(
	eventType eventstore.EventType,
	aggregateType eventstore.AggregateType,
	data []byte,
	creationDate time.Time,
	opts ...eventOption,
) *repository.Event {
	e := &repository.Event{
		Seq:           15,
		CreationDate:  creationDate,
		Typ:           eventType,
		AggregateType: aggregateType,
		Data:          data,
		Version:       "v1",
		AggregateID:   "agg-id",
		ResourceOwner: sql.NullString{String: "ro-id", Valid: true},
		InstanceID:    "instance-id",
		ID:            "event-id",
		EditorUser:    "editor-user",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

type eventOption func(e *repository.Event)

func withVersion(v eventstore.Version) eventOption {
	return func(e *repository.Event) {
		e.Version = v
	}
}

func baseEvent(*testing.T) eventstore.Event {
	return &eventstore.BaseEvent{}
}

func getEvent(event *repository.Event, mapper func(eventstore.Event) (eventstore.Event, error)) func(t *testing.T) eventstore.Event {
	return func(t *testing.T) eventstore.Event {
		e, err := mapper(event)
		if err != nil {
			t.Fatalf("mapper failed: %v", err)
		}
		return e
	}
}

type wantReduce struct {
	aggregateType eventstore.AggregateType
	sequence      uint64
	executer      *testExecuter
	err           func(error) bool
}

func assertReduce(t *testing.T, stmt *handler.Statement, err error, projection string, want wantReduce) {
	t.Helper()
	if want.err == nil && err != nil {
		t.Errorf("unexpected error of type %T: %v", err, err)
		return
	}
	if want.err != nil && want.err(err) {
		return
	}
	if stmt.Aggregate.Type != want.aggregateType {
		t.Errorf("wrong aggregate type: want: %q got: %q", want.aggregateType, stmt.Aggregate.Type)
	}

	if stmt.Sequence != want.sequence {
		t.Errorf("wrong sequence: want: %d got: %d", want.sequence, stmt.Sequence)
	}
	if stmt.Execute == nil {
		want.executer.Validate(t)
		return
	}
	err = stmt.Execute(t.Context(), want.executer, projection)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want.executer.Validate(t)
}
