package pseudo

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix    = eventstore.EventType("pseudo.")
	TimestampEventType = eventTypePrefix + "timestamp"
)

var _ eventstore.Event = (*TimestampEvent)(nil)

type TimestampEvent struct {
	Timestamp   time.Time
	InstanceIDs []string
}

func (t TimestampEvent) Aggregate() eventstore.Aggregate {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) EditorService() string {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) EditorUser() string {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) Type() eventstore.EventType {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) Sequence() uint64 {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) CreationDate() time.Time {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) PreviousAggregateSequence() uint64 {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) PreviousAggregateTypeSequence() uint64 {
	panic("TimestampEvent is not a real event")
}

func (t TimestampEvent) DataAsBytes() []byte {
	panic("TimestampEvent is not a real event")
}

func NewTimestampEvent(
	timestamp time.Time,
	instanceIDs ...string,
) *TimestampEvent {
	return &TimestampEvent{
		Timestamp:   timestamp,
		InstanceIDs: instanceIDs,
	}
}
