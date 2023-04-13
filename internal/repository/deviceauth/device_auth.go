package deviceauth

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix   eventstore.EventType = "device.authorization."
	AddedEventType                         = eventTypePrefix + "added"
	ApprovedEventType                      = eventTypePrefix + "approved"
	CanceledEventType                      = eventTypePrefix + "canceled"
	RemovedEventType                       = eventTypePrefix + "removed"
)

type AddedEvent struct {
	*eventstore.BaseEvent

	ClientID   string
	DeviceCode string
	UserCode   string
	Expires    time.Time
	Scopes     []string
	State      domain.DeviceAuthState
}

func (e *AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
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
		clientID, deviceCode, userCode, expires, scopes, domain.DeviceAuthStateInitiated}
}

type ApprovedEvent struct {
	*eventstore.BaseEvent

	Subject string
}

func (e *ApprovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
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

type CanceledEvent struct {
	*eventstore.BaseEvent
	Reason domain.DeviceAuthCanceled
}

func (e *CanceledEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *CanceledEvent) Data() any {
	return e
}

func (e *CanceledEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCanceledEvent(ctx context.Context, aggregate *eventstore.Aggregate, reason domain.DeviceAuthCanceled) *CanceledEvent {
	return &CanceledEvent{eventstore.NewBaseEventForPush(ctx, aggregate, CanceledEventType), reason}
}

type RemovedEvent struct {
	*eventstore.BaseEvent

	ClientID   string
	DeviceCode string
	UserCode   string
}

func (e *RemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
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
			ctx, aggregate, RemovedEventType,
		),
		clientID, deviceCode, userCode,
	}
}
