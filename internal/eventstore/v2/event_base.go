package eventstore

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/service"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type BaseEvent struct {
	aggregateID   string        `json:"-"`
	aggregateType AggregateType `json:"-"`
	EventType     EventType     `json:"-"`

	resourceOwner     string    `json:"-"`
	aggregateVersion  Version   `json:"-"`
	sequence          uint64    `json:"-"`
	previouseSequence uint64    `json:"-"`
	creationDate      time.Time `json:"-"`

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

func (e *BaseEvent) AggregateID() string {
	return e.aggregateID
}
func (e *BaseEvent) AggregateType() AggregateType {
	return e.aggregateType
}
func (e *BaseEvent) ResourceOwner() string {
	return e.resourceOwner
}
func (e *BaseEvent) AggregateVersion() Version {
	return e.aggregateVersion
}
func (e *BaseEvent) Sequence() uint64 {
	return e.sequence
}
func (e *BaseEvent) PreviousSequence() uint64 {
	return e.previouseSequence
}
func (e *BaseEvent) CreationDate() time.Time {
	return e.creationDate
}

func BaseEventFromRepo(event *repository.Event) *BaseEvent {
	return &BaseEvent{
		aggregateID:      event.AggregateID,
		aggregateType:    AggregateType(event.AggregateType),
		aggregateVersion: Version(event.Version),
		EventType:        EventType(event.Type),
		creationDate:     event.CreationDate,
		sequence:         event.Sequence,
		resourceOwner:    event.ResourceOwner,
		Service:          event.EditorService,
		User:             event.EditorUser,
	}
}

func NewBaseEventForPush(ctx context.Context, typ EventType) *BaseEvent {
	svcName := service.FromContext(ctx)
	return &BaseEvent{
		User:      authz.GetCtxData(ctx).UserID,
		Service:   svcName,
		EventType: typ,
	}
}
