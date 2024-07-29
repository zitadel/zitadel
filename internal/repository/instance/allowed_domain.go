package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	allowedDomainPrefix           = "allowed_domains."
	AllowedDomainAddedEventType   = instanceEventTypePrefix + allowedDomainPrefix + "added"
	AllowedDomainRemovedEventType = instanceEventTypePrefix + allowedDomainPrefix + "removed"
)

type AllowedDomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain"`
}

func (e *AllowedDomainAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAllowedDomainAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowedDomain string,
) *AllowedDomainAddedEvent {
	event := &AllowedDomainAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AllowedDomainAddedEventType,
		),
		Domain: allowedDomain,
	}
	return event
}

func (e *AllowedDomainAddedEvent) Payload() interface{} {
	return e
}

func (e *AllowedDomainAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type AllowedDomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain"`
}

func (e *AllowedDomainRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAllowedDomainRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowedDomain string,
) *AllowedDomainRemovedEvent {
	event := &AllowedDomainRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AllowedDomainRemovedEventType,
		),
		Domain: allowedDomain,
	}
	return event
}

func (e *AllowedDomainRemovedEvent) Payload() interface{} {
	return e
}

func (e *AllowedDomainRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
