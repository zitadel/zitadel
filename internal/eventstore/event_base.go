package eventstore

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
)

var (
	_ Event = (*BaseEvent)(nil)
)

// BaseEvent represents the minimum metadata of an event
type BaseEvent struct {
	ID        string
	EventType EventType `json:"-"`

	Agg *Aggregate

	Seq                           uint64
	Pos                           float64
	Creation                      time.Time
	previousAggregateSequence     uint64
	previousAggregateTypeSequence uint64

	//User who created the event
	User string `json:"-"`
	//Service which created the event
	Service string `json:"-"`
	Data    []byte `json:"-"`
}

// Position implements Event.
func (e *BaseEvent) Position() float64 {
	return e.Pos
}

// EditorService implements Command
func (e *BaseEvent) EditorService() string {
	return e.Service
}

// EditorUser implements Command
func (e *BaseEvent) EditorUser() string {
	return e.User
}

// Creator implements action
func (e *BaseEvent) Creator() string {
	return e.EditorUser()
}

// Type implements action
func (e *BaseEvent) Type() EventType {
	return e.EventType
}

// Sequence is an upcounting unique number of the event
func (e *BaseEvent) Sequence() uint64 {
	return e.Seq
}

// CreationDate is the the time, the event is inserted into the eventstore
func (e *BaseEvent) CreationDate() time.Time {
	return e.Creation
}

// CreatedAt implements Event
func (e *BaseEvent) CreatedAt() time.Time {
	return e.CreationDate()
}

// Aggregate implements action
func (e *BaseEvent) Aggregate() *Aggregate {
	return e.Agg
}

// Data returns the payload of the event. It represent the changed fields by the event
func (e *BaseEvent) DataAsBytes() []byte {
	return e.Data
}

// Revision implements action
func (e *BaseEvent) Revision() uint16 {
	revision, err := strconv.ParseUint(strings.TrimPrefix(string(e.Agg.Version), "v"), 10, 16)
	logging.OnError(err).Debug("failed to parse event revision")
	return uint16(revision)
}

// Unmarshal implements Event
func (e *BaseEvent) Unmarshal(ptr any) error {
	if len(e.Data) == 0 {
		return nil
	}
	return json.Unmarshal(e.Data, ptr)
}

const defaultService = "zitadel"

// BaseEventFromRepo maps a stored event to a BaseEvent
func BaseEventFromRepo(event Event) *BaseEvent {
	return &BaseEvent{
		Agg:       event.Aggregate(),
		EventType: event.Type(),
		Creation:  event.CreatedAt(),
		Seq:       event.Sequence(),
		Service:   defaultService,
		User:      event.Creator(),
		Data:      event.DataAsBytes(),
		Pos:       event.Position(),
	}
}

// NewBaseEventForPush is the constructor for event's which will be pushed into the eventstore
// the resource owner of the aggregate is only used if it's the first event of this aggregate type
// afterwards the resource owner of the first previous events is taken
func NewBaseEventForPush(ctx context.Context, aggregate *Aggregate, typ EventType) *BaseEvent {
	return &BaseEvent{
		Agg:       aggregate,
		User:      authz.GetCtxData(ctx).UserID,
		Service:   service.FromContext(ctx),
		EventType: typ,
	}
}

func (*BaseEvent) Fields() []*FieldOperation {
	return nil
}
