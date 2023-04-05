package deviceauth

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix   eventstore.EventType = "device.authorization."
	AddedEventType                         = eventTypePrefix + "added"
	DeniedEventType                        = eventTypePrefix + "denied"
	ApprovedEventType                      = eventTypePrefix + "approved"
	RemovedEventType                       = eventTypePrefix + "removed"
)

type AddedEvent struct {
	*eventstore.BaseEvent

	ClientID   string
	DeviceCode string
	UserCode   string
	Expires    time.Time
	Scopes     []string
}

func (e *AddedEvent) Data() any {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return NewAddUniqueConstraints(e.ClientID, e.DeviceCode, e.UserCode)
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID string,
	deviceCode string,
	userCode string,
	expires time.Time,
	scopes []string,
) *AddedEvent {
	return &AddedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, AddedEventType,
		),
		clientID, deviceCode, userCode, expires, scopes}
}

type ApprovedEvent struct {
	*eventstore.BaseEvent

	Subject string
}

func (e *ApprovedEvent) Data() any {
	return e
}

func (e *ApprovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewApprovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	subject string,
) *ApprovedEvent {
	return &ApprovedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, ApprovedEventType,
		),
		subject,
	}
}

type DeniedEvent struct {
	*eventstore.BaseEvent
}

func (e *DeniedEvent) Data() any {
	return e
}

func (e *DeniedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDeniedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *DeniedEvent {
	return &DeniedEvent{eventstore.NewBaseEventForPush(ctx, aggregate, DeniedEventType)}
}

type RemovedEvent struct {
	*eventstore.BaseEvent

	ClientID   string
	DeviceCode string
	UserCode   string
}

func (e *RemovedEvent) Data() any {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return NewRemoveUniqueConstraints(e.ClientID, e.DeviceCode, e.UserCode)
}

func NewRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID, deviceCode, userCode string,
) *RemovedEvent {
	return &RemovedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, ApprovedEventType,
		),
		clientID, deviceCode, userCode,
	}
}
