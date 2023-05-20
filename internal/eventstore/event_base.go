package eventstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

// var _ eventstore.Event = (*BaseEvent)(nil)

// BaseEvent represents the minimum metadata of an event
type BaseEvent struct {
	EventType EventType `json:"-"`

	aggregate Aggregate

	sequence                      uint64
	creationDate                  time.Time
	previousAggregateSequence     uint64
	previousAggregateTypeSequence uint64

	//User who created the event
	User string `json:"-"`
	//Service which created the event
	Service string `json:"-"`
	Data    []byte `json:"-"`
}

// EditorService implements Command
func (e *BaseEvent) EditorService() string {
	return e.Service
}

// EditorUser implements Command
func (e *BaseEvent) EditorUser() string {
	return e.User
}

func (e *BaseEvent) Creator() string {
	return e.EditorUser()
}

// Type implements Command
func (e *BaseEvent) Type() EventType {
	return e.EventType
}

// Sequence is an upcounting unique number of the event
func (e *BaseEvent) Sequence() uint64 {
	return e.sequence
}

// CreationDate is the the time, the event is inserted into the eventstore
func (e *BaseEvent) CreationDate() time.Time {
	return e.creationDate
}

// CreationDate is the the time, the event is inserted into the eventstore
func (e *BaseEvent) CreatedAt() time.Time {
	return e.CreationDate()
}

// Aggregate represents the metadata of the event's aggregate
func (e *BaseEvent) Aggregate() Aggregate {
	return e.aggregate
}

// Data returns the payload of the event. It represent the changed fields by the event
func (e *BaseEvent) DataAsBytes() []byte {
	return e.Data
}

// PreviousAggregateSequence implements EventReader
func (e *BaseEvent) PreviousAggregateSequence() uint64 {
	return e.previousAggregateSequence
}

// PreviousAggregateTypeSequence implements EventReader
func (e *BaseEvent) PreviousAggregateTypeSequence() uint64 {
	return e.previousAggregateTypeSequence
}

func (*BaseEvent) Revision() uint16 {
	return 0
}

func (e *BaseEvent) Unmarshal(ptr any) error {
	return json.Unmarshal(e.Data, ptr)
}

// BaseEventFromRepo maps a stored event to a BaseEvent
func BaseEventFromRepo(event *repository.Event) *BaseEvent {
	return &BaseEvent{
		aggregate: Aggregate{
			ID:            event.AggregateID,
			Type:          event.AggregateType,
			ResourceOwner: event.ResourceOwner.String,
			InstanceID:    event.InstanceID,
			Version:       eventstore.Version(event.Version),
		},
		EventType:                     EventType(event.Type),
		creationDate:                  event.CreationDate,
		sequence:                      event.Sequence,
		previousAggregateSequence:     event.PreviousAggregateSequence,
		previousAggregateTypeSequence: event.PreviousAggregateTypeSequence,
		Service:                       event.EditorService,
		User:                          event.EditorUser,
		Data:                          event.Data,
	}
}

// NewBaseEventForPush is the constructor for event's which will be pushed into the eventstore
// the resource owner of the aggregate is only used if it's the first event of this aggregate type
// afterwards the resource owner of the first previous events is taken
func NewBaseEventForPush(ctx context.Context, aggregate *Aggregate, typ EventType) *BaseEvent {
	return &BaseEvent{
		aggregate: Aggregate(*aggregate),
		User:      authz.GetCtxData(ctx).UserID,
		Service:   service.FromContext(ctx),
		EventType: typ,
	}
}
