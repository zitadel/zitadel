package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

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
	projectionName   string
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
		t.Errorf("wrong aggregate type: want: %q got: %q", want.aggregateType, stmt.AggregateType)
	}

	if stmt.PreviousSequence != want.previousSequence {
		t.Errorf("wrong previous sequence: want: %d got: %d", want.previousSequence, stmt.PreviousSequence)
	}

	if stmt.Sequence != want.sequence {
		t.Errorf("wrong sequence: want: %d got: %d", want.sequence, stmt.Sequence)
	}
	if stmt.Execute == nil {
		want.executer.Validate(t)
		return
	}
	err = stmt.Execute(want.executer, want.projectionName)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want.executer.Validate(t)
}
