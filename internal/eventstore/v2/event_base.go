package eventstore

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/service"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

//BaseEvent represents the minimum metadata of an event
type BaseEvent struct {
	EventType EventType

	aggregate Aggregate

	sequence     uint64
	creationDate time.Time

	//User is the user who created the event
	User string `json:"-"`
	//Service is the service which created the event
	Service string `json:"-"`
}

// EditorService implements EventPusher
func (e *BaseEvent) EditorService() string {
	return e.Service
}

//EditorUser implements EventPusher
func (e *BaseEvent) EditorUser() string {
	return e.User
}

//Type implements EventPusher
func (e *BaseEvent) Type() EventType {
	return e.EventType
}

//Sequence is an upcounting unique number of the event
func (e *BaseEvent) Sequence() uint64 {
	return e.sequence
}

//CreationDate is the the time, the event is inserted into the eventstore
func (e *BaseEvent) CreationDate() time.Time {
	return e.creationDate
}

//Aggregate represents the metadata of the event's aggregate
func (e *BaseEvent) Aggregate() Aggregate {
	return e.aggregate
}

//BaseEventFromRepo maps a stored event to a BaseEvent
func BaseEventFromRepo(event *repository.Event) *BaseEvent {
	return &BaseEvent{
		aggregate: Aggregate{
			ID:            event.AggregateID,
			Typ:           AggregateType(event.AggregateType),
			ResourceOwner: event.ResourceOwner,
			Version:       Version(event.Version),
		},
		EventType:    EventType(event.Type),
		creationDate: event.CreationDate,
		sequence:     event.Sequence,
		Service:      event.EditorService,
		User:         event.EditorUser,
	}
}

//NewBaseEventForPush is the constructor for event's which will be pushed into the eventstore
// the resource owner of the aggregate is only used if it's the first event of this aggregateroot
// afterwards the resource owner of the first previous events is taken
func NewBaseEventForPush(ctx context.Context, aggregate *Aggregate, typ EventType) *BaseEvent {
	svcName := service.FromContext(ctx)
	event := &BaseEvent{
		aggregate: *aggregate,
		User:      authz.GetCtxData(ctx).UserID,
		Service:   svcName,
		EventType: typ,
	}

	return event
}
